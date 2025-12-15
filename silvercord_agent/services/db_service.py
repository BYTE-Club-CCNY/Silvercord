import os
import chromadb
import cohere
from openai import OpenAI
from dotenv import load_dotenv

load_dotenv()

# Global clients (initialized once)
_chroma_client = None
_cohere_client = None
_openai_client = None

# Configuration
COHERE_KEY = os.getenv("COHERE_KEY")
OPENAI_KEY = os.getenv("OPEN_AI_KEY")
SCRIPT_DIR = os.path.dirname(os.path.abspath(__file__))
CHROMA_DB_PATH = os.path.join(os.path.dirname(SCRIPT_DIR), "chroma_db")


def init_clients():
    """Initialize all clients (ChromaDB, Cohere, OpenAI) as singletons."""
    global _chroma_client, _cohere_client, _openai_client

    if _chroma_client is None:
        _chroma_client = chromadb.PersistentClient(path=CHROMA_DB_PATH)

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
