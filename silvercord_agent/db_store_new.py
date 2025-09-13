import chromadb
import requests
from bs4 import BeautifulSoup
from langchain_text_splitters import RecursiveCharacterTextSplitter

chroma_client = chromadb.Client()
headers = {
    'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36',
    'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8',
    'Accept-Language': 'en-US,en;q=0.5',
    'Accept-Encoding': 'gzip, deflate',
    'Connection': 'keep-alive',
}

def extract_page_html(url: str):
    try:
        response = requests.get(url, headers=headers, timeout=10)
        response.raise_for_status()
        return response.text
    except requests.RequestException as e:
        print(f"Error fetching URL due to error: {e}")
        return None

def extract_prof_comments(html):
    soup = BeautifulSoup(html, 'html.parser')
    comments = soup.find_all('div', class_='Comments__StyledComments-dzzyvm-0 jpfwLX')
    commentsList = []
    return comments

if __name__ == '__main__':
    html = extract_page_html('https://www.ratemyprofessors.com/professor/2380866')
    print(extract_prof_comments(html))