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
	"net/http"
	"strings"
)

type Response struct {
	methods []string
	mime    string
	writer  http.ResponseWriter
}

func (r *Response) write(c int, e error) {
	if c == StatusMethodNotAllowed {
		allow := strings.Join(r.methods, ",")
		r.writer.Header().Set("Allow", allow)
	}

	body := e.Error()

	r.writer.Header().Set("Content-Type", r.mime)
	r.writer.WriteHeader(c)
	r.writer.Write([]byte(body))
}
