import os
import chromadb
import cohere
from dotenv import load_dotenv

load_dotenv()
COHERE_KEY = os.getenv("COHERE_KEY")
chroma_client = chromadb.PersistentClient(path="./chroma_db")
co = cohere.ClientV2(api_key=COHERE_KEY)
embedding_model = "embed-v4.0"

def retrieve_data(query_in: str):
    embed_query = co.embed(
        texts=[query_in],
        input_type="search_query",
        model=embedding_model,
        embedding_types=["float"]
    ).embeddings.float

    collection = chroma_client.get_collection("professor_reviews")
    results = collection.query(
        query_embeddings=[embed_query[0]],
        n_results = 3
    )
    return results

if __name__ == "__main__":
    print(retrieve_data("Professor Douglas Troeger"))