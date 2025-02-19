# Vocab Importer

## Setup

1. Install dependencies:
```bash
npm install
```

2. Set up environment variables:
- Create a `.env` file in the project root
- Add your Mistral API key:
```
MISTRAL_API_KEY=your_mistral_api_key_here
```

## Running the Application

1. Start the Node.js backend:
```bash
# For development with hot-reload
npm run dev

# For production
npm start
```

2. Open `index.html` in a web browser

## Features
- Generate vocabulary words for a given topic using Mistral AI
- Export generated words to JSON
- Import vocabulary from JSON files

## Dependencies
- Express.js
- Cors
- Mistral AI JavaScript SDK
