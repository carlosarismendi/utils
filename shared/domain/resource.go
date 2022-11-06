package domain

import (
	"fmt"
	"time"

	"github.com/ansel1/merry"
	"github.com/google/uuid"
)

type Updateable struct {
	UpdateAt time.Time `json:"updatedAt"`
}

type Deleteable struct {
	DeletedAt time.Time `json:"deletedAt"`
}

type Resource struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
}

func NewResource(id string, createdAt time.Time) (*Resource, error) {
	if id == "" {
		id = uuid.New().String()
	} else {
		_, err := uuid.Parse(id)
		if err != nil {
			return nil, merry.New(fmt.Sprintf("Invalid field ID: it must be a valid uuid. The value received is '%s'.", id))
		}
	}

	return &Resource{
		ID:        id,
		CreatedAt: createdAt,
	}, nil
}
