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
)

type Server struct {
	network           string
	host              string
	port              string
	domain            string
	readTimeout       int
	writeTimeout      int
	idleTimeout       int
	caCertificate     string
	serverCertificate string
	serverKey         string
	directory         string
	mime              string
	methods           []string
	executables       map[string][]string
	log               *slog.Logger
}

func (s *Server) run() error {
	socket := &httpsh.Socket{
		Network: s.network,
		Host:    s.host,
		Port:    s.port,
		Log:     s.log,
	}

	listener, err := socket.Listen()
	if err != nil {
		return err
	}

	defer listener.Close()

	handler := &httpsh.Handler{
		Directory:   s.directory,
		Mime:        s.mime,
		Methods:     s.methods,
		Executables: s.executables,
		Log:         s.log,
	}

	caCertificate, err := os.ReadFile(s.caCertificate)
	if err != nil {
		return err
	}

	caPool := x509.NewCertPool()

	caPool.AppendCertsFromPEM(caCertificate)

	tls := &tls.Config{
		RootCAs:            caPool,
		ServerName:         s.domain,
		ClientAuth:         tls.RequireAndVerifyClientCert,
		ClientCAs:          caPool,
		ClientSessionCache: tls.NewLRUClientSessionCache(10),
		MinVersion:         tls.VersionTLS13,
		MaxVersion:         tls.VersionTLS13,
	}

	server := &httpsh.Server{
		Listener:     listener,
		Handler:      handler,
		TLS:          tls,
		ReadTimeout:  s.readTimeout,
		WriteTimeout: s.writeTimeout,
		IdleTimeout:  s.idleTimeout,
		Certificate:  s.serverCertificate,
		Key:          s.serverKey,
		Log:          s.log,
	}

	err = server.Run()
	if err != nil {
		return err
	}

	return nil
}
