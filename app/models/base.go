package models

import (
	"time"

	"github.com/google/uuid"
)

type Model struct {
	ID        string    `json:"id,omitempty" bson:"id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

func (m *Model) BeforeCreate() {
	m.ID = uuid.New().String()
	m.CreatedAt = time.Now().UTC()
	m.UpdatedAt = time.Now().UTC()
}

func (m *Model) BeforeUpdate() {
	m.UpdatedAt = time.Now().UTC()
}
