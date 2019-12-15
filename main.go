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

	"github.com/godbus/dbus/v5"
)

const (
	dbusPath      = "/org/fedoraproject/FirewallD1"
	dbusInterface = "org.fedoraproject.FirewallD1"
)

func main() {
	conn, err := dbus.SystemBus()
	if err != nil {
		log.Panicln(err)
	}
	defer conn.Close()

	s, err := createServer(conn, dbusPath, dbusInterface)
	if err != nil {
		log.Panicln(err)
	}

	fw, err := createFirewall()
	if err != nil {
		log.Panicln(err)
	}

	fd, err := createFirewallDirect()
	if err != nil {
		log.Panicln(err)
	}

	err = s.export(fw)
	if err != nil {
		log.Panicln(err)
	}
	defer s.unExport(fw)

	err = s.export(fd)
	if err != nil {
		log.Panicln(err)
	}
	defer s.unExport(fd)

	err = s.publish()
	if err != nil {
		log.Panicln(err)
	}
	defer s.unPublish()

	select {}
}
