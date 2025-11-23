ALTER TABLE events 
DROP COLUMN IF EXISTS attachment_url;

ALTER TABLE event_comments 
DROP COLUMN IF EXISTS attachment_url;
