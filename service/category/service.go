package category

import (
	"context"
	"github.com/ServiceWeaver/weaver"
	"go.opentelemetry.io/otel/trace"
)

type Category struct {
	weaver.AutoMarshal
	ID   int64
	Name string
}

type T interface {
	Get(ctx context.Context, id int64) (Category, error)
	Create(ctx context.Context, category Category) error
	Update(ctx context.Context, id int64, category Category) error
}

type impl struct {
	weaver.Implements[T]
	cache categoryCache
}

func (s *impl) Init(context.Context) error {
	cache, err := weaver.Get[categoryCache](s)
	if err != nil {
		return err
	}
	s.cache = cache
	return nil
}

func (s *impl) LogWithTrace(ctx context.Context) weaver.Logger {
	span := trace.SpanFromContext(ctx)
	return s.Logger().With(
		"spanID", span.SpanContext().SpanID().String(),
		"traceID", span.SpanContext().TraceID().String())
}

func (s *impl) Get(ctx context.Context, id int64) (Category, error) {
	cate, err := s.cache.Get(ctx, id)
	if err != nil {
		s.LogWithTrace(ctx).Error("cache Get err", err, "id", id)
	}
	return cate, nil
}

func (s *impl) Create(ctx context.Context, category Category) error {
	if err := s.cache.Add(ctx, category.ID, category); err != nil {
		s.LogWithTrace(ctx).Error("cache Add err", err, "id", category.ID)
	}
	return nil
}

func (s *impl) Update(ctx context.Context, id int64, category Category) error {
	if err := s.cache.Add(ctx, id, category); err != nil {
		s.LogWithTrace(ctx).Error("cache Add err", err, "id", id)
	}
	return nil
}
