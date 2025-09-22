# Arista CLAB Lab

## Setup Steps
1. Download ubuntu image 
2. Create new VM with image in VirtualBox 

╭──────────────────────────────┬─────────────────────┬─────────┬───────────────────╮
│             Name             │      Kind/Image     │  State  │   IPv4/6 Address  │
├──────────────────────────────┼─────────────────────┼─────────┼───────────────────┤
│ clab-evpn-rdma-fabric-gpu1   │ linux               │ running │ 172.20.20.11      │
│                              │ alpine:3.19         │         │ 3fff:172:20:20::b │
├──────────────────────────────┼─────────────────────┼─────────┼───────────────────┤
│ clab-evpn-rdma-fabric-gpu2   │ linux               │ running │ 172.20.20.3       │
│                              │ alpine:3.19         │         │ 3fff:172:20:20::3 │
├──────────────────────────────┼─────────────────────┼─────────┼───────────────────┤
│ clab-evpn-rdma-fabric-gpu3   │ linux               │ running │ 172.20.20.2       │
│                              │ alpine:3.19         │         │ 3fff:172:20:20::2 │
├──────────────────────────────┼─────────────────────┼─────────┼───────────────────┤
│ clab-evpn-rdma-fabric-gpu4   │ linux               │ running │ 172.20.20.8       │
│                              │ alpine:3.19         │         │ 3fff:172:20:20::8 │
├──────────────────────────────┼─────────────────────┼─────────┼───────────────────┤
│ clab-evpn-rdma-fabric-leaf1  │ ceos                │ running │ 172.20.20.7       │
│                              │ ceosimage:4.34.2.1f │         │ 3fff:172:20:20::7 │
├──────────────────────────────┼─────────────────────┼─────────┼───────────────────┤
│ clab-evpn-rdma-fabric-leaf2  │ ceos                │ running │ 172.20.20.4       │
│                              │ ceosimage:4.34.2.1f │         │ 3fff:172:20:20::4 │
├──────────────────────────────┼─────────────────────┼─────────┼───────────────────┤
│ clab-evpn-rdma-fabric-leaf3  │ ceos                │ running │ 172.20.20.10      │
│                              │ ceosimage:4.34.2.1f │         │ 3fff:172:20:20::a │
├──────────────────────────────┼─────────────────────┼─────────┼───────────────────┤
│ clab-evpn-rdma-fabric-leaf4  │ ceos                │ running │ 172.20.20.5       │
│                              │ ceosimage:4.34.2.1f │         │ 3fff:172:20:20::5 │
├──────────────────────────────┼─────────────────────┼─────────┼───────────────────┤
│ clab-evpn-rdma-fabric-spine1 │ ceos                │ running │ 172.20.20.9       │
│                              │ ceosimage:4.34.2.1f │         │ 3fff:172:20:20::9 │
├──────────────────────────────┼─────────────────────┼─────────┼───────────────────┤
│ clab-evpn-rdma-fabric-spine2 │ ceos                │ running │ 172.20.20.6       │
│                              │ ceosimage:4.34.2.1f │         │ 3fff:172:20:20::6 │
╰──────────────────────────────┴─────────────────────┴─────────┴───────────────────╯

```bash
# 1) watch tx packet counter
sudo docker exec -it clab-evpn-rdma-fabric-gpu4 sh -lc '
  cat /sys/class/net/eth1/statistics/tx_packets ;
  ping -c3 -W1 10.10.10.101 || true ;
  cat /sys/class/net/eth1/statistics/tx_packets
'

# 2) force a gratuitous ARP + broadcast
sudo docker exec -it clab-evpn-rdma-fabric-gpu4 sh -lc '
  ip neigh flush all ;
  arping -c 3 -I eth1 10.10.10.104 || true
'  # if arping missing, apk add iputils-arping
```

```bash
# deploy lab
sudo containerlab deploy -t lab.clab.yml # may need to add --reconfigure
# inspect lab 
sudo containerlab inspect -t lab.clab.yml
# destroy lab
sudo containerlab destroy -t lab.clab.yml
# view topo 
sudo containerlab graph -t lab.clab.yml
```

## Creds 
- user: admin
- pass: admin

## Verify 
```text
show bgp summary
show ip route 10.0.0.1
show ip route 10.0.0.2
show bgp evpn summary
show bgp evpn route-type mac-ip
show vxlan vtep
show vxlan address-table
```

### From GPU1 
```bash
sudo docker exec -it clab-evpn-rdma-fabric-gpu1 sh -lc 'ping -c3 10.10.10.104'
```

## Debug 
### Restart docker 
```bash
sudo snap restart docker
```

## MTU ISSUES
```bash
sudo docker exec -it clab-evpn-rdma-fabric-gpu1 sh -lc 'ip link set eth1 mtu 500'
sudo docker exec -it clab-evpn-rdma-fabric-gpu2 sh -lc 'ip link set eth1 mtu 500'
sudo docker exec -it clab-evpn-rdma-fabric-gpu3 sh -lc 'ip link set eth1 mtu 500'
sudo docker exec -it clab-evpn-rdma-fabric-gpu4 sh -lc 'ip link set eth1 mtu 500'
```

```
leaf1#show vxlan config-sanity
! Your configuration contains warnings. This does not mean misconfigurations. But you may wish to re-check your configurations.
Category                            Result  Detail
---------------------------------- -------- ----------------------------------
Local VTEP Configuration Check       FAIL
  Flood List                         FAIL   No flood list configured
  Flood List                         FAIL   No remote VTEP in VLAN 10
```
