# Frontend Technical Specs

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

- GET `dashboard/last_study_sesssion`
- GET `dashboard/study_progress`
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

- POST `/study-session`






