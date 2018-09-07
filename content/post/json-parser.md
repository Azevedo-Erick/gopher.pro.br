+++
date = "2018-09-06T08:19:09-03:00"
description = "Tutorial de como ensinar o parser JSON para trabalhar com dados em novos formatos"
tags = ["Golang", "Desenvolvimento"]
title = "JSON, criando seu próprio Marshal e Unmarshal"
topics = []

+++

{{< youtube ALlzgQfXVFo >}}

Continuando a conversa sobre interfaces e sobre manipulação de JSON tem um recurso muito útil que usamos para ajudar o nosso sistema a falar melhor com o PostgresQL e pREST.

O formato de data padrão do Postgres é incompatível com o formato padrão do Go e o nosso sistema usa muito JSON para tanto mandar como receber informações do banco de dados.

Ter que lembrar toda hora de fazer o parser da data para o formato correto simplesmente não é pratico, é muito melhor ensinar o Go como lidar com data e hora no bom e velho formato ISO 8601.

Alias recomendo muito sempre usar ISO 8601 UTC para tudo no backend e apenas mudar para o formato e timezone local quando for exibir para o usuário. Mas isso é uma história para outro dia.

O código fonte do package Go [esta disponível no GitHub](https://github.com/nuveo/dbtime).

### MarshalJSON

Agora vamos ver como o código funciona, no exemplo abaixo criamos uma struct `Time` e implementamos uma função MarshalJSON para ela. Na função `main` instanciamos uma struct com alguns dados, em seguida usamos a função `json.MarshalIndent` que percorre a struct e quando ela encontrar o campo Time vai usar a função que definimos e não a default do sistema.

```go
package main

import (
	"encoding/json"
	"fmt"
	"time"
)

type Time struct {
	time.Time
}

const layout = "2006-01-02T15:04:05.999999"

func (t Time) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, t.Time.Format(layout))), nil
}

func main() {
	data := struct {
		Name string
		Time Time
	}{
		Name: "teste",
		Time: Time{time.Now().Add(time.Millisecond * time.Duration(54321))},
	}

	json, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(json))
}
```

Exemplo de retorno usando `The Go Playground` veja que o formato da data obedece a nossa função.

```json
{
	"Name": "teste",
	"Time": "2009-11-10T23:00:54.321"
}
```

### UnmarshalJSON

No proximo exemplo vamos fazer o contrario, agora definimos uma função UnmarshalJSON para nossa struct. Veja o exemplo.


```go
package main

import (
	"encoding/json"
	"fmt"
	"time"
)

type Time struct {
	time.Time
}

const layout = "2006-01-02T15:04:05.999999"

func (t *Time) UnmarshalJSON(b []byte) (err error) {
	if b[0] == '"' && b[len(b)-1] == '"' {
		b = b[1 : len(b)-1]
	}
	if string(b) == `null` {
		*t = Time{}
		return
	}
	t.Time, err = time.Parse(layout, string(b))
	return
}

func main() {
	data := struct {
		Name string
		Time Time
	}{}

	b := []byte(`{"Name": "teste", "Time": "2009-11-10T23:00:54.321"}`)

	err := json.Unmarshal(b, &data)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(data.Time.String())
}
```

Quando o sistema recebe o array de bytes para fazer parse e popular a struct ele percorre os dados e ao encontrar o campo Time o pacote usa a nossa função `UnmarshalJSON` no lugar da default do sistema, assim conseguimos ler corretamente os dados mesmo não sendo o padrão do sistema.

Alem de formatar os dados corretamente também podemos usar esse recurso para outras tarefas, por exemplo é possível percorrer os campos verificando as informações de um determinado tipo e retornar um erro caso encontre algum dado invalido, é uma forma de validação de dados que trabalha internamente no parser do tipo e pode ser muito útil, só devemos tomar cuidado porque pode acabar ocultando de onde o erro esta vindo, então escreva boas mensagens de erro. No exemplo do manual do proprio pacote json o exemplo é um contador que conta animais em um Zoo.

## Links úteis

- Código fonte dos exemplos de hoje
    1. https://github.com/nuveo/dbtime
    2. https://play.golang.org/p/kN-_1lLoggl
    3. https://play.golang.org/p/QaHxm0At5_V
    4. https://play.golang.org/p/owTOkyHi8up
- [Repositório do nosso grupo](https://github.com/go-br/estudos)
- [E você encontra mais exemplos aqui](https://github.com/go-br)
- [Pagina do grupo de estudos](https://gopher.pro.br)
