CREATE TABLE Conversations (
    conversation_id SERIAL PRIMARY KEY,
    participant_ids TEXT NOT NULL, -- JSON or comma-separated user IDs
    is_group BOOLEAN NOT NULL DEFAULT FALSE, -- TRUE for group chats, FALSE for personal chats
    last_message_id INT,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);