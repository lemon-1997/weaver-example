package product

import (
	"context"
	"github.com/ServiceWeaver/weaver"
	"github.com/ServiceWeaver/weaver/weavertest"
	"testing"
)

func TestProduct(t *testing.T) {
	ctx := context.Background()
	root := weavertest.Init(ctx, t, weavertest.Options{
		SingleProcess: true,
		Config: `
			["lemon\\service\\product\\T"]
			dsn = "test.db"`,
	})
	service, err := weaver.Get[T](root)
	if err != nil {
		t.Fatal(err)
	}

	id, err := service.Create(context.Background(), Product{
		Name:        "create",
		Description: "nothing",
		Price:       10,
		CategoryId:  1,
	})
	if err != nil {
		t.Fatal(err)
	}

	if err = service.Update(context.Background(), id, Product{
		Name:        "update",
		Description: "nothing",
		Price:       100,
		CategoryId:  2,
	}); err != nil {
		t.Fatal(err)
	}

	list, err := service.List(context.Background(), []int64{id})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(list)

	if err = service.Delete(context.Background(), id); err != nil {
		t.Fatal(err)
	}
}
