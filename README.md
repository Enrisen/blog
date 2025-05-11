# TechInsights Blog

A modern blog platform built with Go, focusing on technology-related content. This web application allows users to create, edit, view, and delete blog posts with category tagging and user registration functionality.


## Features

- Blog post management (create, read, update, delete)
- Category tagging for posts
- Rich text editing with Summernote //comming soon
- User registration and authentication
- View count tracking for posts
- Responsive design using Tailwind CSS

## Setup Requirements

- Go (1.16 or later)
- PostgreSQL database
- [golang-migrate](https://github.com/golang-migrate/migrate) for database migrations

## Environment Setup

1. Clone the repository
2. Set up the PostgreSQL database
3. Configure the environment variables:

```bash
# PostgreSQL connection string
export it_blog_DB_DSN="postgres://username:password@localhost:5432/blog_db?sslmode=disable"
```

4. Generate TLS certificates for HTTPS:
   - The application requires TLS certificates for secure connections
   - Place your certificate and key files in the `./tls/` directory as:
     - `./tls/cert.pem` (certificate)
     - `./tls/key.pem` (private key)
   - For development, you can generate self-signed certificates

## Database Setup

The application uses PostgreSQL with the following schema:

- `users` - Stores user account information
- `categories` - Stores blog post categories
- `posts` - Stores blog post content and metadata
- `post_categories` - Links posts to categories (many-to-many relationship)
- `images` - Stores image metadata (for future use)

To set up the database:

1. Create a PostgreSQL database
2. Run the migrations:

```bash
make db/migrations/up
```

To create new migrations:

```bash
make db/migrations/new name=migration_name
```

## Running the Application

To start the server:

```bash
make run
```

This will:
1. Format the code (`go fmt`)
2. Run code verification (`go vet`)
3. Start the web server on HTTPS port 4000: [https://localhost:4000](https://localhost:4000)

To connect to the PostgreSQL database directly:

```bash
make db/psql
```

## Application Architecture

The TechInsights Blog is built using a clean architecture approach:

- **HTTP Server**: Uses Go's standard library HTTP server with TLS support
- **Routing**: Custom routing with middleware for authentication, CSRF protection, and logging
- **Templates**: Server-side rendering with Go's HTML templating
- **Database**: PostgreSQL with prepared statements and connection pooling
- **Security**: HTTPS, CSRF protection, secure cookies, and password hashing

## Pages and Navigation

### Main Pages

- **Home Page**: [https://localhost:4000/](https://localhost:4000/)
  - Displays all blog posts with the most recent at the top
  - Features categories, recent posts, and newsletter signup

- **Blog Post View**: [https://localhost:4000/blog/view/{id}](https://localhost:4000/blog/view/{id})
  - Displays a single blog post with full content
  - Shows related posts and categories
  - Tracks view count (increments on each view)

- **Create Post**: [https://localhost:4000/blog/create](https://localhost:4000/blog/create)
  - Form to create a new blog post with validation
  - Supports adding title, content, and categories
  - Uses Summernote rich text editor for content

- **Edit Post**: [https://localhost:4000/blog/edit/{id}](https://localhost:4000/blog/edit/{id})
  - Form to edit an existing blog post
  - Pre-filled with current post data
  - Allows updating title, content, and categories

### User Management

- **User Registration**: [https://localhost:4000/user/register](https://localhost:4000/user/register)
  - Form to register a new user account
  - Validates email and password requirements
  - Stores securely hashed passwords

- **User Login**: [https://localhost:4000/user/login](https://localhost:4000/user/login)
  - Authentication form
  - Session-based authentication
  - Secure cookie storage

## Rich Text Editor

The blog will use Summernote (v0.9.0) as its rich text editor for creating and editing posts:

- Full formatting capabilities (bold, italic, lists, headings, etc.)
- Image embedding support
- Code block formatting
- Multiple language support

The editor is integrated into the create and edit post forms, allowing for rich content creation.

## Project Structure

- `cmd/web/`: Main application code and server setup
  - `main.go`: Application entry point and configuration
  - `routes.go`: HTTP route definitions
  - `handlers.go`: HTTP request handlers
  - `middleware.go`: HTTP middleware components
  - `templates.go`: Template rendering logic
  - `server.go`: HTTP server configuration
- `internal/data/`: Data models and database interactions
  - `blog.go`: Blog post data model and operations
  - `user.go`: User data model and authentication
  - `errors.go`: Custom error types
- `internal/validator/`: Form validation logic
- `migrations/`: SQL database migration files
- `ui/html/`: HTML templates for the web pages
- `ui/static/`: Static assets (CSS, JavaScript, images)
- `plugins/`: Third-party plugins (Summernote editor)
- `tls/`: TLS certificates for HTTPS

## Database Schema Details

### Users Table
```sql
CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);
```

### Posts Table
```sql
CREATE TABLE posts (
    post_id SERIAL PRIMARY KEY,
    author_id INT NOT NULL,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    excerpt TEXT NULL,
    view_count INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_author FOREIGN KEY(author_id) REFERENCES users(user_id) ON DELETE RESTRICT
);
```

### Categories Table
```sql
CREATE TABLE categories (
    category_id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

### Post Categories Table (Many-to-Many)
```sql
CREATE TABLE post_categories (
    post_id INT NOT NULL,
    category_id INT NOT NULL,
    PRIMARY KEY (post_id, category_id),
    CONSTRAINT fk_post FOREIGN KEY(post_id) REFERENCES posts(post_id) ON DELETE CASCADE,
    CONSTRAINT fk_category FOREIGN KEY(category_id) REFERENCES categories(category_id) ON DELETE CASCADE
);
```

## Security Features

- **HTTPS**: All traffic is encrypted using TLS
- **Password Hashing**: User passwords are hashed using bcrypt
- **CSRF Protection**: Cross-Site Request Forgery protection on all forms
- **Secure Cookies**: Session cookies are secure and HTTPS-only
- **Input Validation**: All user input is validated before processing
- **Prepared Statements**: SQL injection protection
- **Authentication Middleware**: Protected routes require authentication

## Troubleshooting

### Common Issues

1. **Database Connection Errors**:
   - Verify your PostgreSQL service is running
   - Check the `it_blog_DB_DSN` environment variable is correctly set
   - Ensure the database exists and the user has appropriate permissions

2. **TLS Certificate Issues**:
   - Ensure certificate files exist in the `./tls/` directory
   - Check that the certificate and key files are valid
   - For development, you can generate self-signed certificates

3. **Migration Errors**:
   - Make sure golang-migrate is installed
   - Check that the migration files are in the correct format
   - Verify database permissions for creating tables

## Development

To format the Go code:

```bash
make fmt
```

To run code verification:

```bash
make vet
```

To connect to the PostgreSQL database:

```bash
make db/psql
```

## Future Enhancements

- Add image upload capabilities
- Implement search functionality
- Add comment system for blog posts
- Create an admin dashboard
- Add pagination for blog posts
- Implement tagging system in addition to categories

## License

All rights reserved.
