package models

// Ingredient representa um ingrediente utilizado em receitas.
// @Description Modelo para gerenciamento de ingredientes.
type Ingredient struct {
	// ID é o identificador único do ingrediente.
	ID uint `gorm:"primaryKey" json:"id"`
	// Name é o nome do ingrediente.
    Name string `gorm:"unique;not null" json:"name" example:"Farinha de trigo."`
}

