import json
import logging
import os
import requests
from dotenv import load_dotenv
from datetime import datetime

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
    
    def structure_transcript(self, transcript: str) -> list:
        # Prompt to extract all conversation cases
        prompt = f"""Analyze the entire transcript and systematically extract information for EACH CONVERSATION CASE:

For EACH case (Case 0, Case 1, etc.), extract:
1. Situation description
2. Full conversation details
3. Specific question for that case

IMPORTANT GUIDELINES:
- Skip initial test instructions
- Extract EXACTLY as they appear in the original transcript
- Provide a structured output for EACH case
- Include all cases present in the transcript
- Maintain original language

Transcript: {transcript}
"""
        
        try:
            # Prepare the request payload
            payload = {
                "contents": [{
                    "parts": [{"text": prompt}]
                }],
                "generationConfig": {
                    "temperature": 0.1,  # Extremely low for precise extraction
                    "maxOutputTokens": 4096  # Increased to accommodate multiple cases
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
            
            # Save raw Gemini response to file
            self._save_gemini_response(response_data)
            
            # Log the raw response for debugging
            logging.info(f"Raw Gemini Response: {generated_text}")
            
            # Try to parse the response as JSON
            try:
                # Remove any code block formatting
                clean_text = generated_text.strip('```json\n```')
                structured_data = json.loads(clean_text)
                
                # Ensure we have a list of cases
                if not isinstance(structured_data, list):
                    structured_data = [structured_data]
                
                # Standardize the output
                cases = []
                for case in structured_data:
                    cases.append({
                        "situation": case.get("situation", ""),
                        "conversation": case.get("conversation", ""),
                        "question": case.get("question", "")
                    })
                
                # Save cases to a text file
                self._save_cases_to_file(cases, response_data)
                
                # Print formatted cases to console
                self._print_formatted_cases(cases)
                
                return cases
            
            except (json.JSONDecodeError, ValueError) as json_error:
                # Fallback manual extraction
                logging.warning(f"JSON parsing failed: {json_error}")
                
                # Manual case extraction
                cases = self._extract_cases(response_data)
                
                # Save cases to a text file
                self._save_cases_to_file(cases, response_data)
                
                # Print formatted cases to console
                self._print_formatted_cases(cases)
                
                return cases
        
        except Exception as e:
            logging.error(f"Error generating structured transcript: {e}")
            logging.exception("Full error traceback:")
            
            return []
    
    def _save_gemini_response(self, response_data: dict):
        """
        Save the raw Gemini API response to a file in the transcripts folder
        """
        try:
            # Ensure the transcripts directory exists
            transcripts_dir = os.path.join(os.path.dirname(__file__), 'transcripts')
            os.makedirs(transcripts_dir, exist_ok=True)
            
            # Generate filename with timestamp
            timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
            filename = f"gemini_raw_response_{timestamp}.json"
            filepath = os.path.join(transcripts_dir, filename)
            
            # Write raw response to file
            with open(filepath, 'w', encoding='utf-8') as f:
                json.dump(response_data, f, ensure_ascii=False, indent=2)
            
            logging.info(f"Raw Gemini response saved to {filepath}")
        
        except Exception as e:
            logging.error(f"Error saving Gemini response to file: {e}")
    
    def _extract_cases(self, response_data: dict) -> list:
        """
        Extract cases from the Gemini API response
        """
        try:
            # Directly use the parts from the response
            cases = []
            for part in response_data.get('parts', []):
                if isinstance(part, dict):
                    cases.append({
                        'situation': part.get('Situation description', ''),
                        'conversation': part.get('Full conversation details', ''),
                        'question': part.get('Specific question', '')
                    })
            
            return cases
        
        except Exception as e:
            logging.error(f"Error extracting cases: {e}")
            return []