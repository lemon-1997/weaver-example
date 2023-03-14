package api

import (
	"github.com/ServiceWeaver/weaver"
	"github.com/lemon-1997/weaver/service/category"
	"github.com/lemon-1997/weaver/service/product"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"net/http"
)

type Server struct {
	root     weaver.Instance
	product  product.T
	category category.T
}

func NewServer(root weaver.Instance) (*Server, error) {
	productSvc, err := weaver.Get[product.T](root)
	if err != nil {
		return nil, err
	}
	categorySvc, err := weaver.Get[category.T](root)
	if err != nil {
		return nil, err
	}
	s := &Server{
		root:     root,
		product:  productSvc,
		category: categorySvc,
	}
	s.registerHandler()
	return s, nil
}

func (s *Server) Run(addr string) error {
	lis, err := s.root.Listener("lemon", weaver.ListenerOptions{LocalAddress: addr})
	if err != nil {
		return err
	}
	s.root.Logger().Debug("listener available", "addr", lis)
	return http.Serve(lis, otelhttp.NewHandler(http.DefaultServeMux, "http"))
}
