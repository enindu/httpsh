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
	"errors"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"sync"
)

type Handler struct {
	Directory   string
	Mime        string
	Methods     []string
	Executables map[string][]string
	Log         *log.Logger
	mutex       sync.Mutex
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	response := &Response{
		writer:  w,
		request: r,
		mime:    h.Mime,
		methods: h.Methods,
		log:     h.Log,
	}

	err := os.Chdir(h.Directory)
	if err != nil {
		response.error(http.StatusBadRequest, errChangeDirectory)
		return
	}

	if !slices.Contains(h.Methods, r.Method) {
		response.error(http.StatusMethodNotAllowed, errMethodNotAllowed)
		return
	}

	if r.URL.Path != "/" {
		response.error(http.StatusForbidden, errAccessDenied)
		return
	}

	queries := r.URL.Query()
	if len(queries) < 1 {
		response.error(http.StatusBadRequest, errQueryInvalid)
		return
	}

	executable, options, err := h.program(queries)
	if err != nil {
		response.error(http.StatusBadRequest, err)
		return
	}

	arguments, err := h.arguments(queries, options)
	if err != nil {
		response.error(http.StatusBadRequest, err)
		return
	}

	line := h.line(executable, arguments)
	stdout, stderr, err := h.command(line)
	if err != nil {
		response.error(http.StatusBadRequest, errors.New(stderr))
		return
	}

	response.write(http.StatusOK, stdout)
}

func (h *Handler) command(l string) (string, string, error) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	stdout := &bytes.Buffer{}
	defer stdout.Reset()

	stderr := &bytes.Buffer{}
	defer stderr.Reset()

	command := exec.Command("sh", "-c", l)
	command.Stdout = stdout
	command.Stderr = stderr

	err := command.Run()
	if err != nil {
		return "", stderr.String(), err
	}

	return stdout.String(), "", nil
}

func (h *Handler) line(e string, a []string) string {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	buffer := bytes.Buffer{}
	defer buffer.Reset()

	buffer.WriteString(e)
	buffer.WriteString(" ")

	for _, v := range a {
		buffer.WriteString(v)
		buffer.WriteString(" ")
	}

	return buffer.String()
}

func (h *Handler) arguments(q map[string][]string, o []string) (a []string, e error) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if len(q["a"]) > 0 {
		for _, v := range q["a"] {
			if len(v) < 3 {
				return nil, errArgumentsInvalid
			}

			switch v[:2] {
			case "d_":
				directory := filepath.Join("./", v[2:])
				info, err := os.Stat(directory)
				if err != nil {
					return nil, errTargetNotFound
				}

				if !info.IsDir() {
					return nil, errTargetNotDirectory
				}

				a = append(a, directory)
			case "f_":
				file := filepath.Join("./", v[2:])
				info, err := os.Stat(file)
				if err != nil {
					return nil, errTargetNotFound
				}

				if info.IsDir() {
					return nil, errTargetNotFile
				}

				a = append(a, file)
			case "o_":
				if !slices.Contains(o, v[2:]) {
					return nil, errOptionNotFound
				}

				a = append(a, v[2:])
			case "t_":
				if !strings.HasPrefix(v[2:], "'") || !strings.HasSuffix(v[2:], "'") {
					return nil, errTextInvalid
				}

				a = append(a, v[2:])
			default:
				return nil, errArgumentsInvalid
			}
		}
	}

	return a, nil
}

func (h *Handler) program(q map[string][]string) (string, []string, error) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if len(q["e"]) != 1 {
		return "", nil, errOneExecutableAllowed
	}

	options, ok := h.Executables[q["e"][0]]
	if !ok {
		return "", nil, errExecutableNotFound
	}

	return q["e"][0], options, nil
}
