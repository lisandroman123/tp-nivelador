import socket

# TODO: Complete with a short-read/short-write tolerant implementation


def recv_all(socket: socket.socket, size):
    buf = b""
    read = 0
    while read < size:
        bytes = socket.recv(size - read)
        sz = len(bytes)
        if sz > 0:
            read += sz
            buf += bytes
        if sz == 0:            
            return b""
    return buf

def send_all(socket: socket.socket, bytes):
    sent = 0
    sz_sent = 0 
    while (sent < len(bytes)):
        sz_sent = socket.send(bytes[sz_sent:])
        if sz_sent == 0:
            return "error"
        sent+=sz_sent
    return sent
        
