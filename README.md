```
      _                _ _
     (_)              (_) |
  ___ _ _ __ ___ _   _ _| |_
 / __| | '__/ __| | | | | __|
| (__| | | | (__| |_| | | |_
 \___|_|_|  \___|\__,_|_|\__|

```

![CI](https://github.com/ehazlett/circuit/workflows/CI/badge.svg)

# Circuit
Circuit is a container management application for [containerd](https://github.com/containerd/containerd).

It can be used imperitively to connect/disconnect containers to/from networks and also run as a daemon to
listen for containerd events and connect/disconnect containers automatically.

Circuit can also provide basic restart capabilities for containers.  By adding the `io.circuit.restart` label,
the daemon will monitor containers and restart if they exit.

Circuit can be used with a [CoreDNS](https://coredns.io/) [Plugin](https://github.com/ehazlett/circuit-coredns)
to provide DNS responses for containers.  This is most useful with the [macvlan](https://github.com/containernetworking/plugins/tree/master/plugins/main/macvlan)
CNI plugin.  Circuit uses [NATS](https://nats.io/) to provide basic clustering to enable container IP resolution
across a fleet of Circuit nodes.

# Usage
The daemon and cli is combined in a single binary.

## Daemon
To run the daemon, use the `server` subcommand:

```
$> circuit server
```

This will start the GRPC server on port `8080` by default.

## CLI
To use the CLI start the server and then use the various subcommands:

Circuit network definitions are simply CNI specs.  To create a network for use with Circuit, use the `create` command.

As an example, you can create a bridge network using the following config as `bridge.json`:

```
{
    "cniVersion": "0.3.1",
    "name": "ctr0",
    "type": "bridge",
    "bridge": "ctr0",
    "isDefaultGateway": true,
    "forceAddress": false,
    "ipMasq": true,
    "hairpinMode": true,
    "ipam": {
        "type": "host-local",
        "subnet": "10.255.0.0/16"
    }
}
```

Create the network in Circuit:

```
$> circuit network create ctr0 bridge.json
```

You can then list networks:

```
$> circuit network ls
NAME      TYPE
ctr0      bridge
```

Run a container with no external networking:

```
$> ctr run -t docker.io/library/alpine:latest shell sh
/ # ip a s
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
    inet 127.0.0.1/8 scope host lo
       valid_lft forever preferred_lft forever
    inet6 ::1/128 scope host
       valid_lft forever preferred_lft forever
```

Now connect the `shell` container to the `ctr0` network:

Note: make sure to have the CNI plugins [installed](https://github.com/containernetworking/plugins/releases) to `/opt/containerd/bin` (can be changed with `--cni-path`).

```
$> circuit network connect shell ctr0
connected shell to ctr0 with ip=10.255.0.2
```

Confirm that the container has the interface:

```
/ # ip a s
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
    inet 127.0.0.1/8 scope host lo
       valid_lft forever preferred_lft forever
    inet6 ::1/128 scope host
       valid_lft forever preferred_lft forever
3: eth0@if45: <BROADCAST,MULTICAST,UP,LOWER_UP,M-DOWN> mtu 1500 qdisc noqueue state UP
    link/ether 22:e0:21:33:ec:18 brd ff:ff:ff:ff:ff:ff
    inet 10.255.0.3/16 brd 10.255.255.255 scope global eth0
       valid_lft forever preferred_lft forever
    inet6 fe80::20e0:21ff:fe33:ec18/64 scope link
       valid_lft forever preferred_lft forever
```

You can also see all of the container IPs:

```
$> circuit network ips shell
NETWORK   IP              INTERFACE
ctr0      10.255.0.3      eth0
```

## Automatic Networking
Circuit can run as a daemon and use containerd events to automatically connect and disconnect
containers.

Note: currently automatic connection is limited to a single network.

To enable automatic connecting, use the `io.circuit.network` label when creating the container:

```
$> ctr run -t --label io.circuit.network=ctr0 docker.io/library/alpine:latest shell sh
```

There should already be an additional interface in the container:

```
/ # ip a s
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
    inet 127.0.0.1/8 scope host lo
       valid_lft forever preferred_lft forever
    inet6 ::1/128 scope host
       valid_lft forever preferred_lft forever
3: eth0@if46: <BROADCAST,MULTICAST,UP,LOWER_UP,M-DOWN> mtu 1500 qdisc noqueue state UP
    link/ether de:6a:45:37:7a:5d brd ff:ff:ff:ff:ff:ff
    inet 10.255.0.4/16 brd 10.255.255.255 scope global eth0
       valid_lft forever preferred_lft forever
    inet6 fe80::dc6a:45ff:fe37:7a5d/64 scope link tentative
       valid_lft forever preferred_lft forever
```

You should also see a log message for the connect event:
```
DEBU[0004] task start: container=shell pid=23217
INFO[0005] connected shell to ctr0 with ip 10.255.0.5
```

# Clustering
Circuit can be configured to use [NATS](https://nats.io/) so that when querying for the container IP
(i.e `circuit network ips <name>`) the Circuit nodes will query each other internally and return
all known IPs of containers with that name.  Note: you will need to setup a NATS host separately.

To form a Circuit cluster, simply configure the Circuit server to connect to NATS:

```
$> circuit server --nats-addr 1.2.3.4:4222
```

You can then list all available nodes:

```
$> circuit cluster nodes
NAME
white-rabbit
```

You can then query for container IPs across all Circuit nodes:

```
$> circuit network ips shell
NETWORK   IP             INTERFACE
ctr0      10.10.214.28   eth0
```

This will query the Circuit cluster for all available IPs of containers with the name `shell`.

# API
There is a GRPC API that the CLI uses for management.  This can also be used in third party applications for more control
over container network management.
