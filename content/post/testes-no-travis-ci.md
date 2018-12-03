+++
title = "Testando código Golang no Travis-CI usando Docker"
description = "Como testar seu codigo Golang no Travis-CI usando Docker."
tags = ["Golang"]
date = "2018-10-17T06:45:14Z"
+++

{{< youtube 7Y6xW7i1EXo >}}

# Testando o código

Essa é uma daquelas gambiarras que salva nossa pele de vez em quando. Eu precisava rodar um código com um ambiente muito especifico e a única forma de montar o ambiente desse código com todas as peculiaridades dele é criar um containers [docker](https://www.docker.com), dessa forma consigo dizer exatamente a versão do que esta instalado e ter certeza que o ambiente esta perfeito. Mas daí eu precisava rodar os testes que são executados pelo [Travis-CI](https://travis-ci.org) e claro que não conseguia montar o ambiente idêntico ao do container no Travis. A solução foi rodar os testes dentro do container via `docker run` e pegar o [errorlevel](https://en.wikipedia.org/wiki/Exit_status) no retorno.

Exemplo de como rodar os testes dentro de um container:

Primeiro criei um `Dockerfile` para gerar a imagem contendo o ambiente que quero testar.

```dockerfile
FROM golang:stretch

# RUN apt update -y && apt install -y "some library heavily dependent of the Linux version, like libmagick++-dev"
RUN mkdir -p /go/src/github.com/crgimenes/tgwd
COPY ./ /go/src/github.com/crgimenes/tgwd
WORKDIR /go/src/github.com/crgimenes/tgwd
RUN go build
CMD ["./tgwd"]
```

Daí criei a imagem normalmente com o `docker build`

```console
docker build -t tgwd .
```

E para testar localmente é só usar `docker run` como no exemplo:

```console
docker run tgwd bash -c "go test ./..."
```

Para testar no Travis preparei um arquivo de configuração `.travis.yml` como este exemplo:

```yml
sudo: false

dist: xenial

language: go

services:
  - docker

os:
  - linux

go:
  - "1.11"

before_install:
  - docker build -t tgwd .

script:
  - docker run tgwd bash -c "go test ./..."
```

Claro que o ambiente go aqui nem era mesmo necessário porque tudo esta dentro da imagem e.o travis na verdade esta testando o código rodando lá dentro, mas eu já deixei assim por pura força do habito.

Eu não gosto muito de rodar tudo em container mas tenho que admitir que é bem interessante quando o ambiente que o sistema precisa é muito alienígena.

- [Código fonte](https://github.com/crgimenes/tgwd)
- [Cesar Gimenes](https://crg.eti.br)
