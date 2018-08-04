+++
date = "2018-08-02T20:05:18-03:00"
title = "Acessando Dropbox com Golang"
description = ""
tags = ["Golang", "Programação"]
+++

Precisei criar um [pacote Go para ao Dropbox](github.com/crgimenes/dropbox) e fazer algumas operações basicas, listar arquivos, upload e download, para isso usei o [dropbox-sdk-go-unofficial](https://github.com/dropbox/dropbox-sdk-go-unofficial).

{{< youtube 9wR3nSc2GMo >}}

## Conectando Dropbox

Para conectar o dropbox precisamos das credenciais de acesso, a forma mais fácil é criar um token, para isso entre em [Dropbox developers apps](https://www.dropbox.com/developers/apps) e crie uma aplicação, então dentro do painel da aplicação crie o token.

## Configurando o sistema

Todas as funções precisam das configurações com as credenciais de acesso e outros parâmetros úteis como por exemplo o nível de log. Então criamos uma função para retornar uma instancia da struct config.

```go
func NewConfig(token string) (config dropbox.Config) {
	config = dropbox.Config{
		Token:    token,
		LogLevel: dropbox.LogOff, // logging level. Default is off
	}
	return
}
```

### Exemplo

```go
config := NewConfig("token aqui")
```

## Listando arquivos

Para facilitar as coisas criamos uma struct que chamamos de Node, um node pode ser um arquivo ou um diretório, é uma forma mais comum de visualizar sistemas de arquivos do que a originalmente usada pelo pacote.

```go
// Node contains metadata to files and folders
type Node struct {
	IsFolder       bool
	Name           string
	Size           uint64
	Rev            string
	ServerModified time.Time
}
```

Para listar o diretório raiz não envia uma "/" como seria comum, no lugar envie uma string vazia.

```go
func List(config dropbox.Config, path string) (nodes []Node, err error) {
	f := files.New(config)
	lfa := files.NewListFolderArg(path)
	lfr, err := f.ListFolder(lfa)
	if err != nil {
		return
	}
	for _, v := range lfr.Entries {
		var n Node
		switch fm := v.(type) {
		case *files.FileMetadata:
			n = parseFileMetadata(fm)
		case *files.FolderMetadata:
			n = parseFolderMetadata(fm)
		}
		nodes = append(nodes, n)
	}
	return
}
```

### Exemplo

Listando arquivos e diretórios

```go
nodes, err := dropbox.List(config, "")
if err != nil {
	log.Fatal(err)
}
for k, v := range nodes {
	fmt.Printf("%v %v\n", k, v.Name)
}
```

## Fazendo Upload de arquivos

Para fazer upload é necessário indicar o caminho completo de destino incluindo a raiz e o nome do arquivo.

A nossa função já controla a sessão de upload de maneira a lidar com as limitações de envio do em uma única sessão da API do Dropbox e também tenta manter o consumo de memória baixo, aqui nos meus testes o melhor resultado foi enviando os arquivos em partes de 1Mb, mas isso pode variar dependendo das condições de rede. O máximo tamanho máximo das partes dos arquivos enviadas de uma vez que a API permite é 150Mb.

### Exemplo

```go
err := dropbox.Upload(config, "origem", "/destino")
if err != nil {
	log.Fatal(err)
}
```

## Fazendo Download de arquivos

Assim como a função Upload, também precisamos tomar cuidado com o consumo de RAM durante o download, mas não é necessário fazer controle de sessão, não temos as mesmas limitações, então podemos usar a função copy para copiar o stream de dados vindo da API para o arquivo de destino. 

```go
err := dropbox.Download(config, "/origem", "destino")
if err != nil {
	log.Fatal(err)
}
```

O código fonte com exemplos e testes esta em [github.com/crgimenes/dropbox](github.com/crgimenes/dropbox).

---

O grupo de estudos de Go se reune todas as quintas-feiras ás 22h00, para participar [entre no canal de Go no slack https://invite.slack.golangbridge.org/ e procure por #brazil

Links úteis

- [Repositório do nosso grupo](https://github.com/go-br/estudos)
- [E você encontra mais exemplos aqui](https://github.com/go-br)
- [Pagina do grupo de estudos](https://gopher.pro.br)
