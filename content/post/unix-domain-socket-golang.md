+++
title = "Unix Domain Socket com Golang"
description = "Como usar Unix Domain Socket com Go"
tags = ["Golang"]
date = "2018-12-06T06:45:14Z"
+++

{{< youtube vOxx0xcNQGg >}}

# Unix domain socket

Unix Domain Sockets ou [IPC](https://en.wikipedia.org/wiki/Inter-process_communication) socket é uma forma muito pratica e segura de trocar informações entre processos. Essa forma de IPC usa um arquivo como endereço/name space no lugar de um IP e uma porta como seria em uma comunicação via rede.

Uma coisa importante para ter em mente é que como vamos usar um arquivo o servidor é responsável por ele, se não existir ele sera criado automaticamente mas se não existir você vai receber um erro com algo como "bind: address already in use" que significa que o arquivo já existe e o servidor não tem como reaproveitar um arquivo que já existe, o correto é fazer shutdown elegantemente e fechar e apagar o arquivo antes de derrubar o servidor. E dependendo do sistema pode ser interessante verificar se o arquivo já existe e apagar antes de subir o servidor. 

Apesar da facilidade, como usamos um arquivo como endereço não da para usar para trocar informação entre maquinas diferentes, e quem fica responsável por manter essa comunicação é o kernel, o arquivo é apenas um name space, nenhum byte vai ser mesmo escrito no arquivo ele vai ocupar zero espaço de disco, toda a comunicação acontece na RAM e gerenciada pelo kernel.

Outra coisa muito importante é que usar Unix domain socket é um recurso padrão de qualquer ambiente [POSIX](https://en.wikipedia.org/wiki/POSIX) mas não esta presente por padrão no Windows.

## Exemplos

### Servidor

A conexão de cliente e servidor é muito parecida a que estamos acostumados quando conectamos via TCP/IP, basicamente ficamos ouvindo o `namespace` indicado pelo arquivo e esperamos por conexões como ficaríamos ouvindo uma porta TCP. Daí quando essa conexão chega usamos a função `Accept` do pacote net e passamos para uma `goroutine` com o handler dessa conexão. 

```go
package main

import (
	"fmt"
	"io"
	"net"
)

func main() {
	l, err := net.Listen("unix", "/tmp/echo.sock")
	if err != nil {
		panic(err)
	}

	for {
		f, err := l.Accept()
		if err != nil {
			panic(err)
		}

		go func(c io.ReadWriter) {
			for {
				buf := make([]byte, 512)
				n, err := c.Read(buf)
				if err != nil {
					return
				}

				fmt.Printf("echo: %s\n", buf[:n])
				_, err = c.Write(buf[:n])
				if err != nil {
					panic(err)
				}
			}
		}(f)
	}
}
```

### Cliente

No nosso exemplo de cliente criamos duas linhas de processamento, uma é a goroutine que vai ficar reatando tudo que vier pela conexão e a outra é a própria função main que vai ficar em loop enviando dados para o servidor. Não precisa ser assim, dependendo de como você quer que seu sistema funcione você pode por exemplo enviar uma mensagem para o servidor e subir uma goroutine para tratar apenas do timeout. O único cuidado é que essas funções bloqueiam o processamento.

```go
package main

import (
	"fmt"
	"io"
	"net"
	"time"
)

func main() {
	f, err := net.Dial("unix", "/tmp/echo.sock")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	go func(r io.Reader) {
		buf := make([]byte, 1024)
		for {
			n, errf := r.Read(buf)
			if errf != nil {
				panic(err)
			}
			fmt.Printf("recebido: %s\n", buf[:n])
		}
	}(f)

	for {
		data := []byte("olá mundo")
		fmt.Printf("enviando: %s\n", data)
		_, err = f.Write([]byte("olá mundo"))
		if err != nil {
			panic(err)
		}

		time.Sleep(time.Duration(400) * time.Millisecond)
	}
}
```

## Servidor usando netcat

Para testes podemos também usar o [netcat](https://en.wikipedia.org/wiki/Netcat) 

```console
nc -lU /tmp/echo.sock && rm /tmp/echo.sock
```

## Cliente usando netcat

```console
nc -U /tmp/echo.sock
```

## Exemplos

### Executando o servidor echo

```console
go run echo/server/main.go
```

Em outro console execute

```console
go run echo/client/main.go
```

---

Implementação do gonf usando npipe no Windows e Unix Domain Socket
https://github.com/gofn/gofn/blob/master/provision/docker_windows.go 
https://github.com/gofn/gofn/blob/master/provision/docker_unix.go

- [Código fonte de hoje](https://github.com/go-br/estudos/tree/master/unixDomainSocket)
- [Repositório do nosso grupo](https://github.com/go-br/estudos)
- [E você encontra mais exemplos aqui](https://github.com/go-br)
- [Pagina do grupo de estudos](https://gopher.pro.br)

Nossos encontros ocorrem todas as quintas-feiras ás 22h00, para participar [entre no canal de Go no Telegram](https://t.me/joinchat/CS0GhBfKbyqZkpl31RRxJQ)