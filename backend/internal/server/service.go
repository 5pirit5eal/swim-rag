package server

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/5pirit5eal/swim-rag/internal/config"
	"github.com/5pirit5eal/swim-rag/internal/models"
	"github.com/5pirit5eal/swim-rag/internal/pdf"
	"github.com/5pirit5eal/swim-rag/internal/rag"
	"github.com/go-chi/httplog/v2"
	"github.com/google/uuid"
)

type RAGService struct {
	// Background context for the server
	ctx context.Context
	// Database client used for storing and querying documents
	db *rag.RAGDB
	// Configuration for the RAG server
	cfg config.Config
}

// Initializes a new RAG service with the given configuration.
// It loads the database password from Google Secret Manager and initializes
// the database connection and LLM client.
// It returns a pointer to the RAGService and an error if any occurred during
// initialization.
func NewRAGService(ctx context.Context, cfg config.Config) (*RAGService, error) {
	slog.Info("Initializing RAG server with config", "cfg", slog.AnyValue(cfg))
	db, err := rag.NewGoogleAIStore(ctx, cfg)
	if err != nil {
		return nil, err
	}

	slog.Info("Creating database connection successfully")

	return &RAGService{
		ctx: ctx,
		cfg: cfg,
		db:  db,
	}, nil
}

// Closes the database connection and LLM client.
// It is important to call this method when the service is no longer needed
// to release resources and avoid memory leaks.
func (rs *RAGService) Close() {
	slog.Info("Closing RAG server...")
	if err := rs.db.Store.Close(); err != nil {
		slog.Error("Error closing database connection", "err", err.Error())
	}
	slog.Info("RAG server closed successfully")
}

// DonatePlanHandler handles the HTTP request to donate a training plan to the database.
// It parses the request, stores the documents and their embeddings in the
// database, and responds with a success message.
// @Summary Add a new training plan
// @Description Upload and store a new swim training plan in the RAG system
// @Tags plans
// @Accept json
// @Produce json
// @Param plan body models.DonatePlanRequest true "Training plan data"
// @Success 200 {string} string "Plan added successfully"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /add [post]
func (rs *RAGService) DonatePlanHandler(w http.ResponseWriter, req *http.Request) {
	logger := httplog.LogEntry(req.Context())
	logger.Info("Adding documents to the database...")
	// Parse HTTP request from JSON.

	dpr := &models.DonatePlanRequest{}

	err := models.GetRequestJSON(req, dpr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if the table is filled
	if len(dpr.Table) == 0 {
		http.Error(w, "Table is empty", http.StatusBadRequest)
		return
	}

	desc := &models.Description{}
	// Check if description is empty and generate one if needed
	if dpr.Description == "" || dpr.Title == "" {
		// Generate a description for the plan
		desc, err := rs.db.Client.DescribeTable(req.Context(), &dpr.Table)
		if err != nil {
			logger.Error("Error when generating description with LLM", httplog.ErrAttr(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		switch {
		case dpr.Title != "":
			desc.Title = dpr.Title
			fallthrough
		case dpr.Description != "":
			desc.Text = dpr.Description
		}
	} else {
		// Generate metadata with improve plan
		m, err := rs.db.Client.GenerateMetadata(req.Context(), &models.Plan{Title: dpr.Title, Description: dpr.Description, Table: dpr.Table})
		if err != nil {
			logger.Error("Error when generating metadata with LLM", httplog.ErrAttr(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		desc.Meta = m
	}

	// Create a donated plan
	plan := &models.DonatedPlan{
		UserID:      dpr.UserID,
		PlanID:      uuid.NewString(),
		CreatedAt:   time.Now().Format(time.DateTime),
		Title:       desc.Title,
		Description: desc.Text,
		Table:       dpr.Table,
	}

	// Store the plan in the database
	err = rs.db.AddDonatedPlan(req.Context(), plan, desc.Meta)
	if err != nil {
		logger.Error("Failed to store plan in the database", httplog.ErrAttr(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Scraping completed successfully"))
}

// QueryHandler handles the RAG query request.
// It parses the request, queries the RAG, generating or choosing a plan, and returns the result as JSON.
// @Summary Query training plans
// @Description Query the RAG system for relevant training plans based on input
// @Tags query
// @Accept json
// @Produce json
// @Param query body models.QueryRequest true "Query parameters"
// @Success 200 {object} models.RAGResponse "Query results"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /query [post]
func (rs *RAGService) QueryHandler(w http.ResponseWriter, req *http.Request) {
	logger := httplog.LogEntry(req.Context())
	logger.Info("Querying the database...")
	// Parse HTTP request from JSON.

	qr := &models.QueryRequest{}
	err := models.GetRequestJSON(req, qr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	answer, err := rs.db.Query(req.Context(), qr.Content, qr.Filter, qr.Method)
	if err != nil {
		if strings.HasPrefix(err.Error(), "unsupported method:") {
			http.Error(w, "Method may only be 'choose' or 'generate', invalid choice.", http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Recalculate the sums of the rows to be sure they are correct
	answer.Table.UpdateSum()
	logger.Debug("Updated the table sums...", "sum", answer.Table[len(answer.Table)-1].Sum)

	logger.Info("Answer generated successfully")
	models.WriteResponseJSON(w, http.StatusOK, answer)
}

// PlanToPDFHandler handles the Plan to PDF export request.
// @Summary Export training plan to PDF
// @Description Generate and download a PDF version of a training plan
// @Tags export
// @Accept json
// @Produce json
// @Param plan body models.PlanToPDFRequest true "Training plan data to export"
// @Success 200 {object} models.PlanToPDFResponse "PDF export response with URI"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /export-pdf [post]
func (rs *RAGService) PlanToPDFHandler(w http.ResponseWriter, req *http.Request) {
	logger := httplog.LogEntry(req.Context())
	logger.Info("Exporting table to PDF...")

	// Parse HTTP request from JSON.
	qr := &models.PlanToPDFRequest{}
	err := models.GetRequestJSON(req, qr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert the table to PDF
	planPDF, err := pdf.PlanToPDF(&models.Plan{
		Title:       qr.Title,
		Description: qr.Description,
		Table:       qr.Table,
	})
	if err != nil {
		logger.Error("Table generation failed", httplog.ErrAttr(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Upload the PDF to cloud storage
	uri, err := pdf.UploadPDF(req.Context(), rs.cfg.Bucket.Name, pdf.GenerateFilename(), planPDF)
	if err != nil {
		logger.Error("PDF upload failed", httplog.ErrAttr(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	answer := &models.PlanToPDFResponse{URI: uri}

	logger.Info("Answer generated successfully")
	models.WriteResponseJSON(w, http.StatusOK, answer)
}
