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
	"flag"
	"log/slog"

	"github.com/spf13/viper"
)

var (
	ca    *bool = flag.Bool("ca", false, "Create CA key and certificate")
	serve *bool = flag.Bool("serve", false, "Run server")
)

func main() {
	flag.Parse()

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

	switch {
	case *ca:
		ca := &CA{
			bits:            viper.GetInt("ca.bits"),
			years:           viper.GetInt("ca.years"),
			months:          viper.GetInt("ca.months"),
			days:            viper.GetInt("ca.days"),
			keyFile:         viper.GetString("ca.key_file"),
			certificateFile: viper.GetString("ca.certificate_file"),
		}

		err := ca.generate()
		if err != nil {
			log.Error("main", "message", err)
			return
		}
	case *serve:
		serve := &Serve{
			network:               viper.GetString("server.network"),
			host:                  viper.GetString("server.host"),
			port:                  viper.GetString("server.port"),
			domain:                viper.GetString("server.domain"),
			readTimeout:           viper.GetInt("server.read_timeout"),
			writeTimeout:          viper.GetInt("server.write_timeout"),
			idleTimeout:           viper.GetInt("server.idle_timeout"),
			caCertificateFile:     viper.GetString("server.ca_certificate_file"),
			serverKeyFile:         viper.GetString("server.server_key_file"),
			serverCertificateFile: viper.GetString("server.server_certificate_file"),
			directory:             viper.GetString("server.directory"),
			mime:                  viper.GetString("server.mime"),
			methods:               viper.GetStringSlice("server.methods"),
			executables:           viper.GetStringMapStringSlice("server.executables"),
			log:                   log,
		}

		err = serve.run()
		if err != nil {
			log.Error("main", "message", err)
			return
		}
	default:
		flag.PrintDefaults()
	}
}
