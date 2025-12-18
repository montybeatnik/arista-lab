from netmiko import ConnectHandler
import requests
import json

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
            print(f"failed to connect: {e}")

    # def apply_configuration(self, device, config):
    #     # Implement the logic to apply the configuration to the device
    #     device_config = {
    #         "device_type": "arista_eos",
    #         "ip": device.ip_address,
    #         "username": device.username,
    #         "password": device.password,
    #         "use_keys": False,
    #         "allow_agent": False,
    #         "conn_timeout": 10,
    #         "read_timeout": 30 
    #     }
    #     with ConnectHandler(**device_config) as net_connect:
    #         net_connect.send_config_set(config.splitlines())

    def apply_configuration(self, device, config):
        url = f"https://{device.ip_address}/command-api"
        auth = (device.username, device.password)
        headers = {"Content-Type": "application/json"}

        commands = ["enable", "configure"] + config.splitlines() + ["end", "wr mem"]

        print(f"applying cfg to {device.hostname}")
        print(f"{commands=}")
        payload = {
            "jsonrpc": "2.0",
            "method": "runCmds",
            "params": {
                "format": "json",
                "timestamps": False,
                "autoComplete": False,
                "expandAliases": False,
                "cmds": commands,
                "version": 1
            },
            "id": "1"
        }

        response = requests.post(url, auth=auth, headers=headers, data=json.dumps(payload), verify=False)

        if response.status_code == 200:
            print(f"Configuration applied successfully to {device.ip_address}")
            print(f"##### DEBUG #####{response.content=}")
        else:
            print(f"Failed to apply configuration to {device.ip_address}: {response.text}")