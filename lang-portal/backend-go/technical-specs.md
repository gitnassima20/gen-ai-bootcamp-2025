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

### Dashboard

- [x] GET `api/v1/dashboard/last-study-sesssion`
  - **Response Body**:

  ```json
  {
    "id": 1,
    "name": "Japanese Verbs Practice",
    "created_at": "2025-02-17T10:30:00Z",
    "group_id": 2,
    "group_name": "Verbs Group"
  }
  ```

- [x] GET `api/v1/dashboard/study-progress`
  - **Response Body**:

  ```json
  {
    "total_words_studied": 150,
    "total_available_words": 400
  }
  ```

- [x] GET `api/v1/dashboard/quick-stats`
  - **Response Body**:

  ```json
  {
    "success-rate": 25,
    "total_study_sessions": 4,
    "total_active_groups": 2,
    "current_streak": 7,
  }
  ```

### Study Activities

- [x] GET `api/v1/study-activities`
  - **Response Body**:

  ```json
  {
    "items": [
      {
        "id": 1,
        "name": "Japanese Verbs Practice",
        "thumbnail_url": "https://example.com/verbs-thumbnail.jpg",
        "total_sessions": 15,
        "last_used": "2025-02-16T14:30:00Z"
      },
      {
        "id": 2,
        "name": "Adjectives Mastery",
        "thumbnail_url": "https://example.com/adjectives-thumbnail.jpg",
        "total_sessions": 10,
        "last_used": "2025-02-15T09:45:00Z"
      }
    ],
    "total_activities": 2
  }
  ```

- [x] GET `api/v1/study-activities/:id`
  - **Response Body**:

  ```json
  {
    "id": 1,
    "name": "Japanese Verbs Practice",
    "thumbnail_url": "https://example.com/thumbnail.jpg"
  }
  ```

- [x] POST `api/v1/study-sessions`
  - **Request Params**:
    - group_id int
    - study_activity_id int

  - **Response Body**:

  ```json
  {
    "id": 1,
    "group_id": 1
  }
  ```

### Words

- GET `api/v1/words`
  - **Query Parameters**:
    - `page` (optional, default: 1)
    - `words_per_page` (optional, default: 100)
  - **Response Body**:

  ```json
  {
    "items": [
      {
        "id": 1,
        "kanji": "食べる",
        "romaji": "taberu",
        "english": "to eat",
        "correct_count": 15,
        "wrong_count": 5
      },
      {
        "id": 2,
        "kanji": "読む",
        "romaji": "yomu",
        "english": "to read", 
        "correct_count": 10,
        "wrong_count": 3
      }
    ],
    "total_count": 250,
    "current_page": 1,
    "total_pages": 3
  }
  ```

- GET `api/v1/words/:id`
  - **Response Body**:

  ```json
  {
    "id": 1,
    "kanji": "食べる",
    "romaji": "taberu",
    "english": "to eat",
    "correct_count": 15,
    "wrong_count": 5,
    "groups": [
      {
        "id": 1,
        "name": "Verbs Group"
      },
      {
        "id": 2,
        "name": "Beginner Vocabulary"
      }
    ]
  }
  ```

### Groups of Words

- [x] GET `api/v1/groups`
  - **Query Parameters**:
    - `page` (optional, default: 1)
    - `groups_per_page` (optional, default: 100)
  - **Response Body**:

  ```json
  {
    "items": [
      {
        "id": 1,
        "name": "Verbs Group",
        "word_count": 50
      },
      {
        "id": 2,
        "name": "Adjectives Group",
        "word_count": 40
      }
    ],
    "total_count": 10,
    "current_page": 1,
    "total_pages": 1
  }
  ```

- [x] GET `api/v1/groups/:id`
  - **Response Body**:

  ```json
  {
    "id": 1,
    "name": "Verbs Group",
    "total_word_count": 50
  }
  ```

- [x] GET `api/v1/groups/:id/words`
  - **Query Parameters**:
    - `page` (optional, default: 1)
    - `words_per_page` (optional, default: 100)
  - **Response Body**:

  ```json
  {
    "items": [
      {
        "id": 1,
        "kanji": "食べる",
        "romaji": "taberu",
        "english": "to eat",
        "correct_count": 15,
        "wrong_count": 5
      },
      {
        "id": 2,
        "kanji": "読む",
        "romaji": "yomu",
        "english": "to read",
        "correct_count": 10,
        "wrong_count": 3
      }
    ],
    "total_words": 50,
    "current_page": 1,
    "total_pages": 1
  }
  ```

