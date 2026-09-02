package client_protocol

import (
	"encoding/binary"
	"net"
	"os"
	"time"

	"github.com/7574-sistemas-distribuidos/tp-nivelador/src/logger"
	"github.com/7574-sistemas-distribuidos/tp-nivelador/src/safe_socket"
)

const CONNECTION_ATTEMPTS_MAX = 3
const CONNECTION_ATTEMPS_DELAY_MS = 200
const ECHO_CLIENT_FIRST_RESPONSE_SIZE = 2

type NetworkPacket struct {
	AgencyId   string
	MessageId  int
	PacketSize uint16
}

type ClientConfig struct {
	ServerHost string
	ServerPort string
	AgencyId   string
}

type Client struct {
	conn   net.Conn
	config ClientConfig
}

func NewClient(config ClientConfig) (*Client, error) {
	conn, err := connectToServer(config.ServerHost, config.ServerPort)
	if err != nil {
		logger.Warn("connect-to-server", logger.Fail)
		return nil, err
	}

	client := &Client{conn: conn, config: config}
	return client, nil
}

func connectToServer(host, port string) (net.Conn, error) {
	const action = "connect-to-server"
	var err error
	var conn net.Conn

	logger.Info(action, logger.InProgress)
	for i := range CONNECTION_ATTEMPTS_MAX {
		conn, err = net.Dial("tcp", host+":"+port)
		if err != nil {
			logger.Warn(action, logger.Fail, "attempt", i)
			time.Sleep(CONNECTION_ATTEMPS_DELAY_MS * time.Millisecond)
			continue
		}

		logger.Info(action, logger.Success)
		break
	}

	return conn, err
}

func sendInfoToServer(line string, client *Client) error {
	sz := make([]byte, 2)
	binary.BigEndian.PutUint16(sz, uint16(len(line)))

	if err := safe_socket.SendAll(client.conn, sz); err != nil {
		logger.Error("send-message-sz-error", logger.Fail)
		return err
	}
	logger.Info("send-message-sz", logger.Success)
	if err := safe_socket.SendAll(client.conn, []byte(line)); err != nil {
		logger.Error("send-message-line-error", logger.Fail)
		return err
	}
	logger.Info("send-message-line", logger.Success)

	responseBuffer, err := safe_socket.RecvAll(client.conn, ECHO_CLIENT_FIRST_RESPONSE_SIZE)
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

func waitForFinalResponse(archive *os.File, client *Client) (string, error) {
	if err := safe_socket.SendAll(client.conn, []byte("EOF")); err != nil {
		logger.Error("send-message-line-error", logger.Fail)
		return "", err
	}

	responseBuffer, err := safe_socket.RecvAll(client.conn, ECHO_CLIENT_FIRST_RESPONSE_SIZE)
	if err != nil {
		logger.Error("recv-response", logger.Fail)
		return "", err
	}
	size := int(binary.BigEndian.Uint16(responseBuffer))
	logger.Info("recv-sz", logger.Success)
	payload, err := safe_socket.RecvAll(client.conn, size)

	return string(payload), nil

}
