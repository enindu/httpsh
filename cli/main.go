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
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"log/slog"
	"os"

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
	viper.SetConfigType("toml")
	viper.AddConfigPath("$HOME/.config/httpsh/cli")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		log.Error(err.Error())
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
			log.Error(err.Error())
			return
		}

		server := &Request{
			bits:         viper.GetInt("cert.request.server.bits"),
			country:      viper.GetStringSlice("cert.request.server.country"),
			organization: viper.GetStringSlice("cert.request.server.organization"),
			unit:         viper.GetStringSlice("cert.request.server.unit"),
			locality:     viper.GetStringSlice("cert.request.server.locality"),
			province:     viper.GetStringSlice("cert.request.server.province"),
			domain:       viper.GetString("cert.request.server.domain"),
			years:        viper.GetInt("cert.request.server.years"),
			months:       viper.GetInt("cert.request.server.months"),
			days:         viper.GetInt("cert.request.server.days"),
			key:          viper.GetString("cert.request.server.key"),
			certificate:  viper.GetString("cert.request.server.certificate"),
		}

		_, err = server.generate(bundle)
		if err != nil {
			log.Error(err.Error())
			return
		}

		client := &Request{
			bits:         viper.GetInt("cert.request.client.bits"),
			country:      viper.GetStringSlice("cert.request.client.country"),
			organization: viper.GetStringSlice("cert.request.client.organization"),
			unit:         viper.GetStringSlice("cert.request.client.unit"),
			locality:     viper.GetStringSlice("cert.request.client.locality"),
			province:     viper.GetStringSlice("cert.request.client.province"),
			domain:       viper.GetString("cert.request.client.domain"),
			years:        viper.GetInt("cert.request.client.years"),
			months:       viper.GetInt("cert.request.client.months"),
			days:         viper.GetInt("cert.request.client.days"),
			key:          viper.GetString("cert.request.client.key"),
			certificate:  viper.GetString("cert.request.client.certificate"),
		}

		_, err = client.generate(bundle)
		if err != nil {
			log.Error(err.Error())
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
			log.Error(err.Error())
			return
		}
	default:
		flag.PrintDefaults()
	}
}

func certificate(f string, s *x509.Certificate, c *x509.Certificate, public *rsa.PublicKey, private *rsa.PrivateKey) (*x509.Certificate, error) {
	certificate, err := x509.CreateCertificate(rand.Reader, s, c, public, private)
	if err != nil {
		return nil, err
	}

	file, err := os.OpenFile(f, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	block := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certificate,
	}

	err = pem.Encode(file, block)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func key(b int, f string) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	key, err := rsa.GenerateKey(rand.Reader, b)
	if err != nil {
		return nil, nil, err
	}

	file, err := os.OpenFile(f, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return nil, nil, err
	}

	defer file.Close()

	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}

	err = pem.Encode(file, block)
	if err != nil {
		return nil, nil, err
	}

	return key, &key.PublicKey, nil
}
