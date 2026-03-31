package store

import (
	"context"
	"database/sql"
	"fmt"

	"content-public-api/model"
)

type ContentStore struct {
	db *sql.DB
}

func NewContentStore(db *sql.DB) *ContentStore {
	return &ContentStore{db: db}
}

func (s *ContentStore) GetDoneContent(ctx context.Context, cursor *int64, limit int) ([]model.Content, error) {
	query := `
		SELECT id, execution_id, title, short_description, message, status, type,
		       sub_type, category, sub_category, image_url, image_prompt, created, last_updated
		FROM content
		WHERE status = 'DONE'`

	args := []any{}

	if cursor != nil {
		query += " AND id < ?"
		args = append(args, *cursor)
	}

	query += fmt.Sprintf(" ORDER BY id DESC LIMIT %d", limit)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []model.Content
	for rows.Next() {
		var c model.Content
		err := rows.Scan(
			&c.ID, &c.ExecutionID, &c.Title, &c.ShortDescription, &c.Message,
			&c.Status, &c.Type, &c.SubType, &c.Category, &c.SubCategory,
			&c.ImageURL, &c.ImagePrompt, &c.Created, &c.LastUpdated,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, c)
	}

	return results, rows.Err()
}
