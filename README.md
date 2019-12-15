docker-firewalld-shim
=====================
This Go program is a compatibility shim to trigger the recreation
of docker iptables to use with netfilter-persistent and systemd.

[![Build Status](https://travis-ci.org/mback2k/docker-firewalld-shim.svg?branch=master)](https://travis-ci.org/mback2k/docker-firewalld-shim)
[![GoDoc](https://godoc.org/github.com/mback2k/docker-firewalld-shim?status.svg)](https://godoc.org/github.com/mback2k/docker-firewalld-shim)
[![Go Report Card](https://goreportcard.com/badge/github.com/mback2k/docker-firewalld-shim)](https://goreportcard.com/report/github.com/mback2k/docker-firewalld-shim)

Installation
------------
You basically have two options to install this Go program package:

1. If you have Go installed and configured on your PATH, just do the following go get inside your GOPATH to get the latest version:

```
go get -u github.com/mback2k/docker-firewalld-shim
```

2. If you do not have Go installed and just want to use a released binary,
then you can just go ahead and download a pre-compiled Linux amd64 binary from the [Github releases](https://github.com/mback2k/docker-firewalld-shim/releases).

Finally put the docker-firewalld-shim binary onto your PATH and make sure it is executable.

Usage
-----
The following is an example of a systemd.service unit file which runs this shim as a daemon:

```
[Unit]
Description=Docker FirewallD Shim
Wants=network.target
After=dbus.service
Before=docker.service
Conflicts=firewalld.service
PartOf=iptables.service ip6tables.service ebtables.service ipset.service netfilter-persistent.service
ReloadPropagatedFrom=iptables.service ip6tables.service ebtables.service ipset.service netfilter-persistent.service

[Service]
Type=dbus
BusName=org.fedoraproject.FirewallD1
ExecStart=/usr/local/sbin/docker-firewalld-shim
KillMode=mixed
Restart=on-failure
PrivateTmp=true
ProtectHome=true
ProtectSystem=full

[Install]
WantedBy=docker.service
```

This daemon service replaces and conflicts with the original `firewalld.service`,
since it provides a subset of the firewalld interface in order to allow the docker
daemon to detect firewalld as running and passthrough all iptables commands to it.

In order to be able to actually run and use this tool as a firewalld simulation,
you need to deploy the dbus system policy for firewalld to `/etc/dbus-1/system.d/`:

Just store the following file from the firewalld project repository:
```
https://raw.githubusercontent.com/firewalld/firewalld/v0.8.0/config/FirewallD.conf
```
as:
```
/etc/dbus-1/system.d/FirewallD.conf
```
and then reload the dbus service via `systemctl reload dbus`.

Disclaimer
----------
This tool is meant as an easy way to make docker recreate its iptables rules
without restarting it and all containers. This tool does not fully implement
the firewalld interface specification and is not tested for any other purpose.

I personally use this tool to make docker recreate its rules after I have
deployed changes to the firewall rules through iptables-restore via Ansible.

Credits
-------
This tool was developed by inspecting the following files from the docker/libnetwork repository:

* https://github.com/docker/libnetwork/blob/bump_19.03/iptables/firewalld.go
* https://github.com/docker/libnetwork/blob/bump_19.03/iptables/iptables.go

This tool was developed by inspecting the following files from firewalld/firewalld repository:

* https://github.com/firewalld/firewalld/blob/v0.8.0/src/firewall/core/fw.py
* https://github.com/firewalld/firewalld/blob/v0.8.0/src/firewall/core/fw_direct.py
* https://github.com/firewalld/firewalld/blob/v0.8.0/src/firewall/core/fw_zone.py
* https://github.com/firewalld/firewalld/blob/v0.8.0/src/firewall/core/ipXtables.py
* https://github.com/firewalld/firewalld/blob/v0.8.0/src/firewall/core/prog.py

Thanks to Open Source it was possible to develop this small helper tool!

License
-------
Copyright (C) 2018 - 2019, Marc Hoersken <info@marc-hoersken.de>

This software is licensed as described in the file LICENSE, which
you should have received as part of this software distribution.

All trademarks are the property of their respective owners.
