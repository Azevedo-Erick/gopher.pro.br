+++
title = "Funções"
date = "2017-15-09T19:46:45-03:00"
description = "Um resumo sobre funções em Go."
tags = ["golang"]
+++

# Funções

Descreve o funcionamento de funções em Golang

* No Go não há metodos, só funções
* Retorno multiplo - Facilita a escrita e torna o código mais limpo
* Retorno assinado - Facilita a leitura do código
* Os parametros por padrão são passados por valor
* Utilizando ponteiros é possível passar valores por referencia

```go
package main

import (
	"fmt"
)

// Retorno simples
func soma(x int, y int) int {
	return x + y
}

// Retorno duplo
func troca(x string, y string) (string, string) {
	return y, x
}

//Retorno assinado
func divide(x, y int) (resultado, resto int) { // os dois retornos são inteiros nesse exemplo
	resto = x % y
	resultado = x / y
	return
}

// função que recebe uma função como parâmetro
func executaFuncao(f func(string) string, valor string) {
	aux := f(valor)
	fmt.Printf(aux)
}

func printValorByRef(valor *string) {
	fmt.Printf("Valor por referencia = %v\r\n", *valor)
}

//funcao variadica - recebe como argumentos infinitos parâmetros do tipo inteiro
func variadicaInteiros(v ...int) {
	fmt.Println(v)
}

//funcao variadica - recebe como argumentos infinitos parâmetros do tipo interface
func variadicaInterface(v ...interface{}) {
	fmt.Println(v)
}

func main() {
	fmt.Printf("Funções!\r\n")

	fmt.Printf("Soma 1+1 = %v\r\n", soma(1, 1))

	b, a := troca("A", "B")
	fmt.Printf("troca A, B = %v, %v\r\n", b, a)

	resu, rest := divide(5, 2)
	fmt.Printf("A divisão de 5 por 2 é = %v\r\n", resu)
	fmt.Printf("O resto da divisão de 5 por 2 é = %v\r\n", rest)

	// função anonima que vamos passar para printFunc
	ola := func(v string) string {
		return "Olá " + v + "!\r\n"
	}

	executaFuncao(ola, "Cesar")

	valor := "Esse valor não vai ser copiado, só estamos passando o ponteiro"
	printValorByRef(&valor)

	//retorno da função variadica, podendo passar infinitos parâmentros do tipo inteiro
	variadicaInteiros(1, 2, 3, 4, 5)

	//retorno da função variadica, podendo passar infinitos parâmentros do tipo interface
	variadicaInterface(1, "dois", 3, "quatro", 5, true, false)

}

}
```
[Playground](https://play.golang.org/p/zpAGhv4ZB-H)