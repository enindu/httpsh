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
	raw := net.JoinHostPort(s.Host, s.Port)

	address, err := net.ResolveTCPAddr(s.Network, raw)
	if err != nil {
		return nil, err
	}

	listener, err := net.ListenTCP(s.Network, address)
	if err != nil {
		return nil, err
	}

	s.Log.Println("socket is listening on", address)

	return listener, nil
}
