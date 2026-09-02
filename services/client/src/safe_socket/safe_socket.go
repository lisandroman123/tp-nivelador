package safe_socket

import "io"

func SendAll(socket io.Writer, bytes []byte) error {
	sent := 0
	for sent < len(bytes) {
		sz_sent, err := socket.Write(bytes[sent:])
		sent += sz_sent
		if err != nil {
			return err
		}
	}

	return nil
}

func RecvAll(socket io.Reader, size int) ([]byte, error) {
	buff := make([]byte, size)
	read := 0
	for read < size {
		sz_read, err := socket.Read(buff[read:])
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
