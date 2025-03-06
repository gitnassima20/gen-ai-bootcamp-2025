import streamlit as st
import requests

st.set_page_config(
    page_title="Setup",
    page_icon="⚙️",
    layout="centered",
)

st.title("⚙️ Setting up the Writing Practice Session !")
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
        response = requests.get("http://localhost:8080/api/v1/groups/:id/words")
        response.raise_for_status()
        st.session_state.words = response.json().get("words", [])
    except requests.exceptions.RequestException as e:
        st.error(f"Failed to fetch words: {e}")

if not st.session_state.words:
    fetch_words_per_group()
    
# Button to generate sentence
# if st.session_state.stage == "setup":
    # if st.button("Generate Sentence"):
    #     if st.session_state.words:
    #         ##todo add prompt related code here
    #     else:
    #         st.error("No words available. Try refreshing the page.")