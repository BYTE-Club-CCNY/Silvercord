import bs4
import os
from langchain_chroma import Chroma
from langchain_openai import OpenAIEmbeddings
from langchain_community.document_loaders import WebBaseLoader
from langchain.text_splitter import RecursiveCharacterTextSplitter
from dotenv import load_dotenv

load_dotenv()
OPEN_AI_KEY = os.getenv('OPEN_AI_KEY')
url = "https://www.ratemyprofessors.com/professor/432142"
CHROMA_PATH = "chroma"

embed = OpenAIEmbeddings(
    api_key=OPEN_AI_KEY,
    model="text-embedding-3-large"
)

def pipeline(url):
    import asyncio
    documents = asyncio.run(load_docs(url))
    chunks = split_text(documents)
    chroma_store(chunks)

async def load_docs(url):
    loader = WebBaseLoader(web_paths=[url])
    docs = []
    async for doc in loader.alazy_load():
        docs.append(doc)
    assert len(docs) == 1
    return docs

def split_text(documents: list):
    text_splitter = RecursiveCharacterTextSplitter(
        chunk_size=1100,
        chunk_overlap=100,
        length_function=len,
    )
    chunks = text_splitter.split_documents(documents)
    document = chunks[10]
    print(document.page_content)
    print("<------------------------------>\n")
    print(document.metadata) # note that this stores the actual url as source
    return chunks

def chroma_store(chunks: list):
    db = Chroma(
        collection_name="silvercord",
        embedding_function=embed,
        persist_directory=CHROMA_PATH
    )
    texts = [chunk.page_content for chunk in chunks]
    metadatas = [chunk.metadata for chunk in chunks]
    db.add_texts(texts=texts, metadatas=metadatas)
    print("#3")
    print(f"Saved {len(chunks)} chunks to {CHROMA_PATH}.")

if __name__ == "__main__":
    pipeline(url)
