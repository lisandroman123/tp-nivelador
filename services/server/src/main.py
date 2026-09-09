import os
import sys

import logger
import server
import signal

SERVER_HOST = os.environ["SERVER_HOST"]
SERVER_PORT = int(os.environ["SERVER_PORT"])
MINIMUM_CLIENTS_NEEDED = int(os.environ["AGENCY_QUORUM_MIN"])


        

def main():
    logger.init()
    s = server.Server(SERVER_HOST, SERVER_PORT, MINIMUM_CLIENTS_NEEDED)
    def handle_sigterm(signum, frame):
        print(f"SEÑAL RECIBIDA: {signum}", flush=True)
        s.shutdown()
    signal.signal(signal.SIGTERM, handle_sigterm)
    try:        
        s.run()
    except Exception as e:
        logger.error("server-run", logger.LogResult.fail, "err", e)
        return 1
    return 0


if __name__ == "__main__":
    sys.exit(main())
