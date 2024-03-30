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
	"bytes"
	"net/http"
)

type Response struct {
	contentType string
	buffer      *bytes.Buffer
	writer      http.ResponseWriter
}

func newResponse(t string, w http.ResponseWriter) *Response {
	w.Header().Set("Content-Type", t)

	return &Response{
		contentType: t,
		buffer:      &bytes.Buffer{},
		writer:      w,
	}
}

func (r *Response) writeError(e error) {
	body := e.Error()

	r.buffer.WriteString(body)

	response := r.buffer.Bytes()

	r.writer.WriteHeader(http.StatusBadRequest)
	r.writer.Write(response)
}
