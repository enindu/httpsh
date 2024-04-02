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
	Listener *net.TCPListener
	Handler  *Handler
	Read     int
	Write    int
	Idle     int
	Log      *log.Logger
}

func (s *Server) Run() error {
	server := &http.Server{
		Handler:           s.Handler,
		ReadTimeout:       time.Duration(s.Read) * time.Second,
		ReadHeaderTimeout: 0,
		WriteTimeout:      time.Duration(s.Write) * time.Second,
		IdleTimeout:       time.Duration(s.Idle) * time.Second,
		ErrorLog:          s.Log,
	}

	go server.Serve(s.Listener)

	s.Log.Printf("server is started")
	s.wait()
	return s.stop(server)
}

func (s *Server) wait() {
	wait := make(chan os.Signal, 1)

	signal.Notify(wait, syscall.SIGINT)
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

	s.Log.Printf("server is stopped")
	return nil
}
