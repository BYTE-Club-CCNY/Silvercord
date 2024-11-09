import os
from dotenv import load_dotenv
from langchain_openai import OpenAIEmbeddings
from langchain_chroma import Chroma
from langchain_anthropic import ChatAnthropic
from langchain.chains import create_history_aware_retriever, create_retrieval_chain
from langchain.prompts import ChatPromptTemplate, MessagesPlaceholder
from langchain.chains.combine_documents import create_stuff_documents_chain
# for chatting testing:
from langchain_core.messages import HumanMessage, SystemMessage


load_dotenv()
API_KEY = os.getenv("OPEN_AI_KEY")
CHROMA_PATH = "chroma"

# below is the RAG chain
# we must initialize the embedding func, read the DB, init the retriever and model that we work with on the client end
# build context + chat history + retriever + prompt -> RAG CHAIN
def build_rag_chain(api_key):
    embed = OpenAIEmbeddings(
        api_key=api_key,
        model="text-embedding-3-large"
    )
    db = Chroma(
        collection_name="silvercord",
        embedding_function=embed,
        persist_directory=CHROMA_PATH,
    )
    retriever = db.as_retriever(
        search_type="similarity_score_threshold",
        search_kwargs={"k": 8, "score_threshold": 0.35},
    )
    model = ChatAnthropic(
        model="claude-3-sonnet-20240229",
    )
    context = (
        "Given a chat history and the latest user question "
        "which might reference context in the chat history, "
        "formulate a standalone question which can be understood "
        "without the chat history. Do NOT answer the question, just "
        "reformulate it if needed and otherwise return it as is."
    )
    context_history = ChatPromptTemplate(
        [
            ("system", context),
            MessagesPlaceholder("chat_history"),
            ("human", "{input}"),
        ]
    )
    history_aware_retriever = create_history_aware_retriever(
        model, retriever, context_history
    )
    query = (
        "You are an assistant for question-answering tasks. "
        "Use the given pieces of retrieved context to answer "
        "the question given by the user. If you don't know the "
        "answer, just say that you do not know. Use 5 sentences "
        "maximum and keep the answer concise. Your answer must strictly "
        "be based on the data and context given to you. "
        "Special case: If asked about Gertner in the context, try to talk like he does "
        "in your response. Here's how he sounds like usually: zere iz alwayz another courze."
        "\n\n"
        "{context}"
    )
    full_prompt = ChatPromptTemplate(
        [
            ("system", query),
            MessagesPlaceholder("chat_history"),
            ("human", "{input}"),
        ]
    )
    qna_chain = create_stuff_documents_chain(model, full_prompt)
    rag_chain = create_retrieval_chain(history_aware_retriever, qna_chain)
    return rag_chain


def process_query(prof_name):
    rag_chain = build_rag_chain(API_KEY)
    user_prompt = f"How is Professor {prof_name}'s course?"
    history = []
    result = rag_chain.invoke({"input": user_prompt, "chat_history": history})
    answer = result.get("answer", "No response available")
    return answer


# below can be used for personal use for testing AND by the discord bot
# works almost just like the db_store.py function for ChromaDB
# takes argument for prof name, and does vector lookup into our chromadb
# example: python llm.py Troeger
# searches for Troeger in the DB, get relevant docs with the retriever, then creates a relevant response
if __name__ == "__main__":
    import sys
    profName = sys.argv[1] if len(sys.argv) > 1 else "Unknown"
    response = process_query(profName)
    print(f"```{response}```")
