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
	"os"
)

type Handler struct {
	ContentType string
	Directory   string
	Whitelist   map[string][]string
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	response := newResponse(h.ContentType, w)

	err := os.Chdir(h.Directory)
	if err != nil {
		response.writeError(err)
		return
	}
}
