/*
	docker-firewalld-shim - Shim to trigger recreation of docker iptables
	Copyright (C) 2018 - 2019, Marc Hoersken <info@marc-hoersken.de>

	This program is free software: you can redistribute it and/or modify
	it under the terms of the GNU General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	This program is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU General Public License for more details.

	You should have received a copy of the GNU General Public License
	along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"os/exec"

	"github.com/godbus/dbus/v5"
)

const (
	dbusInterface = "org.fedoraproject.FirewallD1"
)

type firewall struct {
	methods  map[string]string
	commands map[string]string
}

func (fw *firewall) export(conn *dbus.Conn) error {
	return conn.ExportWithMap(fw, fw.methods, dbusPath, dbusInterface)
}

func (fw *firewall) GetDefaultZone() (string, *dbus.Error) {
	return "docker-firewalld-shim", nil
}

func createFirewall() (*firewall, error) {
	cmds := map[string]string{
		"ipv4": "iptables",
		"ipv6": "ip6tables",
		"eb":   "ebtables",
	}
	for ipv, name := range cmds {
		path, err := exec.LookPath(name)
		if err != nil {
			return nil, err
		}
		cmds[ipv] = path
	}
	return &firewall{
		methods:  map[string]string{"GetDefaultZone": "getDefaultZone"},
		commands: cmds,
	}, nil
}
