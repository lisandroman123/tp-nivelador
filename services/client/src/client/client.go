package client

import (
	"bufio"
	"fmt"
	"os"

	client_protocol "github.com/lisandroman123/tp-nivelador/src/client_protocol"
	"github.com/lisandroman123/tp-nivelador/src/logger"
)

const MAX_AMOUNT_OF_BYTES = 16777215
const MAX_AMOUNT_OF_BYTES_IN_LINE = 1024

type Client struct {
	protocol    *client_protocol.Protocol
	input_file  string
	output_file string
	batch_size  int
}

func readMessagesFromInput(c *Client) error {

	file, err := os.Open(c.input_file)

	if err != nil {
		logger.Error("open archive", logger.Fail)
		return err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	batch_size := c.batch_size * MAX_AMOUNT_OF_BYTES_IN_LINE
	sz_b := 0
	batch := make([]byte, 0, batch_size)
	lines_read := 0
	for scanner.Scan() {
		line := scanner.Bytes()
		sz_l := len(line) + 1

		if lines_read >= c.batch_size ||
			sz_b+sz_l > MAX_AMOUNT_OF_BYTES {

			if sz_b > 0 {
				r, err := c.protocol.SendInfoToServer(batch)
				if err != nil || r != "ACK" {
					return err
				}

				batch = batch[:0]
				sz_b = 0
				lines_read = 0
			}
		}

		batch = append(batch, line...)
		batch = append(batch, '\n')

		sz_b += sz_l
		lines_read++
	}

	if len(batch) > 0 {
		_, err := c.protocol.SendInfoToServer(batch)
		if err != nil {
			return err
		}
	}
	c.protocol.SendEnfOfFile()
	if err := scanner.Err(); err != nil {
		logger.Error("scanning archive", logger.Fail)
		return err
	}
	/**
		ESTA PARTE TIENE QUE ESTAR EN OTRA FUNCION
	**/
	archive, err := os.Create(c.output_file)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	defer archive.Close()
	logger.Info("create-succed", logger.Success)
	for {
		winners, err := c.protocol.ReceiveFinalResponse()
		if err != nil {
			return err
		}
		if len(winners) > 0 {
			err = saveMessagesFromServer(archive, winners)
			if err != nil {
				return err
			}
		} else {
			break
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

func (c *Client) Run() error {
	const mainAction = "test-echo-server"
	readMessagesFromInput(c)
	return nil
}

func New(input_file string, output_file string, batch int, p *client_protocol.Protocol) *Client {
	return &Client{
		protocol:    p,
		input_file:  input_file,
		output_file: output_file,
		batch_size:  batch,
	}

}

func (c *Client) Close() error {
	c.protocol.Close()
	return nil
}
