package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"main.go/app"
	"main.go/models"
)

// @Summary      Buscar todas as receitas
// @Description  Buscar todas as receitas cadastradas
// @Tags         recipe
// @Produce      json
// @Success      200  {array}   models.Recipe
// @Failure      404  "Not Found"
// @Failure      500  "Internal Server Error"
// @Router       /recipe/ [get]
func GetAllRecipesHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var recipes []models.Recipe

		// Retorna as receitas e ingredientes associados a elas da tabela ingredients_recipes
		result := app.DB.Preload("IngredientsRecipes.Ingredient").Find(&recipes)

		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				http.Error(w, "Not Found", http.StatusNotFound)
				return
			} else {
				fmt.Printf("Error querying recipes: %v\n", result.Error)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}

		// Transforma structs das receitas para JSON
		recipesJson, err := json.Marshal(recipes)
		if err != nil {
			http.Error(w, "Error encoding recipes to JSON", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(recipesJson)
	}
}

// @Summary      Buscar receita pelo ID
// @Description  Buscar receita pelo ID
// @Tags         recipe
// @Produce      json
// @Param		 id path int true "ID da receita"
// @Success      200  {array}   models.Recipe
// @Failure      404  "Not Found"
// @Failure      500  "Internal Server Error"
// @Router       /recipe/{id} [get]
func GetRecipeByIdHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		var recipe models.Recipe

		result := app.DB.Preload("IngredientsRecipes.Ingredient").Where("id = ?", id).First(&recipe)

		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				http.Error(w, "Recipe not found", http.StatusNotFound)
				return
			} else {
				fmt.Printf("Error querying recipe: %v\n", result.Error)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}

		// Transforma struct da receita para JSON
		recipeJson, err := json.Marshal(recipe)
		if err != nil {
			http.Error(w, "Error encoding recipe to JSON", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(recipeJson)
	}
}

// @Summary      Buscar receita pelo nome
// @Description  Buscar receita pelo nome sem case sensitive e convertendo '-' para espaços
// @Tags         recipe
// @Produce      json
// @Param		 name path string true "Nome da receita"
// @Success      200  {array}   models.Recipe
// @Failure      404  "Not Found"
// @Failure      500  "Internal Server Error"
// @Router       /recipe/name/{name} [get]
func GetRecipeByNameHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")
		name = strings.ToLower(strings.ReplaceAll(name, "-", " ")) // Substitui hífens por espaço antes de fazer o select

		var recipe models.Recipe

		// Query que seleciona pelo atributo name, comparando ambas Strings em minúsculo
		result := app.DB.Preload("IngredientsRecipes.Ingredient").Where("name LIKE LOWER(?)", name).First(&recipe)

		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				http.Error(w, "Recipe not found", http.StatusNotFound)
				return
			} else {
				fmt.Printf("Error querying recipe: %v\n", result.Error)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}

		// Transforma struct da receita para JSON
		recipeJson, err := json.Marshal(recipe)
		if err != nil {
			http.Error(w, "Error encoding recipe to JSON", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(recipeJson)
	}
}

// @Summary      Criar nova receita
// @Description  Criar nova receita
// @Tags         recipe
// @Accept       json
// @Security Token 
// @Produce      text/plain
// @Param		 recipe body models.Recipe true "Nova receita"
// @Success      201  {string} string "Recipe created!"
// @Failure      400  "Invalid JSON"
// @Router       /recipe [post]
func CreateRecipeHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var recipe models.Recipe

		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()

		err := decoder.Decode(&recipe)
		if err != nil || recipe.UserID == 0 || recipe.Name == "" || recipe.Instructions == "" {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		result := app.DB.Create(&recipe)
		if result.Error != nil {
			http.Error(w, "Recipe already exists or data is incorrect", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Recipe created!"))
	}
}

// @Summary      Atualizar receita
// @Description  Atualizar receita pelo ID
// @Tags         recipe
// @Accept       json
// @Security Token 
// @Produce      text/plain
// @Param		 id path int true "ID da receita"
// @Param		 recipe body models.Recipe true "Receita atualizada"
// @Success      200  {string}   string "Recipe updated!"
// @Failure      400  "Invalid JSON"
// @Failure      404  "Not Found"
// @Failure      500  "Internal Server Error"
// @Router       /recipe/{id} [put]
func UpdateRecipeHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		var reqRecipe models.Recipe

		// Transforma o JSON do body da request em uma struct do modelo Recipe, sem o ID
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&reqRecipe)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		var recipe models.Recipe

		// Seleciona a receita já que se pretende atualizar
		result := app.DB.Where("id = ?", id).First(&recipe)

		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				fmt.Println("Recipe not found")
				http.Error(w, "Not Found", http.StatusNotFound)
				return
			} else {
				fmt.Printf("Error querying recipe: %v\n", result.Error)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}

		// Atualiza seus atributos com os valores da struct da request
		recipe.Name = reqRecipe.Name
		recipe.Instructions = reqRecipe.Instructions
		app.DB.Save(&recipe)

		w.Header().Set("Content-type", "text/plain")
		w.Write([]byte("Recipe updated!"))
	}
}

