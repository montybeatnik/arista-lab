from jinja2 import Environment, FileSystemLoader
from domain.configuration import Configuration
from infrastructure.arista_device_connector import AristaDeviceConnector

class ConfigurationService:
    def __init__(self, template_env, device_repository, device_connector):
        self.template_env = template_env
        self.device_repository = device_repository
        self.device_connector = device_connector

    def generate_configurations(self):
        devices = self.device_repository.get_devices()
        configurations = []
        for device in devices:
            loopback_ip = self.device_connector.get_loopback_ip(device)
            isis_net = self.convert_to_isis_net(loopback_ip)
            template = self.template_env.get_template("config.j2")
            config = template.render(device=device, isis_net=isis_net)
            configurations.append(Configuration(device, config))
        return configurations

    def convert_to_isis_net(self, loopback_ip):
        # Implement the logic to convert the loopback IP to an ISIS NET address
        # For example:
        octets = loopback_ip.split(".")
        isis_net = f"49.0001.{octets[0]:02x}{octets[1]:02x}.{octets[2]:02x}{octets[3]:02x}.00"
        return isis_net