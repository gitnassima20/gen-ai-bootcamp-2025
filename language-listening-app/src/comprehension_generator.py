from transformers import pipeline
import random

class ComprehensionGenerator:
    def __init__(self, model_name='rinna/japanese-gpt2-medium'):
        """
        Initialize text generation and summarization model
        
        Args:
            model_name (str): Hugging Face model for text generation
        """
        self.generator = pipeline('text-generation', model=model_name)
        self.summarizer = pipeline('summarization', model=model_name)

    def generate_questions(self, text, num_questions=3):
        """
        Generate listening comprehension questions
        
        Args:
            text (str): Source text
            num_questions (int): Number of questions to generate
        
        Returns:
            list: Generated comprehension questions
        """
        # Summarize text first to extract key points
        summary = self.summarizer(text, max_length=150, min_length=50, do_sample=False)[0]['summary_text']
        
        # Generate questions based on summary
        questions = []
        for _ in range(num_questions):
            prompt = f"次の文章に基づいて理解度を測る質問を生成してください: {summary}"
            question = self.generator(prompt, max_length=100, num_return_sequences=1)[0]['generated_text']
            questions.append(question.split(prompt)[-1].strip())
        
        return questions

    def generate_multiple_choice(self, text, num_questions=3, choices=4):
        """
        Generate multiple-choice comprehension questions
        
        Args:
            text (str): Source text
            num_questions (int): Number of questions to generate
            choices (int): Number of answer choices
        
        Returns:
            list: Multiple-choice comprehension questions
        """
        questions = []
        key_sentences = text.split('.')[:5]  # Use first few sentences as basis
        
        for sentence in key_sentences:
            if len(sentence.split()) > 5:  # Ensure meaningful sentence
                correct_answer = sentence.strip()
                wrong_answers = [
                    f"Incorrect statement about {sentence.split()[0]}"
                    for _ in range(choices - 1)
                ]
                
                # Shuffle answers
                all_choices = [correct_answer] + wrong_answers
                random.shuffle(all_choices)
                
                question = {
                    'question': f"Which statement best describes the text?",
                    'choices': all_choices,
                    'correct_answer': correct_answer
                }
                questions.append(question)
        
        return questions[:num_questions]