- [x] GET `api/v1/groups/:id/study-sessions`
  - **Query Parameters**:
    - `page` (optional, default: 1)
    - `sessions_per_page` (optional, default: 100)
  - **Response Body**:

  ```json
  {
    "items": [
      {
        "id": 1,
        "name": "Verbs Practice Session",
        "start_time": "2025-02-16T14:30:00Z",
        "end_time": "2025-02-16T15:15:00Z",
        "total_words_reviewed": 20,
      },
      {
        "id": 2,
        "name": "Advanced Verbs Study",
        "start_time": "2025-02-15T10:00:00Z",
        "end_time": "2025-02-15T11:00:00Z",
        "total_words_reviewed": 25,
      }
    ],
    "total_count": 15,
    "current_page": 1,
    "total_pages": 1
  }
  ```

### Study Sessions

- [x] GET `api/v1/study-sessions`
  - **Query Parameters**:
    - `page` (optional, default: 1)
    - `sessions_per_page` (optional, default: 100)
  - **Response Body**:

  ```json
  {
    "items": [
      {
        "id": 1,
        "activity_name": "Verbs Practice",
        "group_name": "Verbs Group",
        "start_time": "2025-02-16T14:30:00Z",
        "end_time": "2025-02-16T15:15:00Z",
        "total_words_reviewed": 20
      },
      {
        "id": 2,
        "activity_name": "Adjectives Study",
        "group_name": "Adjectives Group", 
        "start_time": "2025-02-15T10:00:00Z",
        "end_time": "2025-02-15T11:00:00Z",
        "total_words_reviewed": 25
      }
    ],
    "total_sessions": 50,
    "current_page": 1,
    "total_pages": 1
  }
  ```

- [x] GET `api/v1/study-sessions/:id`
  - **Response Body**:

  ```json
  {
    "id": 1,
    "activity_name": "Verbs Practice",
    "group_name": "Verbs Group",
    "start_time": "2025-02-16T14:30:00Z",
    "end_time": "2025-02-16T15:15:00Z",
    "total_words_reviewed": 20
  }
  ```

- [x] GET `api/v1/study-sessions/:id/words`
  - **Query Parameters**:
    - `page` (optional, default: 1)
    - `words_per_page` (optional, default: 100)
  - **Response Body**:

  ```json
  {
    "items": [
      {
        "id": 1,
        "kanji": "食べる",
        "romaji": "taberu",
        "english": "to eat",
        "correct_count": 5,
        "wrong_count": 2
      },
      {
        "id": 2,
        "kanji": "読む",
        "romaji": "yomu", 
        "english": "to read",
        "correct_count": 10,
        "wrong_count": 3,
      }
    ],
    "total_words": 20,
    "current_page": 1,
    "total_pages": 1
  }
  ```

### Settings

- POST `/reset-history`
  - **Response Body**:

  ```json
  {
    "success": true,
    "message": "Study history has been fully reset"
  }
  ```

- POST `/full-reset`
  - **Response Body**:

  ```json
  {
    "success": true,
    "message": "System has been fully reset"
  }
  ```

- [x] POST `/study-sessions/:id/words/:word-id/review`
  - **Request Params/Body**:
    - id (study-session-id) int
    - word-id int

    ```json
    {
        "correct": true
    }
    ```

  - **Response Body**:

  ```json
  {
    "success": true,
    "word_id": 1,
    "study_session_id": 1,
    "correct": true,
    "created_at": "2025-02-16T14:30:00Z"
  }
  ```

## Mage Tasks

Mage is a task runner for Go.
Lets list out possible tasks we need for our lang portal.

- **Task 1**: Initialize the sqlite3 database called `langportal.db`
- **Task 2**: Run migrations found in `migrations/`
- **Task 3**: Run seed data found in `seed/`
- **Task 4**: Start the server: `go run cmd/server/main.go`

## Environment Setup for CGO and SQLite

### Prerequisites

- MinGW-64 (GCC for Windows)
- Go 1.21 or later
- MSYS2 with MinGW64 toolchain

### Environment Configuration

Before running Go commands that require CGO (such as SQLite operations), set the following environment variables in powershell:

```bash
# Enable CGO
export CGO_ENABLED=1

# Set C and C++ compilers
export CC=x86_64-w64-mingw32-gcc
export CXX=x86_64-w64-mingw32-g++

# Add MinGW64 bin directory to PATH
export PATH=$PATH:/c/msys64/mingw64/bin
```

### Verification

Verify GCC installation:

```bash
gcc --version
```

Expected output should include:

```
gcc.exe (Rev2, Built by MSYS2 project) 14.2.0
Copyright (C) 2024 Free Software Foundation, Inc.
```