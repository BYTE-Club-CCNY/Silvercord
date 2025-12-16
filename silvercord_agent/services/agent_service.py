from datetime import datetime
from difflib import SequenceMatcher
from services.db_service import get_clients
from rmp import get_professor_url
from db_store import process_professor

embedding_model = "embed-v4.0"


def fuzzy_match(query_name: str, stored_names: list, threshold: float = 0.7):
    """
    Fuzzy match professor names to handle typos and partial matches.

    Tries to match in two ways:
    1. Full name exact match
    2. Individual first/last name match

    Example: query_name="Troeger" matches "Douglas Troeger"
    """
    best_match = None
    best_score = 0

    query_clean = query_name.strip().lower()

    for full_name in stored_names:
        full_name_lower = full_name.lower()

        score = SequenceMatcher(None, query_clean, full_name_lower).ratio()
        if score > best_score and score >= threshold:
            best_score = score
            best_match = full_name

        name_parts = full_name.split()
        if len(name_parts) >= 2:
            first_name = name_parts[0].lower()
            last_name = name_parts[-1].lower()

            first_score = SequenceMatcher(None, query_clean, first_name).ratio()
            if first_score > best_score and first_score >= threshold:
                best_score = first_score
                best_match = full_name

            last_score = SequenceMatcher(None, query_clean, last_name).ratio()
            if last_score > best_score and last_score >= threshold:
                best_score = last_score
                best_match = full_name

    return best_match, best_score


def retrieve_data(query_in: str, prof_name_input: str, n_results: int = 10):
    """Retrieve relevant professor reviews from ChromaDB using RAG."""
    chroma_client, cohere_client, _ = get_clients()
    collection = chroma_client.get_or_create_collection("professor_reviews")

    try:
        data = collection.get(include=["metadatas"])
        stored_professors = list(set(meta['professor_name'] for meta in data['metadatas'])) if data['metadatas'] else []
    except Exception as e:
        print(f"[DEBUG] Error getting collection data: {e}")
        stored_professors = []

    matched_prof, score = fuzzy_match(prof_name_input, stored_professors)

    if matched_prof:
        prof_name = matched_prof
    else:
        prof_name = None

    embed_query = cohere_client.embed(
        texts=[query_in],
        input_type="search_query",
        model=embedding_model,
        embedding_types=["float"]
    ).embeddings.float

    where_condition = {"professor_name": prof_name} if prof_name else None

    results = collection.query(
        query_embeddings=[embed_query[0]],
        n_results=n_results,
        include=["documents", "metadatas", "distances"],
        where=where_condition
    )
    return results, prof_name


def build_rag_response(query: str, professor: str):
    """Build RAG response for professor query."""
    _, _, openai_client = get_clients()

    results, matched_prof = retrieve_data(query, professor, n_results=8)
    professor_url = None

    if not matched_prof or not results['documents'][0]:
        professor_url = get_professor_url(professor)
        if professor_url:
            process_professor(professor_url, professor)
            results, matched_prof = retrieve_data(query, professor, n_results=8)
        else:
            return "I couldn't find information about this professor.", None

    context_text = "\n\n".join([
        f"Review: {doc}"
        for doc in results['documents'][0]
    ])

    response = openai_client.chat.completions.create(
        model="gpt-4o-mini",
        messages=[
            {
                "role": "system",
                "content": """You are an assistant helping students learn about professors based on reviews.
                Instructions:
                - Answer based ONLY on the reviews provided
                - Be concise (2-3 sentences max)
                - If the reviews don't contain relevant information, say so
                - Focus on teaching style, workload, and helpfulness"""
            },
            {
                "role": "user",
                "content": f"""Context (student reviews):
                {context_text}
                Question: {query}"""
            }
        ],
        temperature=0.7,
        max_tokens=150
    )

    return response.choices[0].message.content, professor_url


def get_professor_info(professor_name: str, question: str = "How is the following professor?"):
    """
    Get professor information using RAG.

    Args:
        professor_name: Name of the professor
        question: Query about the professor

    Returns:
        dict with name, link, response, processed_at
    """
    answer, professor_url = build_rag_response(question, professor_name)

    return {
        "name": professor_name,
        "link": professor_url,
        "response": answer,
        "processed_at": datetime.utcnow().isoformat()
    }


def get_break_info(year: str):
    """
    Get academic calendar break information.

    Args:
        year: Academic year (e.g., "2025-2026")

    Returns:
        dict with name, link, response, processed_at
    """
    return {
        "name": "Calendar",
        "link": None,
        "response": "This command is under maintenance. Stay tuned!",
        "processed_at": datetime.utcnow().isoformat()
    }
