"""Download and process YouTube transcripts for language learning."""

import os
import json
from pathlib import Path
from typing import Dict, List, Optional, TypedDict

from youtube_transcript_api import YouTubeTranscriptApi
from youtube_transcript_api.formatters import JSONFormatter

from config import DATA_DIR, JPLT_N3_VIDEO_ID


class TranscriptSegment(TypedDict):
    """A single segment of a transcript with timing and text."""
    text: str
    start: float
    duration: float


def fetch_youtube_transcript(video_id: str) -> List[TranscriptSegment]:
    """Fetch transcript for a YouTube video.
    
    Args:
        video_id: YouTube video ID (from URL)
        
    Returns:
        List of transcript segments with text and timing
    """
    try:
        # Fetch available transcripts
        transcript_list = YouTubeTranscriptApi.list_transcripts(video_id)
        
        # Try Japanese first, fallback to English if not available
        try:
            transcript = transcript_list.find_transcript(['ja', 'en'])
            print(f"Found transcript in: {transcript.language}")
            
            # Fetch the actual transcript data
            segments = transcript.fetch()
            
            # Convert transcript segments to our format
            result = []
            for seg in segments:
                try:
                    # Handle both dictionary and object access
                    if hasattr(seg, 'text'):  # It's an object
                        text = seg.text
                        start = seg.start
                        duration = seg.duration
                    else:  # It's a dictionary
                        text = seg.get('text', '')
                        start = seg.get('start', 0)
                        duration = seg.get('duration', 0)
                    
                    text = text.strip()
                    if text:  # Only include non-empty segments
                        result.append({
                            'text': text,
                            'start': float(start),
                            'duration': float(duration)
                        })
                        
                except Exception as seg_error:
                    print(f"Warning: Could not process segment: {seg_error}")
                    continue
            
            if not result:
                print("Warning: No valid segments found in transcript")
                
            return result
            
        except Exception as e:
            print(f"Error finding or processing transcript: {e}")
            return []
            
    except Exception as e:
        print(f"Error listing available transcripts: {e}")
        return []


def save_transcript(transcript: List[TranscriptSegment], output_path: Path) -> None:
    """Save transcript to a JSON file."""
    output_path.parent.mkdir(parents=True, exist_ok=True)
    
    with open(output_path, 'w', encoding='utf-8') as f:
        json.dump(transcript, f, ensure_ascii=False, indent=2)
    
    print(f"Saved transcript to {output_path}")


def main():
    # Create data directory if it doesn't exist
    data_dir = Path(DATA_DIR)
    data_dir.mkdir(exist_ok=True)
    
    # Path for the transcript JSON
    transcript_path = data_dir / f"{JPLT_N3_VIDEO_ID}.json"
    
    # Fetch and save the transcript
    print(f"Fetching transcript for video: {JPLT_N3_VIDEO_ID}")
    transcript = fetch_youtube_transcript(JPLT_N3_VIDEO_ID)
    
    if transcript:
        save_transcript(transcript, transcript_path)
        print(f"Found {len(transcript)} transcript segments.")
    else:
        print("No transcript found or an error occurred.")


if __name__ == "__main__":
    main()
