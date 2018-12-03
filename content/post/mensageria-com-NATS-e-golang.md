+++
title = "Um sistema de mensageria extremamente rápido com NATS e Golang"
description = "NATS é um sistema de mensagens escrito em Go, muito rápido e fácil de usar. Usado para conectar sistemas, seja IoT seja aplicações mobile ou onde precisar de grande desempenho na troca de mensagens."
tags = ["Golang"]
date = "2018-10-19T06:45:14Z"
+++

{{< youtube NBi0r7QOJSs >}}

# mensageria com NATS

Para quem precisa que vários sistemas se comuniquem entre si e precisa de velocidade e simplicidade o [NATS](https://nats.io) é uma grande ajuda. Uma das características mais interessantes é a incrível velocidade e um protocolo simples que pode ser facilmente implementado faz com que o NATS seja usado desde para IoT com bibliotecas cliente escritas até mesmo para Arduino quando para comunicação entre serviços muito maiores. Também é muito usado para aplicações mobile por pelo mesmo motivo, protocolo simples e rápido. É uma alternativa muito interessante ao [MQTT](https://mqtt.org) e é implementado em Go. :D

## Exemplo de cliente NATS

```go
package main

import (
	"fmt"

	"github.com/nats-io/nats"
)

func main() {
	// Connect to a server
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Simple Publisher
	err = nc.Publish("teste", []byte("Hello World"))
	if err != nil {
		fmt.Println(err)
		return
	}

	// Simple Async Subscriber
	_, err = nc.Subscribe("teste", func(m *nats.Msg) {
		fmt.Printf("Received a message: %s\n", string(m.Data))
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	halt()
}

func halt() {
	select {}
}
```

## Codigo fonte de hoje

- [Servidor](https://github.com/nats-io/gnatsd)
- [Cliente](https://github.com/nats-io/go-nats)
- [Site do NATS)(https://nats.io)

- [Repositório do nosso grupo](https://github.com/go-br/estudos)
- [E você encontra mais exemplos aqui](https://github.com/go-br)
- [Pagina do grupo de estudos](https://gopher.pro.br)

Nossos encontros ocorrem todas as quintas-feiras ás 22h00, para participar [entre no canal de Go no slack](https://invite.slack.golangbridge.org/) e procure por #brazil