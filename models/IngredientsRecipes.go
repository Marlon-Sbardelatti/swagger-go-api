package models

// IngredientsRecipes representa a associação de um ingrediente a uma receita.
// @Description Modelo para relacionar um ingrediente da tabela ingredients a uma receita.
type IngredientsRecipes struct {
	// RecipeID é o ID da receita à qual o ingrediente foi adicionado.
	RecipeID uint `gorm:"primaryKey" swaggertype:"integer"`
	// IngredientID é o ID do ingrediente adicionado.
	IngredientID uint `gorm:"primaryKey" json:"ingredient_id" swaggertype:"integer"`
	// Quantity é a quantidade do ingrediente adicionado.
	Quantity string `gorm:"not null" json:"quantity" swaggertype:"string" example:"200g"`
	// Ingredient é o objeto do ingrediente adicionado.
	Ingredient Ingredient `gorm:"foreignKey:IngredientID;constraint:OnDelete:CASCADE" json:"ingredient"`
}
