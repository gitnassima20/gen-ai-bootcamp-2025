import os
import logging
from youtube_transcript_api import YouTubeTranscriptApi

logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

def get_transcript(video_id, languages=None):
    try:
        # Fetch the Japanese transcript
        if languages is None:
            transcript = YouTubeTranscriptApi.get_transcript(video_id)
        else:
            transcript = YouTubeTranscriptApi.get_transcript(video_id, languages=languages)
        return transcript
    except Exception as e:
        print(f"Error downloading transcript: {e}")
        return None

class YouTubeTranscriptDownloader:
    @staticmethod
    def get_transcript(video_id, languages=None):
        """
        Download transcript for a given YouTube video ID
        
        Args:
            video_id (str): YouTube video ID
            languages (list): List of languages to fetch transcript for
        
        Returns:
            list: Transcript text with timestamps
        """
        try:
            transcript = get_transcript(video_id, languages)
            return transcript
        except Exception as e:
            print(f"Error downloading transcript: {e}")
            return None

    @staticmethod
    def extract_text_from_transcript(transcript):
        """
        Extract clean text from transcript
        
        Args:
            transcript (list): Transcript with timestamps
        
        Returns:
            str: Concatenated transcript text
        """
        if not transcript:
            return ""
        return " ".join([entry['text'] for entry in transcript])

    @staticmethod
    def save_transcript_to_file(video_id, transcript, languages=None):
        """
        Save transcript to a text file in the src/transcripts directory
        
        Args:
            video_id (str): YouTube video ID
            transcript (list): Transcript with timestamps
            languages (list, optional): List of languages used for the transcript
        
        Returns:
            str: Path to the saved transcript file
        """
        if not transcript:
            logger.warning(f"No transcript available for video {video_id}")
            return None
        
        # Extract text from transcript
        transcript_text = YouTubeTranscriptDownloader.extract_text_from_transcript(transcript)
        
        # Create transcripts directory if it doesn't exist
        transcripts_dir = os.path.join(os.path.dirname(__file__), 'transcripts')
        os.makedirs(transcripts_dir, exist_ok=True)
        
        # Create filename with video ID and optional language
        filename = f"{video_id}"
        if languages:
            filename += f"_{'-'.join(languages)}"
        filename += ".txt"
        
        # Full path to save the transcript
        file_path = os.path.join(transcripts_dir, filename)
        
        # Write transcript to file
        try:
            with open(file_path, 'w', encoding='utf-8') as f:
                f.write(transcript_text)
            logger.info(f"Transcript saved to {file_path}")
            return file_path
        except Exception as e:
            logger.error(f"Error saving transcript: {e}")
            return None
