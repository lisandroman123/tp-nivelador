package client_protocol

import (
	"encoding/binary"
	"os"

	"github.com/lisandroman123/tp-nivelador/src/logger"
	safe_socket "github.com/lisandroman123/tp-nivelador/src/safe_socket"
)

const ECHO_CLIENT_FIRST_RESPONSE_SIZE = 2
const HEADER_SIZE = 4
const stream = 0x01
const ENDOFSTREAM = 0x02
const ACK = 0x03
const ACK_M = "ACK"

type Protocol struct {
	connection *safe_socket.Connection
}

func (p *Protocol) SendInfoToServer(line string) (string, error) {
	header := make([]byte, 4)
	binary.BigEndian.PutUint16(header[0:1], stream)
	binary.BigEndian.PutUint16(header[1:4], uint16(len(line)))

	if err := p.connection.SendAll(header); err != nil {
		logger.Error("send-message-sz-error", logger.Fail)
		return "", err
	}
	logger.Info("send-message-sz", logger.Success)
	if err := p.connection.SendAll([]byte(line)); err != nil {
		logger.Error("send-message-line-error", logger.Fail)
		return "", err
	}
	logger.Info("send-message-line", logger.Success)

	responseBuffer, err := p.connection.RecvAll(HEADER_SIZE)
	messageType := responseBuffer[0]
	msg_response := string(responseBuffer[1:4])
	//deserialize(responseBuffer)

	if err != nil || messageType != ACK || ACK_M != msg_response {
		logger.Error("recv-response", logger.Fail)
		return "", err
	}

	return "ACK", nil
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

func (p *Protocol) SendEnfOfFile() error {
	header := make([]byte, HEADER_SIZE)

	header[0] = ENDOFSTREAM
	copy(header[1:], "EOS")

	return p.connection.SendAll(header)
}

// func Serialize(line string) {

// }
