CREATE TABLE devices (
    id SERIAL PRIMARY KEY,
    ip_address VARCHAR(15) UNIQUE NOT NULL,
    hostname VARCHAR(255),
    loopback_ip VARCHAR(15),
    username VARCHAR(15),
    password VARCHAR(15)
);

ALTER TABLE devices
ADD COLUMN infrastructure_interfaces TEXT[];