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
    
    def _call_gemini_api(self, prompt: str) -> str:
        """
        Call Gemini API to generate content based on the given prompt.
        
        Args:
            prompt (str): The prompt to send to the Gemini API
        
        Returns:
            str: The generated response from the API
        """
        try:
            headers = {
                'Content-Type': 'application/json'
            }
            
            data = {
                "contents": [{
                    "parts": [{
                        "text": prompt
                    }]
                }],
                "generationConfig": {
                    "temperature": 0.7,
                    "maxOutputTokens": 800
                }
            }
            
            response = requests.post(
                self.api_endpoint, 
                headers=headers, 
                json=data
            )
            
            response.raise_for_status()
            
            # Extract the text from the response
            response_json = response.json()
            generated_text = response_json['candidates'][0]['content']['parts'][0]['text']
            
            logging.info(f"Raw Gemini Response: {generated_text}")
            
            return generated_text
        
        except Exception as e:
            logging.error(f"Error calling Gemini API: {e}")
            return "{}"  # Return an empty JSON string to prevent breaking the process

    def structure_transcript(self, transcript: str) -> tuple:
        # Prompt to extract all conversation cases
        prompt = f"""Analyze the entire transcript and systematically extract information for EACH CONVERSATION CASE:

        For each case, provide:
        1. Situation description
        2. Full conversation details
        3. Specific question to be answered

        Return the result as a valid JSON array.
        """

        try:
            # Call Gemini API to extract structured data
            structured_response = self._call_gemini_api(prompt + "\n\nTranscript:\n" + transcript)
            
            try:
                # Try parsing the response as JSON
                structured_data = json.loads(structured_response)
                
                # Standardize the output
                cases = []
                multiple_choice_options = []
                processed_cases = set()  # To track unique cases
                
                for case in structured_data:
                    # Extract case details with more flexible parsing
                    situation = case.get('Situation description', case.get('situation', '')).strip()
                    conversation = case.get('Full conversation details', case.get('conversation', '')).strip()
                    question = case.get('Specific question', case.get('question', '')).strip()
                    
                    # Create a unique key for the case
                    case_key = f"{situation}_{conversation}_{question}"
                    
                    # Skip only if ALL fields are empty
                    if not situation or not conversation or not question:
                        logging.warning(f"Skipping completely empty case")
                        continue
                    
                    # Skip if this exact case has been processed before
                    if case_key in processed_cases:
                        logging.info(f"Skipping duplicate case: {case_key}")
                        continue
                    
                    # Prepare the prompt for generating multiple-choice options
                    mc_prompt = f"""
                    Based on the following Japanese listening comprehension case, generate multiple-choice options:

                    Situation: {situation}
                    Conversation: {conversation}
                    Specific Question: {question}

                    Instructions:
                    1. Create 4 multiple-choice options for the specific question in japanese
                    2. Ensure the correct answer is derived from the conversation
                    3. Make the other options plausible but incorrect
                    4. Format the response as a JSON with the following structure:
                    {{
                        "correct_answer": "...",
                        "options": ["option1", "option2", "option3", "option4"]
                    }}
                    """
                    
                    # Try to generate multiple-choice options, but don't fail the entire process if it doesn't work
                    try:
                        mc_response = self._call_gemini_api(mc_prompt)
                        multiple_choice_data = json.loads(mc_response)
                    except Exception as mc_error:
                        logging.warning(f"Could not generate multiple-choice options: {mc_error}")
                        multiple_choice_data = {
                            "correct_answer": "",
                            "options": []
                        }
                    
                    # Add the case and its multiple-choice options
                    case_data = {
                        "situation": situation,
                        "conversation": conversation,
                        "question": question
                    }
                    cases.append(case_data)
                    multiple_choice_options.append(multiple_choice_data)
                    
                    # Mark this case as processed
                    processed_cases.add(case_key)
                
                # Log the filtering results with more context
                logging.info(f"Filtered cases: Total input cases={len(structured_data)}, Remaining valid cases={len(cases)}")
                
                return cases, multiple_choice_options
            
            except (json.JSONDecodeError, ValueError) as json_error:
                # Fallback manual extraction
                logging.error(f"JSON Parsing error: {json_error}")
                logging.error(f"Problematic response: {structured_response}")
                
                # Print formatted cases to console
                self._print_formatted_cases(cases)
                
                return cases, []
        
        except Exception as e:
            logging.error(f"Error generating structured transcript: {e}")
            
            return [], []

    def generate_multiple_choice_options(self, extracted_cases: list) -> list:
        """
        Generate multiple-choice options for each case
        
        Args:
            extracted_cases (list): List of extracted cases from the transcript
        
        Returns:
            list: A list of dictionaries containing multiple-choice options
        """
        try:
            multiple_choice_options = []
            
            for case in extracted_cases:
                # Prepare the prompt for generating multiple-choice options
                prompt = f"""
                Based on the following Japanese listening comprehension case, generate multiple-choice options:

                Situation: {case.get('situation', '')}
                Conversation: {case.get('conversation', '')}
                Specific Question: {case.get('question', '')}

                Instructions:
                1. Create 4 multiple-choice options for the specific question in japanese
                2. Ensure the correct answer is derived from the conversation
                3. Make the other options plausible but incorrect
                4. Format the response as a JSON with the following structure:
                {{
                    "question": "...",
                    "options": [
                        {{"text": "Option 1", "is_correct": false}},
                        {{"text": "Option 2", "is_correct": false}},
                        {{"text": "Option 3", "is_correct": false}},
                        {{"text": "Option 4", "is_correct": true}}
                    ]
                }}
                """
                
                # Use Gemini to generate options
                generation_response = self._generate_with_gemini(
                    prompt, 
                    temperature=0.3,  # Lower temperature for more precise answers
                    max_tokens=300
                )
                
                # Parse the JSON response
                options_data = json.loads(generation_response)
                
                multiple_choice_options.append(options_data)
            
            return multiple_choice_options
        
        except Exception as e:
            logging.error(f"Error generating multiple-choice options: {e}")
            return []