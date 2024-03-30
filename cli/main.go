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

package main

import (
	"fmt"
	"net"

	"github.com/enindu/wsh"
)

func main() {
	socket := &wsh.Socket{
		Address: net.JoinHostPort(wsh.Host, wsh.Port),
		Network: wsh.Network,
	}

	listener, err := socket.Listener()
	if err != nil {
		fmt.Printf("\r%v\n", err)
		return
	}

	defer listener.Close()

	handler := &wsh.Handler{
		Commands:  wsh.Commands,
		Directory: wsh.Directory,
		Methods:   wsh.Methods,
		Mime:      wsh.Mime,
	}
}
