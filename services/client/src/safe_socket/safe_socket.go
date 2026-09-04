package safe_socket

import (
	"io"
	"net"
	"time"

	"github.com/lisandroman123/tp-nivelador/src/logger"
)

type ClientConfig struct {
	ServerHost string
	ServerPort string
	AgencyId   string
}

type Connection struct {
	conn   net.Conn
	config ClientConfig
}

const CONNECTION_ATTEMPTS_MAX = 3
const CONNECTION_ATTEMPS_DELAY_MS = 200

func NewClient(config ClientConfig) (*Connection, error) {
	conn, err := connectToServer(config.ServerHost, config.ServerPort)
	if err != nil {
		logger.Warn("connect-to-server", logger.Fail)
		return nil, err
	}

	client := NewConnection(conn, config)
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

func (c *Connection) SendAll(bytes []byte) error {
	sent := 0
	for sent < len(bytes) {
		sz_sent, err := c.conn.Write(bytes[sent:])
		sent += sz_sent
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Connection) RecvAll(size int) ([]byte, error) {
	buff := make([]byte, size)
	read := 0
	for read < size {
		sz_read, err := c.conn.Read(buff[read:])
		if sz_read > 0 {
			read += sz_read
		}

		if err != nil {
			if err == io.EOF && read < size {
				return nil, io.ErrUnexpectedEOF
			}
			return nil, err
		}
	}

	return buff, nil

}

func (c *Connection) Close() error {
	return c.conn.Close()
}

func NewConnection(conn net.Conn, config ClientConfig) *Connection {
	return &Connection{
		conn:   conn,
		config: config,
	}
}
