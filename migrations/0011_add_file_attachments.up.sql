ALTER TABLE events 
ADD COLUMN attachment_url VARCHAR(500);

ALTER TABLE event_comments 
ADD COLUMN attachment_url VARCHAR(500);
