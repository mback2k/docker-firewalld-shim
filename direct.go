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
	"context"
	"log"
	"os/exec"

	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/prop"
)

type firewallDirect struct {
	ctx      context.Context
	commands map[string]string
}

func (fd *firewallDirect) getMethods() map[string]string {
	return map[string]string{"Passthrough": "passthrough"}
}

func (fd *firewallDirect) getName() string {
	return dbusInterface + ".direct"
}

func (fd *firewallDirect) Passthrough(ipv string, args []string) (string, *dbus.Error) {
	path, ok := fd.commands[ipv]
	if !ok {
		return ipv, prop.ErrInvalidArg
	}
	defer log.SetPrefix("")
	log.SetPrefix("> ")
	log.Println(path, args)
	cmd := exec.CommandContext(fd.ctx, path, args...)
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

func createFirewallDirect(ctx context.Context) (*firewallDirect, error) {
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
	return &firewallDirect{
		ctx:      ctx,
		commands: cmds,
	}, nil
}
