import os
import chromadb
import requests
import cohere
from bs4 import BeautifulSoup
from dotenv import load_dotenv
import re

load_dotenv()
COHERE_KEY = os.getenv("COHERE_KEY")
chroma_client = chromadb.PersistentClient(path="./chroma_db")
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
        print(f"Error fetching URL due to error: {e}")
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
    try:
        collection = chroma_client.create_collection(name="professor_reviews")
    except Exception:
        collection = chroma_client.get_collection(name="professor_reviews")
    
    existing = collection.get(
        where={"professor_id": prof_id_in},
        limit=1
    )
    
    if existing['ids']:
        print(f"Professor ID {prof_id_in} already exists in database. Skipping.")
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
    
    print(f"Stored {len(data_list)} chunks for {prof_name_in} (ID: {prof_id_in})")
    return collection

if __name__ == '__main__':
    url = 'https://www.ratemyprofessors.com/professor/432142'
    prof_name = "Douglas Troeger"
    
    print(f"Fetching data for {prof_name}...")
    html = extract_page_html(url)
    
    if not html:
        print("Failed to fetch page")
        exit(1)
    
    prof_id = url.split('/')[-1]
    reviews = extract_prof_reviews(html, prof_name)
    
    if not reviews:
        print("No reviews found")
        exit(1)
    
    print(f"Found {len(reviews)} reviews")
    
    all_chunks = []
    for review in reviews:
        chunks = chunk_text(review['content'])
        for chunk in chunks:
            all_chunks.append({
                'content': chunk,
                'commentId': review['commentId'],
                'source': review['source'],
                'professor_name': prof_name
            })
    
    print(f"Created {len(all_chunks)} chunks from {len(reviews)} reviews")
    
    chunk_texts = [chunk['content'] for chunk in all_chunks]
    embeddings = embed(chunk_texts)
    
    if embeddings:
        vector_store(embeddings, all_chunks, prof_id, prof_name)
        print("Successfully stored all chunks!")
    else:
        print("Failed to generate embeddings")