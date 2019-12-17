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
	"os"
	"os/signal"
	"syscall"

	"github.com/godbus/dbus/v5"
)

const (
	dbusPath      = "/org/fedoraproject/FirewallD1"
	dbusInterface = "org.fedoraproject.FirewallD1"
)

func handleSignals(ctx context.Context, stop chan<- bool) {
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGHUP)
	signal.Notify(sigs, syscall.SIGINT)
	signal.Notify(sigs, syscall.SIGTERM)

	go func() {
		defer signal.Reset(syscall.SIGHUP)
		defer signal.Reset(syscall.SIGINT)
		defer signal.Reset(syscall.SIGTERM)

		defer close(sigs)

	loop:
		for {
			select {
			case <-ctx.Done():
				// returning not to leak the goroutine
				break loop
			case sig := <-sigs:
				switch sig {
				case syscall.SIGINT:
					stop <- true
				case syscall.SIGTERM:
					stop <- true
				case syscall.SIGHUP:
					stop <- false
				}
			}
		}
	}()
}

func workerLoop(ctx context.Context, stop <-chan bool) {
	conn, err := dbus.SystemBusPrivate(dbus.WithContext(ctx))
	if err != nil {
		log.Panicln(err)
	}
	defer conn.Close()

	err = conn.Auth(nil)
	if err != nil {
		log.Panicln(err)
	}

	err = conn.Hello()
	if err != nil {
		log.Panicln(err)
	}

	s, err := createServer(conn, dbusPath, dbusInterface)
	if err != nil {
		log.Panicln(err)
	}

	fw, err := createFirewall()
	if err != nil {
		log.Panicln(err)
	}

	fd, err := createFirewallDirect(ctx)
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

loop:
	for {
		select {
		case <-ctx.Done():
			// returning not to leak the goroutine
			break loop
		case stop := <-stop:
			if stop {
				break loop
			} else {
				log.Println("Sending signal Reloaded ...")
				err := s.signalReloaded()
				if err != nil {
					log.Panicln(err)
				}
			}
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stop := make(chan bool)

	handleSignals(ctx, stop)
	workerLoop(ctx, stop)

	close(stop)
}
