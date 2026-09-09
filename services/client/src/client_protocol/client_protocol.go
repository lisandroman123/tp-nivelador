package client_protocol

import (
	"github.com/lisandroman123/tp-nivelador/src/logger"
	safe_socket "github.com/lisandroman123/tp-nivelador/src/safe_socket"
)

const ECHO_CLIENT_FIRST_RESPONSE_SIZE = 2
const HEADER_SIZE = 5
const SERVER_HEADER_SIZE = 4
const STREAM = 0x01
const ENDOFSTREAM = 0x02
const ACK = 0x03
const ACK_M = "ACK"

type Protocol struct {
	connection *safe_socket.Connection
}

func (p *Protocol) SendInfoToServer(line []byte) (string, error) {
	header := make([]byte, HEADER_SIZE)
	header[0] = STREAM
	size := len(line)

	header[1] = byte(size >> 16)
	header[2] = byte(size >> 8)
	header[3] = byte(size)

	header[4] = byte(p.connection.GetAgencyId()[0])

	message := make([]byte, 0, len(header)+len(line))
	message = append(message, header...)
	message = append(message, line...)

	if err := safe_socket.SendAll(p.connection.GetConn(), message); err != nil {
		logger.Error("send-message-line-error", logger.Fail)
		return "", err
	}
	logger.Info("send-message-line", logger.Success)

	responseBuffer, err := safe_socket.RecvAll(p.connection.GetConn(), SERVER_HEADER_SIZE)
	if err != nil {
		return "", err
	}
	messageType := responseBuffer[0]
	msg_response := string(responseBuffer[1:4])

	if messageType != ACK || ACK_M != msg_response {
		logger.Error("recv-response", logger.Fail)
		return "", err
	}

	return "ACK", nil
}

func (p *Protocol) ReceiveFinalResponse() (string, error) {
	r, err := safe_socket.RecvAll(p.connection.GetConn(), SERVER_HEADER_SIZE)
	if err != nil {
		return "", err
	}
	messageType := r[0]
	size := int(r[1])<<16 | int(r[2])<<8 | int(r[3])
	if messageType == ENDOFSTREAM {
		return "", nil
	} else {
		payload, err := safe_socket.RecvAll(p.connection.GetConn(), size)
		if err != nil {
			return "", err
		}
		return string(payload), nil
	}
}

func New(connection *safe_socket.Connection) *Protocol {
	return &Protocol{
		connection: connection,
	}
}

func (p *Protocol) SendEnfOfFile() error {
	header := make([]byte, HEADER_SIZE)

	header[0] = ENDOFSTREAM
	header[1] = 0
	header[2] = 0
	header[3] = 0
	header[1] = byte(p.connection.GetAgencyId()[0])

	return safe_socket.SendAll(p.connection.GetConn(), header)
}

func (p *Protocol) Close() {
	p.connection.Close()
}

// func Serialize(line string) {

// }
