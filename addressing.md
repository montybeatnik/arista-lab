# Addressing plan (quick reference)

## Loopback0 (overlay):

- spine1: 10.0.0.1/32
- spine2: 10.0.0.2/32
- leaf1: 10.0.0.11/32
- leaf2: 10.0.0.12/32
- leaf3: 10.0.0.13/32
- leaf4: 10.0.0.14/32
- Underlay /31s (per link):

| Link                     | Spine IP      | Leaf IP       |
| ------------------------ | ------------- | ------------- |
| spine1:Eth1 – leaf1:Eth1 | 172.16.1.0/31 | 172.16.1.1/31 |
| spine2:Eth1 – leaf1:Eth2 | 172.16.2.0/31 | 172.16.2.1/31 |
| spine1:Eth2 – leaf2:Eth1 | 172.16.1.2/31 | 172.16.1.3/31 |
| spine2:Eth2 – leaf2:Eth2 | 172.16.2.2/31 | 172.16.2.3/31 |
| spine1:Eth3 – leaf3:Eth1 | 172.16.1.4/31 | 172.16.1.5/31 |
| spine2:Eth3 – leaf3:Eth2 | 172.16.2.4/31 | 172.16.2.5/31 |
| spine1:Eth4 – leaf4:Eth1 | 172.16.1.6/31 | 172.16.1.7/31 |
| spine2:Eth4 – leaf4:Eth2 | 172.16.2.6/31 | 172.16.2.7/31 |

ASNs:

Spines: 65000

Leaves: 65101–65104 (leaf1..leaf4 respectively)

VLAN/VNI: VLAN 10 ↔ VNI 1010
