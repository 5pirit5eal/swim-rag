# System Patterns: Swim Training Plan Generator

## System Architecture

The application follows a typical three-tier web architecture, hosted on Google Cloud Platform.

```plaintext
+---------------------+                                   gRPC/HTTPS                                +---------------------+
| User's Browser      |<--------------------------------------------------------------------------->| Cloud Run Backend   | <--> PostgreSQL
| (Vue.js)            |                                                                             | (Go)                |
|                     |     HTTPS (with IAP)      +---------------------+                           +---------------------+
+---------------------+<------------------------> | Cloud Run Frontend  |                                       |
          ^                                       | (with IAP)          |                                       | (Verifies ID Token)
          | <gcp-api>                             +---------------------+                                       |
          | web-component                                    |                                                  |
          |                                                  | (Passes ID Token in                              |
          |                                                  |  Authorization Header)                           |
          |                                                  |                                                  |
          +--------------------------------------------------+--------------------------------------------------+
```

### Components

1.  **User's Browser (Frontend Client):**
    *   Implemented using Vue.js.
    *   Interacts with the Cloud Run Frontend service.
    *   Handles user interface and client-side logic.

2.  **Cloud Run Frontend (Web Server):**
    *   Hosts the Single-Page Application (SPA) built with Vue.js.
    *   Protected by Identity-Aware Proxy (IAP) for authentication (planned for V2+).
    *   Communicates with the Cloud Run Backend, passing ID tokens in the Authorization header for authenticated requests (planned for V2+).

3.  **Cloud Run Backend (API Server):**
    *   Implemented in Go, utilizing Langchaingo for AI/LLM interactions.
    *   Provides gRPC/HTTPS endpoints for the frontend and potentially other clients.
    *   Handles business logic, data processing, and interaction with the database.
    *   Verifies ID tokens for authenticated requests (planned for V2+).

4.  **PostgreSQL Database (Data Store):**
    *   Managed PostgreSQL service on Google Cloud.
    *   Uses the `pgvector` extension for storing and querying embeddings related to training plans.
    *   Stores training plan data, user data (planned for V2+), and other application-specific information.

## Key Technical Decisions

*   **Microservices-like Architecture:** Separation of frontend and backend into distinct Cloud Run services allows for independent scaling and deployment.
*   **Serverless Compute:** Utilization of Cloud Run for both frontend and backend leverages serverless benefits like auto-scaling and pay-per-use.
*   **Go for Backend:** Chosen for its performance, concurrency features, and suitability for building robust APIs. Langchaingo facilitates integration with Generative AI models.
*   **Vue.js for Frontend:** A progressive JavaScript framework for building user interfaces.
*   **PostgreSQL with pgvector:** Provides robust relational data storage combined with efficient vector similarity search capabilities, crucial for Retrieval Augmented Generation (RAG) of training plans.
*   **Identity-Aware Proxy (IAP):** Google Cloud's solution for managing access to applications based on user identity, simplifying authentication and authorization (planned for V2+).
*   **Infrastructure as Code (IaC):** Terraform is used for managing cloud infrastructure (evident from `deployments/0-infra/` and `deployments/deploy/` directories).
*   **Containerization:** The backend is containerized using Docker (evident from `backend/Dockerfile`), facilitating consistent deployments on Cloud Run.

## Component Relationships & Data Flow

1.  **User Interaction:** The user interacts with the Vue.js application in their browser.
2.  **Frontend to Backend:** The Vue.js app makes HTTPS requests to the Cloud Run Frontend service. This service, in turn, makes gRPC/HTTPS requests to the Cloud Run Backend service. For authenticated actions (V2+), an ID token obtained via IAP is passed in the `Authorization` header.
3.  **Backend Logic:** The Go backend processes requests, interacts with the PostgreSQL database (e.g., for RAG using embeddings stored with `pgvector`), and potentially calls external Generative AI services via Langchaingo.
4.  **Data Persistence:** Data, including training plans, embeddings, and user information (V2+), is stored in the PostgreSQL database.

## Future Considerations (from Roadmap)

*   **V2:** Introduction of user authentication (IAP), user-specific data, history, feedback mechanisms, and expanded export options.
*   **V3:** Multimodal input (PDF, images), community features (sharing, boards).
*   **V4:** User statistics dashboard.
