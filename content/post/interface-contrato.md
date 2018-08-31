+++
date = "2018-08-30T20:05:18-03:00"
title = "interface como um contrato"
description = "Dicas de como usar interface em Golang como um contrato para structs"
tags = ["Golang", "Desenvolvimento"]
+++

{{< youtube ck9rGF0tfRE >}}

Falamos desse mesmo tema na nossa lista de exemplos do grupo de estudos sobre interfaces mas vamos isso foi la no longínquo ano de 2016 então vamos tentar fazer um novo exemplo para mostrar como usar uma interface como um contrato que permite você passar tipos diferentes para uma função ou aceitar tipos diferentes em uma interface desde que ele aceite implemente o `contrato` que em outras palavras significa implementar as funções descritas na interface.

No exemplo abaixo temos um pequeno programa que implementa a `struct melancia`, e dai usa essa struct como parâmetro e também para iniciar uma variável. Funciona mas usar tipos concretos dessa forma engessa o programa, não é possível interferir nas funções que serão executadas uma vez que não é possível alterar a função que vai ser chamada.

```go
// usando tipos concretos
package main

import (
	"fmt"
)

type melancia struct {
}

func (m melancia) nome() string {
	return "melancia"
}

func vitaminaDe(m melancia) {
	fmt.Println("vitamina de", m.nome())
}

func main() {
	var x melancia
	x = melancia{}
	nome := x.nome()
	fmt.Println(nome)

	vitaminaDe(x)
}
```

Agora vamos fazer uma pequena alteração em nosso programa e em vez de usar a struct diretamente vamos usar uma interface tanto para declarar a variável como para passar o o parâmetro para a função.

Note a interface `fruta` no exemplo abaixo e que ela implementa a função `nome()`, qualquer struct que implemente essa função pode ser usada como parâmetros para a função tornando o programa bem mais flexível, permitindo receber diferentes tipos de dados.

```go
// usando interface
package main

import (
	"fmt"
)

type fruta interface {
	nome() string
}

type melancia struct {
}

func (b melancia) nome() string {
	return "melancia"
}

func vitaminaDe(f fruta) {
	fmt.Println("vitamina de", f.nome())
}

func main() {
	var x fruta
	x = melancia{}
	nome := x.nome()
	fmt.Println(nome)

	vitaminaDe(x)
}
```

Mais um exemplo usando as duas structs. Como as duas structs satisfazem nosso contrato não podemos passar qualquer uma. Isso é ótimo por exemplo quando se quer fazer mock para testes de um pacote, se você usar uma interface na hora de receber os parâmetros na hora do teste é só substituir a instancia da struct original por uma mais adequada para o teste.

```go
package main

import (
	"fmt"
)

type fruta interface {
	nome() string
}

type melancia struct {
}

func (m melancia) nome() string {
	return "melancia"
}

type banana struct {
}

func (b banana) nome() string {
	return "banana"
}

func vitaminaDe(f fruta) {
	fmt.Println("vitamina de", f.nome())
}

func main() {
	var m fruta
	m = melancia{}
	nome := m.nome()
	fmt.Println(nome)

	vitaminaDe(m)

	b := banana{}
	nome = b.nome()
	fmt.Println(nome)

	vitaminaDe(b)
}
```

E por fim um pequeno truque, sabia que se sua struct implementar uma função `String()` quando você usar ela com o pacote `fmt` em funções como `Print` por exemplo o sistema vai automaticamente usar a sua função para converter a instancia da sua struct para string. Muito bom para compor logs melhores.

## Links úteis

- Código fonte dos exemplos de hoje
    1. https://play.golang.org/p/7YnT6sB2PCb
	2. https://play.golang.org/p/HtM3kYwygBK
	3. https://play.golang.org/p/bpaYKqNMiOR
	4. https://play.golang.org/p/FA8t_26kwAZ
- [Repositório do nosso grupo](https://github.com/go-br/estudos)
- [E você encontra mais exemplos aqui](https://github.com/go-br)
- [Pagina do grupo de estudos](https://gopher.pro.br)
