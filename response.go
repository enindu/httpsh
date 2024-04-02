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

package httpsh

import (
	"log/slog"
	"net/http"
	"strings"
)

type Response struct {
	writer  http.ResponseWriter
	request *http.Request
	mime    string
	methods []string
	log     *slog.Logger
}

func (r *Response) error(c int, e error) (int, error) {
	if e.Error() == "" {
		e = errUnknown
	}

	r.log.Error("invalid request", "address", r.request.RemoteAddr, "uri", r.request.RequestURI)
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
