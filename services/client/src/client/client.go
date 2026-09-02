package client

import (
	"bufio"
	"fmt"
	"os"

	client_protocol "github.com/7574-sistemas-distribuidos/tp-nivelador/src/client/client_protocol"
	"github.com/7574-sistemas-distribuidos/tp-nivelador/src/logger"
)

const AMOUNT_OF_INPUT_ARCHIVES = 1

func readMessagesFromInput(client *client_protocol.Client) error {

	for i := range AMOUNT_OF_INPUT_ARCHIVES {
		n := i
		fmt.Println("*************************************")
		arch_name := fmt.Sprintf("/input/input-%v.csv", n)
		file, err := os.Open(arch_name)

		if err != nil {
			logger.Error("open archive", logger.Fail)
			return err
		}

		defer file.Close()
		id := 0
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			err := client_protocol.sendInfoToServer(line, client)
			if err {
				return err
			}
			fmt.Println(line)
			id += 1
		}
		archive, err := os.Create("./output/output.txt")
		if err != nil {
			fmt.Println("Error:", err)
			return err
		}
		defer archive.Close()
		logger.Info("create-succed", logger.Success)
		winners := client_protocol.waitForFinalResponse(archive, client)
		err = saveMessagesFromServer(archive, winners)
		/*enviar final del archivo y espera a terminar conexion*/
		if err := scanner.Err(); err != nil {
			logger.Error("scanning archive", logger.Fail)
		}
	}

	return nil
}

func saveMessagesFromServer(archive *os.File, msg string) error {
	/*crear un archivo dentro de la carpeta output con toda la info del servidor*/
	_, err := archive.WriteString(msg + "\n")
	if err != nil {
		fmt.Println("Error al escribir:", err)
		return err
	}
	return nil
}

func (client *client_protocol.Client) Run() error {
	const mainAction = "test-echo-server"
	defer client.conn.Close()
	readMessagesFromInput(client)

	return nil
}
