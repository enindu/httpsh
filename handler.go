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
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
)

type Handler struct {
	Commands  map[string][]string
	Directory string
	Methods   []string
	Mime      string
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	response := &Response{
		methods: h.Methods,
		mime:    h.Mime,
		writer:  w,
	}

	err := os.Chdir(h.Directory)
	if err != nil {
		response.send(StatusBadRequest, err)
		return
	}

	if !slices.Contains(h.Methods, r.Method) {
		response.send(StatusMethodNotAllowed, ErrorMethodNotAllowed)
		return
	}

	queries := r.URL.Query()
	if len(queries) < 1 {
		response.send(StatusUnprocessableContent, ErrorQueryInvalid)
		return
	}

	if len(queries["e"]) != 1 {
		response.send(StatusUnprocessableContent, ErrorOneExecutableAllowed)
		return
	}

	options, ok := h.Commands[queries["e"][0]]
	if !ok {
		response.send(StatusUnprocessableContent, ErrorExecutableNotFound)
		return
	}

	arguments := []string{}

	if len(queries["a"]) > 0 {
		for _, v := range queries["a"] {
			if len(v) < 3 {
				response.send(StatusUnprocessableContent, ErrorArgumentsInvalid)
				return
			}

			switch v[:2] {
			case "d|":
				directory := filepath.Join("./", v[2:])

				info, err := os.Stat(directory)
				if err != nil {
					response.send(StatusBadRequest, err)
					return
				}

				if !info.IsDir() {
					response.send(StatusUnprocessableContent, ErrorTargetNotDirectory)
					return
				}

				arguments = append(arguments, directory)
			case "f|":
				file := filepath.Join("./", v[2:])

				info, err := os.Stat(file)
				if err != nil {
					response.send(StatusBadRequest, err)
					return
				}

				if info.IsDir() {
					response.send(StatusUnprocessableContent, ErrorTargetNotFile)
					return
				}

				arguments = append(arguments, file)
			case "o|":
				if !slices.Contains(options, v[2:]) {
					response.send(StatusUnprocessableContent, ErrorOptionNotFound)
					return
				}

				arguments = append(arguments, v[2:])
			case "t|":
				if !strings.HasPrefix(v[2:], "'") || !strings.HasSuffix(v[2:], "'") {
					response.send(StatusUnprocessableContent, ErrorTextInvalid)
					return
				}

				arguments = append(arguments, v[2:])
			default:
				response.send(StatusUnprocessableContent, ErrorArgumentsInvalid)
				return
			}
		}
	}

	reader, writer := io.Pipe()

	defer writer.Close()

	tailer := strings.Join(arguments, " ")
	line := fmt.Sprintf("%s %s", queries["e"][0], tailer)
	command := exec.Command("sh", "-c", line)

	command.Stdout = writer
	command.Stderr = writer

	go write(w, reader)

	err = command.Run()
	if err != nil {
		response.send(StatusBadRequest, err)
		return
	}
}

func write(w http.ResponseWriter, r *io.PipeReader) {
	io.Copy(w, r)
	r.Close()
}
