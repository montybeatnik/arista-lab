from dataclasses import dataclass
from domain.device import Device

@dataclass
class Configuration:
    device: Device
    config: str