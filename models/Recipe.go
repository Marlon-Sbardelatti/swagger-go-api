package models

// Recipe representa uma receita criada por um usuário.
// @Description Modelo para gerenciamento de receitas.
type Recipe struct {
	// ID é o identificador único da receita.
	ID uint `gorm:"primaryKey" json:"id"`
	// UserID é o identificador do usuário que criou a receita.
	UserID uint `gorm:"not null" json:"user_id"`
	// Name é o nome da sua receita.
	Name string `gorm:"unique;not null" json:"name" swaggertype:"string" example:"bolo de chocolate"`
	// Instructions representa as instruções sobre o modo de preparo da receita.
	Instructions string `gorm:"not null" json:"instructions" swaggertype:"string" example:"Em uma tigela adicione a farinha, o açucar e o cacau em pó." `
	// IngredientsRecipes representa o conjunto de ingredientes que pertence à receita.
    IngredientsRecipes []IngredientsRecipes `gorm:"foreignKey:RecipeID;constraint:OnDelete:CASCADE" json:"ingredients"`
}
