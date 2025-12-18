from dataclasses import dataclass

@dataclass
class Device:
    ip_address: str
    username: str
    password: str