package main

import (
	"bufio"   // NewReader + ReadString leem teclado e rede até o \n
	"fmt"     // Printf/Print para o terminal; Fprintln para escrever na conexão
	"net"     // Dial inicia o handshake TCP com o servidor
	"os"      // Stdin é o teclado do terminal
	"strings" // TrimSpace e ToUpper para limpar e comparar comandos
)

const (
	HOST = "127.0.0.1"
	PORT = "5555"
)

func main() {
	// Dial conecta em HOST:PORT. Se o servidor estiver offline, falha aqui.
	connection, err := net.Dial("tcp", HOST+":"+PORT)
	if err != nil {
		fmt.Printf("Não foi possivel conectar ao servidor: %v\n", err)
		return
	}
	defer connection.Close() // fecha o socket ao sair do main (SAIR, erro ou Ctrl+C)

	fmt.Printf("Conexao estabelecida com o servidor em %s:%s\n", HOST, PORT)
	fmt.Printf("Envie seus comandos (SOM, SUB, MUL, DIV, SAIR) ou digite SAIR para encerrar a conexao. ex: SOM 10 20\n")

	// Dois Readers separados: misturar teclado e socket no mesmo buffer seria incorreto.
	readKeyboard := bufio.NewReader(os.Stdin)
	readNetwork := bufio.NewReader(connection)

	for {
		fmt.Printf("> ")
		// ReadString bloqueia até o Enter; o \n marca o fim do comando no teclado.
		input, err := readKeyboard.ReadString('\n')
		if err != nil {
			break
		}

		command := strings.TrimSpace(input) // tira \r\n do Windows e espaços
		if command == "" {
			continue // linha vazia não vai para a rede
		}

		// Fprintln escreve o comando e um \n — o mesmo delimitador que o Scanner do servidor usa.
		fmt.Fprintln(connection, command)

		if strings.ToUpper(command) == "SAIR" {
			fmt.Println("Conexão encerrada pelo usuário.")
			break // o servidor fecha sem responder; não chamar ReadString aqui
		}

		// Espera a linha de resposta. Sem \n do servidor, esta chamada ficaria bloqueada.
		response, err := readNetwork.ReadString('\n')
		if err != nil {
			fmt.Printf("Conexão perdida com o servidor.\n")
			break
		}

		fmt.Print(strings.TrimSpace(response) + "\n")
	}
}
