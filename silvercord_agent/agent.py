import os
import chromadb
import cohere
from dotenv import load_dotenv
from difflib import SequenceMatcher
from openai import OpenAI

load_dotenv()
COHERE_KEY = os.getenv("COHERE_KEY")
OPENAI_KEY = os.getenv("OPEN_AI_KEY")
chroma_client = chromadb.PersistentClient(path="./chroma_db")
co = cohere.ClientV2(api_key=COHERE_KEY)
oai = OpenAI(api_key=OPENAI_KEY)
embedding_model = "embed-v4.0"


def fuzzy_match(query_name: str, stored_names: list, threshold: float = 0.7):
    '''
    below tries to ensure to try to match for two cases:
    user puts in full name: regular case, tries to match exact name
    user puts in just 1 word (either first name or last name), and sequence matches with
    each full name's first name then last name individually
    EXAMPLE: query_name: Troeger; full_name in stored_names: Douglas Troeger -> MATCH
    '''

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
    collection = chroma_client.get_collection("professor_reviews")

    data = collection.get(include=["metadatas"])
    stored_professors = list(set(meta['professor_name'] for meta in data['metadatas']))

    matched_prof, score = fuzzy_match(prof_name_input, stored_professors)

    if matched_prof:
        print(f"MATCHED '{prof_name_input}' to '{matched_prof}' (score: {score:.2f})")
        prof_name = matched_prof
    else:
        print(f"No match, score: {score:.2f}")
        prof_name = None

    embed_query = co.embed(
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
    results, matched_prof = retrieve_data(query, professor, n_results=8)

    if not results['documents'][0]:
        return "No information found for this professor."

    context_text = "\n\n".join([
        f"Review: {doc}"
        for doc in results['documents'][0]
    ])

    response = oai.chat.completions.create(
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

    return response.choices[0].message.content


if __name__ == "__main__":
    professor = " Troeger"
    question = "How is the following professor?"

    answer = build_rag_response(question, professor)
    print(f"\nQuestion: {question}")
    print(f"Professor: {professor}")
    print(f"\nAnswer:\n{answer}")