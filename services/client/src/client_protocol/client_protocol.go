package client_protocol

import (
	"encoding/binary"
	"os"

	"github.com/lisandroman123/tp-nivelador/src/logger"
	safe_socket "github.com/lisandroman123/tp-nivelador/src/safe_socket"
)

const ECHO_CLIENT_FIRST_RESPONSE_SIZE = 2

type Protocol struct {
	connection *safe_socket.Connection
}

func (p *Protocol) SendInfoToServer(line string) error {
	sz := make([]byte, 2)
	binary.BigEndian.PutUint16(sz, uint16(len(line)))

	if err := p.connection.SendAll(sz); err != nil {
		logger.Error("send-message-sz-error", logger.Fail)
		return err
	}
	logger.Info("send-message-sz", logger.Success)
	if err := p.connection.SendAll([]byte(line)); err != nil {
		logger.Error("send-message-line-error", logger.Fail)
		return err
	}
	logger.Info("send-message-line", logger.Success)

	responseBuffer, err := p.connection.RecvAll(ECHO_CLIENT_FIRST_RESPONSE_SIZE)
	if err != nil || string(responseBuffer) != "ACK" {
		logger.Error("recv-response", logger.Fail)
		return err
	}

	if err != nil {
		logger.Error("error-writing", logger.Fail)
		return err
	}
	return nil
}

func (p *Protocol) WaitForFinalResponse(archive *os.File) (string, error) {
	if err := p.connection.SendAll([]byte("EOF")); err != nil {
		logger.Error("send-message-line-error", logger.Fail)
		return "", err
	}

	responseBuffer, err := p.connection.RecvAll(ECHO_CLIENT_FIRST_RESPONSE_SIZE)
	if err != nil {
		logger.Error("recv-response", logger.Fail)
		return "", err
	}
	size := int(binary.BigEndian.Uint16(responseBuffer))
	logger.Info("recv-sz", logger.Success)
	payload, err := p.connection.RecvAll(size)

	return string(payload), nil

}

func New(connection *safe_socket.Connection) *Protocol {
	return &Protocol{
		connection: connection,
	}
}
