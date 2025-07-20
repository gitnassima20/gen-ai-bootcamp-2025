"""Quiz functionality for the language learning app."""

from typing import Dict, List, Optional, TypedDict
from pathlib import Path
from typing import Dict, List, Optional, Tuple, Any
from pathlib import Path
import json
import random
import time
from datetime import datetime, timedelta

# For LLM integration

__all__ = ['Quiz', 'QuizQuestion']

from mistralai.client import MistralClient
import json
from typing import Dict, Any, List, Optional

class QuestionGenerator:
    """Handles generating questions and answers using Mistral AI."""
    
    def __init__(self, api_key: Optional[str] = None):
        self.client = MistralClient(api_key=api_key) if api_key else None
        self.model = "mistral-medium"  # or "mistral-large-latest" for better quality
    
    def _call_mistral(self, messages: List[Dict[str, str]]) -> str:
        """Make a call to Mistral AI API."""
        if not self.client:
            raise ValueError("Mistral API key not provided")
            
        try:
            # Convert messages to the format expected by Mistral client
            chat_messages = []
            for msg in messages:
                role = "user" if msg["role"] == "user" else "assistant"
                chat_messages.append({"role": role, "content": msg["content"]})
            
            chat_response = self.client.chat(
                model=self.model,
                messages=chat_messages,
                temperature=0.7,
            )
            return chat_response.choices[0].message.content
        except Exception as e:
            print(f"Error calling Mistral API: {e}")
            raise
    
    def generate_question(self, conversation: str) -> Dict[str, Any]:
        """Generate a comprehension question based on the conversation."""
        if not self.client:
            return self._get_default_question()
            
        try:
            # First, generate the question and options
            system_prompt = """You are a Japanese language teaching assistant. Generate a multiple-choice 
                question in English about the following Japanese conversation. The question should test 
                the listener's comprehension of the main topic or key details. Include 4 answer choices 
                where only one is correct. Format your response as a JSON object with 'question', 
                'options' (array), and 'correct_index' (0-3) fields. Only return the JSON object, 
                no other text."""
                
            messages = [
                {"role": "system", "content": system_prompt},
                {"role": "user", "content": f"Conversation in Japanese:\n{conversation}"}
            ]
            
            response = self._call_mistral(messages)
            
            # Clean up the response to extract JSON
            try:
                # Try to find JSON in the response
                start = response.find('{')
                end = response.rfind('}') + 1
                result = json.loads(response[start:end])
            except (json.JSONDecodeError, ValueError) as e:
                print(f"Error parsing JSON from Mistral: {e}")
                return self._get_default_question()
            
            # Add explanation
            explanation_prompt = f"""Explain why the correct answer is right and the others are wrong 
                in one concise sentence. Return only the explanation, no other text.
                
                Question: {result.get('question', '')}
                Correct answer: {result.get('options', [])[result.get('correct_index', 0)]}
                Other options: {[opt for i, opt in enumerate(result.get('options', [])) if i != result.get('correct_index', 0)]}"""
                
            explanation = self._call_mistral([
                {"role": "user", "content": explanation_prompt}
            ])
            
            result['explanation'] = explanation.strip()
            return result
            
        except Exception as e:
            print(f"Error generating question with Mistral: {e}")
            return self._get_default_question()
    
    def _get_default_question(self) -> Dict[str, Any]:
        """Return a default question if LLM generation fails."""
        return {
            'question': 'What is the main topic of this conversation?',
            'options': [
                'A casual conversation between friends',
                'A business meeting',
                'A classroom discussion',
                'A news report'
            ],
            'correct_index': 0,
            'explanation': 'The conversation appears to be a casual exchange between friends.'
        }

class QuizQuestion(TypedDict):
    """A single quiz question with options and correct answer."""
    question: str
    options: List[str]
    correct_index: int
    explanation: str

