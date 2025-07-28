package models

// DonatePlanRequest represents the request body for donating a training plan
// @Description Request payload for donating a swim training plan to the system
type DonatePlanRequest struct {
	UserID      string `json:"user_id" example:"user123" binding:"required"`
	Title       string `json:"title,omitempty" example:"Advanced Freestyle Training"`
	Description string `json:"description,omitempty" example:"A comprehensive training plan for improving freestyle technique"`
	Table       Table  `json:"table" binding:"required"`
	// v3: add other table modalities
	// Image 	 string `json:"image,omitempty"`
	// URI 		 string `json:"uri,omitempty"`
}

// QueryRequest represents the request body for querying the RAG system
// @Description Request payload for querying swim training plans from the RAG system
type QueryRequest struct {
	Content string         `json:"content" example:"I need a training plan for improving my freestyle technique" binding:"required"` // Content describes what kind of training plan is needed
	Filter  map[string]any `json:"filter,omitempty"`                                                                                 // Filter allows filtering plans by metadata like difficulty or stroke type
	Method  string         `json:"method,omitempty" example:"generate" validate:"oneof=choose generate"`                             // Method can be either 'choose' (select existing plan) or 'generate' (create new plan)
}

// RAGResponse represents the response after a query to the RAG system
// @Description Response containing a generated or selected swim training plan
type RAGResponse struct {
	Title       string `json:"title" example:"Advanced Freestyle Training"`
	Description string `json:"description" example:"A comprehensive training plan for improving freestyle technique"`
	Table       Table  `json:"table"`
}

// ChooseResponse represents the response when a plan is chosen rather than generated
// @Description Response when selecting an existing plan from the system
type ChooseResponse struct {
	Idx         int    `json:"index" example:"1"`
	Description string `json:"description" example:"Selected plan based on your requirements"`
}

// PlanToPDFRequest represents the request for PDF export
// @Description Request payload for exporting a training plan to PDF format
type PlanToPDFRequest struct {
	Title       string `json:"title" example:"Advanced Freestyle Training" binding:"required"`
	Description string `json:"description" example:"A comprehensive training plan for improving freestyle technique" binding:"required"`
	Table       Table  `json:"table" binding:"required"`
}

// PlanToPDFResponse represents the response from PDF export
// @Description Response containing the URI to the generated PDF file
type PlanToPDFResponse struct {
	URI string `json:"uri" example:"https://storage.googleapis.com/bucket/plans/plan_123.pdf"`
}
