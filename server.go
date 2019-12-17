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
	"errors"
	"fmt"

	"github.com/godbus/dbus/v5"
)

var (
	errRequestNameReplyInQueue = errors.New("name in queue")
	errRequestNameReplyExists  = errors.New("name exists")
	errRequestNameReplyUnknown = errors.New("unknown request name reply")

	errReleaseNameReplyNonExistent = errors.New("name non existent")
	errReleaseNameReplyNotOwner    = errors.New("name not owner")
	errReleaseNameReplyUnknown     = errors.New("unknown release name reply")
)

type service interface {
	getMethods() map[string]string
	getName() string
}

type server struct {
	conn *dbus.Conn
	path dbus.ObjectPath
	name string
}

func (s *server) signalReloaded() error {
	return s.conn.Emit(s.path, s.name+".Reloaded")
}

func (s *server) export(svc service) error {
	return s.conn.ExportWithMap(svc, svc.getMethods(), s.path, svc.getName())
}

func (s *server) unExport(svc service) error {
	return s.conn.ExportWithMap(nil, svc.getMethods(), s.path, svc.getName())
}

func (s *server) publish() error {
	reply, err := s.conn.RequestName(s.name, dbus.NameFlagDoNotQueue)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	switch reply {
	case dbus.RequestNameReplyAlreadyOwner:
		return nil
	case dbus.RequestNameReplyPrimaryOwner:
		return nil
	case dbus.RequestNameReplyInQueue:
		defer s.unPublish()
		return fmt.Errorf("%w", errRequestNameReplyInQueue)
	case dbus.RequestNameReplyExists:
		defer s.unPublish()
		return fmt.Errorf("%w", errRequestNameReplyExists)
	default:
		defer s.unPublish()
		return fmt.Errorf("%w", errRequestNameReplyUnknown)
	}
}

func (s *server) unPublish() error {
	reply, err := s.conn.ReleaseName(s.name)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	switch reply {
	case dbus.ReleaseNameReplyReleased:
		return nil
	case dbus.ReleaseNameReplyNonExistent:
		return fmt.Errorf("%w", errReleaseNameReplyNonExistent)
	case dbus.ReleaseNameReplyNotOwner:
		return fmt.Errorf("%w", errReleaseNameReplyNotOwner)
	default:
		return fmt.Errorf("%w", errReleaseNameReplyUnknown)
	}
}

func createServer(conn *dbus.Conn, path dbus.ObjectPath, name string) (*server, error) {
	return &server{
		conn: conn,
		path: path,
		name: name,
	}, nil
}
