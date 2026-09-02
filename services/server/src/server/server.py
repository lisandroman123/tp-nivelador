import socket
import logger
import safe_socket

_ECHO_SERVER_MESSAGE_SIZE = 1024
_ECHO_SERVER_FIRST_MSG_SIZE = 2

class Server:
    def __init__(self, server_host: str, server_port: int) -> None:
        self.server_host = server_host
        self.server_port = server_port

    def _handle_client(self, client_socket):
        action = "handle-client"
        message_amount = 0
        try:
            logger.info(action, logger.LogResult.in_progress)
            while message_amount < 3:                                      
                client_message = safe_socket.recv_all(
                    client_socket, _ECHO_SERVER_FIRST_MSG_SIZE
                )                                         
                if not client_message:
                    logger.info(
                        action,
                        logger.LogResult.success,
                        "messages-amount",
                        message_amount,
                    )
                    return
                message_amount += 1
                safe_socket.send_all(client_socket, client_message)

            self._process(client_socket)                

        except Exception as e:
            logger.error(
                action, logger.LogResult.fail, "messages-amount", message_amount
            )
            raise e

    def _process(self,client_socket):
        while True:
            sz_payload=2
            client_message = safe_socket.recv_all(client_socket, sz_payload)
            client_message = int.from_bytes(client_message, byteorder='big')
            print(f"client_message: {client_message}")
            payload = safe_socket.recv_all(client_socket, client_message)                
            print(f"client_message_payload: {payload}")    

            size = len(payload).to_bytes(2, byteorder="big")
            safe_socket.send_all(client_socket, size)
            print(f"sent_message_sz: {len(payload) }")
            safe_socket.send_all(client_socket,payload)
            print(f"sent_message_payload: {payload}")


    def run(self):
        action = "accept-connection"
        with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as server_socket:
            server_socket.bind((self.server_host, self.server_port))
            server_socket.listen()
            while True:
                try:
                    logger.info(action, logger.LogResult.in_progress)
                    client_socket, _ = server_socket.accept()
                except Exception as e:
                    logger.error(action, logger.LogResult.fail)
                    raise e
                logger.info(action, logger.LogResult.success)

                self._handle_client(client_socket)
