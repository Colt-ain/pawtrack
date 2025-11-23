CREATE TABLE permissions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE user_permissions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    permission_id INTEGER NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    granted_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, permission_id)
);

CREATE INDEX idx_user_permissions_user_id ON user_permissions(user_id);
CREATE INDEX idx_user_permissions_permission_id ON user_permissions(permission_id);

-- Seed all permissions
INSERT INTO permissions (name, description) VALUES
-- Dog Permissions
('DOGS_CREATE', 'Create new dogs'),
('DOGS_VIEW_OWN', 'View own dogs'),
('DOGS_VIEW_ASSIGNED', 'View assigned dogs'),
('DOGS_VIEW_ALL', 'View all dogs'),
('DOGS_UPDATE_OWN', 'Update own dogs'),
('DOGS_UPDATE_ALL', 'Update any dog'),
('DOGS_DELETE_OWN', 'Delete own dogs'),
('DOGS_DELETE_ALL', 'Delete any dog'),

-- Event Permissions
('EVENTS_CREATE_OWN', 'Create events for own dogs'),
('EVENTS_CREATE_ASSIGNED', 'Create events for assigned dogs'),
('EVENTS_CREATE_ALL', 'Create events for any dog'),
('EVENTS_VIEW_OWN', 'View events of own dogs'),
('EVENTS_VIEW_ASSIGNED', 'View events of assigned dogs'),
('EVENTS_VIEW_ALL', 'View all events'),
('EVENTS_DELETE_OWN', 'Delete events of own dogs'),
('EVENTS_DELETE_ALL', 'Delete any event'),

-- Event Comment Permissions
('EVENT_COMMENTS_CREATE_OWN', 'Create comments on own dog events'),
('EVENT_COMMENTS_CREATE_ASSIGNED', 'Create comments on assigned dog events'),
('EVENT_COMMENTS_VIEW_OWN', 'View comments on own dog events'),
('EVENT_COMMENTS_VIEW_ASSIGNED', 'View comments on assigned dog events'),
('EVENT_COMMENTS_UPDATE_AUTHORED', 'Update own comments'),
('EVENT_COMMENTS_DELETE_AUTHORED', 'Delete own comments'),
('EVENT_COMMENTS_DELETE_ALL', 'Delete any comment'),

-- Consultant Note Permissions
('CONSULTANT_NOTES_CREATE', 'Create notes for assigned dogs'),
('CONSULTANT_NOTES_VIEW_OWN', 'View own notes'),
('CONSULTANT_NOTES_VIEW_ALL', 'View all notes'),
('CONSULTANT_NOTES_UPDATE_OWN', 'Update own notes'),
('CONSULTANT_NOTES_DELETE_OWN', 'Delete own notes'),
('CONSULTANT_NOTES_DELETE_ALL', 'Delete any note'),

-- Consultant Permissions
('CONSULTANTS_SEARCH', 'Search for consultants'),
('CONSULTANTS_INVITE', 'Invite consultants'),
('CONSULTANTS_PROFILE_UPDATE', 'Update own consultant profile'),
('CONSULTANTS_INVITES_ACCEPT', 'Accept invitations'),

-- User Permissions
('USERS_VIEW_OWN', 'View own profile'),
('USERS_VIEW_ALL', 'View all users'),
('USERS_UPDATE_OWN', 'Update own profile'),
('USERS_UPDATE_ALL', 'Update any user'),
('USERS_DELETE_ALL', 'Delete any user');
