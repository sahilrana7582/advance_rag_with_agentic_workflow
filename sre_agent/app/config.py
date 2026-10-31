import os

from dotenv import load_dotenv


load_dotenv()


class Settings:
    def __init__(self):
        self.sre_api_url = os.getenv(
            "SRE_API_URL",
            "http://localhost:8080/api/v1",
        ).rstrip("/")

        self.sre_api_timeout = float(
            os.getenv("SRE_API_TIMEOUT", "10")
        )

        self.openai_api_key = os.getenv("OPENAI_API_KEY", "")

        self.openai_model = os.getenv(
            "OPENAI_MODEL",
            "gpt-6-luna",
        )


settings = Settings()