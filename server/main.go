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

package main

import (
	"log"

	"github.com/enindu/httpsh"
	"github.com/spf13/viper"
)

func main() {
	logger := log.Default()
	socket := &httpsh.Socket{
		Network: viper.GetString("network"),
		Host:    viper.GetString("host"),
		Port:    viper.GetString("port"),
		Log:     logger,
	}

	listener, err := socket.Listen()
	if err != nil {
		log.Printf("%v", err)
		return
	}

	defer listener.Close()

	handler := &httpsh.Handler{
		Directory:   viper.GetString("directory"),
		Mime:        viper.GetString("mime"),
		Methods:     viper.GetStringSlice("methods"),
		Executables: viper.GetStringMapStringSlice("executables"),
		Log:         logger,
	}

	server := &httpsh.Server{
		Listener: listener,
		Handler:  handler,
		Read:     viper.GetInt("read"),
		Write:    viper.GetInt("write"),
		Idle:     viper.GetInt("idle"),
		Log:      logger,
	}

	err = server.Run()
	if err != nil {
		logger.Printf("%v", err)
		return
	}
}
