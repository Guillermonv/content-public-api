package store

import (
	"context"
	"database/sql"

	"content-public-api/model"
)

type ContentStore struct {
	db *sql.DB
}

func NewContentStore(db *sql.DB) *ContentStore {
	return &ContentStore{db: db}
}

func (s *ContentStore) SearchContent(ctx context.Context, q string, page, pageSize int) ([]model.Content, error) {
	offset := (page - 1) * pageSize
	like := "%" + q + "%"

	rows, err := s.db.QueryContext(ctx, `
		SELECT id, execution_id, title, short_description, message, status,
		       category, sub_category, image_url, image_prompt, slug, created, last_updated
		FROM content
		WHERE status = 'DONE'
		  AND (title LIKE ? OR short_description LIKE ? OR slug LIKE ? OR message LIKE ?)
		ORDER BY
		  CASE
		    WHEN title             LIKE ? THEN 1
		    WHEN short_description LIKE ? THEN 2
		    WHEN slug              LIKE ? THEN 3
		    WHEN message           LIKE ? THEN 4
		    ELSE 5
		  END,
		  id DESC
		LIMIT ? OFFSET ?`,
		like, like, like, like, like, like, like, like, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []model.Content
	for rows.Next() {
		var c model.Content
		err := rows.Scan(
			&c.ID, &c.ExecutionID, &c.Title, &c.ShortDescription, &c.Message,
			&c.Status, &c.Category, &c.SubCategory,
			&c.ImageURL, &c.ImagePrompt, &c.Slug, &c.Created, &c.LastUpdated,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, c)
	}
	return results, rows.Err()
}

func (s *ContentStore) GetDoneContent(ctx context.Context, page, pageSize int) ([]model.Content, error) {
	offset := (page - 1) * pageSize

	rows, err := s.db.QueryContext(ctx, `
		SELECT id, execution_id, title, short_description, message, status,
		       category, sub_category, image_url, image_prompt, slug, created, last_updated
		FROM content
		WHERE status = 'DONE'
		ORDER BY id DESC
		LIMIT ? OFFSET ?`, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []model.Content
	for rows.Next() {
		var c model.Content
		err := rows.Scan(
			&c.ID, &c.ExecutionID, &c.Title, &c.ShortDescription, &c.Message,
			&c.Status, &c.Category, &c.SubCategory,
			&c.ImageURL, &c.ImagePrompt, &c.Slug, &c.Created, &c.LastUpdated,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, c)
	}
	return results, rows.Err()
}

func (s *ContentStore) GetContentBySlug(ctx context.Context, slug string) (*model.Content, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT id, execution_id, title, short_description, message, status,
		       category, sub_category, image_url, image_prompt, slug, created, last_updated
		FROM content
		WHERE status = 'DONE' AND slug = ?
		LIMIT 1`, slug)

	var c model.Content
	err := row.Scan(
		&c.ID, &c.ExecutionID, &c.Title, &c.ShortDescription, &c.Message,
		&c.Status, &c.Category, &c.SubCategory,
		&c.ImageURL, &c.ImagePrompt, &c.Slug, &c.Created, &c.LastUpdated,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}
