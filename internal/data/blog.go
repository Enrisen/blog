package data

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/Enrisen/blog/internal/validator"
)

type Post struct {
	ID         int64     `json:"id"`
	AuthorID   int64     `json:"author_id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	Excerpt    string    `json:"excerpt"`
	ViewCount  int       `json:"view_count"`
	CreatedAt  time.Time `json:"created_at"`
	Categories []string  `json:"categories"` // Keep this for individual post display if needed
}

// Category represents a single category in the database.
type Category struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type BlogModel struct {
	DB *sql.DB
}

func (m *BlogModel) GetAll() ([]*Post, error) {
	query := `
		SELECT p.post_id, p.author_id, p.title, p.content, p.excerpt, p.view_count, p.created_at
		FROM posts p
		ORDER BY p.created_at DESC`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []*Post{}

	for rows.Next() {
		p := &Post{}
		err := rows.Scan(
			&p.ID,
			&p.AuthorID,
			&p.Title,
			&p.Content,
			&p.Excerpt,
			&p.ViewCount,
			&p.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Get categories for this post
		p.Categories, err = m.getPostCategories(p.ID)
		if err != nil {
			return nil, err
		}

		posts = append(posts, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

// CreatePost inserts a new blog post into the database
func (m *BlogModel) CreatePost(authorID int64, title, content string, categories []string) (int64, error) {
	// Create an excerpt from the content (first 150 characters)
	excerpt := content
	if len(content) > 150 {
		excerpt = strings.TrimSpace(content[:150]) + "..."
	}

	// Start a transaction
	tx, err := m.DB.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	// Insert the post
	query := `
		INSERT INTO posts (author_id, title, content, excerpt, view_count, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING post_id`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var postID int64
	err = tx.QueryRowContext(
		ctx,
		query,
		authorID,
		title,
		content,
		excerpt,
		0, // Initial view count is 0
		time.Now(),
	).Scan(&postID)

	if err != nil {
		return 0, err
	}

	// If categories are provided, add them
	if len(categories) > 0 {
		for _, category := range categories {
			// First check if the category exists
			var categoryID int64
			err = tx.QueryRowContext(
				ctx,
				`SELECT category_id FROM categories WHERE name = $1`,
				category,
			).Scan(&categoryID)

			if err == sql.ErrNoRows {
				// Category doesn't exist, create it
				err = tx.QueryRowContext(
					ctx,
					`INSERT INTO categories (name) VALUES ($1) RETURNING category_id`,
					category,
				).Scan(&categoryID)
				if err != nil {
					return 0, err
				}
			} else if err != nil {
				return 0, err
			}

			// Link the post to the category
			_, err = tx.ExecContext(
				ctx,
				`INSERT INTO post_categories (post_id, category_id) VALUES ($1, $2)`,
				postID,
				categoryID,
			)
			if err != nil {
				return 0, err
			}
		}
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return postID, nil
}

func (m *BlogModel) getPostCategories(postID int64) ([]string, error) {
	query := `
		SELECT c.name
		FROM categories c
		JOIN post_categories pc ON c.category_id = pc.category_id
		WHERE pc.post_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := []string{}

	for rows.Next() {
		var category string
		err := rows.Scan(&category)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

// Get retrieves a single blog post by ID
func (m *BlogModel) Get(id int64) (*Post, error) {
	query := `
		SELECT p.post_id, p.author_id, p.title, p.content, p.excerpt, p.view_count, p.created_at
		FROM posts p
		WHERE p.post_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var post Post

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&post.ID,
		&post.AuthorID,
		&post.Title,
		&post.Content,
		&post.Excerpt,
		&post.ViewCount,
		&post.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	// Get categories for this post
	post.Categories, err = m.getPostCategories(post.ID)
	if err != nil {
		return nil, err
	}

	// Increment view count
	_, err = m.DB.ExecContext(ctx,
		`UPDATE posts SET view_count = view_count + 1 WHERE post_id = $1`,
		id)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

// UpdatePost updates an existing blog post in the database
func (m *BlogModel) UpdatePost(id int64, title, content string, categories []string) error {
	// Create an excerpt from the content (first 150 characters)
	excerpt := content
	if len(content) > 150 {
		excerpt = strings.TrimSpace(content[:150]) + "..."
	}

	// Start a transaction
	tx, err := m.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update the post
	query := `
		UPDATE posts
		SET title = $1, content = $2, excerpt = $3
		WHERE post_id = $4`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err = tx.ExecContext(
		ctx,
		query,
		title,
		content,
		excerpt,
		id,
	)

	if err != nil {
		return err
	}

	// Delete existing category associations
	_, err = tx.ExecContext(
		ctx,
		`DELETE FROM post_categories WHERE post_id = $1`,
		id,
	)
	if err != nil {
		return err
	}

	// If categories are provided, add them
	if len(categories) > 0 {
		for _, category := range categories {
			// First check if the category exists
			var categoryID int64
			err = tx.QueryRowContext(
				ctx,
				`SELECT category_id FROM categories WHERE name = $1`,
				category,
			).Scan(&categoryID)

			if err == sql.ErrNoRows {
				// Category doesn't exist, create it
				err = tx.QueryRowContext(
					ctx,
					`INSERT INTO categories (name) VALUES ($1) RETURNING category_id`,
					category,
				).Scan(&categoryID)
				if err != nil {
					return err
				}
			} else if err != nil {
				return err
			}

			// Link the post to the category
			_, err = tx.ExecContext(
				ctx,
				`INSERT INTO post_categories (post_id, category_id) VALUES ($1, $2)`,
				id,
				categoryID,
			)
			if err != nil {
				return err
			}
		}
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

// DeletePost deletes a blog post and its category associations
func (m *BlogModel) DeletePost(id int64) error {
	// Start a transaction
	tx, err := m.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Delete category associations first (due to foreign key constraints)
	_, err = tx.ExecContext(
		ctx,
		`DELETE FROM post_categories WHERE post_id = $1`,
		id,
	)
	if err != nil {
		return err
	}

	// Delete the post
	_, err = tx.ExecContext(
		ctx,
		`DELETE FROM posts WHERE post_id = $1`,
		id,
	)
	if err != nil {
		return err
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

// GetAllCategories retrieves all categories from the database.
func (m *BlogModel) GetAllCategories() ([]*Category, error) {
	query := `
		SELECT category_id, name
		FROM categories
		ORDER BY name ASC`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := []*Category{}

	for rows.Next() {
		var category Category
		err := rows.Scan(
			&category.ID,
			&category.Name,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, &category)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

// ValidateBlogPost validates the blog post form data
func ValidateBlogPost(v *validator.Validator, title, content string, categories []string) {
	// Validate title
	v.Check(validator.NotBlank(title), "title", "Title cannot be empty")
	v.Check(validator.MaxLength(title, 100), "title", "Title cannot be more than 100 characters")

	// Validate content
	v.Check(validator.NotBlank(content), "content", "Content cannot be empty")

	// Validate categories - require at least one category
	if len(categories) == 0 {
		v.AddError("category", "Please select or add a category")
	} else {
		for i, category := range categories {
			if !validator.NotBlank(category) {
				v.AddError(
					"category_"+string(rune('a'+i)),
					"Category cannot be empty",
				)
			}
			if !validator.MaxLength(category, 50) {
				v.AddError(
					"category_"+string(rune('a'+i)),
					"Category cannot be more than 50 characters",
				)
			}
		}
	}
}
