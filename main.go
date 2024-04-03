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
	"errors"

	"github.com/spf13/viper"
)

var (
	errChangeDirectory      error = errors.New("can not change directory")
	errMethodNotAllowed     error = errors.New("method is not allowed")
	errAccessDenied         error = errors.New("access is denied")
	errQueryInvalid         error = errors.New("query is invalid")
	errOneExecutableAllowed error = errors.New("one executable allowed")
	errExecutableNotFound   error = errors.New("executable is not found")
	errArgumentsInvalid     error = errors.New("arguments are invalid")
	errTargetNotFound       error = errors.New("target is not found")
	errTargetNotDirectory   error = errors.New("target is not a directory")
	errTargetNotFile        error = errors.New("target is not a file")
	errOptionNotFound       error = errors.New("option is not found")
	errTextInvalid          error = errors.New("text is invalid")
	errUnknown              error = errors.New("unknown error")
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.httpsh/")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
