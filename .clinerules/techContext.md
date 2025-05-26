# Technical Context: Swim Training Plan Generator

## Technologies Used

*   **Frontend:** Vue.js (Planned, as per `README.md`)
*   **Backend:** Go (Version 1.24.1, as per `backend/go.mod`)
    *   **AI/LLM Integration:** `github.com/tmc/langchaingo`
    *   **Web Framework/Routing:** `github.com/go-chi/chi/v5`
    *   **PDF Generation:** `github.com/johnfercher/maroto/v2`
    *   **Web Scraping:** `github.com/gocolly/colly`
*   **Database:** PostgreSQL (Managed Google Cloud service)
    *   **Vector Extension:** `pgvector` (implied by `README.md` and `github.com/pgvector/pgvector-go` in `go.mod`)
*   **Cloud Platform:** Google Cloud Platform (GCP)
    *   **Compute:** Cloud Run (for both frontend and backend)
    *   **Authentication:** Identity-Aware Proxy (IAP)
    *   **Storage:** Cloud Storage (implied by `cloud.google.com/go/storage` in `go.mod`)
    *   **Secrets Management:** Secret Manager (implied by `cloud.google.com/go/secretmanager` in `go.mod`)
*   **Containerization:** Docker (as per `backend/Dockerfile`)
*   **Infrastructure as Code (IaC):** Terraform (evident from `deployments/` directory structure and `.tf` files)
*   **CI/CD:** Google Cloud Build (evident from `*.cloudbuild.yaml` files in `backend/`)

## Development Setup & Key Dependencies (from `backend/go.mod`)

The backend is a Go module: `github.com/5pirit5eal/swim-rag`.

### Key Backend Dependencies:
*   `cloud.google.com/go/cloudsqlconn`: For connecting to Cloud SQL (PostgreSQL).
*   `cloud.google.com/go/secretmanager`: For managing secrets.
*   `cloud.google.com/go/storage`: For interacting with Google Cloud Storage.
*   `github.com/georgysavva/scany/v2`: For scanning database rows into Go structs.
*   `github.com/go-chi/chi/v5`: HTTP router.
*   `github.com/go-chi/httplog/v2`: HTTP request logging.
*   `github.com/go-chi/render`: JSON and XML response rendering.
*   `github.com/gocolly/colly`: Web scraping framework.
*   `github.com/golobby/dotenv`: For loading environment variables from `.env` files.
*   `github.com/google/uuid`: For generating UUIDs.
*   `github.com/invopop/jsonschema`: For generating JSON schemas from Go types.
*   `github.com/jackc/pgx/v5`: PostgreSQL driver and toolkit.
*   `github.com/johnfercher/maroto/v2`: PDF generation library.
*   `github.com/stretchr/testify`: Testing utilities.
*   `github.com/tmc/langchaingo`: Langchain for Go, for LLM interactions.
*   `google.golang.org/genai`: Google Generative AI SDK.

## Tool Usage Patterns

*   **Backend Development:** Go, Docker for containerization.
*   **Deployment:** Google Cloud Build for CI/CD, Terraform for IaC.
*   **Frontend Development:** Vue.js (TBD).

## Technical Constraints & Considerations

*   The system is designed for Google Cloud Platform.
*   The backend relies on several Google Cloud services (Cloud Run, PostgreSQL, IAP, Secret Manager, Storage).
*   Frontend is not yet implemented.
