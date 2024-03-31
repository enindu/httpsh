// This file is part of wsh.
//
// wsh is free software: you can redistribute it and/or modify it under the
// terms of the GNU General Public License as published by the Free Software
// Foundation, either version 3 of the License, or (at your option) any later
// version.
//
// wsh is distributed in the hope that it will be useful, but WITHOUT ANY
// WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR
// A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with
// wsh. If not, see <https://www.gnu.org/licenses/>.

package wsh

import (
	"fmt"
	"log"
	"net"
)

type Socket struct {
	Network string
	Host    string
	Port    string
	Log     *log.Logger
}

func (s *Socket) Listen() (*net.TCPListener, error) {
	address := net.JoinHostPort(s.Host, s.Port)
	tcpAddress, err := net.ResolveTCPAddr(s.Network, address)
	if err != nil {
		return nil, fmt.Errorf("listener: %w", err)
	}

	tcpListener, err := net.ListenTCP(s.Network, tcpAddress)
	if err != nil {
		return nil, fmt.Errorf("listener: %w", err)
	}

	s.Log.Println("socket is listening on", address)
	return tcpListener, nil
}
