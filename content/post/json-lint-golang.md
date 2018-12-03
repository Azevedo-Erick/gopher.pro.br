+++
title = "Um JSON lint em Golang"
description = "Usamos o retorno de erro do json.Unmarshal para gerar uma mensagem de erro mais útil e completa com direito a indicar o erro com uma setinha e tudo."
tags = ["Golang"]
date = "2018-11-30T06:45:14Z"
+++

{{< youtube VXVZE5SYMiM >}}

# Um JSON lint em Golang

Esse é um pequeno utilitário de linha de comando para validar e formatar JSON que também pode ser usado como pacote, a ideia inicial era criar um parser para validar o JSON mas não foi necessário, o próprio Go traz todas as informações que precisamos no erro inclusive o offset do erro e dai foi simples mostrar os erros de forma mais completa inclusive indicando com uma seta exatamente onde esta o erro.

```go
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/crgimenes/goconfig"
	"github.com/gosidekick/jsonlint"
)

func printError(a ...interface{}) {
	_, err := fmt.Fprintf(os.Stderr, "\x1b[91m%v\033[0;00m\n", a...)
	if err != nil {
		fmt.Println(err)
	}
}

func printIndicator(a ...interface{}) {
	_, err := fmt.Fprintf(os.Stderr, "\x1b[96m%v\033[0;00m\n", a...)
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	type configFlags struct {
		Input  string `json:"i" cfg:"i" cfgDefault:"stdin" cfgHelper:"input from"`
		Output string `json:"o" cfg:"o" cfgDefault:"stdout" cfgHelper:"output to"`
	}

	cfg := configFlags{}
	goconfig.PrefixEnv = "JSON_LINT"
	err := goconfig.Parse(&cfg)
	if err != nil {
		printError(err)
		os.Exit(-1)
	}
	var j []byte

	if cfg.Input == "stdin" {
		j, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			printError(err)
			os.Exit(-1)
		}
	} else {
		j, err = ioutil.ReadFile(cfg.Input)
		if err != nil {
			printError(err)
			os.Exit(-1)
		}
	}
	var m interface{}
	err = json.Unmarshal(j, &m)
	if err != nil {
		out, offset := jsonlint.ParseJSONError(j, err)
		printError(out)
		if offset > 0 {
			out = jsonlint.GetErrorJSONSource(j, offset)
			printIndicator(out)
		}
		os.Exit(-1)
	}
	j, err = json.MarshalIndent(m, "", "\t")
	if err != nil {
		printError(err)
		os.Exit(-1)
	}
	if cfg.Output == "stdout" {
		fmt.Println(string(j))
		return
	}
	err = ioutil.WriteFile(cfg.Output, j, 0644)
	if err != nil {
		printError(err)
	}
}
```

Aqui temos o código do nosso utilitário que também serve de exemplo de como utilizar nosso pacote de tratamento de erro, basicamente basta passar o erro da função `json.Unmarshal` para a função `jsonlint.ParseJSONError` e o retorno vai ser uma string contendo uma descrição mais detalhada do erro e o offset de onde esse erro ocorreu, você pode passar o JSON que originou o erro e esse offset para a função `jsonlint.GetErrorJSONSource` e ela vai gerar uma string contendo a parte do código com problemas junto com uma seta apontando para o primeiro caracter que o parser não conseguiu processar.

Esse utilitário assim como muitos outros esta sendo construído no nosso repositório do [GoSidekick](https://github.com/gosidekick) a ideia é criar um conjunto de utilitários que usamos no nosso dia a dia desenvolvendo código Go e também outras linguagens de programação.

- [Código fonte de hoje](https://github.com/gosidekick/jsonlint)
- [Repositório do nosso grupo](https://github.com/go-br/estudos)
- [E você encontra mais exemplos aqui](https://github.com/go-br)
- [Pagina do grupo de estudos](https://gopher.pro.br)

Nossos encontros ocorrem todas as quintas-feiras ás 22h00, para participar [entre no canal de Go no slack](https://invite.slack.golangbridge.org/) e procure por #brazil