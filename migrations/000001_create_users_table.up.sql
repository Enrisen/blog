-- Users Table
-- Stores the single admin user's login information.
CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,            -- Simple auto-incrementing ID for the user
    username VARCHAR(50) UNIQUE NOT NULL,  -- Username for login
    password_hash VARCHAR(255) NOT NULL   -- Store a secure hash of the password
);