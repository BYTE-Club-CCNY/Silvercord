import os
import re

import chromadb
import cohere
import requests
from bs4 import BeautifulSoup
from dotenv import load_dotenv

load_dotenv()
COHERE_KEY = os.getenv("COHERE_KEY")
CHROMA_DB_PATH = os.getenv("CHROMA_DB_PATH")
chroma_client = chromadb.PersistentClient(path=CHROMA_DB_PATH)
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
CHUNK_SIZE = 400
CHUNK_OVERLAP = 50


def extract_page_html(url_in: str):
    try:
        response = requests.get(url_in, headers=headers, timeout=10)
        response.raise_for_status()
        return response.text
    except requests.RequestException as e:
        return None


def extract_prof_reviews(html_in, prof_name_in):
    soup = BeautifulSoup(html_in, 'html.parser')
    
    review_containers = soup.find_all('div', class_=re.compile(r'Comments__StyledComments'))
    
    if not review_containers:
        review_containers = soup.find_all('div', attrs={'data-testid': re.compile(r'comment')})
    
    reviews_list = []
    seen_reviews = set()
    
    for i, div in enumerate(review_containers):
        review = div.get_text(separator=' ', strip=True)
        
        review = re.sub(r'\s+', ' ', review)
        
        if review and len(review) > 20 and review not in seen_reviews:
            seen_reviews.add(review)
            reviews_list.append({
                'content': review,
                'commentId': i+1,
                'source': 'ratemyprofessor',
                'professor_name': prof_name_in
            })
    
    return reviews_list


def chunk_text(text, chunk_size=CHUNK_SIZE, overlap=CHUNK_OVERLAP):
    chunks = []
    start = 0
    text_len = len(text)
    
    while start < text_len:
        end = start + chunk_size
        chunk = text[start:end]
        
        if len(chunk.strip()) > 50:
            chunks.append(chunk)
        
        start += chunk_size - overlap
    
    return chunks if chunks else [text]


def embed(content):
    if not content:
        return None
    
    embed_response = co.embed(
        texts=content,
        model=embedding_model,
        input_type=input_type,
        embedding_types=["float"]
    )
    return embed_response


def vector_store(embeddings_in, data_list, prof_id_in, prof_name_in):
    collection = chroma_client.get_or_create_collection(name="professor_reviews")
    
    existing = collection.get(
        where={"professor_id": prof_id_in},
        limit=1
    )
    
    if existing['ids']:
        return collection
    
    collection.add(
        ids=[f"prof_{prof_id_in}_chunk_{i}" for i in range(len(data_list))],
        documents=[c['content'] for c in data_list],
        embeddings=embeddings_in.embeddings.float,
        metadatas=[{
            'chunk_id': i,
            'source': c['source'],
            'professor_name': prof_name_in,
            'professor_id': prof_id_in,
            'original_comment_id': c.get('commentId', -1)
        } for i, c in enumerate(data_list)],
    )
    
    return collection


def process_professor(prof_url, prof_name_in):
    html = extract_page_html(prof_url)

    if not html:
        return False

    prof_id = prof_url.split('/')[-1]
    reviews = extract_prof_reviews(html, prof_name_in)

    if not reviews:
        return False

    all_chunks = []
    for review in reviews:
        chunks = chunk_text(review['content'])
        for chunk in chunks:
            all_chunks.append({
                'content': chunk,
                'commentId': review['commentId'],
                'source': review['source'],
                'professor_name': prof_name_in
            })

    chunk_texts = [chunk['content'] for chunk in all_chunks]
    embeddings = embed(chunk_texts)

    if embeddings:
        vector_store(embeddings, all_chunks, prof_id, prof_name_in)
        return True
    else:
        return False


if __name__ == "__main__":
    url = 'https://www.ratemyprofessors.com/professor/2380866'
    prof_name = "Erik Grimmelman"
    process_professor(url, prof_name)
