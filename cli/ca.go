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
	"crypto/x509"
	"math/big"
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

func (c *CA) generate() (*Bundle, error) {
	private, public, err := key(c.bits, c.key)
	if err != nil {
		return nil, err
	}

	serial, err := strconv.ParseInt(time.Now().Format("20060102150405"), 10, 0)
	if err != nil {
		return nil, err
	}

	template := &x509.Certificate{
		SignatureAlgorithm:    x509.SHA512WithRSA,
		SerialNumber:          big.NewInt(serial),
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(c.years, c.months, c.days),
		KeyUsage:              x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	certificate, err := certificate(c.certificate, template, template, public, private)
	if err != nil {
		return nil, err
	}

	return &Bundle{
		private:     private,
		public:      public,
		certificate: certificate,
	}, nil
}
