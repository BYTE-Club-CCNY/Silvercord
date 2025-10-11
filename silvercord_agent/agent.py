import os
import chromadb
import cohere
from dotenv import load_dotenv

load_dotenv()
COHERE_KEY = os.getenv("COHERE_KEY")
chroma_client = chromadb.PersistentClient(path="./chroma_db")
co = cohere.ClientV2(api_key=COHERE_KEY)
embedding_model = "embed-v4.0"

def retrieve_data(query_in: str, n_results: int = 10):
    embed_query = co.embed(
        texts=[query_in],
        input_type="search_query",
        model=embedding_model,
        embedding_types=["float"]
    ).embeddings.float

    collection = chroma_client.get_collection("professor_reviews")
    results = collection.query(
        query_embeddings=[embed_query[0]],
        n_results=n_results,
        include=["documents", "metadatas", "distances"]
    )
    return results

if __name__ == "__main__":
    results = retrieve_data("How is Professor Douglas Troeger's teaching?")
    
    print(f"\nRetrieved {len(results['documents'][0])} chunks:\n")
    for i, (doc, metadata, distance) in enumerate(zip(
        results['documents'][0], 
        results['metadatas'][0], 
        results['distances'][0]
    )):
        print(f"--- Result {i+1} (distance: {distance:.3f}) ---")
        print(f"Professor: {metadata.get('professor_name', 'Unknown')}")
        print(f"Chunk: {doc[:200]}...")
        print()