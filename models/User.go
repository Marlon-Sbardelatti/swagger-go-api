package models

// User representa um usuário do sistema.
// @Description Modelo para gerenciar os usuários do sistema.
type User struct {
	// ID é o código de identificação único do usuário.
	ID uint `gorm:"primaryKey"`
	// Username é o nome único do usuário no sistema.
	Username string `gorm:"unique;not null" example:"seunome"`
	// Email é o email único do usuário.
	Email string `gorm:"unique;not null" example:"seuemail@gmail.com"`
	// Password é a senha de entrada do usuário no sistema.
	Password string `gorm:"not null"`
}

// UserLoginRequest representa as informações de login do usuário no sistema.
// @Description Modelo para gerenciar a entrada dos usuários no sistema.
type UserLoginRequest struct {
	// Email é o email único que permite a identificação do usuário no sistema.
	Email string `example:"seuemail@gmail.com"`
	// Password é a senha cadastrada do usuário com o e-mail informado.
	Password string
}
