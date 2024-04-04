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
	"crypto/tls"
	"crypto/x509"
	"log/slog"
	"os"

	"github.com/enindu/httpsh"
	"github.com/spf13/viper"
)

func main() {
	logger := slog.Default()
	socket := &httpsh.Socket{
		Network: viper.GetString("network"),
		Host:    viper.GetString("host"),
		Port:    viper.GetString("port"),
		Log:     logger,
	}

	listener, err := socket.Listen()
	if err != nil {
		logger.Error("main", "message", err.Error())
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

	caCertificate, err := os.ReadFile(viper.GetString("ca_certificate"))
	if err != nil {
		logger.Error("main", "message", err.Error())
		return
	}

	caPool := x509.NewCertPool()
	caPool.AppendCertsFromPEM(caCertificate)

	config := &tls.Config{
		RootCAs:            caPool,
		ServerName:         viper.GetString("host"),
		ClientAuth:         tls.RequireAndVerifyClientCert,
		ClientCAs:          caPool,
		ClientSessionCache: tls.NewLRUClientSessionCache(10),
		MinVersion:         tls.VersionTLS13,
		MaxVersion:         tls.VersionTLS13,
	}

	server := &httpsh.Server{
		Listener:    listener,
		Handler:     handler,
		Config:      config,
		Read:        viper.GetInt("read"),
		Write:       viper.GetInt("write"),
		Idle:        viper.GetInt("idle"),
		Certificate: viper.GetString("server_certificate"),
		Key:         viper.GetString("server_key"),
		Log:         logger,
	}

	err = server.Run()
	if err != nil {
		logger.Error("main", "message", err.Error())
		return
	}
}
