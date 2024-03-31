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
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
)

type Handler struct {
	BaseDirectory      string
	ContentType        string
	AllowedMethods     []string
	AllowedExecutables map[string][]string
	Log                *log.Logger
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	response := &Response{
		responseWriter: w,
		request:        r,
		contentType:    h.ContentType,
		allowedMethods: h.AllowedMethods,
		log:            h.Log,
	}

	err := os.Chdir(h.BaseDirectory)
	if err != nil {
		response.write(StatusBadRequest, err)
		return
	}

	if !slices.Contains(h.AllowedMethods, r.Method) {
		response.write(StatusMethodNotAllowed, ErrorMethodNotAllowed)
		return
	}

	if r.URL.Path != "/" {
		response.write(StatusForbidden, ErrorAccessDenied)
		return
	}

	queries := r.URL.Query()
	if len(queries) < 1 {
		response.write(StatusUnprocessableContent, ErrorQueryInvalid)
		return
	}

	if len(queries["e"]) != 1 {
		response.write(StatusUnprocessableContent, ErrorOneExecutableAllowed)
		return
	}

	options, ok := h.AllowedExecutables[queries["e"][0]]
	if !ok {
		response.write(StatusUnprocessableContent, ErrorExecutableNotFound)
		return
	}

	arguments := []string{}
	if len(queries["a"]) > 0 {
		for _, v := range queries["a"] {
			if len(v) < 3 {
				response.write(StatusUnprocessableContent, ErrorArgumentsInvalid)
				return
			}

			switch v[:2] {
			case "d_":
				directory := filepath.Join("./", v[2:])
				info, err := os.Stat(directory)
				if err != nil {
					response.write(StatusBadRequest, err)
					return
				}

				if !info.IsDir() {
					response.write(StatusUnprocessableContent, ErrorTargetNotDirectory)
					return
				}

				arguments = append(arguments, directory)
			case "f_":
				file := filepath.Join("./", v[2:])
				info, err := os.Stat(file)
				if err != nil {
					response.write(StatusBadRequest, err)
					return
				}

				if info.IsDir() {
					response.write(StatusUnprocessableContent, ErrorTargetNotFile)
					return
				}

				arguments = append(arguments, file)
			case "o_":
				if !slices.Contains(options, v[2:]) {
					response.write(StatusUnprocessableContent, ErrorOptionNotFound)
					return
				}

				arguments = append(arguments, v[2:])
			case "t_":
				if !strings.HasPrefix(v[2:], "'") || !strings.HasSuffix(v[2:], "'") {
					response.write(StatusUnprocessableContent, ErrorTextInvalid)
					return
				}

				arguments = append(arguments, v[2:])
			default:
				response.write(StatusUnprocessableContent, ErrorArgumentsInvalid)
				return
			}
		}
	}

	reader, writer := io.Pipe()

	defer writer.Close()

	tailer := strings.Join(arguments, " ")
	execution := fmt.Sprintf("%s %s", queries["e"][0], tailer)
	command := exec.Command("sh", "-c", execution)
	command.Stdout = writer
	command.Stderr = writer

	go response.copy(reader)

	command.Run()
}
