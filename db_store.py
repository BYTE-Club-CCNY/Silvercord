import bs4
import sys
import os
os.environ['USER_AGENT'] = 'myagent'
import requests
from langchain_chroma import Chroma
from langchain_openai import OpenAIEmbeddings
from langchain_community.document_loaders import WebBaseLoader
from langchain.text_splitter import RecursiveCharacterTextSplitter
from dotenv import load_dotenv
from backend.rmp import get_professor_url

load_dotenv()

CHROMA_PATH = "chroma"
OPEN_AI_KEY = os.getenv('OPEN_AI_KEY')
embed = OpenAIEmbeddings(
    api_key=OPEN_AI_KEY,
    model="text-embedding-3-large"
)

url1 = "https://www.ratemyprofessors.com/professor/2380866" 

def pipeline(url):
    import asyncio
    documents = asyncio.run(load_docs(url))
    chunks = split_text(documents)
    # print(f"{chunks=}")
    # gotta be a better way of handling this exception !!
    try:    
        if chroma_store(chunks, url):
            print("ChromaDB Store Successful!")
        else:
            print("Nothing was stored.")
    except Exception as e:
        print("Error occurred when vector storing: ", e)
        
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
    chunk_size = len(chunks)-1
    document = chunks[chunk_size]
    #print(document.page_content)
    #print("<------------------------------>\n")
    #print(document.metadata) # note that this stores the actual url as source
    return chunks

def chroma_store(chunks: list, url: str):
    db = Chroma(
        collection_name="silvercord",
        embedding_function=embed,
        persist_directory=CHROMA_PATH
    )
    texts = [chunk.page_content for chunk in chunks]
    metadatas = [chunk.metadata for chunk in chunks]
    # print(f"{texts=}")
    # print(f"{metadatas=}")
    # query the db with the url
    # if it exists, don't update, but if it doesn't exist, add it
    if db.get(include=["metadatas"], where={"source": url}):
        print(f"URL: {url} already exists in the database.")
        return False
    else:
        db.add_texts(texts=texts, metadatas=metadatas)
        return True
    #print(f"Saved {len(chunks)} chunks to {CHROMA_PATH}.")

if __name__ == "__main__":
    url = sys.argv[1] if len(sys.argv) > 1 else "default_url"
    pipeline(url)
