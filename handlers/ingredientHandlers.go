package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"main.go/app"
	"main.go/models"
)

// @Summary      Buscar todos os ingredientes
// @Description  Buscar todos os ingredientes cadastrados
// @Tags         ingredient
// @Produce      json
// @Success      200  {array}   models.Ingredient
// @Failure      404    "Not Found"
// @Failure      500    "Internal Server Error"
// @Router       /ingredient [get]
func GetAllIngredientsHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var ingredients []models.Ingredient

		result := app.DB.Find(&ingredients)

		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				http.Error(w, "Not Found", http.StatusNotFound)
				return
			} else {
				fmt.Printf("Error querying ingredients: %v\n", result.Error)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}

		if len(ingredients) == 0 {
			http.Error(w, "No ingredients found", http.StatusNotFound)
			return
		}

		ingredientsJson, err := json.Marshal(ingredients)

		if err != nil {
			http.Error(w, "Error encoding ingredients to JSON", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(ingredientsJson)
	}
}

// @Summary      Buscar ingrediente pelo ID
// @Description  Buscar ingrediente pelo ID
// @Tags         ingredient
// @Produce      json
// @Param		 id path int true "ID do ingrediente"
// @Success      200  {array}   models.Ingredient
// @Failure      404  "Not Found"
// @Failure      500  "Internal Server Error"
// @Router       /ingredient/{id} [get]
func GetIngredientByIdHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		var ingredient models.Ingredient

		result := app.DB.Where("id = ?", id).First(&ingredient)

		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				http.Error(w, "Ingredient not Found", http.StatusNotFound)
				return
			} else {
				fmt.Printf("Error querying ingredients: %v\n", result.Error)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}

		ingredientJson, err := json.Marshal(ingredient)

		if err != nil {
			http.Error(w, "Error encoding ingredient to JSON", http.StatusInternalServerError)
			return
		}

		log.Println(ingredientJson)
		w.Header().Set("Content-Type", "application/json")
		w.Write(ingredientJson)

	}
}

// @Summary      Buscar ingrediente pelo nome
// @Description  Buscar ingrediente pelo nome sem case sensitive e convertendo '-' para espaços
// @Tags         ingredient
// @Produce      json
// @Param		 name path string true "Nome do ingrediente"
// @Success      200  {array}   models.Ingredient
// @Failure      404  "Not Found"
// @Failure      500  "Internal Server Error"
// @Router       /ingredient/name/{name} [get]
func GetIngredientByNameHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")
		name = "%" + name + "%"

		var ingredient []models.Ingredient

		result := app.DB.Where("LOWER(name) LIKE ?", strings.ToLower(name)).Find(&ingredient)

		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				http.Error(w, "Ingredient not Found", http.StatusNotFound)
				return
			} else {
				fmt.Printf("Error querying ingredients: %v\n", result.Error)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}

		ingredientJson, err := json.Marshal(ingredient)

		if err != nil {
			http.Error(w, "Error encoding ingredient to JSON", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(ingredientJson)

	}
}

// @Summary      Criar novo ingrediente
// @Description  Criar novo ingrediente
// @Tags         ingredient
// @Security Token 
// @Accept       json
// @Produce      text/plain
// @Param		 ingredient body models.Ingredient true "Novo ingrediente"
// @Success      201  {string}   string "Ingredient created!"
// @Failure      400  "Invalid JSON"
// @Router       /ingredient [post]
func CreateIngredientHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var ingredient models.Ingredient

		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()

		err := decoder.Decode(&ingredient)
		if err != nil || ingredient.Name == "" {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		result := app.DB.Create(&ingredient)

		if result.Error != nil {
			http.Error(w, "Ingredient already exists or data is incorrect", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Ingredient created!"))
	}
}

// @Summary      Atualizar ingrediente
// @Description  Atualizar ingrediente pelo ID
// @Tags         ingredient
// @Accept       json
// @Security Token 
// @Produce      text/plain
// @Param		 id path int true "ID do ingrediente"
// @Param		 ingredient body models.Ingredient true "Ingrediente atualizado"
// @Success      200  {string}   string "Ingredient updated!"
// @Failure      400  "Invalid JSON"
// @Failure      404  "Not Found"
// @Failure      500  "Internal Server Error"
// @Router       /ingredient/{id} [put]
func UpdateIngredientHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var reqIngredient models.Ingredient

		// Transforma body da request para uma struct, sem o ID
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()

		err := decoder.Decode(&reqIngredient)
		if err != nil || reqIngredient.Name == "" {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Cria struct para armazenar dados do atual ingrediente
		var ingredient models.Ingredient

		result := app.DB.Where("id = ?", id).First(&ingredient)

		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				fmt.Println("Ingredient not found")
				http.Error(w, "Not Found", http.StatusNotFound)
				return
			} else {
				fmt.Printf("Error querying ingredient: %v\n", result.Error)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}

		// Struct com o ingrediente recebe nome da struct da request
		ingredient.Name = reqIngredient.Name
		app.DB.Save(&ingredient)

		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Ingredient updated!"))
	}
}

// @Summary      Deletar ingrediente
// @Description  Deletar ingrediente pelo ID
// @Tags         ingredient
// @Produce      text/plain
// @Security Token 
// @Param		 id path int true "ID do ingrediente"
// @Success      200  {string}   string "Ingredient deleted!"
// @Failure      404  "Not Found"
// @Failure      500  "Internal Server Error"
// @Router       /ingredient/{id} [delete]
func DeleteIngredientHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		var ingredient models.Ingredient

		result := app.DB.Where("id = ?", id).Delete(&ingredient)

		if result.Error != nil {
			fmt.Printf("Error querying user: %v\n", result.Error)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if result.RowsAffected == 0 {
			fmt.Println("User not found")
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Ingredient deleted!"))
	}
}
