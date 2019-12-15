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
	"github.com/godbus/dbus/v5"
)

type firewall struct{}

func (fw *firewall) getMethods() map[string]string {
	return map[string]string{"GetDefaultZone": "getDefaultZone"}
}

func (fw *firewall) getName() string {
	return dbusInterface
}

func (fw *firewall) GetDefaultZone() (string, *dbus.Error) {
	return "docker-firewalld-shim", nil
}

func createFirewall() (*firewall, error) {
	return &firewall{}, nil
}
