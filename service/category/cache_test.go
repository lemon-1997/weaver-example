package category

import (
	"context"
	"github.com/ServiceWeaver/weaver"
	"github.com/ServiceWeaver/weaver/weavertest"
	"reflect"
	"testing"
)

func TestCache(t *testing.T) {
	ctx := context.Background()
	root := weavertest.Init(ctx, t, weavertest.Options{SingleProcess: true})
	cache, err := weaver.Get[categoryCache](root)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		id      int64
		want    Category
		wantErr bool
		fun     func(c categoryCache) error
	}{
		{
			name:    "Add",
			id:      1,
			want:    Category{ID: 1, Name: "1"},
			wantErr: false,
			fun: func(c categoryCache) error {
				return cache.Add(ctx, 1, Category{ID: 1, Name: "1"})
			},
		},
		{
			name:    "Update",
			id:      2,
			want:    Category{ID: 2, Name: "2"},
			wantErr: false,
			fun: func(c categoryCache) error {
				if err = cache.Add(ctx, 2, Category{ID: 2, Name: "1"}); err != nil {
					return err
				}
				if err = cache.Add(ctx, 2, Category{ID: 2, Name: "2"}); err != nil {
					return err
				}
				return nil
			},
		},
		{
			name:    "Remove",
			id:      1,
			wantErr: true,
			fun: func(c categoryCache) error {
				return cache.Remove(ctx, 1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err = tt.fun(cache); err != nil {
				t.Fatal(err)
			}
			got, err := cache.Get(context.Background(), tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}
