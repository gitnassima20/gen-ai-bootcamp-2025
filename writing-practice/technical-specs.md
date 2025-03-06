# Writing Practice

## Initialization Step

The app will fetch ``http://localhost:8080/api/v1/groups/:id/words`` from @backend_go and store the response in memory.

## Page Stages:

How the app should behave from the user's perspective.

### Setup Stage

When a user first's start the app,
They will see a single button called 'Generate Sentence',
When they press the button the app will generate a simple sentence using the Sentence Generator LLM prompt, and the state will be updated to Practice Stage.

### Practice Stage

When on practice stage,
user will be able to see the generated sentence in english,
Also they will see an upload button under the generated sentence,
They will see submit after uploading their file for review
When they press the submit button, the uploaded image will be passed to the Grading System and transite to Review Stage.

### Review Stage

When a user is on the review stage,
They will still see the english sentence,
The upload field will be gone,
They will see a review of the output from the Grading System:

- Transcription of the image
- Translation of Transcription
- Grading:
  - a letter score using S rank
  - a description (areas for improvement, accuracy level..)

There will be a button called 'Next Sentence' which will take the user to the Setup Stage.

## Sentence Generator LLM Prompt

Generate a simple sentence using the word {{word}}

- Use JPLT2 Grammar

## Grading System

The Grading System will do the following:

- Transcribe the image using MangaOCR.
- Use an LLM to produce literal translation of the transcription.
- Use another LLM to grade.
- Return the result to the frontend app.