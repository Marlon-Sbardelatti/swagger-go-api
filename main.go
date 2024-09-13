package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	// "github.com/swaggo/http-swagger"
	// "github.com/swaggo/http-swagger/swaggerFiles"
	"main.go/app"
	"main.go/db"
	// "main.go/docs"
	"main.go/routes"
)

// @title           Cookbook API
// @version         1.0
// @description     Crie suas próprias receitas e faça um incrível livro de receitas!

// @contact.name   Sofia Sousa Lindner
// @contact.email  sslindner@furb.br

// @license.name  MIT
// @license.url   https://github.com/Marlon-Sbardelatti/go-rest-api/blob/master/LICENSE

// @host      localhost:3000
// @BasePath  /

// @SecurityDefinitions.apikey Token
// @In header
// @Name Authorization

func main() {
	// Inicializa conexão com banco e cria DAO
	db := db.InitDB()
	app := &app.App{DB: db}

	// Cria o router e registra as rotas do servidor
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	// r.Get("/swagger/*", httpSwagger.WrapHandler(swaggerFiles.Handler))
	routes.RegisterRoutes(r, app)

	log.Println("Server running on Port 3000")
	http.ListenAndServe(":3000", r)
}
