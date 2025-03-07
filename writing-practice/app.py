import streamlit as st
from manga_ocr import MangaOcr


pg = st.navigation([st.Page(page="00_setup.py", url_path='setup'),
                    st.Page(page="01_practice.py", url_path='practice_writing'),
                    st.Page(page="02_review.py", url_path='review')], 
                   position="sidebar")
pg.run()

# Initialize the MangaOCR model
@st.cache_resource
def load_ocr_model():
    return MangaOcr()

# Store OCR model in session state if not already set
if "ocr_model" not in st.session_state:
    st.session_state.ocr_model = load_ocr_model()

# Initialize session state
if "words" not in st.session_state:
    st.session_state.words = []
if "sentence" not in st.session_state:
    st.session_state.sentence = None
if "stage" not in st.session_state:
    st.session_state.stage = "setup"