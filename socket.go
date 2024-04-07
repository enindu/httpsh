// This file is part of httpsh.
//
// httpsh is free software: you can redistribute it and/or modify it under the
// terms of the GNU General Public License as published by the Free Software
// Foundation, either version 3 of the License, or (at your option) any later
// version.
//
// httpsh is distributed in the hope that it will be useful, but WITHOUT ANY
// WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR
// A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with
// httpsh. If not, see <https://www.gnu.org/licenses/>.

package httpsh

import (
	"log/slog"
	"net"
)

type Socket struct {
	Network string
	Host    string
	Port    string
	Log     *slog.Logger
}

func (s *Socket) Listen() (*net.TCPListener, error) {
	address, err := net.ResolveTCPAddr(s.Network, net.JoinHostPort(s.Host, s.Port))
	if err != nil {
		return nil, err
	}

	listener, err := net.ListenTCP(s.Network, address)
	if err != nil {
		return nil, err
	}

	s.Log.Info("socket.listen", "message", "socket is listening", "network", s.Network, "host", s.Host, "port", s.Port)
	return listener, nil
}
