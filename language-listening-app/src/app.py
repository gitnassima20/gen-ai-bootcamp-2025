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
        """
        Process and structure YouTube transcript
        
        Args:
            video_id (str): YouTube video ID
            language_level (str): Language difficulty level
        
        Returns:
            dict: Structured transcript content
        """
        # Download full transcript
        full_transcript = self.transcript_downloader.get_transcript(video_id)
        
        if not full_transcript:
            return None
        
        # Extract text
        full_text = self.transcript_downloader.extract_text_from_transcript(full_transcript)
        
        # Split transcript into sections (simplified approach)
        sections = self._split_transcript_into_sections(full_text)
        
        # Generate additional questions if not provided in transcript
        generated_questions = self.comprehension_generator.generate_multiple_choice(
            sections['conversation'], num_questions=3
        )
        
        # Print structured content before insertion
        print(f"""
===== Structured Transcript Content =====
Video ID: {video_id}
Language Level: {language_level}

--- Introduction ---
{sections['introduction']}

--- Conversation ---
{sections['conversation']}

--- Generated Questions ---
{generated_questions}
===============================
""")
        
        # Store structured content
        self.vector_store.insert_structured_transcript(
            video_id=video_id,
            language_level=language_level,
            introduction=sections['introduction'],
            conversation=sections['conversation'],
            questions=generated_questions
        )
        
        return sections

    def _split_transcript_into_sections(self, full_text):
        """
        Split transcript into introduction, conversation, and questions sections
        
        Args:
            full_text (str): Complete transcript text
        
        Returns:
            dict: Structured transcript sections
        """
        # Very basic splitting logic - you might want to improve this
        sentences = full_text.split('.')
        
        return {
            'introduction': '. '.join(sentences[:3]) + '.',
            'conversation': '. '.join(sentences[3:-3]) + '.',
            'questions': '. '.join(sentences[-3:]) + '.'
        }

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
    st.title("Structured Language Listening Comprehension App")
    
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
            st.subheader("Introduction")
            st.write(structured_content['introduction'])
            
            st.subheader("Conversation")
            st.write(structured_content['conversation'])
            
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
