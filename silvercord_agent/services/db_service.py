import os

import chromadb
import cohere
from dotenv import load_dotenv
from openai import OpenAI

load_dotenv()

# Global clients (initialized once)
_chroma_client = None
_cohere_client = None
_openai_client = None
_professor_collection = None

# Configuration
COHERE_KEY = os.getenv("COHERE_KEY")
OPENAI_KEY = os.getenv("OPEN_AI_KEY")
CHROMA_DB_PATH = os.getenv("CHROMA_DB_PATH")


def init_clients():
    """Initialize all clients (ChromaDB, Cohere, OpenAI) as singletons."""
    global _chroma_client, _cohere_client, _openai_client, _professor_collection

    if _chroma_client is None:
        _chroma_client = chromadb.PersistentClient(path=CHROMA_DB_PATH)

    if _professor_collection is None: 
        _professor_collection = _chroma_client.get_or_create_collection(
                name="professor_reviews",
                metadata={"hnsw:space": "cosine"}
        )

    if _cohere_client is None:
        if COHERE_KEY:
            _cohere_client = cohere.ClientV2(api_key=COHERE_KEY)

    if _openai_client is None:
        if OPENAI_KEY:
            _openai_client = OpenAI(api_key=OPENAI_KEY)

    return _chroma_client, _cohere_client, _openai_client


def get_clients():
    """Get initialized clients. Will initialize if not already done."""
    if _chroma_client is None:
        init_clients()
    return _chroma_client, _cohere_client, _openai_client


def check_health():
    """Check health of all services for health endpoint."""
    chroma_client, cohere_client, openai_client = get_clients()

    status = {
        'chromadb': 'disconnected',
        'cohere': 'missing',
        'openai': 'missing'
    }

    # Check ChromaDB
    try:
        chroma_client.heartbeat()
        status['chromadb'] = 'connected'
    except Exception:
        status['chromadb'] = 'disconnected'

    # Check Cohere
    if cohere_client is not None:
        status['cohere'] = 'configured'

    # Check OpenAI
    if openai_client is not None:
        status['openai'] = 'configured'

    return status
