+++
title = "Maps"
date = "2020-10-22T19:46:45-03:00"
description = "Exemplo do uso de maps em Go."
tags = ["golang"]
+++

# Maps
```go
package main

import "fmt"

func main() {
	//declarando um map
	var m = map[string]int{"Arroz": 1, "Feijão": 2}
	fmt.Println(m)

	//declarando map com make
	m1 := make(map[string]int)

	m1["Arroz"] = 1
	m1["Feijão"] = 2

	fmt.Println(m1)

	//verificar tamanho de um map
	fmt.Println(len(m1))

	//deletando elemento de um map
	delete(m1, "Feijão")
	fmt.Println(m1)
}
```
[Playground](https://play.golang.org/p/eF5_SbYO_d-)