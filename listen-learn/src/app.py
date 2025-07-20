"""Streamlit UI for the language learning app using FAISS."""

import os
import json
import time
from pathlib import Path
from typing import Dict, List, Optional, Tuple, Any
from dotenv import load_dotenv

# Load environment variables from .env file
load_dotenv()

import streamlit as st
import numpy as np
import faiss
from sentence_transformers import SentenceTransformer
from streamlit_option_menu import option_menu

from config import DATA_DIR, JPLT_N3_VIDEO_ID
from quiz import Quiz
from audio_player import AudioPlayer

# FAISS index and metadata files
FAISS_INDEX_FILE = DATA_DIR / "faiss_index.bin"
TRANSCRIPT_METADATA_FILE = DATA_DIR / "transcript_metadata.pkl"

# Set page config
st.set_page_config(
    page_title="Japanese Learning App",
    layout="wide",
    initial_sidebar_state="expanded"
)

def load_faiss_index():
    """Load the FAISS index and metadata."""
    try:
        if not FAISS_INDEX_FILE.exists() or not TRANSCRIPT_METADATA_FILE.exists():
            st.error("FAISS index not found. Please run process.py first.")
            st.stop()
            
        # Load the FAISS index
        index = faiss.read_index(str(FAISS_INDEX_FILE))
        
        # Load the metadata
        with open(TRANSCRIPT_METADATA_FILE, 'rb') as f:
            metadatas = pickle.load(f)
            
        # Initialize the sentence transformer
        model = SentenceTransformer('all-MiniLM-L6-v2')
        
        return index, model, metadatas
        
    except Exception as e:
        st.error(f"Error loading FAISS index: {e}")
        st.stop()

def show_home(index, model, metadatas):
    """Show the home page with search functionality."""
    st.header("🔍 Search Transcript")
    
    # Add a search box with clear instructions
    col1, col2 = st.columns([3, 1])
    with col1:
        query = st.text_input(
            "Search for phrases or topics in Japanese or English:",
            placeholder="Try searching in Japanese or English...",
            label_visibility="collapsed"
        )
    with col2:
        st.write("")
        st.caption("Tip: Search in either language")
    
    if query:
        with st.spinner("🔍 Searching..."):
            try:
                # Encode the query
                query_embedding = model.encode(
                    [query], 
                    show_progress_bar=False, 
                    convert_to_numpy=True
                )
                
                # Search the index
                k = 5  # Number of results to return
                distances, indices = index.search(query_embedding, k)
                
                # Display results
                if len(indices) > 0 and len(indices[0]) > 0:
                    st.subheader("🎯 Matching Segments")
                    
                    # Add a clear button to play all results
                    if st.button("▶️ Play All Results"):
                        for idx in indices[0]:
                            if idx < len(metadatas):
                                meta = metadatas[idx]
                                st.session_state.audio_player.play_segment(meta)
                    
                    # Display each result with audio and text
                    for i, idx in enumerate(indices[0], 1):
                        if idx < len(metadatas):  # Check if index is valid
                            meta = metadatas[idx]
                            with st.expander(
                                f"{i}. ⏱️ {int(meta['start_time']//60)}:{int(meta['start_time']%60):02d} - "
                                f"{meta['text'][:50]}..."
                            ):
                                # Show audio player
                                if 'audio_player' in st.session_state:
                                    st.session_state.audio_player.play_segment(meta)
                                
                                # Display the full text
                                st.write(meta['text'])
                                
                                # Add a copy button
                                if st.button(f"📋 Copy Text {i}", key=f"copy_{i}"):
                                    st.session_state.clipboard = meta['text']
                                    st.toast("Text copied to clipboard!", icon="📋")
                                
                                st.caption(f"Duration: {meta['duration']:.1f} seconds")
                else:
                    st.info("No matching segments found. Try a different search term.")
                    
            except Exception as e:
                st.error(f"Error performing search: {e}")
                st.exception(e)

