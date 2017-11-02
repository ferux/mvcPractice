package model

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Modeler interface {
	Create(*sqlx.DB, Modeler) (uuid.UUID, error)
	Read(*sqlx.DB) ([]Modeler, error)
	ReadUUID(*sqlx.DB, uuid.UUID) (Modeler, error)
	Update(*sqlx.DB, Modeler) error
	Delete(*sqlx.DB, uuid.UUID) error
}
