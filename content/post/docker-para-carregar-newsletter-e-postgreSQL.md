+++
title = "Usando docker para carregar o Newsletter e PostgreSQL"
description = "Colocando nosso aplicativo em containers docker."
tags = ["Golang"]
date = "2018-10-12T06:45:14Z"
+++

{{< youtube jmCE2S65zNQ >}}

# Usando Docker

Nesse episódio do grupo de estudos vimos como colocar uma aplicação para rodar dentro de um container [Docker](https://www.docker.com), também criamos um Makefile para ajudar a executar os comandos docker e todo o básico para rodar nosso newsletter em um container.

Nosso exemplo de Makefile

```makefile
.PHONY: up up.build rmi stop down purge

up:
	docker-compose up

up.build:
	docker-compose up --build

down:
	docker-compose down --remove-orphans

stop:
	docker-compose stop

rmi:
	docker system prune -f

purge:
	docker system prune -fa
```

Esse é nosso exemplo de Dockerfile já usando múltiplos estágios para só levar para a imagem final o executável e o mínimo para funcionar assim gerando uma imagem bem pequena.

```Dockerfile
FROM golang as builder

RUN mkdir -p /go/src/github.com/crgimenes/newsletter
COPY . /go/src/github.com/crgimenes/newsletter
WORKDIR /go/src/github.com/crgimenes/newsletter
RUN go get ./... && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o nl

FROM alpine:3.6
RUN apk add --no-cache ca-certificates
WORKDIR /
COPY --from=builder go/src/github.com/crgimenes/newsletter/newsletter/template /newsletter/template
COPY --from=builder go/src/github.com/crgimenes/newsletter/nl /nl
ENTRYPOINT ["/nl"]
```

E finalmente nosso exemplo de docker-compose que vai ajudar a subir nosso projeto facilmente com as dependências e variáveis tudo pronto para funcionar. No Mac o docker-compose é instalado junto com o Docker mas se não funcionar no seu ambiente de desenvolvimento verifique se esta instalado corretamente.

```yml
version: "2"
services:
    redis:
        image: redis:latest
        ports:
            - 6379:6379
        volumes:
            - ./data/redis:/data
    postgres:
        image: postgres:latest
        environment:
            - POSTGRES_DB=erp
            - POSTGRES_USER=postgres
        ports:
            - 5432:5432
        volumes:
            - ./data/postgres:/var/lib/postgresql/data
    newsletter:
        image: newsletter
        build: ./
        depends_on:
            - postgres
            - redis
        environment:
            - SENDGRID_API_KEY=XXXXXXXXXXXX
        ports:
            - 8080:8080
```

- [Código fonte de hoje](https://github.com/crgimenes/newsletter)
- [Repositório do nosso grupo](https://github.com/go-br/estudos)
- [E você encontra mais exemplos aqui](https://github.com/go-br)
- [Pagina do grupo de estudos](https://gopher.pro.br)

Nossos encontros ocorrem todas as quintas-feiras ás 22h00, para participar [entre no canal de Go no slack](https://invite.slack.golangbridge.org/) e procure por #brazil