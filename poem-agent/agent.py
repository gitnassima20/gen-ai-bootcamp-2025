from transformers import MarianTokenizer, MarianMTModel, BartForConditionalGeneration, BartTokenizer
from diffusers import StableDiffusionPipeline
import pytesseract
import whisper
import yt_dlp
from PIL import Image
import fugashi
from typing import Dict, Any
import json

class MultimediaAgent:
    def __init__(self, 
                 whisper_model: str = "medium",
                 ocr_lang: str = "jpn",
                 sd_model_path: str = "stabilityai/stable-diffusion-2"):
        """
        Initialize multimedia processing agent
        """

        # Whisper Model Setup
        self.whisper_model = whisper.load_model(whisper_model)

        # OCR Configuration
        pytesseract.pytesseract.tesseract_cmd = r'/usr/bin/tesseract'
        self.ocr_lang = ocr_lang

        # Translation Model
        self.translator = MarianMTModel.from_pretrained("Helsinki-NLP/opus-mt-ja-en")
        self.translation_tokenizer = MarianTokenizer.from_pretrained("Helsinki-NLP/opus-mt-ja-en")

        # Summarization Model Setup (BART)
        self.summarizer = BartForConditionalGeneration.from_pretrained("facebook/bart-large-cnn")
        self.summarizer_tokenizer = BartTokenizer.from_pretrained("facebook/bart-large-cnn")

        # Stable Diffusion Setup
        self.sd_model = StableDiffusionPipeline.from_pretrained("stabilityai/stable-diffusion-2")

    def download_media(self, url: str, media_type: str = "audio") -> str:
        """
        Download media from URL
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
        """
        result = self.whisper_model.transcribe(audio_path, language=language)
        return result['text']

    def normalize_text(self, text: str) -> str:
        """
        Normalize Japanese text
        """
        tokenizer = fugashi.Tagger()
        morphemes = tokenizer(text)
        normalized_words = [m.feature.lemma if m.feature.lemma else m.surface for m in morphemes]
        return ''.join(normalized_words)

    def translate_to_english(self, text: str) -> str:
        """
        Translate Japanese to English
        """
        inputs = self.translation_tokenizer(text, return_tensors="pt")
        outputs = self.translator.generate(**inputs)
        return self.translation_tokenizer.decode(outputs[0], skip_special_tokens=True)

    def summarize_text(self, text: str) -> str:
        """
        Summarize the text to ensure it fits within token limits
        """
        inputs = self.summarizer_tokenizer(text, return_tensors="pt", max_length=1024, truncation=True)
        summary_ids = self.summarizer.generate(**inputs, max_length=150, min_length=50, length_penalty=2.0, num_beams=4, early_stopping=True)
        return self.summarizer_tokenizer.decode(summary_ids[0], skip_special_tokens=True)

    def generate_image(self, prompt: str, output_path: str = "output.png") -> str:
        """
        Generate an image using Stable Diffusion
        """
        # Use the Stable Diffusion pipeline to generate the image.
        image = self.sd_model(prompt).images[0]
        
        # Save the image to the specified output path
        image.save(output_path)
        
        return output_path

    def process_workflow(self, url: str) -> Dict[str, Any]:
        """
        Complete multimedia processing workflow
        """
        media_path = self.download_media(url)
        transcription = self.transcribe_audio(media_path)
        normalized_text = self.normalize_text(transcription)

        translated_text = self.translate_to_english(normalized_text)

        # Summarize the translated text before image generation
        summarized_text = self.summarize_text(translated_text)

        image_path = self.generate_image(summarized_text)

        return {
            "original_url": url,
            "media_path": media_path,
            "transcription": transcription,
            "normalized_text": normalized_text,
            "translated_text": translated_text,
            "summarized_text": summarized_text,
            "image_path": image_path
        }

def main():
    agent = MultimediaAgent()
    result = agent.process_workflow("https://www.youtube.com/watch?v=aZds4UrUuko")
    print(json.dumps(result, indent=2, ensure_ascii=False))

if __name__ == "__main__":
    main()
