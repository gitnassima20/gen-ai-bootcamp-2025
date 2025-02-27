# Language Listening Comprehension App

## Overview
This application helps language learners practice listening comprehension by generating exercises from YouTube video transcripts.

## Features
- Download YouTube video transcripts
- Generate multiple-choice and open-ended comprehension questions
- Text-to-Speech functionality
- Vector-based transcript search and similarity

## Setup

### Prerequisites
- Python 3.8+
- pip

### Installation
1. Clone the repository
2. Create a virtual environment
```bash
python -m venv venv
source venv/bin/activate  # On Windows use `venv\Scripts\activate`
```

3. Install dependencies
```bash
pip install -r requirements.txt
```

### Running the App
```bash
streamlit run src/app.py
```

## Dependencies
- YouTube Transcript API
- Vosk (Speech Recognition)
- Pyttsx3 (Text-to-Speech)
- Transformers (LLM)
- Streamlit (Frontend)
- SQLite Vector (Knowledge Base)

## Limitations
- Requires internet connection
- Transcript quality depends on source video
- Limited language support based on model capabilities

## Future Improvements
- Multi-language support
- More advanced question generation
- Adaptive difficulty levels

## Videos Links:

- JPLT N3: <https://www.youtube.com/watch?v=lasKN-LsJwQ>
- JPLT N4: <https://www.youtube.com/watch?v=F4sqJAPyB4o>
- JPLT N5: <https://www.youtube.com/watch?v=sY7L5cfCWno>

## Video Structure:

1. Introduction
2. Conversation
3. Questions