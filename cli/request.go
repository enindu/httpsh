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
	"net"
	"strconv"
	"time"
)

type Request struct {
	bits        int
	years       int
	months      int
	days        int
	domains     []string
	emails      []string
	ips         []string
	key         string
	certificate string
}

func (r *Request) generate(b *Bundle) (*Bundle, error) {
	private, public, err := key(r.bits, r.key)
	if err != nil {
		return nil, err
	}

	serial, err := strconv.ParseInt(time.Now().Format("20060102150405"), 10, 0)
	if err != nil {
		return nil, err
	}

	ips := []net.IP{}
	for _, v := range r.ips {
		ips = append(ips, net.ParseIP(v))
	}

	template := &x509.Certificate{
		SignatureAlgorithm: x509.SHA512WithRSA,
		SerialNumber:       big.NewInt(serial),
		NotBefore:          time.Now(),
		NotAfter:           time.Now().AddDate(r.years, r.months, r.days),
		KeyUsage:           x509.KeyUsageDigitalSignature,
		ExtKeyUsage:        []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		DNSNames:           r.domains,
		EmailAddresses:     r.emails,
		IPAddresses:        ips,
	}

	certificate, err := certificate(r.certificate, template, b.certificate, b.public, b.private)
	if err != nil {
		return nil, err
	}

	return &Bundle{
		private:     private,
		public:      public,
		certificate: certificate,
	}, nil
}
