```
 ________  ___  ________  ________  ___  ___  ___  _________
|\   ____\|\  \|\   __  \|\   ____\|\  \|\  \|\  \|\___   ___\
\ \  \___|\ \  \ \  \|\  \ \  \___|\ \  \\\  \ \  \|___ \  \_|
 \ \  \    \ \  \ \   _  _\ \  \    \ \  \\\  \ \  \   \ \  \
  \ \  \____\ \  \ \  \\  \\ \  \____\ \  \\\  \ \  \   \ \  \
   \ \_______\ \__\ \__\\ _\\ \_______\ \_______\ \__\   \ \__\
    \|_______|\|__|\|__|\|__|\|_______|\|_______|\|__|    \|__|

```

Circuit manages networks for [runc](https://runc.io).

- CNI network management (define and manage CNI networks and connectivity)
- CNI compatible (use CNI plugins)
- Quality of service management for networks and container interfaces
- Load balancing using IPVS

Circuit has been designed for flexibility.  For example, the controller has
been designed to be replaced.  Circuit leverages
[CNI](https://github.com/containernetworking/cni)
for setting up networking using various plugins such as bridge, ptp, etc.
Define multiple CNI networks and connect/disconnect, load balance, etc.

Huge thanks to @jessfraz [netns](https://github.com/jessfraz/netns) for
inspiration :)

# Usage
The following examples assume you have CNI plugins.  Checkout the [CNI docs](https://github.com/containernetworking/cni#running-the-plugins)
on getting started (mainly just `clone` and `./build`).  From there you can
place those binaries somewhere on your `PATH` and it should just work.

## Create a Network

Specify a CNI network config when creating.  For the examples following, we
will assume a network config like so:

```
{
    "cniVersion": "0.2.0",
    "name": "br-sandbox",
    "type": "bridge",
    "bridge": "cni0",
    "ipMasq": true,
    "isGateway": true,
    "ipam": {
        "type": "host-local",
        "subnet": "10.100.10.0/24",
        "routes": [
            {
                "dst": "0.0.0.0/0"
            }
        ]
    }
}
```

```
$> circuit network create /path/to/cni.conf
```

## View Networks
```
$> circuit network list
NAME                TYPE                VERSION             PEERS
local               ipvlan              0.2.0
sandbox             bridge              0.2.0
shell               bridge              0.2.0               10.30.30.2 (19022)
```

## Connect a Container to a Network
```
$> runc list
ID          PID         STATUS      BUNDLE         CREATED
web-00      4668        running     /root/web-00   2016-10-12T18:45:27.787840219Z

$> circuit network connect 4668 sandbox
connected container 4668 to network sandbox
```

## Set QoS for Network
This will set the target rate of 5mbps with a ceiling of 6mbps
```
$> circuit network qos set --rate 5000 --ceiling 6000 sandbox
qos configured for sandbox
```

This will add 50ms latency to the network
```
$> circuit network qos set --delay 50ms sandbox
qos configured for sandbox
```

An example ping from the container with before and after QOS:

```
$> ping 10.254.1.1
64 bytes from 10.254.1.1: icmp_seq=1 ttl=64 time=0.176 ms
64 bytes from 10.254.1.1: icmp_seq=2 ttl=64 time=0.136 ms
64 bytes from 10.254.1.1: icmp_seq=3 ttl=64 time=0.150 ms
64 bytes from 10.254.1.1: icmp_seq=4 ttl=64 time=0.138 ms
64 bytes from 10.254.1.1: icmp_seq=5 ttl=64 time=50.361 ms
64 bytes from 10.254.1.1: icmp_seq=6 ttl=64 time=50.323 ms
64 bytes from 10.254.1.1: icmp_seq=7 ttl=64 time=50.280 ms
64 bytes from 10.254.1.1: icmp_seq=8 ttl=64 time=50.352 ms
```

## Clear QoS for a Network
```
$> circuit network qos reset sandbox
qos reset for sandbox
```

Circuit supports basic load balancing via IPVS.

Note: this is experimental and the implementation may change.

## Create a Load Balancer Service
```
$> circuit lb create demo 192.168.100.235:80
service demo created
```

## Create a Load Balancer Service with Custom Scheduler
```
$> circuit lb create demo-wrr --scheduler wrr 192.168.100.235:80
service demo-wrr created
```
## List Load Balancer Services
```
$> circuit lb list
NAME                ADDR                 PROTOCOL            SCHEDULER
demo                192.168.100.235:80   tcp                 rr
```

## Add Target to Service
```
$> circuit lb add demo 10.254.1.196:80
service demo updated
```

## List Load Balancer Services with Details
```
$> circuit lb list --details
NAME                ADDR                 PROTOCOL            SCHEDULER
demo                192.168.100.235:80   tcp                 rr
  -> 10.254.1.196:80
```

## Remove Target from Service
```
$> circuit lb remove demo 10.254.1.196:80
service demo updated
```

## Remove Service
```
$> circuit lb delete demo
service demo removed
```

## Disconnect Container from Network
```
$> runc list
ID          PID         STATUS      BUNDLE         CREATED
web-00      4668        running     /root/web-00   2016-10-12T18:45:27.787840219Z

$> circuit network disconnect 4668 sandbox
disconnected container 4668 from network sandbox
```

## Delete Network
```
$> circuit network delete sandbox
sandbox deleted
```

# runc Hooks
Circuit also supports runc hooks.  This will automatically create and configure
networks upon start / stop for runc containers.  To setup, simply add Circuit
as `prestart` and `poststop` hooks in a runc config:

```
...

"hooks": {
    "prestart": [
        {
            "path": "/usr/local/bin/circuit",
            "env": [
                "CNI_CONF=/etc/cni/conf.d/bridge.conf",
                "CNI_PATH=/path/to/cni/plugins",
                "PATH=/bin:/usr/bin:/usr/sbin:/sbin"
            ]
        }
    ],
    "poststop": [
        {
            "path": "/usr/local/bin/circuit",
            "env": [
                "CNI_CONF=/etc/cni/conf.d/bridge.conf",
                "CNI_PATH=/path/to/cni/plugins",
                "PATH=/bin:/usr/bin:/usr/sbin:/sbin"
            ]
        }
    ]
},
...
```
