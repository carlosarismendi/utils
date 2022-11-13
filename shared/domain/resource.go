package domain

import (
	"time"
)

type CreatedAt struct {
	CreatedAt time.Time `json:"createdAt"`
}

type UpdatedAt struct {
	UpdateAt time.Time `json:"updatedAt"`
}

type DeletedAt struct {
	DeletedAt time.Time `json:"deletedAt"`
}
