# Swim Training Plan Generator - Backend API Documentation

**Note:** The frontend for this application is currently To Be Determined (TBD). This API documentation outlines the backend endpoints available for future frontend integration or direct API consumption.

All endpoints are relative to the base path of the service.

---

## Endpoints

### 1. Donate Training Plan

*   **Method:** `POST`
*   **Path:** `/add`
*   **Description:** Allows a user to donate a training plan to the system. If the title or description is not provided, the system will attempt to generate them based on the table content. Metadata for the plan will also be generated.
*   **Request Body:** `application/json`
    ```json
    {
        "user_id": "string (optional, for V1 anonymous)",
        "title": "string (optional)",
        "description": "string (optional)",
        "table": [
            {
                "Amount": "integer (Amount of repetitions)",
                "Multiplier": "string (Multiplier for the distance, e.g., 'x' or 'times')",
                "Distance": "integer (Distance in meters)",
                "Break": "string (Break time, typically in seconds)",
                "Content": "string (Content or description of the row)",
                "Intensity": "string (Intensity level of the activity)",
                "Sum": "integer (Total volume or sum for the row, usually Amount * Distance)"
            }
        ]
    }
    ```
*   **Responses:**
    *   `200 OK`: "Scraping completed successfully" (Note: The success message seems to be a misnomer from the code, it should indicate successful plan donation).
    *   `400 Bad Request`: If the request is malformed (e.g., invalid JSON, empty table).
    *   `500 Internal Server Error`: If an error occurs during processing (e.g., LLM interaction failure, database error).
*   **Example Request:**
    ```json
    {
        "user_id": "user123",
        "title": "My Awesome Swim Plan",
        "description": "A great plan for intermediate swimmers.",
        "table": [
            {
                "Amount": 4,
                "Multiplier": "x",
                "Distance": 100,
                "Break": "30s",
                "Content": "Freestyle",
                "Intensity": "Moderate",
                "Sum": 400
            },
            {
                "Amount": 1,
                "Multiplier": "",
                "Distance": 200,
                "Break": "60s",
                "Content": "Cool down",
                "Intensity": "Easy",
                "Sum": 200
            }
        ]
    }
    ```

---

### 2. Query for Training Plan

*   **Method:** `POST`
*   **Path:** `/query`
*   **Description:** Queries the system for a training plan based on user content. The method can be 'generate' (to create a new plan) or 'choose' (to select from existing plans).
*   **Request Body:** `application/json`
    ```json
    {
        "content": "string (User's query or description of desired plan)",
        "filter": "object (optional, key-value pairs for filtering, e.g., {\"difficulty\": \"easy\"})",
        "method": "string (optional, 'generate' or 'choose'. Defaults to 'generate' if not specified or invalid)"
    }
    ```
*   **Response Body (`200 OK`):** `application/json`
    ```json
    {
        "title": "string (Title of the generated/chosen plan)",
        "description": "string (Description of the plan)",
        "table": [ // Same Row structure as in /add request
            {
                "Amount": "integer",
                "Multiplier": "string",
                "Distance": "integer",
                "Break": "string",
                "Content": "string",
                "Intensity": "string",
                "Sum": "integer"
            }
        ]
    }
    ```
*   **Other Responses:**
    *   `400 Bad Request`: If the request is malformed or an unsupported method is provided.
    *   `500 Internal Server Error`: If an error occurs during processing.
*   **Example Request:**
    ```json
    {
        "content": "I want a 2000m swimming plan for beginners focusing on freestyle.",
        "method": "generate"
    }
    ```

---

### 3. Scrape URL for Training Plan Data

*   **Method:** `GET`
*   **Path:** `/scrape`
*   **Description:** Scrapes a given URL for training plan data and attempts to store it in the database.
*   **Query Parameters:**
    *   `url` (string, required): The URL to scrape.
*   **Responses:**
    *   `200 OK`: "Scraping completed successfully"
    *   `400 Bad Request`: If the `url` parameter is missing.
    *   `500 Internal Server Error`: If an error occurs during scraping or database interaction.
*   **Example Request:**
    `GET /scrape?url=https://example.com/my-swim-plan`

---

### 4. Export Training Plan to PDF

*   **Method:** `POST`
*   **Path:** `/export-pdf`
*   **Description:** Converts a given training plan (title, description, table) into a PDF and uploads it to cloud storage, returning the URI of the stored PDF.
*   **Request Body:** `application/json`
    ```json
    {
        "title": "string",
        "description": "string",
        "table": [ // Same Row structure as in /add request
            {
                "Amount": "integer",
                "Multiplier": "string",
                "Distance": "integer",
                "Break": "string",
                "Content": "string",
                "Intensity": "string",
                "Sum": "integer"
            }
        ]
    }
    ```
*   **Response Body (`200 OK`):** `application/json`
    ```json
    {
        "uri": "string (URI of the generated PDF in cloud storage)"
    }
    ```
*   **Other Responses:**
    *   `400 Bad Request`: If the request is malformed.
    *   `500 Internal Server Error`: If an error occurs during PDF generation or upload.
*   **Example Request:**
    ```json
    {
        "title": "My Exported Plan",
        "description": "This is a plan to be exported.",
        "table": [
            {
                "Amount": 8,
                "Multiplier": "x",
                "Distance": 50,
                "Break": "15s",
                "Content": "Kick with board",
                "Intensity": "Easy",
                "Sum": 400
            }
        ]
    }
    ```

---

### 5. Health Check

*   **Method:** `GET`
*   **Path:** `/health`
*   **Description:** Standard health check endpoint.
*   **Responses:**
    *   `200 OK`: "OK"
