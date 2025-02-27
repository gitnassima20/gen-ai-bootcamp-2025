import json
import logging
import os
import re
import requests
from dotenv import load_dotenv

class TranScriptStructure:
    def __init__(self):
        # Load environment variables
        load_dotenv()
        
        # Configure API key
        self.api_key = os.getenv('GOOGLE_API_KEY')
        if not self.api_key:
            raise ValueError("No Google API key found. Please set GOOGLE_API_KEY in .env file.")
        
        # API endpoint for Gemini 2.0 Flash
        self.api_endpoint = f"https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent?key={self.api_key}"
    
    def structure_transcript(self, transcript: str) -> dict:
        # Detailed, structured prompt for Gemini
        prompt = f"""Carefully analyze the following JLPT Listening Comprehension Transcript:

Transcript: {transcript}

IMPORTANT: Provide a detailed, structured JSON response with these exact keys:

{{
    "introduction": "Describe the initial context, setting, and participants of the conversation in detail.",
    "conversation": "Provide a comprehensive summary of the dialogue, including key interactions and main points.",
    "question": "Identify and clearly state the specific comprehension question or challenge in the transcript.",
    "analysis": {{
        "language_skills_tested": "Specify the exact language skills being evaluated",
        "difficulty_level": "Determine the JLPT difficulty level (N5-N1)",
        "key_vocabulary": ["list", "of", "important", "words"],
        "grammatical_points": ["list", "of", "key", "grammar", "structures"]
    }}
}}

Guidelines:
- Be extremely specific and detailed
- Focus on educational insights
- Ensure the JSON is valid and complete
- If no clear question is present, explain why
"""
        
        try:
            # Prepare the request payload
            payload = {
                "contents": [{
                    "parts": [{"text": prompt}]
                }],
                "generationConfig": {
                    "temperature": 0.7,
                    "maxOutputTokens": 2048
                }
            }
            
            # Send request to Gemini API
            response = requests.post(
                self.api_endpoint, 
                json=payload,
                headers={'Content-Type': 'application/json'}
            )
            
            # Check for successful response
            response.raise_for_status()
            
            # Extract the generated text
            response_data = response.json()
            generated_text = response_data['candidates'][0]['content']['parts'][0]['text']
            
            # Log the raw response for debugging
            logging.info(f"Raw Gemini Response: {generated_text}")
            
            # Try to parse the response as JSON
            try:
                # Remove any code block formatting
                clean_text = generated_text.strip('```json\n```')
                structured_data = json.loads(clean_text)
                
                # Validate the structure
                if not all(key in structured_data for key in ['introduction', 'conversation', 'question']):
                    raise ValueError("Missing required keys in JSON response")
                
                return structured_data
            
            except (json.JSONDecodeError, ValueError) as json_error:
                # Fallback parsing if JSON fails
                logging.warning(f"JSON parsing failed: {json_error}")
                logging.warning(f"Attempting manual extraction from text: {generated_text}")
                
                return {
                    "introduction": self._extract_section(generated_text, "introduction"),
                    "conversation": self._extract_section(generated_text, "conversation"),
                    "question": self._extract_section(generated_text, "question")
                }
        
        except Exception as e:
            logging.error(f"Comprehensive error generating structured transcript: {e}")
            # Log the full error details
            logging.exception("Full error traceback:")
            
            return {
                "introduction": f"Error in analysis: {str(e)}",
                "conversation": "",
                "question": ""
            }
    
    def _extract_section(self, text: str, section: str) -> str:
        """
        Manually extract a section from the text if JSON parsing fails
        """
        try:
            # More robust regex to extract section
            pattern = rf'"{section}":\s*"(.*?)"'
            match = re.search(pattern, text, re.DOTALL | re.IGNORECASE)
            return match.group(1).strip() if match else ""
        except Exception as e:
            logging.error(f"Error extracting {section}: {e}")
            return ""