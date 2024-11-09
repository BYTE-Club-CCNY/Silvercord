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

# basic configuration with LLM, chromadb, langchain
load_dotenv()
CHROMA_PATH = "chroma"
OPEN_AI_KEY = os.getenv('OPEN_AI_KEY')
embed = OpenAIEmbeddings(
    api_key=OPEN_AI_KEY,
    model="text-embedding-3-large"
)


# pipeline function that simply takes url as data and stores content into chroma for vector searching for the LLM
# if you're not familiar with building RAGs, the pipeline is as follows:
  # 1. load documents (whichever type)
  # 2. split text into multiple chunks for storing into ChromaDB
  # 3. create a new collection with the LLM embedding function
  # 4. store the text & metadata collected from the embedded docs into chromaDB for LLM lookup 
def pipeline(url):
    import asyncio
    documents = asyncio.run(load_docs(url))
    chunks = split_text(documents)
    try:    
        chroma_store(chunks)
        print("ChromaDB Store Successful!")
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

def chroma_store(chunks: list):
    db = Chroma(
        collection_name="silvercord",
        embedding_function=embed,
        persist_directory=CHROMA_PATH
    )
    texts = [chunk.page_content for chunk in chunks]
    metadatas = [chunk.metadata for chunk in chunks]
    db.add_texts(texts=texts, metadatas=metadatas)
    #print(f"Saved {len(chunks)} chunks to {CHROMA_PATH}.")

# below defines the specification:
# arguments are taken by sys.argvs in terminal
# example is: python db_store.py https://www.ratemyprofessors.com/professor/432142
# above would process the parameter link as the url for chroma storing
if __name__ == "__main__":
    url = sys.argv[1] if len(sys.argv) > 1 else "default_url"
    pipeline(url)
