from typing import Literal, Optional

from pydantic import BaseModel, Field


class Filter(BaseModel):
    """Represents the filter criteria for querying SWIM RAG."""

    freistil: Optional[bool] = Field(
        None, description="Filter for freestyle swimming technique"
    )
    brust: Optional[bool] = Field(
        None, description="Filter for breaststroke swimming technique"
    )
    ruecken: Optional[bool] = Field(
        None, description="Filter for backstroke swimming technique"
    )
    delfin: Optional[bool] = Field(
        None, description="Filter for butterfly swimming technique"
    )
    lagen: Optional[bool] = Field(
        None, description="Filter for individual medley swimming"
    )
    schwierigkeitsgrad: Optional[
        Literal[
            "Nichtschwimmer",
            "Anfaenger",
            "Fortgeschritten",
            "Leistungsschwimmer",
            "Top-Athlet",
        ]
    ] = Field(None, description="Difficulty level filter for the training plan")
    trainingstyp: Optional[
        Literal[
            "Techniktraining",
            "Leistungstest",
            "Grundlagenausdauer",
            "Recovery",
            "Kurzstrecken",
            "Langstrecken",
            "Atemmangel",
            "Wettkampfvorbereitung",
        ]
    ] = Field(None, description="Type of training session filter")


class QueryRequest(BaseModel):
    """Represents a request to query the Swim RAG backend."""

    content: str = Field(
        ..., description="The query string to be processed by the Swim RAG MCP."
    )
    filter: Optional[Filter] = Field(
        None,
        description="Optional filters to apply to the RAG query.",
    )
    method: Literal["generate", "choose"] = Field(
        "generate",
        description="The method to use for the returned trainingsplan, "
        "either generating a new one or returning an existing from the database.",
    )


class Row(BaseModel):
    """Represents a single row within the training plan data."""

    Amount: int
    Multiplier: str = Field(
        "x", description="Multiplier for the amount, e.g., 'x' for repetitions."
    )
    Distance: int = Field(
        ..., description="Distance for the row, typically in meters."
    )
    Break: str = Field(
        ...,
        description="Break time between sets or activities, typically in seconds or meters.",
    )
    Content: str = Field(
        ..., description="Content or description of the set or activity"
    )
    Intensity: str = Field(..., description="Intensity level of the activity")
    Sum: int = Field(..., description="Total distance or amount for the row.")


class QueryResponse(BaseModel):
    """Represents the response from the Swim RAG backend after processing a query."""

    title: str = Field(
        ...,
        description="The status of the query response, e.g., 'success' or 'error'.",
    )
    description: str = Field(
        ..., description="The data returned from the Swim RAG backend."
    )
    table: list[Row] = Field(
        ...,
        description="Training plan data returned from the Swim RAG backend.",
    )


class ExportResponse(BaseModel):
    """Exported pdf url for a plan."""

    uri: str = Field(..., description="The URL of the exported PDF file.")
