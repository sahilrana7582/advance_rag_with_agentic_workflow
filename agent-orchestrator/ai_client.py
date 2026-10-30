"""Minimal AI client: one method that sends a prompt to Gemini and returns the reply."""
from __future__ import annotations

import httpx

from config import GEMINI_API_KEY, GEMINI_MODEL


class AIClient:
    def __init__(self, api_key: str = GEMINI_API_KEY, model: str = GEMINI_MODEL):
        self.api_key = api_key
        self.model = model

    def ask(self, prompt: str, system: str = "You are a helpful assistant.") -> str:
        """Send ``prompt`` to the model and return the text response."""
        if not self.api_key or self.api_key == "your_gemini_api_key":
            raise RuntimeError("GEMINI_API_KEY is not set — copy .env.example to .env and add your key.")

        url = f"https://generativelanguage.googleapis.com/v1beta/models/{self.model}:generateContent"
        payload = {
            "systemInstruction": {"parts": [{"text": system}]},
            "contents": [{"role": "user", "parts": [{"text": prompt}]}],
        }
        resp = httpx.post(
            url,
            headers={"x-goog-api-key": self.api_key},
            json=payload,
            timeout=60,
        )
        if not resp.is_success:
            raise RuntimeError(f"Gemini error {resp.status_code}: {resp.text}")
        data = resp.json()
        try:
            return data["candidates"][0]["content"]["parts"][0]["text"]
        except (KeyError, IndexError) as exc:
            raise RuntimeError(f"Unexpected response from Gemini: {data}") from exc

