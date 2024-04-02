// This file is part of wsh.
//
// wsh is free software: you can redistribute it and/or modify it under the
// terms of the GNU General Public License as published by the Free Software
// Foundation, either version 3 of the License, or (at your option) any later
// version.
//
// wsh is distributed in the hope that it will be useful, but WITHOUT ANY
// WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR
// A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with
// wsh. If not, see <https://www.gnu.org/licenses/>.

package wsh

import (
	"log"
	"net/http"
	"strings"
)

type Response struct {
	writer  http.ResponseWriter
	request *http.Request
	mime    string
	methods []string
	log     *log.Logger
}

func (r *Response) error(c int, e error) (int, error) {
	if e.Error() == "" {
		e = errUnknown
	}

	r.log.Printf("%q %q", r.request.RemoteAddr, r.request.RequestURI)
	return r.write(c, e.Error())
}

func (r *Response) write(c int, s string) (int, error) {
	if c == http.StatusMethodNotAllowed {
		r.writer.Header().Set("Allow", strings.Join(r.methods, ", "))
	}

	r.writer.Header().Set("Content-Type", r.mime)
	r.writer.WriteHeader(c)
	return r.writer.Write([]byte(s))
}
