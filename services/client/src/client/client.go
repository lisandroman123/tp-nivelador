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
				_, err := c.protocol.SendInfoToServer(batch)
				if err != nil {
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

	if err := scanner.Err(); err != nil {
		logger.Error("scanning archive", logger.Fail)
		return err
	}

	if err := c.protocol.SendEndOfFile(); err != nil {
		logger.Error("send-end-of-file-error", logger.Fail)
		return err
	}

	return nil
}

func saveMessagesFromServer(archive *os.File, msg string) error {
	_, err := archive.WriteString(msg + "\n")
	if err != nil {
		fmt.Println("Error al escribir:", err)
		return err
	}
	return nil
}

func receiveWinnersFromServer(c *Client) error {
	archive, err := os.Create(c.output_file)
	if err != nil {
		logger.Error("Error opening output_file", logger.Fail)
		return err
	}
	defer archive.Close()
	logger.Info("create-output_file-succed", logger.Success)
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

func (c *Client) Run() error {
	const mainAction = "test-echo-server"
	r := readMessagesFromInput(c)
	if r != nil {
		fmt.Printf("readMessagesFromInput ERROR: %v\n", r)
		return r
	}

	logger.Info("input-finished", logger.Success)

	save := receiveWinnersFromServer(c)
	if save != nil {
		fmt.Printf("receiveWinnersFromServer ERROR: %v\n", save)
		return save
	}

	logger.Info("client-finished-successfully", logger.Success)
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
