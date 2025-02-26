import os
from youtube_transcript_api import YouTubeTranscriptApi

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
