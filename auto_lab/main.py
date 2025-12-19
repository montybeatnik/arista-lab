import psycopg2
from jinja2 import Environment, FileSystemLoader
from application.configuration_service import ConfigurationService
from infrastructure.postgres_device_repository import PostgresDeviceRepository
from infrastructure.arista_device_connector import AristaDeviceConnector

def main():
    db_connection = psycopg2.connect(
        host="localhost",
        database="device_inventory",
        user="lab",
        password="password"
    )
    device_repository = PostgresDeviceRepository(db_connection)
    template_env = Environment(loader=FileSystemLoader("templates"))
    device_connector = AristaDeviceConnector()
    configuration_service = ConfigurationService(template_env, device_repository, device_connector)

    templates = {
        "isis": "isis.j2",
        "mpls": "mpls.j2",
        "ipv6": "ipv6.j2"
    }

    for tmpl in templates:
        configurations = configuration_service.generate_configurations(templates[tmpl])
        for config in configurations:
            device_connector.apply_configuration(config.device, config.config)


if __name__ == "__main__":
    main()