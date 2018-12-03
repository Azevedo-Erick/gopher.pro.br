+++
title = "Trafegando dados entre middleware http usando contexto em Golang"
description = "Como trocar informações entre os middlewares http do seu sistema usando contexto."
tags = ["Golang"]
date = "2018-11-09T06:45:14Z"
+++

{{< youtube Xyj-dQvfvq0 >}}

# Trafegando dados entre middleware

Já vimos [como funciona um middleware HTTP](https://www.youtube.com/watch?v=aQGRh7ECgwg) e agora vamos ver como passar informação entre os middlewares. Isso é usado por exemplo para passar as credenciais de um usuário para o proximo middleware e qualquer outra informação que seja coletada em algum dos middlewares e você queria passar para frente.

Antes de mais nada vamos fazer um exemplo para mostrar da forma mais clara possível quando os middlewares são executados, isso é muito importante porque erros no entendimento dessa ordem de execução é uma incrível fonte de bugs. No exemplo abaixo colocamos mensagens no terminal toda vez que um middleware é carregado e descarregado. E é sempre bom reforçar isso vai acontecer na mesma ordem em que eles foram registrados.

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func handleMain(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("{\"value\":42}\n"))
	if err != nil {
		fmt.Println("error handleMain", err)
	}
}

func middleware1() negroni.Handler {
	fmt.Println("carregando middleware 1")
	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		fmt.Println("empilhando middleware 1")
		next(w, r)
		fmt.Println("desempilhando middleware 1")
	})
}

func middleware2() negroni.Handler {
	fmt.Println("carregando middleware 2")
	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		fmt.Println("empilhando middleware 2")
		next(w, r)
		fmt.Println("desempilhando middleware 2")
	})
}

func middleware3() negroni.Handler {
	fmt.Println("carregando middleware 3")
	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		fmt.Println("empilhando middleware 3")
		next(w, r)
		fmt.Println("desempilhando middleware 3")
	})
}

func main() {
	n := negroni.Classic()
	n.Use(middleware1())
	n.Use(middleware2())
	n.Use(middleware3())
	fmt.Println("-=-=-=-=-=-=-=-=-=-=-=-=-=-")
	r := mux.NewRouter().StrictSlash(true)
	n.UseHandler(r)

	r.HandleFunc("/", handleMain).Methods("GET")

	fmt.Println("main listen at :8080")
	err := http.ListenAndServe(":8080", n)
	if err != nil {
		fmt.Println(err)
	}
}
```

A maneira canônica de passar informações a diante durante o processamento de uma requisição HTTP é usando contexto. Veja o exemplo.

```go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

type key int

const (
	dataKey key = iota
)

type data struct {
	ValueA string `json:"value_a"`
	ValueB int    `json:"value_b"`
}

func setContextData(r *http.Request, d *data) (ro *http.Request) {
	ctx := r.Context()
	ctx = context.WithValue(ctx, dataKey, d)
	ro = r.WithContext(ctx)
	return
}

func getContextData(r *http.Request) (d data) {
	d = *r.Context().Value(dataKey).(*data)
	return
}

func middleware1() negroni.Handler {
	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		d := data{
			ValueA: "valor A",
			ValueB: 42,
		}
		r = setContextData(r, &d)
		next(w, r)
	})
}

func middleware2() negroni.Handler {
	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		d := getContextData(r)
		d.ValueA += "A"
		r = setContextData(r, &d)
		next(w, r)
	})
}

func middleware3() negroni.Handler {
	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		d := getContextData(r)
		d.ValueA += "A"
		r = setContextData(r, &d)
		next(w, r)
	})
}

func handleMain(w http.ResponseWriter, r *http.Request) {
	d := getContextData(r)
	j, err := json.MarshalIndent(d, "", "\t")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(j)
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	n := negroni.Classic()
	n.Use(middleware1())
	n.Use(middleware2())
	n.Use(middleware3())

	r := mux.NewRouter().StrictSlash(true)
	n.UseHandler(r)

	r.HandleFunc("/", handleMain).Methods("GET")

	fmt.Println("main listen at :8080")
	err := http.ListenAndServe(":8080", n)
	if err != nil {
		fmt.Println(err)
	}
}
```

Tem muito mais exemplos no repositório do nosso [grupo de estudos de Go](https://github.com/go-br/estudos), vale a pena dar uma conferida e também colaborar, estamos sempre precisando de ajuda para ter uma material completo e atualizado.

- [Código fonte de hoje](https://github.com/go-br/estudos/tree/master/http_middleware)
- [Repositório do nosso grupo](https://github.com/go-br/estudos)
- [E você encontra mais exemplos aqui](https://github.com/go-br)
- [Pagina do grupo de estudos](https://gopher.pro.br)

Nossos encontros ocorrem todas as quintas-feiras ás 22h00, para participar [entre no canal de Go no slack](https://invite.slack.golangbridge.org/) e procure por #brazil