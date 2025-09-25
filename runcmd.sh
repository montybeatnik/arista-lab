SW=172.20.20.9          # switch mgmt IP/hostname
USER=admin              # your creds
PASS=admin

curl -sk -u "$USER:$PASS" \
  -H "Content-Type: application/json" \
  -X POST https://$SW/command-api \
  -d '{
        "jsonrpc": "2.0",
        "method": "runCmds",
        "params": {
          "version": 1,
          "format": "json",
          "cmds": ["show ip interface brief"]
        },
        "id": 1
      }' | jq
