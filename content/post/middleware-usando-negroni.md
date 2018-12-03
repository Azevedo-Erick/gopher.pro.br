+++
title = "Como fazer HTTP middleware em Go tanto usando Negroni como usando a biblioteca padrão"
description = "Exempos de como criar middleware para seus serviços HTTP mostrando exemplos com a bilioteca padrão e também com a dobradinha Negroni e Gorilla mux."
tags = ["Golang"]
date = "2018-11-02T06:45:14Z"
+++

{{< youtube aQGRh7ECgwg >}}

Usar [middleware](https://en.wikipedia.org/wiki/Middleware) HTTP são muito úteis para evitar duplicidade de código quando você tem vários endpoints na sua aplicação, por exemplo se você quiser ter certeza que as credenciais do usuário foram verificadas, ou que o conteúdo foi comprimido, e assim por diante.

A coisa mais importante que se deve lembrar é que cada middleware vai ser chamado na ordem que foi registrado, então por exemplo podemos ter um middleware que tem a responsabilidade de preparar o ambiente, como abrir o banco de dados ou preparar o controle de sessão que vem antes do middleware que valida as credenciais do usuário.

# via biblioteca padrão

```go
package main

import (
	"fmt"
	"log"
	"net/http"
)

func handleMain(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("{\"value\":42}\n"))
	if err != nil {
		fmt.Println("error handleMain", err)
	}
}

func handleHealthcheck(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("{\"status\":\"ok\"}\n"))
	if err != nil {
		fmt.Println("error handleHealthcheck", err)
	}
}

func applicationJSON(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		h.ServeHTTP(w, r)
	}
}

func basicAuth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path == "/healthcheck" {
			h.ServeHTTP(w, r)
			return
		}

		user, pass, ok := r.BasicAuth()
		if !ok || user != "admin" || pass != "admin" {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintln(w, `{"error": "Unauthorized"}`)
			return
		}

		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		h.ServeHTTP(w, r)
	}
}

func main() {
	http.HandleFunc("/", applicationJSON(basicAuth(handleMain)))
	http.HandleFunc("/healthcheck", applicationJSON(handleHealthcheck))

	fmt.Println("main listen at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

Nesse exemplo temos dois endpoints e um middleware que vai ajustar o header `Content-Type` para `application/json` em todas as requisições HTTP, dessa forma não precisamos mais nos preocupar com isso, não importa quantos endpoints tenhamos na nossa aplicação todos terão esse mesmo cabeçalho. Ou seja evitamos de duplicar esse código em vários pontos do nosso programa.

# Via negroni e gorilla/mux

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

func handleHealthcheck(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("{\"status\":\"ok\"}\n"))
	if err != nil {
		fmt.Println("error handleHealthcheck", err)
	}
}

func applicationJSON() negroni.Handler {
	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	})
}

func basicAuth() negroni.Handler {
	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		if r.URL.Path == "/healthcheck" {
			next(w, r)
			return
		}

		user, pass, ok := r.BasicAuth()
		if !ok || user != "admin" || pass != "admin" {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintln(w, `{"error": "Unauthorized"}`)
			return
		}

		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		next(w, r)
	})
}

func main() {
	n := negroni.Classic()
	n.Use(applicationJSON())
	n.Use(basicAuth())

	r := mux.NewRouter().StrictSlash(true)
	n.UseHandler(r)

	r.HandleFunc("/", handleMain).Methods("GET")
	r.HandleFunc("/healthcheck", handleHealthcheck).Methods("GET")

	fmt.Println("main listen at :8080")
	err := http.ListenAndServe(":8080", n)
	if err != nil {
		fmt.Println(err)
	}
}
```

Agora temos o mesmo exemplo mas usando a clássica dobradinha [Negroni](https://github.com/urfave/negroni) e [Gorilla mux](https://github.com/gorilla/mux). No nosso repositório do [grupo de estudos de Golang](https://github.com/go-br/estudos) temos mais alguns exemplos.

- [Código fonte de hoje](https://github.com/go-br/estudos/tree/master/http_middleware)
- [Repositório do nosso grupo](https://github.com/go-br/estudos)
- [E você encontra mais exemplos aqui](https://github.com/go-br)
- [Pagina do grupo de estudos](https://gopher.pro.br)

Nossos encontros ocorrem todas as quintas-feiras ás 22h00, para participar [entre no canal de Go no slack](https://invite.slack.golangbridge.org/) e procure por #brazil