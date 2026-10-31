from langchain.tools import tool

from app.clients.sre_api_client import SREAPIClient


@tool
def list_pods(namespace: str) -> dict:
    """
    List all Kubernetes pods in the specified namespace.

    Use this tool when you need to inspect pods,
    check their status, or investigate pod-related issues.

    Args:
        namespace: Kubernetes namespace to list pods from.
    """

    client = SREAPIClient()

    try:
        response = client.get_json(
            f"/{namespace}/pods"
        )

        return response

    finally:
        client.close()