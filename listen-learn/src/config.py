"""Configuration for the language learning app."""
from pathlib import Path

# Base directory (one level up from src/)
BASE_DIR = Path(__file__).parent.parent

# YouTube video to process (JPLT N3 example from README)
JPLT_N3_VIDEO_ID = "lasKN-LsJwQ"  # From: https://www.youtube.com/watch?v=lasKN-LsJwQ

# Directory for storing data (transcripts, audio, vector DB)
DATA_DIR = BASE_DIR / "data"

# Create data directory if it doesn't exist
DATA_DIR.mkdir(exist_ok=True)

# ChromaDB settings
CHROMA_COLLECTION = "jplt_n3"
