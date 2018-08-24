+++
date = "2018-08-23T20:05:18-03:00"
title = "interface vazia"
description = "Dicas de como usar interface vazia e seus riscos"
tags = ["Golang", "Desenvolvimento"]
+++

{{< youtube yBTzaVwqvng >}}

Interface vazia é um tipo que aceita qualquer coisa, você pode passar o que quiser como parâmetros de função ou variáveis do tipo `interface{}`.

Parece muito pratico a primeira vista mas quando usamos `interface{}` estamos jogando pela janela a validação de tipos feita em tempo de compilação e perdemos uma das grandes vantagens de uma linguagem compilada de tipagem forte e estática.

E como a checagem de tipo não vai acontecer em tempo de compilação é sua responsabilidade checar se esta recebendo o tipo certo em tempo de execução.

## Identificando o tipo

Quando queremos saber o tipo de um determinado principalmente quando estamos depurando e tentando ver se o tipo que chegou foi o topo esperado podemos usar a função `fmt.Printf` passando o formato `%T` como no exemplo. 


```go
var value interface{}
value = 1
fmt.Printf("tipo de value: %T\n", value)
```

No nosso exemplo declaramos a variável `value` como `interface{}` e agora podemos passar qualquer valor para ela, como passamos um inteiro a função `Printf` com o formato `%T` vai informar que a variável é do tipo `int`.

Vamos ver um exemplo completo

```go
package main

import (
	"fmt"
)

func main() {
	var value interface{}
	value = 1
	fmt.Printf("tipo de value: %T\n", value)
	value = 3.14
	fmt.Printf("tipo de value: %T\n", value)
	value = "isso é uma string"
	fmt.Printf("tipo de value: %T\n", value)
}
```

Esse pequeno programa deve mostrar o seguinte resultado

```console
tipo de value: int
tipo de value: float64
tipo de value: string
```

## Switch

Inspecionar o tipo de uma variável como vimos é bem simples, ajuda bastante a depurar o código e saber se estamos recebendo o que esperamos, agora vamos fazer nosso código fazer esse trabalho sozinho.

Vamos usar a combinação do switch com o cast `(type)` que vai retornar o tipo da variável, veja o exemplo.

```go
package main

import (
	"fmt"
)

func main() {
	var value interface{}
	value = 1
	switch value.(type) {
	case int:
		fmt.Println("Value é do tipo int")
	case string:
		fmt.Println("Value é do tipo string")
	default:
		fmt.Printf("Tipo %T não implementado\n", value)

	}
}
```

Se a variável `value` conter um `int` nosso programa vai retornar `Value é do tipo int`, se passarmos uma string para `value` o retorno vai ser `Value é do tipo string` e se passarmos qualquer outro valor o programa vai avisar que o suporte para aquele tipo ainda não esta implementado.

Uma vez que já sabemos o tipo da variável podemos fazer cast para o tipo correto e dai e dai usar sem problemas.

```go
switch value.(type) {
case int:
	fmt.Println(value.(int))
case string:
	fmt.Println(value.(string))
```

Se não fizermos essa validação e simplesmente usar o cast para tentar converter a variável para o tipo que queremos corremos o risco do programa quebrar com um `panic` caso o tipo passado não seja o que estamos esperando.

## map[string]interface{}

Um uso muito comum é usar interfaces vazias em conjunto com mapas para converter strings JSON para estruturas que não conhecemos com certeza o formato. Como não sabemos quais campos e nem os tipos dos campos a estrutura contem podemos usar um `map[string]interface{}` que basicamente casa com qualquer estrutura.

```go
package main

import (
	"fmt"
	"encoding/json"
)

func main() {
	b := []byte(`{"Name":"Cesar","Value":10}`)
	
	var m map[string]interface{}
	m = make(map[string]interface{})
	
	err := json.Unmarshal(b, &m)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%#v\n", m)
}
```

Neste exemplo primeiro declaramos um array de bytes com uma estrutura JSON, depois declaramos um mapa de strings e interfaces vazias, e instanciamos na memória. Daí usamos `json.Unmarshal` para ler os array de bytes e popular o mapa com os campos e valores encontrados. Verificamos se houve algum erro porque sempre tem a possibilidade do JSON ser inválido e então usamos mais um pequeno truque para mostrar a estrutura. O formato `%#v` na função `fmt.Printf` para mostrar mais dados e não apenas o valor como aconteceria se usássemos apenas o formato `%v`.

Vamos ver um exemplo mais completo usando `range` para percorrer os campos do JSON e `switch` para desviar o código para o tipo correto. 

```go
package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	b := []byte(`{"Name":"Banana","Value":2.10}`)

	var m map[string]interface{}
	m = make(map[string]interface{})

	err := json.Unmarshal(b, &m)
	if err != nil {
		fmt.Println(err)
		return
	}
	for k, v := range m {
		switch v.(type) {
		case float64:
			fmt.Printf("%v %v\n", k, v.(float64))
		case string:
			fmt.Printf("%v %v\n", k, v.(string))
		default:
			fmt.Printf("Tipo %T não implementado\n", v)

		}
	}
}
```

Mais uma forma de como validar se a interface vazia tem o tipo que você quer é fazendo como no exemplo abaixo.

```go
package main

import "fmt"

func main() {
	var value interface{}
	value = 1

	str, ok := value.(string)
	if !ok {
		fmt.Println("Value não é string")
	}
	fmt.Println(str)
}
```

Aqui usamos `str, ok := value.(string)` de forma que se `value` for do tipo `string` ok será true e str contera o valor propriamente dito já como string, é uma forma simples de testar o tipo sem escrever o `switch case`.

Espero ter esclarecido mais um pouco sobre interface vazia, seus usos e porque pode ser arriscado.  

## Links úteis

- Código fonte dos exemplos de hoje
    1. https://play.golang.org/p/mf317GEVSo3
	2. https://play.golang.org/p/yAVCQjdk56Q
	3. https://play.golang.org/p/ucMSAgbGIEO
	4. https://play.golang.org/p/-eLaKhS-ojh
	5. https://play.golang.org/p/vrelMrzXnN-
	6. https://play.golang.org/p/cA-JutjQWWf
- [Repositório do nosso grupo](https://github.com/go-br/estudos)
- [E você encontra mais exemplos aqui](https://github.com/go-br)
- [Pagina do grupo de estudos](https://gopher.pro.br)
