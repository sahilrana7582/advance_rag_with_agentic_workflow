"""Run an AI call.

Usage:
    python main.py                          # uses the sample prompt below
    python main.py "your question here"     # ask anything
"""
from __future__ import annotations

import sys

from ai_client import AIClient

SAMPLE_PROMPT = "In one sentence, what is a Kubernetes CrashLoopBackOff?"


def main() -> None:
    prompt = " ".join(sys.argv[1:]).strip() or SAMPLE_PROMPT
    client = AIClient()
    print("You:", prompt)
    print("-" * 40)
    try:
        print(client.ask(prompt))
    except RuntimeError as exc:
        print(f"[setup] {exc}")


if __name__ == "__main__":
    main()
