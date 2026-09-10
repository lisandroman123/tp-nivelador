import csv
from collections.abc import Iterator
from .bet import Bet
import threading

_LOTTERY_WINNER_NUMBER = 7574

class ReadWriteLock:
    def __init__(self):
        self.condition = threading.Condition()
        self.readers = 0
        self.writer = False

    def acquire_read(self):
        with self.condition:
            while self.writer:
                self.condition.wait()

            self.readers += 1

    def release_read(self):
        with self.condition:
            self.readers -= 1

            if self.readers == 0:
                self.condition.notify_all()

    def acquire_write(self):
        with self.condition:
            while self.writer or self.readers > 0:
                self.condition.wait()

            self.writer = True

    def release_write(self):
        with self.condition:
            self.writer = False
            self.condition.notify_all()



class Lottery:
    def __init__(self, storage_path) -> None:
        self.storage_path = storage_path
        self.file_lock = ReadWriteLock()

    def has_won(self, bet: Bet) -> bool:
        return bet.number == _LOTTERY_WINNER_NUMBER

    def store_bets(self, bets: list[Bet]) -> None:
        self.file_lock.acquire_write()
        try:
            with open(self.storage_path, "a+") as file:
                writer = csv.writer(file, quoting=csv.QUOTE_MINIMAL)
                for bet in bets:
                    writer.writerow(
                        [
                            bet.agency_id,
                            bet.first_name,
                            bet.last_name,
                            bet.document,
                            bet.birthdate,
                            bet.number,
                        ]
                    )
        finally:
            self.file_lock.release_write()

    def load_bets(self) -> Iterator[Bet]: 
        self.file_lock.acquire_read()
        try:       
            with open(self.storage_path, "r") as file:
                reader = csv.reader(file, quoting=csv.QUOTE_MINIMAL)
                for row in reader:
                    [agency_id, first_name, last_name, document, birthdate, number] = row
                    yield Bet(
                        int(agency_id),
                        first_name,
                        last_name,
                        int(document),
                        birthdate,
                        int(number),
                    )
        finally:
            self.file_lock.release_read()