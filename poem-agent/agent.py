import torch
from transformers import AutoModelForCausalLM, AutoTokenizer, AutoModelForSeq2SeqLM
from typing import List, Dict, Any
import yt_dlp
import whisper
import pytesseract
from PIL import Image
import json
import fugashi

class MultimediaAgent:
    def __init__(self, 
                 qwen_model_path: str = "Qwen/Qwen1.5-1.8B",
                 whisper_model: str = "base",
                 ocr_lang: str = "jpn"):
        """
        Initialize multimedia processing agent with multiple specialized models
        
        Args:
            qwen_model_path: Path to Qwen language model
            whisper_model: Whisper model size for transcription
            ocr_lang: OCR language setting
        """
        # Qwen Model Setup
        self.tokenizer = AutoTokenizer.from_pretrained(qwen_model_path)
        self.qwen_model = AutoModelForCausalLM.from_pretrained(
            qwen_model_path, 
            device_map="auto",
            torch_dtype=torch.float16
        )
        
        # Whisper Model Setup
        self.whisper_model = whisper.load_model(whisper_model)
        
        # OCR Configuration
        pytesseract.pytesseract.tesseract_cmd = r'/usr/bin/tesseract'
        self.ocr_lang = ocr_lang
        
        # Define available tools
        self.tools = [
            {
                "name": "download_media",
                "description": "Download media from URL",
                "parameters": {
                    "url": {"type": "string", "description": "Media URL"},
                    "media_type": {"type": "string", "enum": ["video", "audio"]}
                }
            },
            {
                "name": "transcribe_audio",
                "description": "Transcribe audio to text",
                "parameters": {
                    "audio_path": {"type": "string", "description": "Path to audio file"},
                    "language": {"type": "string", "description": "Language of audio"}
                }
            },
            {
                "name": "ocr_extract",
                "description": "Extract text from image or PDF",
                "parameters": {
                    "file_path": {"type": "string", "description": "Path to image/PDF"},
                    "page_range": {"type": "array", "description": "Optional page range"}
                }
            },
            {
                "name": "normalize_text",
                "description": "Normalize Japanese text",
                "parameters": {
                    "text": {"type": "string", "description": "Text to normalize"}
                }
            }
        ]
    
    def download_media(self, url: str, media_type: str = "audio") -> str:
        """
        Download media from URL using yt-dlp
        
        Args:
            url: Media URL
            media_type: Type of media to download
        
        Returns:
            Path to downloaded media file
        """
        ydl_opts = {
            'format': 'bestaudio/best' if media_type == 'audio' else 'best',
            'outtmpl': f'downloads/%(title)s.%(ext)s'
        }
        
        with yt_dlp.YoutubeDL(ydl_opts) as ydl:
            info = ydl.extract_info(url, download=True)
            return info['requested_downloads'][0]['filepath']
    
    def transcribe_audio(self, audio_path: str, language: str = "ja") -> str:
        """
        Transcribe audio using Whisper
        
        Args:
            audio_path: Path to audio file
            language: Language of audio
        
        Returns:
            Transcribed text
        """
        result = self.whisper_model.transcribe(audio_path, language=language)
        return result['text']
    
    def ocr_extract(self, file_path: str, page_range: List[int] = None) -> str:
        """
        Extract text from image or PDF using OCR
        
        Args:
            file_path: Path to file
            page_range: Optional page range for multi-page documents
        
        Returns:
            Extracted text
        """
        # Simple OCR implementation
        try:
            image = Image.open(file_path)
            text = pytesseract.image_to_string(image, lang=self.ocr_lang)
            return text
        except Exception as e:
            return f"OCR Error: {str(e)}"
    
    def normalize_text(self, text: str) -> str:
        """
        Normalize Japanese text
        
        Args:
            text: Input text to normalize
        
        Returns:
            Normalized text
        """
        # Basic normalization (can be expanded)
        tokenizer = fugashi.Dictionary().create()
    
        # Normalize text:
        # 1. Convert to full-width characters
        # 2. Remove redundant whitespaces
        # 3. Lowercase (if needed)
        # 4. Remove repeated phrases
        
        # Split into morphemes
        morphemes = tokenizer.tokenize(text)
        
        # Extract base forms
        normalized_words = [m.dictionary_form() for m in morphemes]
        
        # Rejoin and clean
        normalized_text = ''.join(normalized_words)
        
        return normalized_text
    
    def process_workflow(self, url: str) -> Dict[str, Any]:
        """
        Complete multimedia processing workflow
        
        Args:
            url: Source media URL
        
        Returns:
            Processed multimedia data
        """
        # Download media
        media_path = self.download_media(url)
        
        # Transcribe audio
        transcription = self.transcribe_audio(media_path)
        
        # Normalize text
        normalized_text = self.normalize_text(transcription)
        
        return {
            "original_url": url,
            "media_path": media_path,
            "transcription": transcription,
            "normalized_text": normalized_text
        }


def main():
    agent = MultimediaAgent()
    result = agent.process_workflow("https://www.youtube.com/watch?v=aZds4UrUuko&list=PLuNFC5RjK4NS0moKF-p1YkKJjcXU3ztvF")
    print(json.dumps(result, indent=2, ensure_ascii=False))

if __name__ == "__main__":
    main()