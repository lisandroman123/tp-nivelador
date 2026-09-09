import safe_socket
import lottery.bet as bet

HEADER_SIZE = 5
STREAM = 0x01
ENDOFSTREAM = 0x02
ACK = 0x03
ACK_M = b"ACK"

class ServerProtocol:
    def reciveMessageFromClient(self, client_socket):        
        client_message = safe_socket.recv_all(client_socket, HEADER_SIZE)        

        msg_type = int.from_bytes(client_message[0:1],"big")
        payload_sz = int.from_bytes(client_message[1:4],"big")
        agency_id = int.from_bytes(client_message[4:5],"big")
        
        if msg_type == STREAM:
            payload = safe_socket.recv_all(client_socket, payload_sz)                                 
            msg_type, msg = self.deserialize(agency_id,payload)            
            return msg_type,msg                 
        elif msg_type == ENDOFSTREAM:
            return 2, []                                    

    def sendMessageToClient(self,client_socket,payload):
        payload = payload.encode("utf-8")
        size = len(payload).to_bytes(3, byteorder="big")
        header = bytes([STREAM]) + size                                   
        message = header + payload
        safe_socket.send_all(client_socket,message)
        print(f"sent_message_payload: {payload}")

    def deserialize(self,agency_id, payload):
        payload = payload.decode("utf-8")
        bets = []
        for line in payload.splitlines():
            split = line.split(",")       
            first_name, last_name,document, birthdate, number = split[0],split[1],split[2],split[3],split[4]
            bets.append(
                bet.Bet(agency_id=agency_id, first_name=first_name, last_name=last_name,document=document,birthdate=birthdate,number=number)
            )
        return 1, bets
                             
    def serialize(self,bet=bet.Bet):
        return f"{bet.first_name},{bet.last_name},{bet.document},{bet.birthdate},{bet.number}"

    def sendACKToClient(self,client_socket):
        #esto deberia estar dentro del serializer capaz lo hago una clase
        message_type = bytes([0x03])
        payload = ACK_M
        message = message_type + payload
        safe_socket.send_all(client_socket, message)

    def sendEndOfStream(self, client_socket):        
        header = bytes([ENDOFSTREAM,0,0,0])                                    
        safe_socket.send_all(client_socket, header)
