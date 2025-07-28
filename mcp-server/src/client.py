import asyncio

from fastmcp import Client
from fastmcp.client.transports import StreamableHttpTransport

from swim_rag_mcp.schemas import QueryRequest

transport = StreamableHttpTransport(url="http://localhost:8000/mcp")
client = Client(transport)


async def call_tool():
    async with client:
        result = await client.call_tool(
            "generate_or_choose_plan",
            dict(
                query=QueryRequest(
                    content="Was ist der beste Trainingsplan?",
                    filter=None,
                    method="generate",
                ).model_dump(),
            ),
        )
        print(result)


asyncio.run(call_tool())
