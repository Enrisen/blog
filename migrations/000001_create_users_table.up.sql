-- Users Table
-- Stores the single admin user's login information.
CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,            -- Simple auto-incrementing ID for the user
    name VARCHAR(100) NOT NULL,            -- User's full name
    email VARCHAR(255) UNIQUE NOT NULL,    -- User's email address
    password_hash VARCHAR(255) NOT NULL,   -- Store a secure hash of the password
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL  -- When the user was created
);