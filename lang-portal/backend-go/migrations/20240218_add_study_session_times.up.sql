-- Create new table with desired schema
CREATE TABLE new_study_sessions AS 
SELECT *, 
       created_at as start_time,
       datetime(created_at, '+45 minutes') as end_time 
FROM study_sessions;

-- Drop the old table
DROP TABLE study_sessions;

-- Rename the new table to the original name
ALTER TABLE new_study_sessions RENAME TO study_sessions;