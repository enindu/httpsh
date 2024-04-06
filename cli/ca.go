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
	"math/big"
	"os"
	"strconv"
	"time"
)

type CA struct {
	bits        int
	years       int
	months      int
	days        int
	key         string
	certificate string
}

func (c *CA) generate() error {
	private, public, err := c.generateKey()
	if err != nil {
		return err
	}

	return c.generateCertificate(private, public)
}

func (c *CA) generateCertificate(private *rsa.PrivateKey, public *rsa.PublicKey) error {
	serial, err := strconv.ParseInt(time.Now().Format("20060102150405"), 10, 0)
	if err != nil {
		return err
	}

	template := &x509.Certificate{
		SerialNumber:          big.NewInt(serial),
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(c.years, c.months, c.days),
		KeyUsage:              x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	certificate, err := x509.CreateCertificate(rand.Reader, template, template, public, private)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(c.certificate, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	block := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certificate,
	}

	return pem.Encode(file, block)
}

func (c *CA) generateKey() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	key, err := rsa.GenerateKey(rand.Reader, c.bits)
	if err != nil {
		return nil, nil, err
	}

	file, err := os.OpenFile(c.key, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
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
