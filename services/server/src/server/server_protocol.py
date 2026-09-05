import safe_socket
import src_frozen.lottery.bet as bet

# @dataclass
# class NetworkPacket:
#     msg_type: int
#     bet: bet.Bet

HEADER_SIZE = 4
STREAM = 0x01
ENDOFSTREAM = 0x02
ACK = 0x03
ACK_M = b"ACK"

#enumerate msg types to receive: 0x01 stream - 0x02 endofstream 0x03 ack


class ServerProtocol:
    def reciveMessageFromClient(self, client_socket):        
        client_message = safe_socket.recv_all(client_socket, HEADER_SIZE)
        header = int.from_bytes(client_message, byteorder='big')

        msg_type = int.from_bytes(header[0:1],"big")
        payload_sz = int.from_bytes(header[1:4],"big")

        if msg_type == STREAM:
            payload = safe_socket.recv_all(client_socket, payload_sz)            
            print(f"client_message_payload: {payload}")         
            msg_type, msg = self.deserialize(payload)        
        elif msg_type == ENDOFSTREAM:
            return 2, "EOS"                    
        return msg_type,msg            

    def sendMessageToClient(self,client_socket,payload):                            
            size = len(payload).to_bytes(2, byteorder="big")
            safe_socket.send_all(client_socket, size)
            print(f"sent_message_sz: {len(payload) }")
            safe_socket.send_all(client_socket,payload)
            print(f"sent_message_payload: {payload}")

    def deserialize(payload):
        split = payload.split(",")       
        first_name, last_name,document, birthdate, number = split[0],split[1],split[2],split[3],split[4]
        return 1,bet.Bet(agency_id=1, first_name=first_name, last_name=last_name,document=document,birthdate=birthdate,number=number)
                             
    def serialize(bet=bet.Bet):
        return f"{bet.first_name},{bet.last_name},{bet.document},{bet.birthdate},{bet.number}"

    def sendACKToClient(client_socket):
        #esto deberia estar dentro del serializer capaz lo hago una clase
        message_type = bytes([0x03])
        payload = ACK_M
        message = message_type + payload
        safe_socket.send_all(client_socket, message)

         