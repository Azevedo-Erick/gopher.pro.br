+++
date = "2018-09-20T08:19:09-03:00"
description = "Primeiros passos usando Protocol Buffers com golang"
tags = ["Golang", "Desenvolvimento", "RPC"]
title = "Protocol Buffers"
topics = []
draft = true

+++

{{< youtube XzSgREcg_sU >}}

Este é o primeiro de uma série de tutoriais rápidos onde pretendo cobrir o uso de [gRPC](https://grpc.io) e vários aspectos como testes, TLS, boas praticas e muito mais.

Para iniciar com gRPC precisamos ir até a base então vamos primeiro falar de [Protocol Buffers](https://developers.google.com/protocol-buffers/).

Protocol Buffers é uma forma simples e agnóstica com relação a linguagem de se definir uma estrutura de dados, como XML só que melhor, mais simples e mais rápido, muito mais rápido.

Também podemos dizer que Protocol Buffers é uma forma de definir como seus dados estão organizados e então você pode usar essa definição para gerar automaticamente código para várias linguagens. Hoje são suportadas C++, Go, JAVA, Python, Ruby, C#, Objective-C, Javascript e PHP é bem possível que a linguagem do seu coração esteja nessa lista, as minhas são logo as duas primeiras :D, mas se você procurar um pouco vai encontrar implementações para outras linguagens não oficialmente suportadas como Lua por exemplo. O suporte a tantas linguagens de programação é possível porque o compilador que le a definição do seu protocolo usa um sistema de plugins para definir o código de saída. [Aqui tem uma lista de plugins de terceiros](https://github.com/protocolbuffers/protobuf/blob/master/docs/third_party.md)

Alem do nosso material tem muita coisa boa na internet explicando muito bem o funcionamento do Protocol Buffers, não deixe de por exemplo ver os vídeos do [Francesc Campoy em JustForFunc](https://www.youtube.com/watch?v=_jQ3i_fyqGA)

Agora vamos ver um exemplo simples de como salvar e recuperar structs em um arquivo. Esse exemplo é derivado do exemplo do Francesc, a principal diferença esta na hora de carregar os dados do arquivo, eu prefiro evitar carregar tudo para RAM e depois fazer o parse dos dados. No lugar disso é melhor ler o arquivo e ir parseando os dados.

Para poder usar o compilador protoc com Go não deixe de instalar o plugin como no exemplo

```console
go get -u github.com/golang/protobuf/protoc-gen-go
````

Não ter o plugin da linguagem que se pretende gerar o código é a falha mais comum quando se usa protocol buffers.


Agora vamos criar um exemplo bem simples de arquivo .proto o user.proto

```proto
syntax = "proto3";

package user;

message User {
    int64 ID = 1;
    string email = 2;
    string name = 3;
}
```

Os arquivos .proto são usados para definir as estruturas/mensagens serializadas pelo protocol buffer, com esse arquivo o compilador protoc pode gerar código para várias linguagens e esse é o grande truque, é rápido porque o código é gerado para tratar dados binários de uma estrutura especifica facilitando muito o trabalho do parser, outros formatos como JSON existem muito mais carga de processamento.

No arquivo de exemplo um detalhe importante é o ID que vem depois do sinal de igualdade, como os dados serializadas serão binários esse ID sera usado para distinguir os campos, você pode adicionar novos campos na ordem que quiser, basta ter um ID diferente.

## Gerando código

Agora que temos o arquivo definindo o formato que os dados serão serializadas podemos usar o protoc para gerar um package contendo nossa struct já em código Go.

```console
protoc --go_out=. user.proto
```

Esse é o comando para gerar o código manualmente mas eu prefiro chamar o protoc via go generate, para isso coloquei dentro do a seguinte linha

```go
//go:generate protoc --go_out=. ./user/user.proto
```

Assim podemos gerar essa dependência e qualquer outra chamando `go generate`

```console
go generate
```

O arquivo `user.pb.go` sera gerado contendo o package `user` e todo o necessário para serializarmos e deserializamos ela.

## Gravando dados serializadas

Depois de gerar o package agora vamos finalmente estudar nosso código de exemplo.

Temos duas funções, uma para adicionar e outra para listar usuários.

A função `add` primeiro cria uma instancia da struct user e popula os campos

```go 
u := &user.User{
	ID:    id,
	Name:  name,
	Email: email,
}
```

Em seguida serializamos essa struct para binário

```go
b, err := proto.Marshal(u)
if err != nil {
	return fmt.Errorf("could not encode task: %v", err)
}
```

Nesse ponto já temos a struct serializada, ou seja ela pode ser gravada ou transmitida e qualquer linguagem que use o mesmo arquivo .proto para gerar o código conseguiria deserializar os dados.

No nosso caso vamos gravar em um arquivo, então vamos abrir um arquivo no modo append, ou se não existir criar um arquivo vazio.

```go
f, err := os.OpenFile(dbPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
if err != nil {
	return fmt.Errorf("could not open %s: %v", dbPath, err)
}
```

Agora antes de gravar a struct no arquivo vamos gravar um inteiro contendo o tamanho da struct, isso é necessário porque o protocol buffers não contem nenhum metadados dizendo o tamanho da struct nem nenhuma informação alem do essencial, esse é parte do motivo por ser tão rápido. 

Então a parte a seguir não tem nada de protocol buffers mas é interessante por si só. Vamos gravar o tamanho no arquivo usando `binary.Write` e isso depende do formato do inteiro na memória o que requer lidar com `endianness`, essa é uma troca que temos que fazer quando queremos velocidade, precisamos descer até a estrutura da plataforma. 

```go
// add record length to file
if err = binary.Write(f, endianness, length(len(b))); err != nil {
	return fmt.Errorf("could not encode length of message: %v", err)
}
```

Finalmente vamos gravar a estrutura

```go
// add rocord to file
_, err = f.Write(b)
if err != nil {
	return fmt.Errorf("could not write task to file: %v", err)
}
```

E terminamos fechando o arquivo

```go
err = f.Close()
if err != nil {
	return fmt.Errorf("could not close file %s: %v", dbPath, err)
}
```

## Lendo dados serializados

Agora os dados estão salvos no disco usando um formato bem simples, o tamanho dos dados seguido payload da estrutura, seguindo pelo tamanho do proximo registro e em seguida pelos seus dados e assim por diante até o fim do arquivo.

A primeira coisa que vamos fazer é abrir o arquivo para leitura

```go
f, err := os.Open(dbPath)
if err != nil {
	return fmt.Errorf("could not open file %s: %v", dbPath, err)
}
defer func() {
	e := f.Close()
	if e != nil {
		fmt.Println(e)
	}
}()
```

Em seguida entramos em um loop em que vamos ler o arquivo e só vamos sair dele quando terminarmos de ler o arquivo `EOF`

Primeiro lemos o inteiro contendo o tamanho do proximo registro, caso `binary.Read` retorne erro ou EOF saímos do loop.

```go
// load record file
var l length
err = binary.Read(f, endianness, &l)
if err != nil {
	if err == io.EOF {
		err = nil
		return
	}
	return fmt.Errorf("could not read file %s: %v", dbPath, err)
}
```

Agora que sabemos o tamanho da struct que esta serializada no arquivo podemos usar `io.ReadFull` para ler exatamente essa quantidade de bytes e para isso criamos um buffer.

```go
// load record
bs := make([]byte, l)
_, err = io.ReadFull(f, bs)
if err != nil {
	return fmt.Errorf("could not read file %s: %v", dbPath, err)
}
```

Nosso buffer agora contem os dados serializados e vamos usar esses dados junto com `proto.Unmarshal` para preencher uma nova instancia de user. 

```go
// Unmarshal
var u user.User
err = proto.Unmarshal(bs, &u)
if err != nil {
	return fmt.Errorf("could not read user: %v", err)
}
```

E por fim exibimos os dados na tela, usamos os getters que o protoc gerou para nos mas eles não são necessários já que Go não permite strings NULL, mas se fosse um ponteiro para a instancia de user esses getters evitariam erro retornando strings vazias.

```go
// Print
fmt.Println("id:", u.GetID())
fmt.Println("name:", u.GetName())
fmt.Println("e-mail:", u.GetEmail())
fmt.Println("------------------")
```

Aos poucos vamos intercalar tópicos mais avançados como esse com tópicos mais simples e principalmente e práticos.

## Links úteis

- [Código fonte de hoje](https://github.com/go-br/estudos/tree/master/protobuf)
- [Protocol Buffers](https://developers.google.com/protocol-buffers/)
- [Repositório do nosso grupo](https://github.com/go-br/estudos)
- [E você encontra mais exemplos aqui](https://github.com/go-br)
- [Pagina do grupo de estudos](https://gopher.pro.br)
