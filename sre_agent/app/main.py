from app.agents.pod_agent import llm_with_tools
from app.tools.pod_tools import list_pods

from langchain.messages import HumanMessage, ToolMessage


def main():
    messages = [
        HumanMessage(
            content="Show me all pods in the taskflow namespace."
        )
    ]

    response = llm_with_tools.invoke(messages)

    messages.append(response)

    for tool_call in response.tool_calls:
        if tool_call["name"] == "list_pods":
            result = list_pods.invoke(tool_call["args"])

            messages.append(
                ToolMessage(
                    content=str(result),
                    tool_call_id=tool_call["id"],
                )
            )

    final_response = llm_with_tools.invoke(messages)

    print("\nSRE Agent:")
    print(final_response.content)


if __name__ == "__main__":
    main()