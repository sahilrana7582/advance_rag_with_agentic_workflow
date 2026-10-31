
from app.clients.sre_api_client import SREAPIClient


def main():
    client = SREAPIClient()

    try:
        result = client.get_json("/health")
        print(result)
    finally:
        client.close()


if __name__ == "__main__":
    main()
