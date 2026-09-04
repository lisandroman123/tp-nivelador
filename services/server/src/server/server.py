import socket
import logger
import server_protocol
import src_frozen.lottery.lottery as lottery_module

class Server:
    protocol: server_protocol.ServerProtocol
    lottery: lottery_module.Lottery

    def __init__(self, server_host: str, server_port: int) -> None:
        self.server_host = server_host
        self.server_port = server_port
        self.protocol = server_protocol.ServerProtocol()
        self.lottery = lottery_module.Lottery()

    def _handle_client(self, client_socket):
        action = "handle-client"        
        logger.info(action, logger.LogResult.in_progress)            
        self._process(client_socket)                

    def _process(self,client_socket):
        while True:            
            msg_type, bet = self.protocol.reciveMessageFromClient(client_socket)
            if msg_type == 2 and bet == "ACK":
                break
            else:
                self.lottery.store_bets([bet])         
        winners = self.calculate_winners()
        self.protocol.sendMessageToClient(client_socket,winners)

    def calculate_winners(self):        
        winners = []
        for bet in self.lottery.load_bets():
            if self.lottery.has_won(bet):
                winners.append(bet)
        return winners
    
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