def show_quiz(quiz: Quiz, audio_player: AudioPlayer):
    """Show the quiz interface with Japanese conversation and English questions."""
    st.header("📝 Japanese Listening Comprehension Quiz")
    
    # Initialize session state for quiz if needed
    if 'quiz_started' not in st.session_state:
        st.session_state.quiz_started = True
        st.session_state.show_feedback = False
    
    # Show progress
    progress = quiz.get_progress()
    
    # Progress bar with score
    col1, col2 = st.columns([4, 1])
    with col1:
        st.progress(int(progress['progress']))
    with col2:
        st.metric("Score", f"{progress['score']}/{progress['total']}")
    
    st.caption(f"Question {progress['current']} of {progress['total']}")
    
    # Get current question
    question = quiz.get_current_question()
    
    # Quiz completion screen
    if not question or quiz.is_complete():
        st.balloons()
        st.success("🎉 Quiz Complete!")
        
        # Calculate score percentage
        score_percent = (quiz.score / progress['total']) * 100 if progress['total'] > 0 else 0
        
        # Show score with emoji based on performance
        st.markdown("### Your Results:")
        col1, col2 = st.columns(2)
        with col1:
            st.metric("Score", f"{quiz.score}/{progress['total']}")
        with col2:
            st.metric("Percentage", f"{score_percent:.0f}%")
        
        if score_percent >= 80:
            st.success("🏆 Excellent! You have a strong understanding of the material!")
        elif score_percent >= 50:
            st.info("👍 Good job! Keep practicing to improve your comprehension!")
        else:
            st.warning("💪 Keep practicing! You'll get better with more practice!")
        
        # Add buttons for next steps
        col1, col2 = st.columns(2)
        with col1:
            if st.button("🔄 Take the Quiz Again", use_container_width=True):
                quiz.reset()
                st.session_state.show_feedback = False
                st.rerun()
        with col2:
            if st.button("🏠 Return to Home", use_container_width=True):
                st.session_state.quiz_started = False
                st.rerun()
        return
    
    # Display the Japanese conversation
    if 'conversation' in question and question['conversation']:
        with st.expander("🎧 Listen to the conversation", expanded=True):
            st.write("**Japanese Conversation:**")
            
            # Show conversation with speaker avatars
            for line in question['conversation']:
                speaker = line.get('speaker', 'A')
                gender = line.get('gender', 'female' if speaker == 'A' else 'male')
                avatar = "👩" if gender == 'female' else "👨"
                
                col1, col2 = st.columns([1, 10])
                with col1:
                    st.markdown(f"**{avatar} {speaker}**")
                with col2:
                    st.markdown(f"{line.get('text', '')}")
            
            # Play the conversation audio
            audio_player.play_conversation(
                question['conversation'],
                f"conversation_{quiz.current_question}"
            )
    
    # Display the English question
    st.subheader("Question:")
    st.write(question['question'])
    
    # Show feedback if answer was submitted
    if st.session_state.show_feedback:
        last_answer = quiz.answers[-1] if quiz.answers else None
        if last_answer and last_answer['question_idx'] == quiz.current_question - 1:
            if last_answer['is_correct']:
                st.success("✅ Correct! " + question.get('explanation', ''))
            else:
                correct_answer = question['options'][question['correct_index']]
                st.error(f"❌ Incorrect. The correct answer is: {correct_answer}")
                if 'explanation' in question:
                    st.info(f"💡 {question['explanation']}")
            
            # Add a button to continue to the next question
            if st.button("Next Question ➡️", key=f"next_{quiz.current_question}"):
                st.session_state.show_feedback = False
                st.rerun()
            return
    
    # Display answer options
    selected = st.radio(
        "Select your answer:",
        question['options'],
        key=f"question_{quiz.current_question}",
        index=None
    )
    
    # Submit button
    if st.button("Submit Answer", disabled=selected is None, key=f"submit_{quiz.current_question}"):
        selected_idx = question['options'].index(selected)
        is_correct = quiz.submit_answer(selected_idx)
        st.session_state.show_feedback = True
        st.rerun()

def initialize_session_state():
    """Initialize session state variables with error handling."""
    status = st.status("🚀 Initializing application...", expanded=True)
    
    try:
        if 'initialized' not in st.session_state:
            with status:
                st.write("📂 Loading search index...")
                index, model, metadatas = load_faiss_index()
                st.session_state.faiss_data = {
                    'index': index,
                    'model': model,
                    'metadatas': metadatas
                }
                st.success("✓ Search index loaded")
                
                # Initialize audio player
                st.write("🔊 Initializing audio player...")
                st.session_state.audio_player = AudioPlayer(DATA_DIR)
                
                # Initialize quiz with Mistral API key
                st.write("📝 Setting up quiz...")
                mistral_api_key = os.getenv("MISTRAL_API_KEY")
                if not mistral_api_key and os.path.exists("../../.env"):
                    # Try to load from parent directory's .env if exists
                    load_dotenv("../../.env")
                    mistral_api_key = os.getenv("MISTRAL_API_KEY")
                
                if not mistral_api_key:
                    st.warning("Mistral API key not found. Using default questions.")
                
                st.session_state.quiz = Quiz(DATA_DIR, mistral_api_key)
                
                if not st.session_state.quiz.load_questions(JPLT_N3_VIDEO_ID):
                    with st.spinner("Generating quiz questions..."):
                        st.session_state.quiz.generate_sample_questions(
                            st.session_state.faiss_data['metadatas'],
                            num_questions=5
                        )
                
                st.success("✓ Quiz ready")
                st.session_state.initialized = True
                st.rerun()
    except Exception as e:
        status.error(f"Error initializing application: {str(e)}")
        st.stop()

# Main App
def main():
    st.title("🇯🇵 Japanese Learning Companion")
    
    # Initialize session state
    initialize_session_state()
    
    # Navigation
    menu = option_menu(
        menu_title=None,
        options=["Home", "Quiz", "About"],
        icons=["house", "book", "info-circle"],
        menu_icon="cast",
        default_index=0,
        orientation="horizontal"
    )
    
    # Get FAISS data
    faiss_data = st.session_state.faiss_data
    
    # Show selected page
    if menu == "Home":
        show_home(faiss_data['index'], faiss_data['model'], faiss_data['metadatas'])
    elif menu == "Quiz":
        show_quiz(st.session_state.quiz, st.session_state.audio_player)
    elif menu == "About":
        st.markdown("""
        ## About Japanese Learning Companion
        
        This app helps you practice Japanese listening comprehension using JPLT N3 materials.
        
        ### Features:
        - Search through transcript segments
        - Listen to text with high-quality TTS
        - Practice with interactive quizzes
        - Track your progress
        
        ### How to use:
        1. Use the **Search** tab to find specific phrases
        2. Try the **Quiz** to test your comprehension
        3. Click the audio icon to hear the text
        """)

if __name__ == "__main__":
    main()
