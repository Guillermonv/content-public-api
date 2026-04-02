package model

import "time"

type Content struct {
	ID               int64      `json:"id"`
	ExecutionID      int64      `json:"execution_id"`
	Title            *string    `json:"title"`
	ShortDescription *string    `json:"short_description"`
	Message          *string    `json:"message"`
	Status           *string    `json:"status"`
	Type             *string    `json:"type"`
	SubType          *string    `json:"sub_type"`
	Category         *string    `json:"category"`
	SubCategory      *string    `json:"sub_category"`
	ImageURL         *string    `json:"image_url"`
	ImagePrompt      *string    `json:"image_prompt"`
	Slug             *string    `json:"slug"`
	Created          *time.Time `json:"created"`
	LastUpdated      *time.Time `json:"last_updated"`
}

type ContentPage struct {
	Data     []Content `json:"data"`
	Page     int       `json:"page"`
	PageSize int       `json:"page_size"`
	HasNext  bool      `json:"has_next"`
}
