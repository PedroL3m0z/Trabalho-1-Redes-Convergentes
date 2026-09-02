# Calculadora TCP

Cliente e servidor TCP em Go que executam operações aritméticas por linha de texto. O servidor escuta em `127.0.0.1:5555`, aceita vários clientes ao mesmo tempo e responde a cada comando com o resultado ou uma mensagem de erro.

Trabalho 1 — Redes Convergentes.

## Requisitos

- [Go](https://go.dev/dl/) 1.22 ou superior

## Estrutura

```
trabalho-redes/
├── go.mod
├── src/
│   ├── server/servidor.go   # escuta TCP e processa os comandos
│   └── client/cliente.go    # conecta, lê o teclado e exibe as respostas
└── README.md
```

## Como executar

Abra dois terminais na raiz do projeto. O servidor precisa estar no ar antes do cliente.

```bash
# terminal 1
go run ./src/server
```

```bash
# terminal 2
go run ./src/client
```

Saída esperada no servidor:

```text
Servidor escutando em 127.0.0.1:5555
```

Saída esperada no cliente:

```text
Conexao estabelecida com o servidor em 127.0.0.1:5555
Envie seus comandos (SOM, SUB, MUL, DIV, SAIR) ou digite SAIR para encerrar a conexao. ex: SOM 10 20
>
```

Para encerrar a sessão, digite `SAIR`. O servidor continua escutando outros clientes até ser interrompido com `Ctrl+C`.

## Protocolo

A comunicação usa TCP, texto UTF-8 e uma mensagem por linha (`\n` como delimitador).

Formato das operações:

```text
OPERACAO NUM1 NUM2
```

| Comando | Significado | Exemplo | Resposta |
|---------|-------------|---------|----------|
| `SOM` | soma | `SOM 10 20` | `Resultado: 30` |
| `SUB` | subtração | `SUB 10 3` | `Resultado: 7` |
| `MUL` | multiplicação | `MUL 4 5` | `Resultado: 20` |
| `DIV` | divisão inteira | `DIV 7 2` | `Resultado: 3` |
| `SAIR` | encerra a conexão | `SAIR` | (fecha o socket) |

Os comandos não diferenciam maiúsculas de minúsculas (`som 10 20` é válido). `NUM1` e `NUM2` precisam ser inteiros.

Erros devolvidos pelo servidor:

| Mensagem | Quando ocorre |
|----------|----------------|
| `ERRO: Operacao invalida.` | operação diferente de `SOM`, `SUB`, `MUL`, `DIV` ou `SAIR` |
| `ERRO: Formato invalido (useSOM: OPERACAO NUM1 NUM2)` | quantidade de argumentos errada ou operandos que não são inteiros |
| `ERRO: Divisao por zero.` | `DIV` com segundo operando igual a `0` |

`SAIR` não espera resposta: o cliente envia o comando e fecha a conexão; o servidor sai do loop daquele cliente e também fecha o socket.

## Como o código funciona

### Servidor

1. `net.Listen("tcp", "127.0.0.1:5555")` abre o socket de escuta.
2. O loop de `Accept()` bloqueia até um cliente conectar.
3. Cada conexão é atendida em uma goroutine (`go resolveClient`), então um cliente não impede os demais.
4. `bufio.Scanner` lê o stream TCP linha a linha. Sem um delimitador, um `Read` cru poderia entregar o comando pela metade.
5. `processCommand` interpreta a linha, valida e calcula.
6. A resposta é escrita de volta com `\n`. Se o comando for `SAIR`, a goroutine encerra e `defer connection.Close()` libera o socket.

### Cliente

1. `net.Dial("tcp", "127.0.0.1:5555")` inicia o handshake TCP com o servidor.
2. Há dois leitores: um para o teclado (`os.Stdin`) e outro para a rede (`connection`).
3. O loop é síncrono: lê o comando, envia com `fmt.Fprintln` (já inclui `\n`), espera a resposta e imprime.
4. Linhas vazias não são enviadas. Depois de `SAIR`, o cliente não tenta ler resposta — o servidor fecha sem escrever nada.

Host e porta estão definidos nas constantes `HOST` e `PORT` dos dois arquivos. Para escutar em outra máquina da rede, altere o servidor para `0.0.0.0` (ou o IP da interface) e o cliente para o IP do servidor.
