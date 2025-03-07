import streamlit as st
from PIL import Image
import io
import google.generativeai as genai
from dotenv import load_dotenv
import os

# Load environment variables
load_dotenv()

# Config
st.set_page_config(
    page_title="Review",
    page_icon="📝",
    layout="centered",
)

st.title("📝 Review Writing")
st.divider()

# Initialize session state
if "sentence" not in st.session_state:
    st.session_state.sentence = None
if "stage" not in st.session_state:
    st.session_state.stage = "review"  # Set stage to review

# Load Gemini API key from .env
GEMINI_API_KEY = os.getenv("GEMINI_API_KEY")

if not GEMINI_API_KEY:
    st.error("Missing Gemini API Key. Please set GEMINI_API_KEY in the .env file.")
else:
    genai.configure(api_key=GEMINI_API_KEY)

    # Define the grading system function
    def grade_transcription(ocr_text, original_sentence):
        prompt = (
            f"Grade the transcription based on the original sentence. "
            f"Here is the OCR transcription: '{ocr_text}'. "
            f"Here is the original sentence: '{original_sentence}'. "
            f"Provide a English Translation for the original sentence."
            f"Provide a letter grade (S, A, B, etc.), a description of accuracy"
           
        )
        model = genai.GenerativeModel("gemini-2.0-flash")
        response = model.generate_content(prompt)
        return response.text if response else "Failed to generate grading result."

    # Ensure we're in the review stage
    if st.session_state.get("stage") != "review":
        st.warning("No submission found. Please go back and submit your handwriting first.")
        if st.button("Go to Practice"):
            st.switch_page("01_practice.py")
    else:
        # Display sentence
        st.subheader("Original Sentence:")
        st.write(st.session_state.sentence)

        # Retrieve and display uploaded image
        if "practice_image" in st.session_state:
            image = Image.open(io.BytesIO(st.session_state.practice_image))
            st.image(image, caption="Your handwritten practice", use_container_width=True)

            # Call the grading function with the OCR text and original sentence
            grading_result = grade_transcription(ocr_text, st.session_state.sentence)

            st.subheader("Grading Result:")
            st.write(grading_result)

            # Next Sentence Button
            if st.button("Next Sentence"):
                st.session_state.sentence = None  # Reset for next round
                st.session_state.stage = "setup"
                st.switch_page("00_setup.py")

        else:
            st.warning("No image found. Please upload your handwriting in the practice stage.")
