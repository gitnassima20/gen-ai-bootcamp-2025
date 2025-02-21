import os
from youtube_transcript_api import YouTubeTranscriptApi

class YouTubeTranscriptDownloader:
    @staticmethod
    def get_transcript(video_id):
        """
        Download transcript for a given YouTube video ID
        
        Args:
            video_id (str): YouTube video ID
        
        Returns:
            list: Transcript text with timestamps
        """
        try:
            transcript = YouTubeTranscriptApi.get_transcript(video_id)
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
