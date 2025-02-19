const express = require('express');
const cors = require('cors');
const axios = require('axios');
const path = require('path');
require('dotenv').config();

const app = express();
const PORT = process.env.PORT || 5000;

// Middleware
app.use(cors());
app.use(express.json());

// Serve static files (including index.html)
app.use(express.static(path.join(__dirname)));

// Vocabulary generation endpoint
app.post('/generate', async (req, res) => {
  try {
    const { topic } = req.body;

    if (!topic) {
      return res.status(400).json({ error: 'No topic provided' });
    }

    // Prompt for generating structured vocabulary
    const prompt = `For the topic "${topic}", generate a list of 10 vocabulary words in JSON format. 
    Each word should have the following structure:
    {
      "kanji": "Japanese kanji representation",
      "romaji": "romaji transliteration",
      "english": "English translation",
      "parts": [
        { "kanji": "individual kanji character", "romaji": ["possible romaji readings"] }
      ]
    }
    Provide the response as a JSON array of these word objects.`;

    const response = await axios.post('https://api.mistral.ai/v1/chat/completions', {
      model: 'mistral-small-latest',
      messages: [{ role: 'user', content: prompt }]
    }, {
      headers: {
        'Authorization': `Bearer ${process.env.MISTRAL_API_KEY}`,
        'Content-Type': 'application/json'
      }
    });

    // Extract and parse the JSON response
    const responseText = response.data.choices[0].message.content;
    
    // Try to parse the response, handling potential formatting issues
    let words;
    try {
      // Remove any markdown code block formatting
      const cleanedText = responseText.replace(/```json?/g, '').replace(/```/g, '').trim();
      words = JSON.parse(cleanedText);
    } catch (parseError) {
      console.error('Failed to parse JSON:', parseError);
      console.error('Raw response:', responseText);
      return res.status(500).json({ 
        error: 'Failed to parse vocabulary', 
        rawResponse: responseText 
      });
    }

    res.json({ words });
  } catch (error) {
    console.error('Error generating vocabulary:', error.response ? error.response.data : error.message);
    res.status(500).json({ error: 'Failed to generate vocabulary' });
  }
});

// Default route to serve index.html
app.get('/', (req, res) => {
  res.sendFile(path.join(__dirname, 'index.html'));
});

// Start the server
app.listen(PORT, () => {
  console.log(`Server running on http://localhost:${PORT}`);
});
