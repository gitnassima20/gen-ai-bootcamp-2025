import streamlit as st
import requests
import random
import google.generativeai as genai
from dotenv import load_dotenv
import os

# Load environment variables
load_dotenv()

# Config
st.set_page_config(
    page_title="Setup",
    page_icon="⚙️",
    layout="centered",
)

st.title("⚙️ Setting up the Study Activity Session !")
st.divider()

# Initialize session state
if "words" not in st.session_state:
    st.session_state.words = []
if "sentence" not in st.session_state:
    st.session_state.sentence = None
if "stage" not in st.session_state:
    st.session_state.stage = "setup"  # Start in setup stage

# Fetch words from backend
def fetch_words_per_group():
    try:
        response = requests.get("http://localhost:8080/api/v1/groups/10/words/raw")
        response.raise_for_status()
        st.session_state.words = response.json().get("words", [])
        print('words: ', st.session_state.words)
    except requests.exceptions.RequestException as e:
        st.error(f"Failed to fetch words: {e}")

if not st.session_state.words:
    fetch_words_per_group()

# Load Gemini API key from .env
GEMINI_API_KEY = os.getenv("GEMINI_API_KEY")

if not GEMINI_API_KEY:
    st.error("Missing Gemini API Key. Please set GEMINI_API_KEY in the .env file.")
else:
    genai.configure(api_key=GEMINI_API_KEY)

    def generate_sentence(english):
        prompt = (
            f"Generate a very simple Japanese sentence using the word '{english}', "
            f"following JLPT N2 grammar rules."
            f"Use at max three words in the sentence leveraging its simplicity"
            f"The Output will contain only the generated sentence, no other extra explanation"
        )
        model = genai.GenerativeModel("gemini-2.0-flash")
        response = model.generate_content(prompt)
        return response.text if response else "Failed to generate sentence."

        # Button to generate sentence
    if st.session_state.stage == "setup":
        if st.button("Generate Sentence"):
            if st.session_state.words:
                random_word = random.choice(st.session_state.words)
                st.session_state.selected_word = random_word  # Store selected word
                english = random_word.get("english", "")
                kanji = random_word.get("kanji", "")

                if english:
                    st.session_state.sentence = generate_sentence(english)
                else:
                    st.error("Could not retrieve a valid word.")
            else:
                st.error("No words available. Try refreshing the page.")

    # Display the generated sentence
    if st.session_state.sentence:
        st.subheader("Generated Sentence:")
        st.write(f"## {st.session_state.sentence}")

        # Display the word used
        if "selected_word" in st.session_state:
            word = st.session_state.selected_word
            st.markdown(f"""
            **Word Used:** {word.get('english', 'Unknown')}  
            **Kanji:** {word.get('kanji', 'N/A')}  
            """)
        
        # Show button to proceed to practice
        if st.button("Next"):
            st.session_state.stage = "practice"
            st.switch_page("01_practice.py")
