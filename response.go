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
	"io"
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

func (r *Response) write(c int, e error) (int, error) {
	if c == http.StatusMethodNotAllowed {
		allow := strings.Join(r.methods, ",")

		r.writer.Header().Set("Allow", allow)
	}

	body := e.Error()

	r.log.Println(r.request.RemoteAddr, r.request.RequestURI, body)
	r.writer.Header().Set("Content-Type", r.mime)
	r.writer.WriteHeader(c)

	return r.writer.Write([]byte(body))
}

func (r *Response) copy(reader *io.PipeReader) (int64, error) {
	length, err := io.Copy(r.writer, reader)
	if err != nil {
		return length, err
	}

	return length, reader.Close()
}
