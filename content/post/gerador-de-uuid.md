+++
title = "Gerador de UUID"
description = "Neste episódio do grupo de estudos de Go mostramos nosso gerador de UUID e também o inicio do desenvolvimento da nossa ferramenta de migration."
tags = ["Golang"]
date = "2018-11-23T06:45:14Z"
+++

{{< youtube AJSjBiuv65c >}}

# Gerador de UUID

Mais um utilitário no nosso repositório do [GoSidekick](https://github.com/gosidekick), um utilitário para gerar [UUID v4](https://en.wikipedia.org/wiki/Universally_unique_identifier) (aquele que é totalmente aleatório).

```go
package main

import (
	"fmt"
	"strings"

	"github.com/crgimenes/goconfig"
	"github.com/google/uuid"
)

func main() {
	type configFlags struct {
		CarriageReturn bool   `json:"cr" cfg:"cr" cfgDefault:"false"`
		LineFeed       bool   `json:"lf" cfg:"lf" cfgDefault:"false"`
		Uppercase      bool   `json:"u" cfg:"u" cfgDefault:"false"`
		Armored        bool   `json:"a" cfg:"a" cfgDefault:"false"`
		ArmorChar      string `json:"ac" cfg:"ac" cfgDefault:"\""`
		NtoGenerate    int    `json:"n" cfg:"n" cfgDefault:"1"`
		IDSeparator    string `json:"ids" cfg:"ids" cfgDefault:""`
	}

	cfg := configFlags{}
	err := goconfig.Parse(&cfg)
	if err != nil {
		fmt.Println(err)
		return
	}

	cr := ""
	lf := ""
	ac := ""
	if cfg.CarriageReturn {
		cr = "\r"
	}
	if cfg.LineFeed {
		lf = "\n"
	}
	if cfg.Armored {
		ac = cfg.ArmorChar
	}
	for n := 0; n < cfg.NtoGenerate; n++ {
		id := uuid.New().String()
		if cfg.Uppercase {
			id = strings.ToUpper(id)
		}
		fmt.Printf("%v%v%v%v%v%v", ac, id, ac, cfg.IDSeparator, cr, lf)
	}
}
```

- [Código fonte de hoje](https://github.com/gosidekick/uuid)
- [Repositório do nosso grupo](https://github.com/go-br/estudos)
- [E você encontra mais exemplos aqui](https://github.com/go-br)
- [Pagina do grupo de estudos](https://gopher.pro.br)

Nossos encontros ocorrem todas as quintas-feiras ás 22h00, para participar [entre no canal de Go no slack](https://invite.slack.golangbridge.org/) e procure por #brazil