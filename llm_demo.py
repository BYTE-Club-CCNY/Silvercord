from openai import OpenAI
import os
from dotenv import load_dotenv

'''
be sure to pip install -r requirements.txt, 
and also supply a .env file with your OpenAI API key
'''
load_dotenv()
def llm_query(context): 
    client = OpenAI( #client object, handles api requests
        api_key=os.environ.get("OPENAI_API_KEY"),
    )

    chat_completion = client.chat.completions.create( #here is where the requests gets processed
                                                      #also note, it's a bit obvious when gpt is used if it doesnt use this exact endpoint on currrent version. the chat endpoint gpt is trained off of is deprecated and doesnt work
        messages=[ 
            {
                "role": "user", #user input
                "content": f"Who am I? Please use the following context: My name is {context}", #here is your question. you would add in here the context and do some prompt engineering to get a desired result
                                    #note about prompt engineering, it is easier to get what you want than to get what you don't want
            }
        ],
        model="gpt-4o-mini", #model to use
    )
    return chat_completion.choices[0].message.content #response from the model
context=input("Who are you \n")
print(llm_query(context)) #print the response

# i want you to try running this file and providing no input (just hit enter) and see what is outputted.
# then put in your name and see the difference in the output

# this is the power of a context, llms are naturally stupid (despite being called artificial intelligence) and need a lot of help to get the right answer