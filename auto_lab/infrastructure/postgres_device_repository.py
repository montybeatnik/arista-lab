import psycopg2
from domain.device import Device

class PostgresDeviceRepository:
    def __init__(self, db_connection):
        self.db_connection = db_connection

    def get_devices(self):
        devices = []
        with self.db_connection.cursor() as cursor:
            cursor.execute("SELECT hostname, ip_address, loopback_ip, username, password, infrastructure_interfaces FROM devices")
            for row in cursor.fetchall():
                devices.append(Device(*row))
        return devices