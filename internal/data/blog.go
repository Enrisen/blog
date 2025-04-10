package data

import (
	"context"
	"database/sql"
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
