CREATE TABLE Messages (
    message_id SERIAL PRIMARY KEY,
    conversation_id INT NOT NULL,
    sender_id INT NOT NULL,
    receiver_id INT, -- Nullable for group chats
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    content TEXT NOT NULL,
    status VARCHAR(20) DEFAULT 'Unread'
);