import os
from dotenv import load_dotenv
from langchain_openai import OpenAIEmbeddings
from langchain_chroma import Chroma
from langchain_anthropic import ChatAnthropic
from langchain.chains import create_history_aware_retriever, create_retrieval_chain
from langchain.prompts import ChatPromptTemplate, MessagesPlaceholder
from langchain.chains.combine_documents import create_stuff_documents_chain
from backend.rmp import get_professor_url
from db_store import pipeline
# for chatting testing:
from langchain_core.messages import HumanMessage, SystemMessage


load_dotenv()
API_KEY = os.getenv("OPEN_AI_KEY")
CHROMA_PATH = "chroma"

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
        search_kwargs={"k": 6, "score_threshold": 0.45},
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
    # do check to see if url is in db
    # Fetch all documents in the database (adjust parameters as needed)
    query_results = db.get(include=["metadatas"], where={"source": get_professor_url(profName)})  # Ensure this fetches all docs, or use pagination if the db is large

    # Define the specific metadata key-value pair to search for
    # target_source = "https://www.ratemyprofessors.com/professor/432142"

    # Check if any document in all_documents matches the target metadata pair
    # matching_docs = [doc for doc in all_documents if doc.metadata.get("source") == target_source]
    
    print(f"{query_results=}")

    # if matching_docs:
    #     print("Document with the specified source exists.")
    #     # Take action based on found documents
    # else:
    #     print("No document found with the specified source.")

    qna_chain = create_stuff_documents_chain(model, full_prompt)
    rag_chain = create_retrieval_chain(history_aware_retriever, qna_chain)
    return rag_chain


def chat(): # this is for personal testing with the LLM
    print("Start asking about professors at CCNY")
    chat_history = []

    while True:
        query = input("You: ")
        if query.lower() == "exit":
            break
        rag_chain = build_rag_chain(API_KEY)
        result = rag_chain.invoke({"input": query, "chat_history": chat_history})
        print(f"AI: {result['answer']}")
        chat_history.append(HumanMessage(content=query))
        chat_history.append(SystemMessage(content=result["answer"]))


def process_query(prof_name):
    rag_chain = build_rag_chain(API_KEY)
    user_prompt = f"How is Professor {prof_name}'s course?"
    history = []
    result = rag_chain.invoke({"input": user_prompt, "chat_history": history})
    answer = result.get("answer", "No response available")
    return answer

if __name__ == "__main__":
    import sys
    profName = sys.argv[1] if len(sys.argv) > 1 else "Unknown"
    
    pipeline(get_professor_url(profName))
    response = process_query(profName)
    print(f"```{response}```")
