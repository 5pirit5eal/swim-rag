package genai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/5pirit5eal/swim-rag/internal/models"
	"github.com/go-chi/httplog/v2"
	"github.com/tmc/langchaingo/schema"
	"google.golang.org/genai"
)

// GeneratePlan generates a plan using the LLM based on the provided query and documents.
func (gc *GoogleGenAIClient) GeneratePlan(ctx context.Context, q string, docs []schema.Document) (*models.RAGResponse, error) {
	logger := httplog.LogEntry(ctx)
	ts, err := models.TableSchema()
	if err != nil {
		return nil, fmt.Errorf("failed to get table schema: %w", err)
	}

	var dc []string
	for _, doc := range docs {
		dc = append(dc, doc.PageContent)
	}

	// Create a RAG query for the LLM with the most relevant documents as context
	query := fmt.Sprintf(ragTemplateStr, ts, q, strings.Join(dc, "\n \n"))
	genCfg := *gc.gcfg
	genCfg.ResponseMIMEType = "application/json"
	answer, err := gc.gc.Models.GenerateContent(ctx, gc.cfg.Model, genai.Text(query), &genCfg)

	if err != nil {
		logger.Error("Error when generating answer with LLM", httplog.ErrAttr(err))
		return nil, fmt.Errorf("error generating answer: %w", err)
	}

	// read description and table from the LLM response
	var p models.RAGResponse
	err = json.Unmarshal([]byte(answer.Text()), &p)
	if err != nil {
		logger.Error("Error parsing LLM response", httplog.ErrAttr(err), "raw_response", answer)
		return nil, fmt.Errorf("error parsing LLM response: %w", err)
	}
	// Add the total to the table if it is not already present
	if !strings.Contains(p.Table[len(p.Table)-1].Content, "Gesamt") {
		p.Table.AddSum()
	}
	// Recalculate the sums of the rows to be sure they are correct
	p.Table.UpdateSum()

	// Add the plan to the response
	logger.Debug("Plan generated successfully")
	return &p, nil
}

// ChoosePlan lets an LLM choose the best fitting plan from the given documents.
// Returns the plan id of the chosen plan
func (gc *GoogleGenAIClient) ChoosePlan(ctx context.Context, q string, docs []schema.Document) (string, error) {
	logger := httplog.LogEntry(ctx)
	var dc string
	for i, doc := range docs {
		dc += fmt.Sprintf("%d: %s \n\n", i, doc.PageContent)
	}

	// Create a RAG query for the LLM with the most relevant documents as context
	query := fmt.Sprintf(choosePlanTemplateStr, q, dc)
	genCfg := *gc.gcfg
	genCfg.ResponseMIMEType = "application/json"
	answer, err := gc.gc.Models.GenerateContent(ctx, gc.cfg.Model, genai.Text(query), &genCfg)
	if err != nil {
		logger.Error("Error when generating answer with LLM", httplog.ErrAttr(err))
		return "", fmt.Errorf("error generating answer: %w", err)
	}
	logger.Debug("Successful answer from LLM", "answer", answer)

	var cr models.ChoiceResult
	err = json.Unmarshal([]byte(answer.Text()), &cr)
	if err != nil {
		logger.Error("Error parsing LLM response", httplog.ErrAttr(err), "raw_response", answer)
		return "", fmt.Errorf("error parsing LLM response: %w", err)
	}
	planID, ok := docs[cr.Idx].Metadata["plan_id"]
	if !ok {
		return "", fmt.Errorf("plan_id not found in Metadata for document at index %d", cr.Idx)
	}
	planIDStr, ok := planID.(string)
	if !ok {
		return "", fmt.Errorf("plan_id is not a string in Metadata for document at index %d", cr.Idx)
	}
	return planIDStr, nil
}

func (gc *GoogleGenAIClient) ImprovePlan(ctx context.Context, plan models.Planable, syncGroup *sync.WaitGroup, c chan<- models.Document, ec chan<- error) {
	if syncGroup != nil {
		defer syncGroup.Done()
	}
	logger := httplog.LogEntry(ctx)
	meta, err := gc.GenerateMetadata(ctx, plan)
	if err != nil {
		logger.Error("Error when generating metadata with LLM", httplog.ErrAttr(err))
		ec <- fmt.Errorf("error generating metadata: %w", err)
		return
	}

	// Create request body by converting the plans into documents
	c <- models.Document{
		Plan: plan,
		Meta: meta,
	}
}

func (gc *GoogleGenAIClient) DescribeTable(ctx context.Context, table *models.Table) (*models.Description, error) {
	logger := httplog.LogEntry(ctx)
	ds, err := models.DescriptionSchema()
	if err != nil {
		logger.Error("Failed in retrieving Schema", httplog.ErrAttr(err))
		return nil, fmt.Errorf("models.MetadataSchema: %w", err)
	}
	// Create a description of the table
	query := fmt.Sprintf(describeTemplateStr, ds, table.String())
	genCfg := *gc.gcfg
	genCfg.ResponseMIMEType = "application/json"
	answer, err := gc.gc.Models.GenerateContent(ctx, gc.cfg.Model, genai.Text(query), &genCfg)
	if err != nil {
		return nil, fmt.Errorf("Models.GenerateContent: %w", err)
	}
	var desc models.Description
	err = json.Unmarshal([]byte(answer.Text()), &desc)
	if err != nil {
		logger.Error("Error parsing LLM response", httplog.ErrAttr(err), "raw_response", answer.Text())
		return nil, fmt.Errorf("error parsing LLM response: %w", err)
	}
	return &desc, nil
}

func (gc *GoogleGenAIClient) GenerateMetadata(ctx context.Context, plan models.Planable) (*models.Metadata, error) {
	logger := httplog.LogEntry(ctx)
	ms, err := models.MetadataSchema()
	if err != nil {
		logger.Error("Failed in retrieving Schema", httplog.ErrAttr(err))
		return nil, fmt.Errorf("models.MetadataSchema: %w", err)
	}
	// Enhance scraped documents with gemini and create meaningful metadata
	genericPlan := plan.Plan()
	query := fmt.Sprintf(metadataTemplateStr, genericPlan.Title, genericPlan.Description, genericPlan.Table.String(), ms)
	genCfg := *gc.gcfg
	genCfg.ResponseMIMEType = "application/json"
	answer, err := gc.gc.Models.GenerateContent(ctx, gc.cfg.Model, genai.Text(query), &genCfg)
	if err != nil {
		logger.Error("Error when generating answer with LLM", httplog.ErrAttr(err))
		return nil, fmt.Errorf("Models.GenerateContent: %w", err)
	}
	logger.Debug("Successful answer from LLM", "answer", answer.Text())

	// Parse the answer as JSON
	var metadata models.Metadata
	err = json.Unmarshal([]byte(answer.Text()), &metadata)
	if err != nil {
		logger.Error("Error parsing LLM response", httplog.ErrAttr(err), "raw_response", answer.Text())
		return nil, fmt.Errorf("JSON unmarshal error: %w with raw response %s", err, answer.Text())
	}

	return &metadata, nil
}
