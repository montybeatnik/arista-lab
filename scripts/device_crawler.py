import subprocess
import re
import psycopg2
from netmiko import ConnectHandler

def get_containerlab_ips(topology_file):
    try:
        output = subprocess.check_output(["sudo", "containerlab", "inspect", "-t", topology_file]).decode("utf-8")
        ips = []
        for line in output.splitlines():
            if "172." in line:
                parts = line.split()
                for part in parts:
                    if part.startswith("172."):
                        ips.append(part)
                        break
        return ips
    except subprocess.CalledProcessError as e:
        print(f"Failed to run containerlab inspect: {str(e)}")
        return []

def connect_to_device(ip, username, password):
    device_config = {
        "device_type": "arista_eos",
        "ip": ip,
        "username": username,
        "password": password,
    }
    try:
        net_connect = ConnectHandler(**device_config)
        return net_connect
    except Exception as e:
        print(f"Failed to connect to {ip}: {str(e)}")
        return None

def get_device_info(net_connect):
    try:
        output = net_connect.send_command("show hostname")
        hostname = output.strip()
        output = net_connect.send_command("show ip int lo0")
        loopback_ip = None
        for line in output.splitlines():
            if "IP Address" in line:
                loopback_ip = line.split(":")[1].strip().split("/")[0]
                break
        print(f"get_device_info(): {hostname=}, {loopback_ip=}")
        return hostname, loopback_ip
    except Exception as e:
        print(f"Failed to retrieve device info: {str(e)}")
        return None, None

def populate_db(db_connection, ip, hostname, loopback_ip):
    try:
        with db_connection.cursor() as cursor:
            cursor.execute("INSERT INTO devices (ip_address, hostname, loopback_ip) VALUES (%s, %s, %s) ON CONFLICT (ip_address) DO UPDATE SET hostname = %s, loopback_ip = %s", (ip, hostname, loopback_ip, hostname, loopback_ip))
        db_connection.commit()
        print(f"Populated DB with device {ip}")
    except Exception as e:
        print(f"Failed to populate DB with device {ip}: {str(e)}")
        db_connection.rollback()

def main():
    db_connection = psycopg2.connect(
        host="localhost",
        database="device_inventory",
        user="username",
        password="password"
    )
    ips = get_containerlab_ips("lab.clab.yml")
    for ip in ips:
        net_connect = connect_to_device(ip, "admin", "admin") #TODO: lab but these should be somewhere else
        if net_connect:
            hostname, loopback_ip = get_device_info(net_connect)
            if hostname and loopback_ip:
                populate_db(db_connection, ip, hostname, loopback_ip)
            net_connect.disconnect()
    db_connection.close()

if __name__ == "__main__":
    main()