import asyncio
import logging
import os

import httpx
from dotenv import load_dotenv
from fastmcp import FastMCP
from fastmcp.exceptions import ToolError
from fastmcp.server.middleware.error_handling import RetryMiddleware
from fastmcp.server.middleware.logging import StructuredLoggingMiddleware
from fastmcp.server.middleware.rate_limiting import (
    RateLimitingMiddleware,
    SlidingWindowRateLimitingMiddleware,
)

from swim_rag_mcp.schemas import ExportResponse, QueryRequest, QueryResponse
from swim_rag_mcp.utils import get_id_token

load_dotenv(".config.env")

URL = os.getenv("SWIM_RAG_API_URL", "http://localhost:8080")
print(f"Using Swim RAG API URL: {URL}")

mcp: FastMCP = FastMCP(
    name="swim-rag-mcp",
    instructions="""
        This is the MCP Server connected to the Swim RAG backend, an application meant for generating and 
        exporting german training plans for swimming. It allows the user to query for a personalized training plan,
        edit it and send the edited plan to the Swim RAG backend for export to a PDF file.
    """,
    exclude_tags={"internal"},
    include_tags={"public"},
    on_duplicate_tools="error",  # Handle duplicate registrations
    on_duplicate_resources="warn",
    on_duplicate_prompts="replace",
    middleware=[
        StructuredLoggingMiddleware(
            include_payloads=True, log_level=logging.INFO
        ),
        RateLimitingMiddleware(burst_capacity=20),
        SlidingWindowRateLimitingMiddleware(
            max_requests=100,
            window_minutes=1,
        ),
        RetryMiddleware(max_retries=3),
    ],
)


@mcp.tool(tags={"public"})
async def generate_or_choose_plan(query: QueryRequest) -> QueryResponse:
    """Query the Swim RAG system with a given german query string.
    It parses the request, queries the RAG, generating or choosing a plan, and returns the result as JSON.
    """
    # Send the request to the Swim RAG backend
    try:
        response = httpx.post(
            url=URL + "/query",
            json=query.model_dump(),
            timeout=60.0,  # Set a timeout for the request
            headers=await get_id_token(URL),  # Get the auth token if available,
        )
        response.raise_for_status()  # Raise an error for bad responses
    except httpx.RequestError as e:
        raise ToolError(f"Request error: {e}")
    except httpx.HTTPStatusError as e:
        raise ToolError(
            f"HTTP error: {e.response.status_code} - {e.response.text}"
        )
    except Exception as e:
        raise ToolError(f"An unexpected error occurred: {e}")
    # Return the response
    return QueryResponse.model_validate_json(response.text)


@mcp.tool(tags={"public"})
async def export_plan(plan: QueryResponse) -> ExportResponse:
    """Export a plan as a PDF file for easier printing and sharing."""
    try:
        response = httpx.post(
            url=URL + "/export-pdf",
            json=plan.model_dump(),
            timeout=60.0,  # Set a timeout for the request
            headers=await get_id_token(URL),  # Get the auth token if available
        )
        response.raise_for_status()  # Raise an error for bad responses
    except httpx.RequestError as e:
        raise ToolError(f"Request error: {e}")
    except httpx.HTTPStatusError as e:
        raise ToolError(
            f"HTTP error: {e.response.status_code} - {e.response.text}"
        )
    except Exception as e:
        raise ToolError(f"An unexpected error occurred: {e}")

    # Return the response
    return ExportResponse.model_validate_json(response.text)


@mcp.tool(tags={"internal"})
async def scrape_plans_from_web(url: str) -> str:
    """Scrape plans from a given URL.

    Args:
        url (str): The URL to scrape plans from.

    Returns:
        str: Confirmation message of the scraping operation.
    """
    # Here you would implement the logic to scrape plans from
    return f"Plans scraped from {url} successfully."


if __name__ == "__main__":
    asyncio.run(mcp.run_async(transport="http"))
