-- Categories Table
-- Stores available post categories.
CREATE TABLE categories (
    category_id SERIAL PRIMARY KEY,        -- Simple auto-incrementing ID for the category
    name VARCHAR(100) UNIQUE NOT NULL,     -- Name of the category (e.g., 'Technology', 'Travel')
    -- slug VARCHAR(120) UNIQUE NOT NULL,  -- Optional: URL-friendly version of the name
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);