// @Summary      Deletar receita
// @Description  Deletar receita pelo ID
// @Tags         recipe
// @Produce      text/plain
// @Security Token 
// @Param		 id path int true "ID da receita"
// @Success      200  {string}   string "Recipe deleted!"
// @Failure      404  "Not Found"
// @Failure      500  "Internal Server Error"
// @Router       /recipe/{id} [delete]
func DeleteRecipeHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		var recipe models.Recipe

		result := app.DB.Where("id = ?", id).Delete(&recipe)

		if result.Error != nil {
			fmt.Printf("Error querying recipe: %v\n", result.Error)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if result.RowsAffected == 0 {
			fmt.Println("Recipe not found")
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-type", "text/plain")
		w.Write([]byte("Recipe deleted!"))
	}
}

// @Summary      Adicionar ingrediente à receita
// @Description  Adicionar ingrediente cadastrado à uma receita criada, passando ambos IDs
// @Tags         ingredients_recipes
// @Accept       json
// @Produce      text/plain
// @Security Token 
// @Param		 id path int true "ID da receita"
// @Param		 reqIngredientRecipe body models.IngredientsRecipes true "Ingrediente adicionado"
// @Success      201  {string}   string "Ingredient added!"
// @Failure      400  "Error adding ingredient to recipe"
// @Failure      500  "Internal Server Error"
// @Router       /recipe/{id}/ingredients [post]
func AddIngredientRecipeHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Converte o parâmetro ID de String para uint
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Error converting id param to uint", http.StatusInternalServerError)
			return
		}

		var reqIngredientRecipe models.IngredientsRecipes

		// Transforma o JSON do body da request em uma struct do modelo IngredientsRecipes, sem o ID
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		err = decoder.Decode(&reqIngredientRecipe)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Nova struct do modelo que recebe o RecipeID do parâmetro e demais atributos da struct da request
		newRecipe := models.IngredientsRecipes{
			RecipeID:     uint(id),
			IngredientID: reqIngredientRecipe.IngredientID,
			Quantity:     reqIngredientRecipe.Quantity,
		}

		result := app.DB.Create(&newRecipe)
		if result.Error != nil {
			http.Error(w, "Error adding ingredient to recipe", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-type", "text/plain")
		w.Write([]byte("Ingredient added!"))
		w.WriteHeader(http.StatusCreated)
	}
}

// @Summary      Remover ingrediente da receita
// @Description  Remover ingrediente cadastrado em uma receita, passando ambos IDs
// @Tags         ingredients_recipes
// @Produce      text/plain
// @Param		 id path int true "ID da receita"
// @Security Token 
// @Param		 ingredient_id path int true "ID do ingrediente"
// @Success      200  {string}   string "Ingredient removed from recipe!"
// @Failure      404  "Not Found"
// @Failure      500  "Internal Server Error"
// @Router       /recipe/{id}/ingredients/{ingredient_id} [delete]
func DeleteIngredientRecipeHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		ingredient_id := chi.URLParam(r, "ingredient_id")

		var ingredientRecipe models.IngredientsRecipes

		// Query bicondicional que seleciona somente linhas que possuam, simultaneamente, os ids da receita e do ingrediente passados
		result := app.DB.Where("recipe_id = ? AND ingredient_id = ?", id, ingredient_id).Delete(&ingredientRecipe)

		if result.Error != nil {
			fmt.Printf("Error querying recipe: %v\n", result.Error)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if result.RowsAffected == 0 {
			fmt.Println("Recipe or ingredient not found")
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-type", "text/plain")
		w.Write([]byte("Ingredient removed from recipe!"))
	}
}
