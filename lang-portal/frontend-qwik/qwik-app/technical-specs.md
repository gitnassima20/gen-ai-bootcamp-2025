# Frontend Technical Specs

We would like to build a Japanese language learning app.

## Role/Profession:

Front-End Developer

## Project Description

### Project Brief

We are building a japanese language learning web-app which serves the following purposes:

- A portal to launch study activities
- to store, group and explore japanese vocabulary
- to review study progress

The web-app is intended for desktop only, so we don't have to concerned with mobile layouts.

### Technical Requirements

- Qwick as the frontend library
- Tailwind CSS as the CSS framework
- Vite.js as the local development server
- Typescript for the programming language
- ShadCN for components

### Front-End Routes

This is a list of routes for our web-app we are building
Each of these routes are a page and we'll describe them
in more details under the pages heading.

- /dashboard
- /study-activities
- /study-activities/:id
- /words
- /words/:id
- /groups
- /groups/:id
- /sessions
- /settings

The default route / should forward to /dashboard

### Global Components

There will be a horizental navigation bar with the following links:

- Dashboard
- Study Activities
- Words
- Word Groups
- Study Sessions
- Settings

### Breadcrumbs

Beneath the navigation there will be breadcrumbs so users can easily see where they are Examples of breadcrumbs

Dashboard
Study Activities > Adventure MUD
Study Activities > Typing Tutor
Words > こんにちは
Word Groups > Core Verbs

## Pages

### Dashboard /

#### Purpose

The purpose of this page is to show the user's progress and quick stats and it acts as a default page in the web-app.

#### Components

- Last Study Session
  - shows last activity used
  - shows when last activity used
  - summarizes wrong vs correct from last activity
  - has a link to the group
- Study Progress
  - shows total words studied/total stored words
  - displays a mastery progress
- Quick Stats
  - success rate eg. 80%
  - total study sessions eg. 4
  - total active groups eg.3
  - study streak 4 days
- Start Studying Button
  - goes to study activities page

#### APIs to Consume

- GET `dashboard/last-study-sesssion`
- GET `dashboard/study-progress`
- GET `dashboard/quick-stats`

### Study Activities `/study-activities`

#### Purpose

The purpose of this page is show a list of all the study activities with thumbnail and its name, button to launch the study activity or view it.

#### Components

- Study Activity Card
  - thumbnail
  - name
  - button to launch the study activity or view it
  - button to view details about the study activity

#### APIs to Consume

- GET `/study-activities`

### Study Activity Details `/study-activities/:id`

#### Purpose

The purpose of this page is show a details about a single study activity.

#### Components

- Name
- Thumbnail
- Launch button
- Study Sessions associated with pagination
  - id
  - study activity name
  - group name
  - start time
  - end time (inferred by last word_review_item submitted)
  - number of review items

#### APIs to Consume

- GET `/study-activities/:id`

### Study Activity Launch `/study-activities/:id/launch`

#### Purpose

The purpose of this page is to launch a study activity.

#### Components

- name of study activity
- launch form
  - select field for group
  - launch button

#### Behavior

After submission of the form, a new tab will open with the study activity url provided, and it will redirect to the study activity page.

#### APIs to Consume

- POST `/study-sessions`

### Words `/words`

#### Purpose

The purpose of this page is to show a list of all the words in our database.

#### Components

- Paginated Word List
  - Columns
    - Japanese
    - Romaji
    - English
    - Correct Count
    - Wrong Count
  - Pagination with 100 per page
  - Clicking the Japanese word will take us to the word show page

#### APIs to Consume

- GET `/words`

### Word Show `/words/:id`

#### Purpose

The purpose of this page is to information about a specific word.

#### Components

- Japanese
- Romaji
- English
- Study Statistics
  - Correct Count
  - Wrong Count
- Word Groups
  - shown as a series of pills eg. tags
  - when group name is clicked it will take us to the group show page

#### APIs to Consume

- GET `/words/:id`

### Word Groups `/groups`

#### Purpose

The purpose of this page is to show a list of all groups in our database.
#### Components

- Paginated Group List
  - Columns
    - Group Name
    - Word Count
  - Pagination with 100 per page
  - Clicking the group name will take us to the group show page

#### APIs to Consume

- GET `/groups`

### Group Show `/groups/:id`

#### Purpose

The purpose of this page is to show information about a specific word group.
#### Components

- Group Name
- Group Statistics
  - Total Word Count
- Words in Group (Paginated List of Words)
  - Should use the same compoent  as the words index page
- Study Sessions (Paginated List of Study Sessions)
  - Should use the same compoent  as the study session index page

#### APIs to Consume

- GET `/groups/:id` (the name and the group's stat)
- GET `/groups/:id/words` 
- GET `/groups/:id/study-sessions`

### Study Sessions `/groups`

#### Purpose

The purpose of this page is to show a list of all study sessions in our database.
#### Components

- Paginated Study Session List
  - Columns
    - Id 
    - Activity Name
    - Group Name
    - Start Time
    - End Time
  - Clicking the study session will take us to the study session show page

#### APIs to Consume

- GET `/study-sessions`

### Study Session Details `/study-sessions/:id`

#### Purpose

The purpose of this page is to show information about a specific study session.
#### Components

- Study Session Details
  - Study Activity Name
  - Group Name
  - Start Time
  - Number of Review Items
- Words Reviewd Items (Paginated List of Words)
  - should use the same component as the words page

#### APIs to Consume

- GET `/study-sessions/:id`
- GET `/study-sessions/:id/words`

### Settings `/settings`

#### Purpose

The purpose of this page is to configure of the study portal

#### Components

- Theme Selection
- Reset History Button
  - this will delete all study sessions and words review items.
- Full Reset
  - this will drop all tables and re-create with seed data

#### APIs to Consume

- POST `/reset-history`
- POST `/full-reset`