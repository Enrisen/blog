# TechInsights Blog

A modern blog platform built with Go, focusing on technology-related content. This web application allows users to create, edit, view, and delete blog posts with category tagging and user registration functionality.

User authentication and validation is not functional so CRUD can be performed by all users.

## Features
- Blog post management (create, read, update, delete)
- User registration (Not completely functional)
- View count tracking for posts (Not functional, using only dummy values)


### What you need

- Go (1.16 or later)
- PostgreSQL database

### Running the Application

To start the server, run:

```bash
make run
```

This will start the web server on [http://localhost:4000](http://localhost:4000).

### Database Setup

The application uses PostgreSQL. Make sure you have the database set up and the connection string available in the environment variable `it_blog_DB_DSN`.

To run database migrations:

```bash
make db/migrations/up
```

## Pages and Navigation

### Main Pages

- **Home Page**: [http://localhost:4000/](http://localhost:4000/)
  - Displays all blog posts with the most recent at the top
  - Features categories, recent posts, and newsletter signup

- **Blog Post View**: [http://localhost:4000/blog/view/{id}](http://localhost:4000/blog/view/{id})
  - Displays a single blog post with full content
  - Shows related posts and categories
  - Tracks view count

- **Create Post**: [http://localhost:4000/blog/create](http://localhost:4000/blog/create)
  - Form to create a new blog post with validation
  - Supports adding title, content, and categories

- **Edit Post**: [http://localhost:4000/blog/edit/{id}](http://localhost:4000/blog/edit/{id})
  - Form to edit an existing blog post
  - Pre-filled with current post data

### User Management

- **User Registration**: [http://localhost:4000/user/register](http://localhost:4000/user/register)
  - Form to register a new user account
  - Validates email and password requirements

## Project Structure

- `cmd/web/`: Main application code and server setup
- `internal/data/`: Data models and database interactions
- `internal/validator/`: Form validation logic
- `migrations/`: SQL database migration files
- `ui/html/`: HTML templates for the web pages
- `ui/static/`: Static assets (CSS, JavaScript, images)

## Development

To format the Go code:

```bash
make fmt
```

To run code verification:

```bash
make vet
```

## License

All rights reserved.