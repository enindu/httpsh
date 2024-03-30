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
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	Handler      *Handler
	IdleTimeout  int
	Listener     *net.TCPListener
	Log          *log.Logger
	ReadTimeout  int
	WriteTimeout int
}

func (s *Server) Run() error {
	server := &http.Server{
		ErrorLog:     s.Log,
		Handler:      s.Handler,
		IdleTimeout:  time.Duration(s.IdleTimeout) * time.Second,
		ReadTimeout:  time.Duration(s.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(s.WriteTimeout) * time.Second,
	}

	address := s.Listener.Addr().String()

	s.Log.Println("server is started on", address)

	go server.Serve(s.Listener)

	wait := make(chan os.Signal, 1)

	signal.Notify(wait, syscall.SIGINT)

	<-wait

	parent := context.Background()
	stop, cancel := context.WithTimeout(parent, time.Minute)

	defer cancel()

	err := server.Shutdown(stop)
	if err != nil {
		return err
	}

	fmt.Printf("\r")
	s.Log.Println("server is stopped")

	return nil
}
