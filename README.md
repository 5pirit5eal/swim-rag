# Swim Training Plan Generator

A web application for generating, recommeding and sharing training plans for swimmers, triathletes and trainers.

## Tech Stack

The web application is hosted in Google Cloud. The following components are involved:

1. Identity-Aware Proxy (IAP) to authenticate and identify users
2. Single-Page Frontend hosted in Cloud Run implemented in React
3. Backend hosted in Cloud Run implemented in Go and Langchaingo
4. PostgreSQL database (managed) with pgvector

```plaintext
+---------------------+                                   gRPC/HTTPS                                +---------------------+
| User's Browser      |<--------------------------------------------------------------------------->| Cloud Run Backend   | <--> PostgreSQL
| (React)             |                                                                             | (Go)                |
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

## Roadmap

### V1: The Prototype

The initial version only contains the core concept of being able to generate individual training plans based on user input (text/form).

- **Generate TP based on user input (free form or configuration menu)**
  - Frontend displays simple one page app with text box and advanced settings
  - Backend provides endpoint for retrieval-augmented generation
  - Database contains training plans scraped from the web (tiny scope) and their embeddings
- **TP are exportable/printable as PDF**
  - Frontend includes Button for exporting PDF
  - Backend allows the user to get PDF exports of Trainingplan data
- **Requests are anonym and no user identification is necessary**
  - Authentication is omitted for v1
  - Any user information is stored encrypted if not configured otherwise
  - Frontend displays the option to not donate results in advanced settings
- **Backend allows user to donate training plans**
  - Backend-only implementation of training plan donation
  - Backend-only implementation of feedback

### V2: User-specific TPs and TP History

The next version takes lessons learned in v1 into consideration and adds user authentication and generation history, making the web application a lot more useful.

- **Generate TP based on user input (free form, configuration) and enable the user to give feedback to adapt results**
  - Builds upon the content of v1 by adding chat history in backend
- **TP are exportable as PDF, excel or odt**
- **TP can be recommended/donated**
  - Form added to frontend to input/donate your favorite training plans
  - URL input to report new scraping targets (after allowListing and configuration)
  - Request user feedback to plan after new login after previously exporting TP
  - Endpoint in backend for saving feedback
  - Notation/Abbreviations of different sources can be processed correctly
- **User authentication and identification**
  - Identity-Aware Proxy for user identification and user data persistence (TPs)
- **TP History can be viewed**
  - TPs are saved associated with user ID (JWT)
  - Frontend is extended with `History` page to view previously created TPs
- **TP can be shared via URL**
- **Enable deletion of user data**
  - Add url to soft-delete all user data at the bottom of the page
  - User data is retrievable or 30 days, then fully deleted
- **Add content required by german/eu law**
  - Impressum

### V3: Multimodal Generation and Community-Sharing

- **Generate TP based on multimodal user input (pdf, image, hand-written notes)**
  - The third version further builds upon v1 and v2 by adding more input options for generating or donating TPs
- **Searchable Community Board for most popular TP**
  - A searchable community forum of previously created and liked TPs.
  - This also allows user interaction, such as comments, likes and tagging.
  - Optional sign-up for a leaderboard of most active users
- **Add anonymous/guest mode**
  - An anonymous mode is added which enables new users to test the application in a similar way to how v1 works.

### V4: Usage Statistics Dashboard

- **Add usage statistics dashboard for each signed-in user**
  - Total distance (km) swam in plans
  - Number of plans exported
  - Number of plans donated
  - Comparative statistics with other user behaviour
