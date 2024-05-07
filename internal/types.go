package internal

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
)

type Server struct {
	Handler    http.Handler
	ServerAddr int
	AdminAddr  int
}

func (s *Server) Start(ctx context.Context) error {
	eg := errgroup.Group{}
	eg.Go(func() error {
		return http.ListenAndServe(fmt.Sprintf(":%d", s.ServerAddr), s.Handler)
	})
	eg.Go(func() error {
		return http.ListenAndServe(fmt.Sprintf(":%d", s.AdminAddr),
			http.FileServer(http.Dir("./databases")))
	})
	eg.Go(func() error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		}
	})
	return eg.Wait()
}
