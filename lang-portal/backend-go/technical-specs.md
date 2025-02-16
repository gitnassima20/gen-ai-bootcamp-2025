# Backend Server Technical Specs

## Business Goal

A language learning school wants to build a prototype of learning portal which will act as three things:

- Inventory of possible vocabulary that can be learned
- Act as a  Learning record store (LRS), providing correct and wrong score on practice vocabulary
- A unified launchpad to launch different learning apps

## Technical Requirments

- Go (1.21 or later)
- SQLite3
- Gin framework
- API REST
- No Auth/Authz

## Database Schema

We have the following tables:

- words - stored vocabulary words
  - id int
  - japanese string
  - romaji string
  - english string
  - parts json

- word_groups - join table for words and groups: many-to-many
  - id int
  - word_id int
  - group_id int

- groups - thematic groups of words
  - id int
  - name string
  - words_count int

- study_sessions - records of study sessions grouping word_review_items
  - id int
  - name string
  - url string
- study_activities - a specific study activity, linking a study session to a group
  - id int
  - group_id int
  - study_activity_id int
  - created_at timestamp

- word_review_items - a record of word practice, determining if the word was correct or not
  - id int
  - word_id int
  - study_session_id int
  - correct boolean
  - created_at timestamp

## API Endpoints

- GET `dashboard/last-study-sesssion`
- GET `dashboard/study-progress`
- GET `dashboard/quick-stats`

- GET `/study-activities`
- GET `/study-activities/:id`

- POST `/study-sessions`
  - params : group_id, study_activity_id

- GET `/words`
  - pagination with 100 items per page
- GET `/words/:id`

- GET `/groups`
  - pagination with 100 items per page
- GET `/groups/:id`
- GET `/groups/:id/words` 
- GET `/groups/:id/study-sessions`

- GET `/study-sessions`
  - pagination with 100 items per page
- GET `/study-sessions/:id`
- GET `/study-sessions/:id/words`

- POST `/reset-history`
- POST `/full-reset`

- POST `/study-sessions/:id/words/:word-id/review`
  - params : correct