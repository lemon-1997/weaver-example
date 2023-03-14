package api

import (
	"encoding/json"
	"github.com/lemon-1997/weaver/service/category"
	"github.com/lemon-1997/weaver/service/product"
	"io"
	"net/http"
	"strconv"
)

type Product struct {
	ID           int64   `json:"id"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	Price        float64 `json:"price"`
	CategoryId   int64   `json:"category_id"`
	CategoryName string  `json:"category_name"`
}

type Category struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (s *Server) getProduct(w http.ResponseWriter, r *http.Request) {
	var ids []int64
	for _, item := range r.URL.Query()["id"] {
		id, err := strconv.ParseInt(item, 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		ids = append(ids, id)
	}
	list, err := s.product.List(r.Context(), ids)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	productList := transProductList(list)
	for _, item := range productList {
		cate, err := s.category.Get(r.Context(), item.CategoryId)
		if err != nil {
			// Todo log
			//http.Error(w, err.Error(), http.StatusBadRequest)
		}
		item.CategoryName = cate.Name
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(productList)
}

func (s *Server) createProduct(w http.ResponseWriter, r *http.Request) {
	var req Product
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if _, err := s.category.Get(r.Context(), req.CategoryId); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id, err := s.product.Create(r.Context(), product.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		CategoryId:  req.CategoryId,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	type ID struct {
		Id int64
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ID{Id: id})
}

func (s *Server) updateProduct(w http.ResponseWriter, r *http.Request) {
	var req Product
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if _, err := s.category.Get(r.Context(), req.CategoryId); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := s.product.Update(r.Context(), req.ID, product.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		CategoryId:  req.CategoryId,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, "ok")
}

func (s *Server) deleteProduct(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = s.product.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, "ok")
}

func (s *Server) getCategory(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	cate, err := s.category.Get(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Category{ID: cate.ID, Name: cate.Name})
}

func (s *Server) createCategory(w http.ResponseWriter, r *http.Request) {
	var req Category
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := s.category.Create(r.Context(), category.Category{
		ID:   req.ID,
		Name: req.Name,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, "ok")
}

func (s *Server) updateCategory(w http.ResponseWriter, r *http.Request) {
	var req Category
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := s.category.Update(r.Context(), req.ID, category.Category{
		ID:   req.ID,
		Name: req.Name,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, "ok")
}

func transProductList(list []product.Product) (res []*Product) {
	for _, item := range list {
		res = append(res, &Product{
			Name:        item.Name,
			Description: item.Description,
			Price:       item.Price,
			CategoryId:  item.CategoryId,
		})
	}
	return
}
