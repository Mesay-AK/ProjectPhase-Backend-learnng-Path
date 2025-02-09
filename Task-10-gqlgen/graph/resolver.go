package graph

import (
	"gorm.io/gorm"
	"taskManager/internal/database"
)

type Resolver struct {
	DB *gorm.DB
}

func NewResolver() *Resolver {
	return &Resolver{
		DB: database.DB, 
	}
}
