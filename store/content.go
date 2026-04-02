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

func (s *ContentStore) GetDoneContent(ctx context.Context, page, pageSize int) ([]model.Content, error) {
	offset := (page - 1) * pageSize

	query := fmt.Sprintf(`
		SELECT id, execution_id, title, short_description, message, status, type,
		       sub_type, category, sub_category, image_url, image_prompt, slug, created, last_updated
		FROM content
		WHERE status = 'DONE'
		ORDER BY id DESC
		LIMIT %d OFFSET %d`, pageSize, offset)

	rows, err := s.db.QueryContext(ctx, query)
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
			&c.ImageURL, &c.ImagePrompt, &c.Slug, &c.Created, &c.LastUpdated,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, c)
	}

	return results, rows.Err()
}
