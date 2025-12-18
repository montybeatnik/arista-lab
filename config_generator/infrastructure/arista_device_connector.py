from netmiko import ConnectHandler

class AristaDeviceConnector:
    def __init__(self):
        pass

    def get_loopback_ip(self, device):
        # Implement the logic to retrieve the loopback IP address using eAPI or Netmiko
        # For example using Netmiko:
        device_config = {
            "device_type": "arista_eos",
            "ip": device.ip_address,
            "username": device.username,
            "password": device.password,
            "use_keys": False,
            "allow_agent": False
        }
        try:
            with ConnectHandler(**device_config) as net_connect:
                output = net_connect.send_command("show ip int lo0")
                # Parse the output to extract the IP address
                for line in output.splitlines():
                    if "IP Address" in line:
                        loopback_ip = line.split(":")[1].strip().split("/")[0]
                        return loopback_ip
        except Exception as e:
            print(device)
            print(f"failed to connect: {e}")

    def apply_configuration(self, device, config):
        # Implement the logic to apply the configuration to the device
        device_config = {
            "device_type": "arista_eos",
            "ip": device.ip_address,
            "username": device.username,
            "password": device.password,
            "use_keys": False,
            "allow_agent": False            
        }
        with ConnectHandler(**device_config) as net_connect:
            net_connect.send_config_set(config.splitlines())