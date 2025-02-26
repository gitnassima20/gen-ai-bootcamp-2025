import streamlit as st
import pyttsx3
from youtube_transcript import YouTubeTranscriptDownloader
from vector_store import VectorStore
from comprehension_generator import ComprehensionGenerator

# Predefined video links by JLPT level
VIDEO_LINKS = {
    "JLPT N5": "sY7L5cfCWno",  # From README
    "JLPT N4": "F4sqJAPyB4o",  # From README
    "JLPT N3": "lasKN-LsJwQ"   # From README
}

class LanguageListeningApp:
    def __init__(self):
        """Initialize app components"""
        self.vector_store = VectorStore()
        self.transcript_downloader = YouTubeTranscriptDownloader()
        self.comprehension_generator = ComprehensionGenerator()
        self.tts_engine = pyttsx3.init()

    def process_structured_transcript(self, video_id, language_level):
        # Fetch the transcript
        full_transcript = self.transcript_downloader.get_transcript(video_id, languages=['ja'])
        
        if not full_transcript:
            return None

        # Initialize sections
        structured_content = []

        # Temporary variables to hold current situation, conversation, and question
        current_situation = None
        current_conversation = None
        current_question = None

        # Logic to structure the transcript
        for entry in full_transcript:
            text = entry['text']
            
            if "男" in text and "女" in text:  # Identify situations
                if current_situation:  # If there is a previous situation, store it
                    structured_content.append({
                        "situation": current_situation,
                        "conversation": current_conversation,
                        "question": current_question
                    })
                current_situation = text  # Update current situation
                current_conversation = None  # Reset conversation and question
                current_question = None
            
            elif "もしもし" in text:  # Check for conversation starters
                current_conversation = text  # Capture conversation
            
            elif "質問" in text:  # Check for question markers
                current_question = text  # Capture question

        # Append the last set if available
        if current_situation:
            structured_content.append({
                "situation": current_situation,
                "conversation": current_conversation,
                "question": current_question
            })

        return structured_content

    def retrieve_similar_content(self, query, search_type='all'):
        """
        Retrieve similar content based on query
        
        Args:
            query (str): Search query
            search_type (str): Section to search
        
        Returns:
            list: Similar content results
        """
        return self.vector_store.search_similar_content(query, search_type)

def main():
    st.title("Language Listening Comprehension App")
    
    app = LanguageListeningApp()
    
    # Video Input Section
    col1, col2 = st.columns(2)
    with col1:
        language_level = st.selectbox("Select Language Level", 
            list(VIDEO_LINKS.keys())
        )
        # Automatically set video ID based on selected language level
        video_id = VIDEO_LINKS[language_level]
        st.write(f"Selected Video: {video_id}")
    
    if st.button("Process Video"):
        with st.spinner("Processing structured transcript..."):
            structured_content = app.process_structured_transcript(video_id, language_level)
        
        if structured_content:
            # Display Sections
            for i, content in enumerate(structured_content):
                st.subheader(f"Section {i+1}")
                st.write(f"Situation: {content['situation']}")
                if content['conversation']:
                    st.write(f"Conversation: {content['conversation']}")
                if content['question']:
                    st.write(f"Question: {content['question']}")
            
            # Semantic Search
            st.subheader("Semantic Content Search")
            search_query = st.text_input("Search similar content")
            if search_query:
                similar_content = app.retrieve_similar_content(search_query)
                for result in similar_content:
                    st.write(f"Video ID: {result['video_id']}")
                    st.write(f"Language Level: {result['language_level']}")
                    st.write("Similarity Scores:")
                    for section, score in result['similarities'].items():
                        st.write(f"- {section}: {score:.2f}")

if __name__ == "__main__":
    main()
