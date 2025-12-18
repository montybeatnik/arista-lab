setup_git:
	git config --global user.email "christopher.t.hern@gmail.com"
	git config --global user.name "Christopher Hern"

test_data_plane:
	sudo docker exec -it clab-evpn-rdma-fabric-gpu1 sh -lc 'ping -c3 10.10.10.104'

fix_acl_perm_issue:
	# 1) Pick a local (ext4) path for clab’s runtime
	mkdir -p ~/.clab-runs
	
	# 2) Ensure the env var is set AND preserved for sudo
	export CLAB_LABDIR_BASE=$HOME/.clab-runs
	sudo -E env | grep CLAB_LABDIR_BASE || echo "env not preserved!"
	
	# 3) Clean any half-created labdir on the mounted repo (optional but tidy)
	rm -rf ~/lab/clab-evpn-rdma-fabric
	
	# 4) Redeploy (note the -E to keep the env var)
	sudo -E containerlab destroy  -t ~/lab/lab.clab.yml || true
	sudo -E containerlab deploy   -t ~/lab/lab.clab.yml --reconfigure

update_python_deps:
	sudo apt install python3.12-venv -y
	sudo apt install python3-pip
	python3 -m venv venv
	source venv/bin/activate
	pip install -r requirements.txt

done_with_py_venv:
	deactivate

db_shell:
	sudo docker-compose exec db psql -U username device_inventory

# db_create_user:
#     docker exec -i db psql -U username device_inventory -c "CREATE ROLE lab WITH PASSWORD 'password';"
#     docker exec -i db psql -U username device_inventory -c "ALTER ROLE lab WITH SUPERUSER;"
# 	CREATE ROLE lab WITH PASSWORD 'password';
# 	ALTER ROLE lab WITH SUPERUSER;
#   ALTER ROLE lab WITH LOGIN;