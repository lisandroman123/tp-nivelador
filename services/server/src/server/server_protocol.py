import safe_socket
import src_frozen.lottery.bet as bet

# @dataclass
# class NetworkPacket:
#     msg_type: int
#     bet: bet.Bet


class ServerProtocol:
    def reciveMessageFromClient(self, client_socket):
        sz_payload=2
        client_message = safe_socket.recv_all(client_socket, sz_payload)
        client_message = int.from_bytes(client_message, byteorder='big')
        print(f"client_message: {client_message}")
        payload = safe_socket.recv_all(client_socket, client_message)                
        print(f"client_message_payload: {payload}")         
        msg_type, msg = self.deserialize(payload)   
        return msg_type,msg            

    def sendMessageToClient(self,client_socket,winners:list[bet.Bet]):
        for winner in winners:
            payload = self.serialize(winner)
            size = len(payload).to_bytes(2, byteorder="big")
            safe_socket.send_all(client_socket, size)
            print(f"sent_message_sz: {len(payload) }")
            safe_socket.send_all(client_socket,payload)
            print(f"sent_message_payload: {payload}")

    def deserialize(payload):
        split = payload.split(",")
        if len(split) > 1: 
            first_name, last_name,document, birthdate, number = split[0],split[1],split[2],split[3],split[4]
            return 1,bet.Bet(agency_id=1, first_name=first_name, last_name=last_name,document=document,birthdate=birthdate,number=number)
        else:
            return 2,"ACK"            
        

    def serialize(bet=bet.Bet):
        return f"{bet.first_name},{bet.last_name},{bet.document},{bet.birthdate},{bet.number}"

#     def mi_serializador_tlv(tipo: int, valor_texto: str) -> bytes:
#     # 1. Convertir el texto a su representación de bytes (Value)
#     value_bytes = valor_texto.encode('utf-8')
#     longitud = len(value_bytes)
    
#     # 2. Validar límites para evitar desbordamientos
#     if tipo < 0 or tipo > 255:
#         raise ValueError("El tipo debe caber en 1 byte (0-255)")
#     if longitud > 65535:
#         raise ValueError("El payload es demasiado grande para 2 bytes (máx 65535)")
        
#     # 3. Serializar manualmente el Tipo (1 byte)
#     t_byte = bytes([tipo])
    
#     # 4. Serializar manualmente la Longitud (2 bytes - Big Endian / Network Byte Order)
#     # Tomamos el número entero y lo dividimos en sus dos bytes componentes
#     l_byte_alto = (longitud >> 8) & 0xFF  # Desplaza 8 bits a la derecha para obtener los 8 bits superiores
#     l_byte_bajo = longitud & 0xFF         # Aplica una máscara para quedarse solo con los 8 bits inferiores
    
#     l_bytes = bytes([l_byte_alto, l_byte_bajo])
    
#     # 5. Concatenar todo en un solo bloque de bytes (T + L + V)
#     paquete_completo = t_byte + l_bytes + value_bytes
    
#     return paquete_completo

# # === PRUEBA DEL SERIALIZADOR MANUAL ===

# # Queremos enviar el tipo 5, y un texto que genera 10 bytes de payload
# tipo_msg = 5
# payload_msg = "Hola Mundo"

# buffer_para_enviar = mi_serializador_tlv(tipo_msg, payload_msg)

# print(f"Paquete serializado a mano: {buffer_para_enviar}")
# print(f"Representación en lista de bytes: {list(buffer_para_enviar)}")


# import socket

# def leer_exactamente(sock: socket.socket, cantidad_bytes: int) -> bytes:
#     """Asegura leer exactamente el número de bytes pedidos, 
#     manejando fragmentación de red."""
#     buffer = bytearray()
#     while len(buffer) < cantidad_bytes:
#         bytes_a_leer = cantidad_bytes - len(buffer)
#         datos = sock.recv(bytes_a_leer)
#         if not datos:
#             raise ConnectionError("El cliente cerró la conexión inesperadamente.")
#         buffer.extend(datos)
#     return bytes(buffer)

# def mi_deserializador_tlv(sock: socket.socket):
#     # FASE 1: Leer la cabecera fija (1 byte de Tipo + 2 bytes de Longitud = 3 bytes)
#     cabecera = leer_exactamente(sock, 3)
    
#     # Extraer el Tipo (está en el byte index 0)
#     tipo = cabecera[0]
    
#     # Extraer la Longitud (bytes en index 1 y 2) usando operaciones de bits inversas
#     # Desplazamos el byte alto 8 bits a la izquierda y le sumamos el byte bajo
#     l_byte_alto = cabecera[1]
#     l_byte_bajo = cabecera[2]
#     longitud = (l_byte_alto << 8) | l_byte_bajo
    
#     # FASE 2: Leer el Valor (exactamente la cantidad de bytes que nos dijo 'longitud')
#     valor_bytes = leer_exactamente(sock, longitud)
    
#     # Reconstruir el texto original
#     valor_texto = valor_bytes.decode('utf-8')
    
#     return tipo, longitud, valor_texto
