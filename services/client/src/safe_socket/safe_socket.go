package safe_socket

import (
	"io"
	"net"
	"time"
)

type ClientConfig struct {
	ServerHost string
	ServerPort string
	AgencyId   string
	InputFile  string
	OutputFile string
	BatchSize  int
}

type Connection struct {
	conn   net.Conn
	config ClientConfig
}

const CONNECTION_ATTEMPTS_MAX = 3
const CONNECTION_ATTEMPS_DELAY_MS = 200

func (c *Connection) GetAgencyId() string {
	return c.config.AgencyId
}

func NewClient(config ClientConfig) (*Connection, error) {
	conn, err := connectToServer(config.ServerHost, config.ServerPort)
	if err != nil {

		return nil, err
	}

	client := NewConnection(conn, config)
	return client, nil
}

func connectToServer(host, port string) (net.Conn, error) {
	const action = "connect-to-server"
	var err error
	var conn net.Conn

	for _ = range CONNECTION_ATTEMPTS_MAX {
		conn, err = net.Dial("tcp", host+":"+port)
		if err != nil {

			time.Sleep(CONNECTION_ATTEMPS_DELAY_MS * time.Millisecond)
			continue
		}

		break
	}

	return conn, err
}

func SendAll(w io.Writer, bytes []byte) error {
	sent := 0
	sz := len(bytes)
	for sent < sz {
		sz_sent, err := w.Write(bytes[sent:])
		sent += sz_sent
		if err != nil {
			return err
		}
	}

	return nil
}

func RecvAll(r io.Reader, size int) ([]byte, error) {
	buff := make([]byte, size)
	read := 0
	for read < size {
		sz_read, err := r.Read(buff[read:])
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

func (c *Connection) GetConn() net.Conn {
	return c.conn
}
