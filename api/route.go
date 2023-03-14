package api

import (
	"github.com/ServiceWeaver/weaver"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"net/http"
)

func (s *Server) registerHandler() {
	instrument := func(label string, fn func(http.ResponseWriter, *http.Request)) http.Handler {
		return weaver.InstrumentHandlerFunc(label, func(w http.ResponseWriter, r *http.Request) {
			span := trace.SpanFromContext(r.Context())
			span.SetAttributes(attribute.String("http.path", r.URL.Path))
			fn(w, r)
		})
	}
	http.Handle("/product", instrument("product", s.productHandler))
	http.Handle("/category", instrument("category", s.categoryHandler))
}

func (s *Server) productHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.getProduct(w, r)
	case http.MethodPost:
		s.createProduct(w, r)
	case http.MethodPut:
		s.updateProduct(w, r)
	case http.MethodDelete:
		s.deleteProduct(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) categoryHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.getCategory(w, r)
	case http.MethodPost:
		s.createCategory(w, r)
	case http.MethodPut:
		s.updateCategory(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
