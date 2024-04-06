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
	"os"
)

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
		Type:  "KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}

	err = pem.Encode(file, block)
	if err != nil {
		return nil, nil, err
	}

	return key, &key.PublicKey, nil
}
