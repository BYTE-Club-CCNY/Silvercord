import os
import chromadb
import requests
import cohere
from bs4 import BeautifulSoup
from dotenv import load_dotenv

load_dotenv()
COHERE_KEY = os.getenv("COHERE_KEY")
chroma_client = chromadb.Client()
co = cohere.ClientV2(api_key=COHERE_KEY)
headers = {
    'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36',
    'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8',
    'Accept-Language': 'en-US,en;q=0.5',
    'Accept-Encoding': 'gzip, deflate',
    'Connection': 'keep-alive',
}

embedding_model = "embed-v4.0"
input_type = "search_document"

def extract_page_html(url: str):
    try:
        response = requests.get(url, headers=headers, timeout=10)
        response.raise_for_status()
        return response.text
    except requests.RequestException as e:
        print(f"Error fetching URL due to error: {e}")
        return None

def extract_prof_reviews(html_in):
    soup = BeautifulSoup(html_in, 'html.parser')
    reviewsFull = soup.find_all('div', class_='Comments__StyledComments-dzzyvm-0 jpfwLX')
    reviews_list = []
    for i, div in enumerate(reviewsFull):
        review = div.get_text(strip=True)
        if review:
            reviews_list.append({
                'content': review,
                'commentId': i+1,
                'source': 'ratemyprofessor'
            })
    return reviews_list

def embed(content):
    embed_response = co.embed(
        texts=content,
        model=embedding_model,
        input_type=input_type,
        embedding_types=["float"]
    )
    return embed_response

def vector_store(embeddings_in, data_list, prof_id_in):
    collection = chroma_client.create_collection(name="professor_reviews")
    collection.add(
        ids=[f"prof_{prof_id_in}_review_{c['commentId']}" for c in data_list],
        documents=[c['content'] for c in data_list],
        embeddings=embeddings_in.embeddings.float,
        metadatas=[{'commentId': c['commentId'], 'source': c['source']} for c in data_list],
    )
    return collection

if __name__ == '__main__':
    url = 'https://www.ratemyprofessors.com/professor/2380866'
    html = extract_page_html(url)
    prof_id = url.split('/')[-1]
    print("prof_id: ", prof_id)
    reviews = extract_prof_reviews(html)
    review_texts= [review['content'] for review in reviews]
    embeddings = embed(review_texts)
    print(vector_store(embeddings, reviews, prof_id))