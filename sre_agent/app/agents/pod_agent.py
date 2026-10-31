from langchain_openai import ChatOpenAI

from app.config import settings
from app.tools.pod_tools import list_pods


llm = ChatOpenAI(
    model=settings.openai_model,
    api_key=settings.openai_api_key,
)

llm_with_tools = llm.bind_tools([list_pods])