import chromadb
from sentence_transformers import SentenceTransformer
import json
import uuid

class VectorStore:
    def __init__(self, persist_directory='./chroma_db'):
        """
        Initialize ChromaDB vector store
        
        Args:
            persist_directory (str): Directory to persist vector database
        """
        # ChromaDB client
        self.client = chromadb.PersistentClient(path=persist_directory)
        
        # Embedding model
        self.model = SentenceTransformer('all-MiniLM-L6-v2')
        
        # Create collections for different content types
        self.introduction_collection = self.client.get_or_create_collection("introductions")
        self.conversation_collection = self.client.get_or_create_collection("conversations")
        self.questions_collection = self.client.get_or_create_collection("questions")

    def insert_structured_transcript(self, video_id, language_level, introduction, conversation, questions):
        """
        Insert structured transcript with embeddings
        
        Args:
            video_id (str): YouTube video ID
            language_level (str): JLPT level or language difficulty
            introduction (str): Video introduction text
            conversation (str): Main conversation text
            questions (list): List of comprehension questions
        """
        # Generate embeddings
        intro_embedding = self.model.encode(introduction).tolist()
        conv_embedding = self.model.encode(conversation).tolist()
        questions_embedding = self.model.encode(' '.join(questions)).tolist()
        
        # Unique IDs
        intro_id = str(uuid.uuid4())
        conv_id = str(uuid.uuid4())
        questions_id = str(uuid.uuid4())
        
        # Metadata
        metadata = {
            'video_id': video_id,
            'language_level': language_level
        }
        
        # Add to collections
        self.introduction_collection.add(
            embeddings=intro_embedding,
            documents=[introduction],
            metadatas=[metadata],
            ids=[intro_id]
        )
        
        self.conversation_collection.add(
            embeddings=conv_embedding,
            documents=[conversation],
            metadatas=[metadata],
            ids=[conv_id]
        )
        
        self.questions_collection.add(
            embeddings=questions_embedding,
            documents=[json.dumps(questions)],
            metadatas=[metadata],
            ids=[questions_id]
        )

    def search_similar_content(self, query, search_type='all', top_k=5):
        """
        Search for similar content across different sections
        
        Args:
            query (str): Search query
            search_type (str): Section to search ('introduction', 'conversation', 'questions', 'all')
            top_k (int): Number of results to return
        
        Returns:
            list: Similar content results
        """
        query_embedding = self.model.encode(query).tolist()
        
        results = []
        
        # Search based on type
        if search_type in ['introduction', 'all']:
            intro_results = self.introduction_collection.query(
                query_embeddings=[query_embedding],
                n_results=top_k
            )
            results.extend(self._format_results(intro_results, 'introduction'))
        
        if search_type in ['conversation', 'all']:
            conv_results = self.conversation_collection.query(
                query_embeddings=[query_embedding],
                n_results=top_k
            )
            results.extend(self._format_results(conv_results, 'conversation'))
        
        if search_type in ['questions', 'all']:
            questions_results = self.questions_collection.query(
                query_embeddings=[query_embedding],
                n_results=top_k
            )
            results.extend(self._format_results(questions_results, 'questions'))
        
        # Sort and return top results
        return sorted(results, key=lambda x: x['similarity'], reverse=True)[:top_k]

    def _format_results(self, query_results, section_type):
        """
        Format ChromaDB query results
        
        Args:
            query_results (dict): ChromaDB query results
            section_type (str): Type of section
        
        Returns:
            list: Formatted results
        """
        formatted_results = []
        
        for i in range(len(query_results['ids'][0])):
            formatted_results.append({
                'video_id': query_results['metadatas'][0][i]['video_id'],
                'language_level': query_results['metadatas'][0][i]['language_level'],
                'section_type': section_type,
                'content': query_results['documents'][0][i],
                'similarity': query_results['distances'][0][i]
            })
        
        return formatted_results

    def get_structured_transcript(self, video_id):
        """
        Retrieve structured transcript for a video
        
        Args:
            video_id (str): YouTube video ID
        
        Returns:
            dict: Structured video content
        """
        # Retrieve from each collection
        intro_results = self.introduction_collection.get(
            where={'video_id': video_id}
        )
        
        conv_results = self.conversation_collection.get(
            where={'video_id': video_id}
        )
        
        questions_results = self.questions_collection.get(
            where={'video_id': video_id}
        )
        
        # Check if results exist
        if not (intro_results['documents'] and conv_results['documents'] and questions_results['documents']):
            return None
        
        return {
            'language_level': intro_results['metadatas'][0]['language_level'],
            'introduction': intro_results['documents'][0],
            'conversation': conv_results['documents'][0],
            'questions': json.loads(questions_results['documents'][0])
        }
