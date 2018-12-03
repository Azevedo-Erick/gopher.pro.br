+++
title = "Usando goto e label em Go"
description = "goto é um recurso injustiçado e hoje vamos mostrar como usar goto no em Golang."
tags = ["Golang"]
date = "2018-10-26T06:45:14Z"
+++

{{< youtube r84yDrU0Bug >}}

# goto e label em Go

A instrução goto tem uma má fama que vem do tempo do BASIC quando era usada indiscriminadamente e acabava tornando o código impossível de ler. Em linguagens modernas entretanto é uma instrução perfeitamente válida e desde que usada com critério pode ajudar a tornar seu código mais limpo.

Alem de goto as instruções break e continue também aceitam labels, isso é muito util para quando por exemplo se quer sair de um *for* aninhado em outro *for* ou especificar para qual dos *fors* aninhados se quer fazer *continue*.

Veja os exemplos de [goto e várias outras coisas que aceitam label no nosso grupo de estudos](https://github.com/go-br/estudos/tree/master/goto)

Go tem algumas regras para usar com goto, continue e break:

- Não se pode saltar para fora de um escopo de função
- Não se pode saltar sobre a declaração de uma variável
- Quando se usa *break* seguido de um label a instrução imediatamente após o label precisa ser um *for*, *switch*, ou *select*.
- Quando se usa *continue* seguido de um label a instrução imediatamente após o label precisa ser um *for*.

- [Código fonte de hoje](https://github.com/go-br/estudos/tree/master/goto)
- [Repositório do nosso grupo](https://github.com/go-br/estudos)
- [E você encontra mais exemplos aqui](https://github.com/go-br)
- [Pagina do grupo de estudos](https://gopher.pro.br)

Nossos encontros ocorrem todas as quintas-feiras ás 22h00, para participar [entre no canal de Go no slack](https://invite.slack.golangbridge.org/) e procure por #brazil