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
	cert   *bool = flag.Bool("cert", false, "Create CA, server, and client certs")
	server *bool = flag.Bool("server", false, "Run server")
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
	case *cert:
		ca := &CA{
			bits:        viper.GetInt("cert.ca.bits"),
			years:       viper.GetInt("cert.ca.years"),
			months:      viper.GetInt("cert.ca.months"),
			days:        viper.GetInt("cert.ca.days"),
			key:         viper.GetString("cert.ca.key"),
			certificate: viper.GetString("cert.ca.certificate"),
		}

		bundle, err := ca.generate()
		if err != nil {
			log.Error("main", "message", err)
			return
		}

		server := &Request{
			bits:        viper.GetInt("cert.request.server.bits"),
			years:       viper.GetInt("cert.request.server.years"),
			months:      viper.GetInt("cert.request.server.months"),
			days:        viper.GetInt("cert.request.server.days"),
			domains:     viper.GetStringSlice("cert.request.server.domains"),
			emails:      viper.GetStringSlice("cert.request.server.emails"),
			ips:         viper.GetStringSlice("cert.request.server.ips"),
			key:         viper.GetString("cert.request.server.key"),
			certificate: viper.GetString("cert.request.server.certificate"),
		}

		_, err = server.generate(bundle)
		if err != nil {
			log.Error("main", "message", err)
			return
		}

		client := &Request{
			bits:        viper.GetInt("cert.request.client.bits"),
			years:       viper.GetInt("cert.request.client.years"),
			months:      viper.GetInt("cert.request.client.months"),
			days:        viper.GetInt("cert.request.client.days"),
			domains:     viper.GetStringSlice("cert.request.client.domains"),
			emails:      viper.GetStringSlice("cert.request.client.emails"),
			ips:         viper.GetStringSlice("cert.request.client.ips"),
			key:         viper.GetString("cert.request.client.key"),
			certificate: viper.GetString("cert.request.client.certificate"),
		}

		_, err = client.generate(bundle)
		if err != nil {
			log.Error("main", "message", err)
			return
		}
	case *server:
		server := &Server{
			network:           viper.GetString("server.network"),
			host:              viper.GetString("server.host"),
			port:              viper.GetString("server.port"),
			domain:            viper.GetString("server.domain"),
			readTimeout:       viper.GetInt("server.read_timeout"),
			writeTimeout:      viper.GetInt("server.write_timeout"),
			idleTimeout:       viper.GetInt("server.idle_timeout"),
			caCertificate:     viper.GetString("server.ca_certificate"),
			serverKey:         viper.GetString("server.server_key"),
			serverCertificate: viper.GetString("server.server_certificate"),
			directory:         viper.GetString("server.directory"),
			mime:              viper.GetString("server.mime"),
			methods:           viper.GetStringSlice("server.methods"),
			executables:       viper.GetStringMapStringSlice("server.executables"),
			log:               log,
		}

		err := server.run()
		if err != nil {
			log.Error("main", "message", err)
			return
		}
	default:
		flag.PrintDefaults()
	}
}
