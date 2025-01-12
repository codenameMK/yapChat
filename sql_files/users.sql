CREATE TABLE Users (
    user_id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    profile_pic TEXT, -- Reference to the profile picture (URL or path)
    last_seen TIMESTAMP
);