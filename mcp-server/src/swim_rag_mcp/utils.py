import os

import requests
import httpx
from google.oauth2 import _id_token_async
from google.auth.transport import _aiohttp_requests


async def get_id_token(url: str) -> dict[str, str] | None:
    """Get the authorization token for the Swim RAG API.

    Args:
        url (str): The URL of the Swim RAG API.

    Returns:
        httpx.BasicAuth | None: The authorization token or None if not applicable.
    """
    if os.getenv("K_SERVICE"):
        # Add authorization headers from the service account in env as the service runs in google cloud run
        request = _aiohttp_requests.Request()
        id_token = await _id_token_async.fetch_id_token(request, url)
        auth = {"X-Serverless-Authorization": f"Bearer {id_token}"}
    else:
        # Expect local proxy of cloud run service
        auth = None
    return auth
