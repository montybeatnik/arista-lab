from dataclasses import dataclass
from typing import List

@dataclass
class Device:
    hostname: str
    ip_address: str
    loopback_ip: str
    username: str
    password: str
    infrastructure_interfaces: List[str]