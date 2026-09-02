package main

import (
	"bufio"   // leitura em buffer; NewScanner junta bytes do TCP até um \n
	"fmt"     // formatação de texto (Printf, Sprintf) para logs e respostas
	"net"     // sockets TCP: Listen, Accept e o tipo Conn
	"strconv" // Atoi converte os operandos de texto para int
	"strings" // TrimSpace, Fields e ToUpper para interpretar o protocolo
)

const (
	HOST = "127.0.0.1" // loopback: só a própria máquina alcança o servidor
	PORT = "5555"

	OPERATION_ERROR   = "ERRO: Formato invalido (useSOM: OPERACAO NUM1 NUM2)"
	OPERATION_INVALID = "ERRO: Operacao invalida."
	DIVISION_BY_ZERO  = "ERRO: Divisao por zero."
)

// processCommand interpreta uma linha do protocolo e devolve a resposta.
// O bool indica se a conexão deve ser encerrada (comando SAIR).
func processCommand(message string) (string, bool) {
	// TrimSpace remove \r/\n e espaços; Fields quebra em tokens ignorando espaços extras.
	parts := strings.Fields(strings.TrimSpace(message))

	if len(parts) == 0 {
		return OPERATION_ERROR, false
	}

	// ToUpper faz som/SAIR serem aceitos da mesma forma que SOM/SAIR.
	operation := strings.ToUpper(parts[0])

	if operation == "SAIR" {
		return "Conexao encerrada.", true
	}

	if operation != "SOM" && operation != "SUB" && operation != "MUL" && operation != "DIV" {
		return OPERATION_INVALID, false
	}

	if len(parts) != 3 {
		return OPERATION_ERROR, false
	}

	// Atoi falha se o token não for um inteiro (ex.: "abc" ou "3.14").
	num1, err1 := strconv.Atoi(parts[1])
	num2, err2 := strconv.Atoi(parts[2])

	if err1 != nil || err2 != nil {
		return OPERATION_ERROR, false
	}

	var result int

	switch operation {
	case "SOM":
		result = num1 + num2
	case "SUB":
		result = num1 - num2
	case "MUL":
		result = num1 * num2
	case "DIV":
		if num2 == 0 {
			return DIVISION_BY_ZERO, false
		}
		result = num1 / num2 // divisão inteira
	default:
		return OPERATION_ERROR, false
	}

	return fmt.Sprintf("Resultado: %d", result), false
}

// resolveClient atende um cliente até SAIR, desconexão ou erro de leitura.
func resolveClient(connection net.Conn) {
	// Close libera o descritor do SO ao sair da função (defer executa no return).
	defer connection.Close()
	address := connection.RemoteAddr().String() // IP:porta do cliente, só para log
	fmt.Printf("Conexao estabelecida com %s\n", address)

	// Scanner lê o stream TCP linha a linha. Sem isso, um Read cru poderia
	// devolver o comando pela metade ("SOM 1" + "0 20\n").
	scanner := bufio.NewScanner(connection)

	for scanner.Scan() {
		message := scanner.Text() // linha já sem o \n
		response, close := processCommand(message)

		if close {
			break
		}

		// Write envia bytes crus; o \n é o delimitador combinado com o cliente.
		connection.Write([]byte(response + "\n"))
	}

	// Scan() retorna false no fim normal e em erro de I/O; Err() distingue os dois.
	if err := scanner.Err(); err != nil {
		fmt.Printf("Erro ao ler dados de %s: %v\n", address, err)
	}

	fmt.Printf("Conexao com %s encerrada.\n", address)
}

func main() {
	// Listen abre o socket de escuta. O servidor espera conexões; não inicia o handshake.
	listener, err := net.Listen("tcp", HOST+":"+PORT)
	if err != nil {
		fmt.Printf("Erro ao iniciar servidor: %v\n", err)
		return
	}
	defer listener.Close() // libera a porta 5555 ao encerrar o processo

	fmt.Printf("Servidor escutando em %s:%s\n", HOST, PORT)

	for {
		// Accept bloqueia até um cliente conectar e devolve um Conn dedicado a ele.
		connection, err := listener.Accept()
		if err != nil {
			fmt.Printf("Erro ao aceitar conexao: %v\n", err)
			continue // erro pontual não derruba o servidor
		}

		// go inicia uma goroutine: cada cliente é atendido em paralelo.
		go resolveClient(connection)
	}
}
