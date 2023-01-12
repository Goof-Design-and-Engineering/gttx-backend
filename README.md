# gttx-backend

Takes the `pocketbase.go` file and compile to single binary via docker. Use the docker image build process to host for free on fly.io

## How 2 Deploy

### If not signed up

flyctl auth signup
flyctl deploy

### if not signed in

flyctl auth signin
flyctl deploy

## How 2 Backup DB

```bash
# this will register a ssh key with your local agent (if you haven't already)
flyctl ssh issue --agent

# proxies connections to a fly VM through a Wireguard tunnel
flyctl proxy 10022:22

# run in a separate terminal to copy the pb_data directory
scp -r -P 10022 root@localhost:/pb/pb_data  /your/local/pb_data
```
