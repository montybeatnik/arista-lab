from dataclasses import dataclass

@dataclass
class Device:
    hostname: str
    ip_address: str
    loopback_ip: str
    username: str
    password: str