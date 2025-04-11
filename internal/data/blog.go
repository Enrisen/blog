package data

import (
	"context"
	"database/sql"
	"strings"
	"time"
)

type Post struct {
	ID         int64     `json:"id"`
	AuthorID   int64     `json:"author_id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	Excerpt    string    `json:"excerpt"`
	ViewCount  int       `json:"view_count"`
	CreatedAt  time.Time `json:"created_at"`
	Categories []string  `json:"categories"`
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
