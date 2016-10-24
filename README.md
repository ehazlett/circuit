```
 ________  ___  ________  ________  ___  ___  ___  _________
|\   ____\|\  \|\   __  \|\   ____\|\  \|\  \|\  \|\___   ___\
\ \  \___|\ \  \ \  \|\  \ \  \___|\ \  \\\  \ \  \|___ \  \_|
 \ \  \    \ \  \ \   _  _\ \  \    \ \  \\\  \ \  \   \ \  \
  \ \  \____\ \  \ \  \\  \\ \  \____\ \  \\\  \ \  \   \ \  \
   \ \_______\ \__\ \__\\ _\\ \_______\ \_______\ \__\   \ \__\
    \|_______|\|__|\|__|\|__|\|_______|\|_______|\|__|    \|__|

```

Circuit manages networks for [runc](https://runc.io).  Features include
but not limited to:

- Network bridge creation
- Virtual ethernet pairs for containers
- Quality of service management for networks and container interfaces
- Load balancing using IPVS

Circuit has been designed for flexibility.  For example, the controller has
been designed to be replaced.  By default, Circuit uses internal bridging
but it could be replaced by an Open vSwitch controller.  The same goes for
the internal IPAM, data service and load balancing.  External load balancer
implementations such as HAProxy or Nginx integration  would be trivial
to write and utilize in Circuit.

# Usage
The following show example usage of Circuit.

## Create a Network

```
$> circuit network create sandbox 10.254.1.0/24
```

## View Networks
```
$> circuit network ls
NAME                SUBNET
sandbox             10.254.1.0/24
```

## Connect a Container to a Network
```
$> runc list
ID          PID         STATUS      BUNDLE         CREATED
web-00      4668        running     /root/web-00   2016-10-12T18:45:27.787840219Z

$> circuit network connect 4668 sandbox
INFO[0000] connected container 4668 to network sandbox
```

## Set QoS for Network
This will set the target rate of 5mbps with a ceiling of 6mbps
```
$> circuit network qos set --rate 5000 --ceiling 6000 sandbox
INFO[0000] qos configured for sandbox
```

This will add 50ms latency to the network
```
$> circuit network qos set --delay 50ms sandbox
INFO[0000] qos configured for sandbox
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
INFO[0000] qos reset for sandbox
```

## Create a Load Balancer Service
```
$> circuit lb create demo 192.168.100.235:80
INFO[0000] service demo created
```

## Create a Load Balancer Service with Custom Scheduler
```
$> circuit lb create demo-wrr --scheduler wrr 192.168.100.235:80
INFO[0000] service demo-wrr created
```
## List Load Balancer Services
```
$> circuit lb ls
NAME                ADDR                 PROTOCOL            SCHEDULER
demo                192.168.100.235:80   tcp                 rr
```

## Add Target to Service
```
$> circuit lb add demo 10.254.1.196:80
INFO[0000] service demo updated
```

## List Load Balancer Services with Details
```
$> circuit lb ls --details
NAME                ADDR                 PROTOCOL            SCHEDULER
demo                192.168.100.235:80   tcp                 rr
  -> 10.254.1.196:80
```

## Remove Target from Service
```
$> circuit lb remove demo 10.254.1.196:80
INFO[0000] service demo updated
```

## Remove Service
```
$> circuit lb delete demo
INFO[0000] service demo removed
```

## Disconnect Container from Network
```
$> runc list
ID          PID         STATUS      BUNDLE         CREATED
web-00      4668        running     /root/web-00   2016-10-12T18:45:27.787840219Z

$> circuit network disconnect 4668 sandbox
INFO[0000] disconnected container 4668 from network sandbox
```

## Delete Network
```
$> circuit network delete sandbox
INFO[0000] sandbox deleted
```

# runc Hooks
Circuit also supports runc hooks.  This will automatically create and configure
networks upon start / stop for runc containers.  To setup, simply add Circuit
as `prestart` and `poststop` hooks in a runc config:

This will automatically create a new network named the same as the container.
It will also remove and cleanup upon container stop.
```
...

"hooks": {
		"prestart": [
			{
				"path": "/usr/local/bin/circuit"
			}
		],
		"poststop": [
			{
				"path": "/usr/local/bin/circuit"
			}
		]
	},

...
```

You can also configure the network name and subnet with environment
variables.  Specify `NETWORK` and `SUBNET` in the config to have Circuit
create the network using the specified name and subnet:

```
...

"hooks": {
		"prestart": [
			{
				"path": "/usr/local/bin/circuit",
				"env": ["NETWORK=demo", "SUBNET=10.254.200.0/24","PATH=/bin:/usr/bin:/usr/sbin:/sbin"]
			}
		],
		"poststop": [
			{
				"path": "/usr/local/bin/circuit",
				"env": ["NETWORK=demo", "PATH=/bin:/usr/bin:/usr/sbin:/sbin"]
			}
		]
	},

...
```