class Quiz:
    """Manages quiz questions and user progress."""
    
    def __init__(self, data_dir: Path, openai_api_key: Optional[str] = None):
        self.data_dir = data_dir
        self.questions: List[QuizQuestion] = []
        self.current_question = 0
        self.score = 0
        self.answers = []
        self.question_generator = QuestionGenerator(openai_api_key) if openai_api_key else None
        
    def load_questions(self, video_id: str) -> bool:
        """Load questions from a JSON file."""
        quiz_file = self.data_dir / f"{video_id}_quiz.json"
        if not quiz_file.exists():
            return False
            
        with open(quiz_file, 'r', encoding='utf-8') as f:
            self.questions = json.load(f)
        return True
    
    def _generate_conversation(self, segments: List[dict], min_duration: int = 30) -> Tuple[str, List[dict]]:
        """Generate a conversation between two people with at least min_duration seconds."""
        conversation = []
        current_duration = 0
        
        for segment in segments:
            if current_duration >= min_duration:
                break
                
            # Split segment into sentences for more natural conversation
            sentences = [s.strip() + ('。' if not s.endswith(('。', '？', '！')) else '') 
                        for s in segment.get('text', '').replace('。', '。\n').split('\n') 
                        if s.strip()]
            
            # Alternate between two speakers
            for j, sentence in enumerate(sentences):
                speaker = "A" if j % 2 == 0 else "B"
                duration = max(1.0, len(sentence) * 0.2)  # More realistic timing
                
                conversation.append({
                    'speaker': speaker,
                    'text': sentence,
                    'duration': duration,
                    'gender': 'female' if speaker == 'A' else 'male'  # Speaker A: female, B: male
                })
                current_duration += duration
                
                if current_duration >= min_duration:
                    break
        
        # Format conversation with speaker labels
        formatted_text = "\n\n".join(
            f"{line['speaker']}: {line['text']}" 
            for line in conversation
        )
        
        return formatted_text, conversation

    def generate_sample_questions(self, metadata: List[dict], num_questions: int = 5) -> None:
        """Generate questions from transcript metadata with realistic conversations."""
        if not metadata:
            return
            
        # Group segments into conversations of at least 30 seconds
        current_segments = []
        current_duration = 0
        generated_count = 0
        
        for segment in metadata:
            if generated_count >= num_questions:
                break
                
            segment_duration = segment.get('duration', 0)
            
            if current_duration + segment_duration >= 30 or not current_segments:
                # Start a new conversation
                current_batch = current_segments if current_segments else [segment]
                conversation_text, conversation_data = self._generate_conversation(
                    current_batch,
                    min_duration=30
                )
                
                # Generate question using LLM if available, otherwise use default
                if self.question_generator:
                    try:
                        qa = self.question_generator.generate_question(conversation_text)
                    except Exception as e:
                        print(f"Error generating question with LLM: {e}")
                        qa = None
                else:
                    qa = None
                
                if not qa:
                    # Fallback to default question
                    qa = {
                        'question': 'What is the main topic of this conversation?',
                        'options': [
                            'A casual conversation between friends',
                            'A business meeting',
                            'A classroom discussion',
                            'A news report'
                        ],
                        'correct_index': 0,
                        'explanation': 'The conversation appears to be a casual exchange between friends.'
                    }
                
                # Add question for this conversation
                self.questions.append({
                    'id': f"q{len(self.questions) + 1}",
                    'japanese_text': conversation_text,
                    'conversation': conversation_data,
                    'question': qa['question'],
                    'options': qa['options'],
                    'correct_index': qa['correct_index'],
                    'explanation': qa['explanation'],
                    'segment': current_batch[0] if current_batch else segment,
                    'timestamp': datetime.now().isoformat()
                })
                
                generated_count += 1
                current_segments = []
                current_duration = 0
                
                # Don't skip the current segment if we just started with it
                if current_duration == 0 and segment_duration < 30:
                    current_segments.append(segment)
                    current_duration += segment_duration
            else:
                current_segments.append(segment)
                current_duration += segment_duration
    
    def get_current_question(self) -> Optional[Dict]:
        """Get the current question."""
        if 0 <= self.current_question < len(self.questions):
            return self.questions[self.current_question]
        return None
    
    def submit_answer(self, answer_index: int) -> bool:
        """Submit an answer and return if it was correct."""
        if not 0 <= self.current_question < len(self.questions):
            return False
            
        is_correct = (answer_index == self.questions[self.current_question]['correct_index'])
        
        self.answers.append({
            'question_idx': self.current_question,
            'answer_idx': answer_index,
            'is_correct': is_correct,
            'timestamp': datetime.now().isoformat()
        })
        
        if is_correct:
            self.score += 1
            
        self.current_question += 1
        return is_correct
    
    def get_progress(self) -> Dict[str, int]:
        """Get quiz progress information."""
        total = len(self.questions)
        current = min(self.current_question, total)
        return {
            'current': current,
            'total': total,
            'score': self.score,
            'progress': (current / total * 100) if total > 0 else 0
        }
    
    def is_complete(self) -> bool:
        """Check if the quiz is complete."""
        return self.current_question >= len(self.questions)
    
    def reset(self) -> None:
        """Reset the quiz state."""
        self.current_question = 0
        self.score = 0
        self.answers = []
