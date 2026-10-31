
import httpx

from app.config import settings


class SREAPIClient:

    def __init__(self):
        self.client = httpx.Client(
            base_url=settings.sre_api_url,
            timeout=settings.sre_api_timeout,
        )

    def get(self, endpoint: str, params: dict | None = None):
        try:
            response = self.client.get(
                endpoint,
                params=params,
            )

            response.raise_for_status()

            return response

        except httpx.TimeoutException as exc:
            raise RuntimeError(
                "SRE API request timed out"
            ) from exc

        except httpx.HTTPStatusError as exc:
            raise RuntimeError(
                f"SRE API returned HTTP {exc.response.status_code}: "
                f"{exc.response.text[:500]}"
            ) from exc

        except httpx.RequestError as exc:
            raise RuntimeError(
                f"Unable to communicate with SRE API: {exc}"
            ) from exc

    def get_json(self, endpoint: str, params: dict | None = None):
        response = self.get(endpoint, params)
        return response.json()

    def get_text(self, endpoint: str, params: dict | None = None):
        response = self.get(endpoint, params)
        return response.text

    def close(self):
        self.client.close()
