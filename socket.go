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
	"net"
)

type Socket struct {
	Network string
	Address string
}

func (s *Socket) Listener() (*net.TCPListener, error) {
	address, err := net.ResolveTCPAddr(s.Network, s.Address)
	if err != nil {
		return nil, fmt.Errorf("listener: %w", err)
	}

	listener, err := net.ListenTCP(s.Network, address)
	if err != nil {
		return nil, fmt.Errorf("listener: %w", err)
	}

	return listener, nil
}
