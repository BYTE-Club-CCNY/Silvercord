import codecs
import json
import os
import sys
from dotenv import load_dotenv
from openai import OpenAI

load_dotenv()
OPENAI_KEY = os.getenv("OPEN_AI_KEY")
SCRIPT_DIR = os.path.dirname(os.path.abspath(__file__))

oai = OpenAI(api_key=OPENAI_KEY)
embedding_model = "embed-v4.0"

sys.stdout = codecs.getwriter("utf-8")(sys.stdout.detach())

def summarize_conversation(messages: str): # reminder: require a type for this later
    
    response = oai.chat.completions.create(
        model="gpt-4o-mini",
        messages=[
            {
                "role": "system",
                "content": """You are an assistant tasked with summarizing long discord conversations & meetings.
                Instructions:
                - Summarize based ONLY on the messages provided
                - Be concise (Up to 5 sentences max)
                - ONLY If an event/emergency is explicitly mentioned: Also return a bulleted list with the following:
                  topic name, key dates and times, people involved, decision/conclusion made.
                """
            },
            {
                "role": "user",
                "content": f"""Context (conversation to be summarized):
                {messages}"""
            }
        ],
        temperature=0.3,
        max_tokens=150
    )

    print(response.choices[0].message.content)

if __name__ == "__main__":
    try:
        input_message = sys.argv[1]
    except IndexError:
        sys.stderr.write("Error: No message argument received from Node.js.\n")
        sys.exit(1)
        
    summarize_conversation(input_message)