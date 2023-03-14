package category

import (
	"context"
	"errors"
	"github.com/ServiceWeaver/weaver"
	"sync"
)

type categoryCache interface {
	Add(context.Context, int64, Category) error
	Get(context.Context, int64) (Category, error)
	Remove(context.Context, int64) error
}

type categoryCacheImpl struct {
	weaver.Implements[categoryCache]
	weaver.WithRouter[categoryCacheRouter]

	cache sync.Map
}

func (c *categoryCacheImpl) Add(_ context.Context, id int64, category Category) error {
	c.cache.Store(id, category)
	return nil
}

func (c *categoryCacheImpl) Get(_ context.Context, id int64) (Category, error) {
	value, ok := c.cache.Load(id)
	if !ok {
		return Category{}, errors.New("record not found")
	}
	cate, ok := value.(Category)
	if !ok {
		return Category{}, errors.New("data error")
	}
	return cate, nil
}

func (c *categoryCacheImpl) Remove(_ context.Context, id int64) error {
	c.cache.Delete(id)
	return nil
}

type categoryCacheRouter struct{}

func (categoryCacheRouter) Add(_ context.Context, key int64, _ Category) int64 { return key }
func (categoryCacheRouter) Get(_ context.Context, key int64) int64             { return key }
func (categoryCacheRouter) Remove(_ context.Context, key int64) int64          { return key }
