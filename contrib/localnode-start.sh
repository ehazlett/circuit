#!/bin/sh
CONF_DIR=/var/lib/circuit
mkdir -p $CONF_DIR
cat <<EOF>$CONF_DIR/ctr0.json
{
  "cniVersion": "0.3.1",
  "name": "ctr0",
  "type": "bridge",
  "bridge": "ctr0",
  "isGateway": true,
  "ipMasq": true,
  "promiscMode": true,
  "ipam": {
    "type": "host-local",
    "subnet": "10.88.0.0/16",
    "routes": [
      { "dst": "0.0.0.0/0" }
    ]
  }
}
EOF
export CONTAINERD_SNAPSHOTTER=native
nohup containerd -c /etc/containerd/config.toml &
exec circuit $@
