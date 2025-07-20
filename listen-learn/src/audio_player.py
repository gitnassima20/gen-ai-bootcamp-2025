"""Audio player component for the language learning app using gTTS with different voices."""

import os
import tempfile
from pathlib import Path
from typing import Optional, Dict, List, Union

import streamlit as st
from gtts import gTTS
import base64
from pydub import AudioSegment
from pydub.effects import speedup

# Voice configurations for different speakers
VOICE_CONFIGS = {
    'A': {
        'lang': 'ja',
        'tld': 'co.jp',  # Japanese accent
        'slow': False,
        'pitch': 1.0,
    },
    'B': {
        'lang': 'ja',
        'tld': 'co.jp',
        'slow': False,
        'pitch': 0.9,  # Slightly different pitch for second speaker
    },
    'default': {
        'lang': 'ja',
        'tld': 'co.jp',
        'slow': False,
    }
}

class AudioPlayer:
    """Handles text-to-speech and audio playback using gTTS and Streamlit's audio component."""
    
    def __init__(self, data_dir: Path):
        self.data_dir = data_dir
        self.audio_dir = data_dir / "audio"
        self.audio_dir.mkdir(exist_ok=True, parents=True)
        self.initialized = True  # gTTS doesn't require initialization
        
    def text_to_speech(self, text: str, filename: str, speaker: str = 'default') -> Optional[Path]:
        """Convert text to speech with specified speaker voice and save as audio file."""
        if not text:
            return None
            
        audio_path = self.audio_dir / f"{filename}_{speaker}.mp3"
        
        # Only generate if not already exists
        if not audio_path.exists():
            try:
                # Get voice configuration for the speaker
                voice_config = VOICE_CONFIGS.get(speaker, VOICE_CONFIGS['default']).copy()
                
                # Generate speech using gTTS with speaker-specific settings
                tts = gTTS(
                    text=text,
                    lang=voice_config['lang'],
                    tld=voice_config.get('tld', 'com'),
                    slow=voice_config.get('slow', False)
                )
                
                # Save the audio file
                tts.save(str(audio_path))
                
                # Note: gTTS has limited voice customization. For more realistic voices,
                # consider using a paid TTS service like Google Cloud TTS or Amazon Polly.
                
            except Exception as e:
                st.error(f"Error generating speech for speaker {speaker}: {e}")
                return None
                
        return audio_path if audio_path.exists() else None
    
    def play_conversation(self, conversation: List[Dict[str, str]], filename: str) -> None:
        """Play a conversation with different voices for each speaker."""
        if not conversation:
            return
            
        # Create a single audio file for the entire conversation
        combined_audio_path = self.audio_dir / f"{filename}_conversation.mp3"
        
        if not combined_audio_path.exists():
            try:
                combined = AudioSegment.empty()
                
                for i, line in enumerate(conversation):
                    speaker = line.get('speaker', 'A')
                    text = line.get('text', '')
                    
                    if not text.strip():
                        continue
                        
                    # Generate audio for this line
                    line_audio_path = self.text_to_speech(text, f"{filename}_line_{i}", speaker=speaker)
                    if not line_audio_path or not line_audio_path.exists():
                        continue
                        
                    # Load the audio and add a small pause between lines
                    audio = AudioSegment.from_mp3(str(line_audio_path))
                    combined += audio
                    combined += AudioSegment.silent(duration=500)  # 0.5s pause between lines
                
                # Save the combined audio
                combined.export(str(combined_audio_path), format="mp3")
                
            except Exception as e:
                st.error(f"Error combining audio files: {e}")
                return
        
        # Play the combined audio
        if combined_audio_path.exists():
            audio_bytes = combined_audio_path.read_bytes()
            audio_base64 = base64.b64encode(audio_bytes).decode()
            
            audio_html = f"""
            <audio autoplay controls>
                <source src="data:audio/mp3;base64,{audio_base64}" type="audio/mp3">
                Your browser does not support the audio element.
            </audio>
            """
            st.components.v1.html(audio_html, height=50)
        else:
            st.warning("Could not generate conversation audio.")
    
    def play_audio(self, text: str, filename: str, autoplay: bool = True, speaker: str = 'default') -> None:
        """Play text as speech using gTTS with specified speaker."""
        audio_path = self.text_to_speech(text, filename, speaker=speaker)
        if audio_path and audio_path.exists():
            audio_bytes = audio_path.read_bytes()
            audio_base64 = base64.b64encode(audio_bytes).decode()
            
            # Create an audio player with autoplay
            audio_html = f"""
            <audio autoplay={'true' if autoplay else 'false'} controls>
                <source src="data:audio/mp3;base64,{audio_base64}" type="audio/mp3">
                Your browser does not support the audio element.
            </audio>
            """
            st.components.v1.html(audio_html, height=50)
        else:
            st.warning("Could not generate audio for the given text.")
    
    def play_segment(self, segment: dict) -> None:
        """Play audio for a transcript segment."""
        if not segment or 'text' not in segment:
            return
            
        # Create a unique filename based on segment content
        filename = f"segment_{hash(segment['text']) & 0xffffffff:08x}"
        self.play_audio(segment['text'], filename)
        
    def initialize_tts(self):
        """For compatibility with existing code."""
        pass
