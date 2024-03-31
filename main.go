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
	"errors"

	"github.com/spf13/viper"
)

const (
	StatusBadRequest           int = 400
	StatusMethodNotAllowed     int = 405
	StatusUnprocessableContent int = 422
)

var (
	ErrorMethodNotAllowed     error = errors.New("method is not allowed")
	ErrorQueryInvalid         error = errors.New("query is invalid")
	ErrorOneExecutableAllowed error = errors.New("one executable allowed")
	ErrorExecutableNotFound   error = errors.New("executable is not found")
	ErrorArgumentsInvalid     error = errors.New("arguments are invalid")
	ErrorTargetNotDirectory   error = errors.New("target is not a directory")
	ErrorTargetNotFile        error = errors.New("target is not a file")
	ErrorOptionNotFound       error = errors.New("option is not found")
	ErrorTextInvalid          error = errors.New("text is invalid")
)

var (
	Network            string
	Host               string
	Port               string
	BaseDirectory      string
	ContentType        string
	AllowedMethods     []string
	AllowedExecutables map[string][]string
	ReadTimeout        int
	WriteTimeout       int
	IdleTimeout        int
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath("$HOME/.wsh/")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	Network = viper.GetString("network")
	Host = viper.GetString("host")
	Port = viper.GetString("port")
	BaseDirectory = viper.GetString("base_directory")
	ContentType = viper.GetString("content_type")
	AllowedMethods = viper.GetStringSlice("allowed_methods")
	AllowedExecutables = viper.GetStringMapStringSlice("allowed_executables")
	ReadTimeout = viper.GetInt("read_timeout")
	WriteTimeout = viper.GetInt("write_timeout")
	IdleTimeout = viper.GetInt("idle_timeout")
}
