test_data_plane:
	sudo docker exec -it clab-evpn-rdma-fabric-gpu1 sh -lc 'ping -c3 10.10.10.104'
