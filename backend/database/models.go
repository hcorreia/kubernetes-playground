// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package database

import (
	"time"

	null "github.com/guregu/null/v6"
)

type Post struct {
	ID        int32       `json:"id"`
	Title     string      `json:"title"`
	Image     null.String `json:"image"`
	Content   string      `json:"content"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}
