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
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	Listener    *net.TCPListener
	Handler     *Handler
	Config      *tls.Config
	Read        int
	Write       int
	Idle        int
	Certificate string
	Key         string
	Log         *slog.Logger
}

func (s *Server) Run() error {
	server := &http.Server{
		Handler:           s.Handler,
		TLSConfig:         s.Config,
		ReadTimeout:       time.Duration(s.Read) * time.Second,
		ReadHeaderTimeout: 0,
		WriteTimeout:      time.Duration(s.Write) * time.Second,
		IdleTimeout:       time.Duration(s.Idle) * time.Second,
		ErrorLog:          slog.NewLogLogger(s.Log.Handler(), slog.LevelError),
	}

	go func() {
		err := server.ServeTLS(s.Listener, s.Certificate, s.Key)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.Log.Error("server.run", "message", err)
		}
	}()

	s.Log.Info("server.run")
	s.wait()
	return s.stop(server)
}

func (s *Server) wait() {
	wait := make(chan os.Signal, 1)

	signal.Notify(wait, syscall.SIGINT)
	s.Log.Info("server.wait")

	<-wait

	fmt.Printf("\r")
}

func (s *Server) stop(server *http.Server) error {
	stop, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	err := server.Shutdown(stop)
	if err != nil {
		return err
	}

	s.Log.Info("server.stop")
	return nil
}
