package main

import (
	"fmt"
	"strconv"
)

func main() {
	login()
}

func login() {
	for {
		fmt.Println("(1) - Criar conta\n(2) - Entrar")
		fmt.Print("Escolha uma opção: ")

		var input string
		fmt.Scanln(&input)

		res, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("\nResposta inválida. Tente novamente:")
			continue
		}

		if res == 1 || res == 2 {
			fmt.Printf("Você escolheu a opção: %d\n", res)
			break
		} else {
			fmt.Println("Opção inválida. Tente novamente:")
		}
	}
}
