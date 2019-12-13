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
	"log"
	"os/exec"

	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/prop"
)

const (
	dbusInterfaceDirect = dbusInterface + ".direct"
)

type firewallDirect struct {
	methods  map[string]string
	firewall *firewall
}

func (fd *firewallDirect) export(conn *dbus.Conn) error {
	return conn.ExportWithMap(fd, fd.methods, dbusPath, dbusInterfaceDirect)
}

func (fd *firewallDirect) Passthrough(ipv string, args []string) (string, *dbus.Error) {
	path, ok := fd.firewall.commands[ipv]
	if !ok {
		return ipv, prop.ErrInvalidArg
	}
	log.SetPrefix("> ")
	log.Println(path, args)
	cmd := exec.Command(path, args...)
	bytes, err := cmd.CombinedOutput()
	output := string(bytes)
	if err != nil {
		log.SetPrefix("! ")
		log.Println(output)
		return err.Error(), prop.ErrInvalidArg
	}
	log.SetPrefix("< ")
	log.Println(output)
	return output, nil
}

func createFirwallDirect(fw *firewall) (*firewallDirect, error) {
	return &firewallDirect{
		firewall: fw,
		methods:  map[string]string{"Passthrough": "passthrough"},
	}, nil
}
