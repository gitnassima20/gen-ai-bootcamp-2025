import streamlit as st
import random
from PIL import Image
import numpy as np
import tempfile
import os
from manga_ocr import MangaOcr

# Initialize the MangaOCR model
@st.cache_resource
def load_ocr_model():
    return MangaOcr()

# Word pairs for simple sentences
word_pairs = [
    {"eng": "I eat", "jap": "私は食べる", "romaji": "Watashi wa taberu"},
    {"eng": "You drink", "jap": "あなたは飲む", "romaji": "Anata wa nomu"},
    {"eng": "She reads", "jap": "彼女は読む", "romaji": "Kanojo wa yomu"},
    {"eng": "He writes", "jap": "彼は書く", "romaji": "Kare wa kaku"},
    {"eng": "We walk", "jap": "私たちは歩く", "romaji": "Watashitachi wa aruku"},
    {"eng": "They run", "jap": "彼らは走る", "romaji": "Karera wa hashiru"},
    {"eng": "I sleep", "jap": "私は寝る", "romaji": "Watashi wa neru"},
    {"eng": "Dog barks", "jap": "犬が吠える", "romaji": "Inu ga hoeru"},
    {"eng": "Cat meows", "jap": "猫が鳴く", "romaji": "Neko ga naku"},
    {"eng": "Bird flies", "jap": "鳥が飛ぶ", "romaji": "Tori ga tobu"},
]

def calculate_score(expected, actual):
    """Calculate a simple similarity score between expected and actual text"""
    if not actual:
        return 0
    
    # Convert to sets of characters for comparison
    expected_chars = set(expected)
    actual_chars = set(actual)
    
    # Calculate intersection and union of characters
    intersection = expected_chars.intersection(actual_chars)
    union = expected_chars.union(actual_chars)
    
    # Calculate Jaccard similarity
    if len(union) == 0:
        return 0
    
    similarity = len(intersection) / len(union)
    return similarity * 100

def main():
    st.title("Japanese Writing Practice")
    st.write("Practice writing simple Japanese sentences and get feedback on your handwriting.")
    
    # Session state for tracking current exercise and history
    if 'current_pair' not in st.session_state:
        st.session_state.current_pair = random.choice(word_pairs)
    
    if 'history' not in st.session_state:
        st.session_state.history = []
    
    if 'ocr_model' not in st.session_state:
        with st.spinner("Loading OCR model... (this may take a moment)"):
            st.session_state.ocr_model = load_ocr_model()
    
    # Display the current English sentence to write
    st.markdown("### Write this sentence in Japanese:")
    st.markdown(f"**{st.session_state.current_pair['eng']}**")
    
    # File uploader for the handwritten image
    uploaded_file = st.file_uploader("Upload a photo of your handwritten Japanese", type=["jpg", "jpeg", "png"])
    
    col1, col2 = st.columns(2)
    
    if uploaded_file is not None:
        # Display the uploaded image
        image = Image.open(uploaded_file)
        col1.image(image, caption="Your submission", use_column_width=True)
        
        # Save image to a temporary file for OCR processing
        with tempfile.NamedTemporaryFile(delete=False, suffix='.jpg') as tmp:
            # Convert image to RGB mode to remove alpha channel
            image_rgb = image.convert('RGB')
            image_rgb.save(tmp, format="JPEG")
            tmp_path = tmp.name
        
        try:
            # Run OCR on the image
            with st.spinner("Analyzing your handwriting..."):
                ocr_text = st.session_state.ocr_model(tmp_path)
            
            # Clean up the temporary file
            os.unlink(tmp_path)
            
            # Display the OCR result
            col2.markdown("### OCR Result:")
            col2.write(f"Detected text: {ocr_text}")
            
            # Calculate and display the score
            expected_text = st.session_state.current_pair['jap']
            score = calculate_score(expected_text, ocr_text)
            
            st.markdown("### Feedback:")
            st.write(f"Correct Japanese: **{expected_text}** ({st.session_state.current_pair['romaji']})")
            
            # Display score with color-coded feedback
            if score >= 80:
                st.success(f"Great job! Your writing clarity score: {score:.1f}%")
            elif score >= 50:
                st.warning(f"Good attempt. Your writing clarity score: {score:.1f}%")
            else:
                st.error(f"Keep practicing! Your writing clarity score: {score:.1f}%")
            
            # Add to history
            st.session_state.history.append({
                "english": st.session_state.current_pair['eng'],
                "expected": expected_text,
                "detected": ocr_text,
                "score": score
            })
            
            # Button to get a new sentence
            if st.button("Next Sentence"):
                st.session_state.current_pair = random.choice(word_pairs)
                st.experimental_rerun()
        
        except Exception as e:
            st.error(f"Error processing image: {e}")
    
    # Display history if available
    if st.session_state.history:
        st.markdown("### Practice History")
        
        # Calculate average score
        avg_score = sum(item["score"] for item in st.session_state.history) / len(st.session_state.history)
        st.write(f"Average score: {avg_score:.1f}%")
        
        # Display history items
        for i, item in enumerate(reversed(st.session_state.history[-5:])):
            with st.expander(f"Practice #{len(st.session_state.history) - i}: {item['english']} ({item['score']:.1f}%)"):
                st.write(f"Expected: {item['expected']}")
                st.write(f"Detected: {item['detected']}")
    
    # Instructions
    with st.expander("How to use this app"):
        st.write("""
        1. You'll see a simple English sentence.
        2. Write this sentence in Japanese on a piece of paper.
        3. Take a clear photo of your writing.
        4. Upload the photo using the uploader above.
        5. The app will analyze your handwriting and give you a score.
        6. Click 'Next Sentence' to practice with a new sentence.
        
        Tips for better results:
        - Write clearly on white paper with black ink
        - Make sure there's good lighting when taking the photo
        - Crop the image to include just your writing
        - Hold the camera directly above the paper to avoid angle distortion
        """)
    
    # Footer
    st.markdown("---")
    st.caption("Japanese Writing Practice App | Language Learning Prototype")

if __name__ == "__main__":
    main()