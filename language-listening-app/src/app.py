import streamlit as st
import pyttsx3
from youtube_transcript import YouTubeTranscriptDownloader
from vector_store import VectorStore
from comprehension_generator import ComprehensionGenerator
from structred_data import TranScriptStructure

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
        self.transcript_analyzer = TranScriptStructure()

    def process_structured_transcript(self, video_id, language_level):
        # Fetch the transcript
        full_transcript = self.transcript_downloader.get_transcript(video_id, languages=['ja'])
        
        if not full_transcript:
            return None
        # Save the transcript to a file
        transcript_file = self.transcript_downloader.save_transcript_to_file(
            video_id, 
            full_transcript, 
            languages=['ja']
        )
        
        # Convert transcript to text
        transcript_text = " ".join([entry['text'] for entry in full_transcript])
        
        # Use TranScriptStructure to analyze the transcript
        structured_response = self.transcript_analyzer.structure_transcript(transcript_text)
        
        return structured_response

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
            # Display the structured transcript as text
            st.text_area("Structured Transcript Analysis", structured_content, height=400)

if __name__ == "__main__":
    main()