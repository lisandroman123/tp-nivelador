import socket
import logger
from .server_protocol import ServerProtocol
import lottery.lottery as lottery_module
import threading


MAX_AMOUNT_OF_BYTES = 16777215

class Server:
    protocol: ServerProtocol
    lottery: lottery_module.Lottery
    threads: list        

    def __init__(self, server_host: str, server_port: int, minimum_clients_needed:int) -> None:
        self.server_host = server_host
        self.server_port = server_port
        self.protocol = ServerProtocol()
        self.lottery = lottery_module.Lottery("./bets.txt")
        self._shutdown = False
        self.threads = []
        self.threads_lock = threading.Lock()
        self.condition = threading.Condition()
        self.finished_clients = 0
        self.minimum_clients_needed=minimum_clients_needed
                        
    def _handle_client(self, client_socket):
        action = "handle-client"        
        logger.info(action, logger.LogResult.in_progress)            
        try:
            self._process(client_socket)                                                               
        finally:
            client_socket.close()

    def _process(self,client_socket):        
        agency_id = -1
        while True:            
            msg_type, bets = self.protocol.reciveMessageFromClient(client_socket)
            if msg_type == 2:
                break
            else:
                if agency_id==-1:
                    agency_id=bets[0].agency_id
                    print(f"agency_id: {agency_id}")
                self.lottery.store_bets(bets)   #if not error                      
                self.protocol.sendACKToClient(client_socket)
                #if error send other msg conection closed
        with self.condition:
            self.finished_clients+=1
            if self._shutdown:
                return
            if self.finished_clients >= self.minimum_clients_needed:
                self.winners = self.calculate_winners()
                self.condition.notify_all()
            else:
                self.condition.wait()
                if self._shutdown:
                    return
        winners_batch = ""
        for winner in self.winners:
            if winner.agency_id == agency_id:
                msg = self.protocol.serialize(winner)
                candidate = winners_batch + msg + "\n"
                if len(candidate.encode("utf-8")) > MAX_AMOUNT_OF_BYTES: 
                    if len(winners_batch) > 0:               
                        self.protocol.sendMessageToClient(client_socket,winners_batch)
                        winners_batch = ""
                    winners_batch += msg + "\n"
                else:
                    winners_batch = candidate
        if len(winners_batch) > 0:
            self.protocol.sendMessageToClient(client_socket, winners_batch)
        self.protocol.sendEndOfStream(client_socket)
        
    def calculate_winners(self):        
        winners = []
        for bet in self.lottery.load_bets():
            if self.lottery.has_won(bet):
                winners.append(bet)
        return winners
    
    def run(self):                
        self.server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        
        self.server_socket.bind((self.server_host, self.server_port))
        self.server_socket.listen()    
        self.acceptor = threading.Thread(
            target=self._accept_connections,
            args=(self.server_socket,)
        )
        self.acceptor.start()
        print("Esperando acceptor...", flush=True)
        self.acceptor.join()
        print("Acceptor terminó", flush=True)

        for thread in self.threads:
            print("Esperando thread cliente...", flush=True)
            thread.join()

        print("Todos los threads terminaron", flush=True)   
        
    def _accept_connections(self, server_socket):        
        while not self._shutdown:
            try:
                client_socket, _ = server_socket.accept()

                thread = threading.Thread(
                    target=self._handle_client,
                    args=(client_socket,)
                )

                self.threads.append(thread)
                thread.start()

            except OSError:
                if self._shutdown:
                    break
                raise
        
    def shutdown(self):
        self._shutdown = True
        print("se envia close al server socket", flush=True)  
        try:
            self.server_socket.shutdown(socket.SHUT_RDWR)
        except OSError:
            pass
        self.server_socket.close()        
        print("server socket close termino", flush=True)  
        with self.condition:
            self.condition.notify_all() 