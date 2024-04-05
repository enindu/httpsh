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
	"log/slog"

	"github.com/spf13/viper"
)

func main() {
	log := slog.Default()

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.config/httpsh/cli")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		log.Error("main", "message", err)
		return
	}

	server := &Server{
		network:           viper.GetString("server.network"),
		host:              viper.GetString("server.host"),
		port:              viper.GetString("server.port"),
		domain:            viper.GetString("server.domain"),
		readTimeout:       viper.GetInt("server.read_timeout"),
		writeTimeout:      viper.GetInt("server.write_timeout"),
		idleTimeout:       viper.GetInt("server.idle_timeout"),
		caCertificate:     viper.GetString("server.ca_certificate"),
		serverCertificate: viper.GetString("server.server_certificate"),
		serverKey:         viper.GetString("server.server_key"),
		directory:         viper.GetString("server.directory"),
		mime:              viper.GetString("server.mime"),
		methods:           viper.GetStringSlice("server.methods"),
		executables:       viper.GetStringMapStringSlice("server.executables"),
		log:               log,
	}

	err = server.run()
	if err != nil {
		log.Error("main", "message", err)
		return
	}
}
