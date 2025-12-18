from dataclasses import dataclass

@dataclass
class Device:
    hostname: str
    ip_address: str
    username: str
    password: str