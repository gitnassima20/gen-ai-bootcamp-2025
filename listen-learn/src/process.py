"""Process and index transcript data for the language learning app using FAISS."""

import json
import pickle
from pathlib import Path
from typing import Dict, List, TypedDict, Tuple

import numpy as np
import faiss
from sentence_transformers import SentenceTransformer

from config import DATA_DIR, JPLT_N3_VIDEO_ID

# FAISS index file
FAISS_INDEX_FILE = DATA_DIR / "faiss_index.bin"
TRANSCRIPT_METADATA_FILE = DATA_DIR / "transcript_metadata.pkl"


class TranscriptSegment(TypedDict):
    """A single segment of a transcript with timing and text."""
    text: str
    start: float
    duration: float


def load_transcript(transcript_path: Path) -> List[TranscriptSegment]:
    """Load transcript from a JSON file."""
    with open(transcript_path, 'r', encoding='utf-8') as f:
        return json.load(f)


def split_into_sections(transcript: List[TranscriptSegment]) -> Dict[str, List[TranscriptSegment]]:
    """Split transcript into Introduction, Conversation, and Questions sections.
    
    For now, this is a simple time-based split. In a real app, you might want to 
    use more sophisticated methods (e.g., silence detection, topic modeling).
    """
    if not transcript:
        return {}
    
    total_duration = transcript[-1]['start'] + transcript[-1]['duration']
    
    # Simple time-based split (adjust these ratios as needed)
    intro_end = total_duration * 0.2  # First 20% as introduction
    conversation_end = total_duration * 0.7  # Next 50% as conversation
    
    sections = {
        'introduction': [],
        'conversation': [],
        'questions': []
    }
    
    for segment in transcript:
        if segment['start'] < intro_end:
            sections['introduction'].append(segment)
        elif segment['start'] < conversation_end:
            sections['conversation'].append(segment)
        else:
            sections['questions'].append(segment)
    
    return sections


def initialize_faiss(embedding_dim: int = 384) -> Tuple[faiss.Index, SentenceTransformer]:
    """Initialize or load FAISS index and sentence transformer."""
    # Initialize the sentence transformer
    model = SentenceTransformer('all-MiniLM-L6-v2')
    
    # Try to load existing index and metadata
    if FAISS_INDEX_FILE.exists() and TRANSCRIPT_METADATA_FILE.exists():
        print("Loading existing FAISS index...")
        index = faiss.read_index(str(FAISS_INDEX_FILE))
        return index, model
    
    # Create a new index if it doesn't exist
    print("Creating new FAISS index...")
    index = faiss.IndexFlatL2(embedding_dim)
    return index, model


def index_transcript(transcript: List[TranscriptSegment], index: faiss.Index, model: SentenceTransformer) -> None:
    """Index transcript segments using FAISS."""
    if not transcript:
        print("No transcript segments to index.")
        return
    
    # Extract texts and metadata
    texts = [segment['text'] for segment in transcript]
    metadatas = [{
        'start_time': segment['start'],
        'duration': segment['duration'],
        'text': segment['text']
    } for segment in transcript]
    
    # Generate embeddings
    print("Generating embeddings...")
    embeddings = model.encode(texts, show_progress_bar=True, convert_to_numpy=True)
    
    # Add to FAISS index
    if index.ntotal == 0:  # If index is empty
        index.add(embeddings)
    else:  # If adding to existing index
        index.add(embeddings)
    
    # Save metadata
    with open(TRANSCRIPT_METADATA_FILE, 'wb') as f:
        pickle.dump(metadatas, f)
    
    # Save FAISS index
    faiss.write_index(index, str(FAISS_INDEX_FILE))
    print(f"Indexed {len(texts)} segments with FAISS.")


def main():
    # Initialize FAISS
    print("Initializing FAISS...")
    index, model = initialize_faiss()
    
    # Load the transcript
    transcript_path = DATA_DIR / f"{JPLT_N3_VIDEO_ID}.json"
    if not transcript_path.exists():
        print(f"Transcript not found at {transcript_path}. Run ingest.py first.")
        return
    
    print(f"Loading transcript from {transcript_path}...")
    transcript = load_transcript(transcript_path)
    
    # Split into sections
    print("Splitting transcript into sections...")
    sections = split_into_sections(transcript)
    
    # Index each section
    for section_name, section_data in sections.items():
        print(f"Indexing {len(section_data)} segments from {section_name}...")
        index_transcript(section_data, index, model)
    
    print("Indexing complete. FAISS index and metadata have been saved.")
    print(f"Index size: {index.ntotal} vectors")
    
    print("Done!")


if __name__ == "__main__":
    main()
