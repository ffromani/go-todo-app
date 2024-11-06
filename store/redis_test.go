package store

import (
	"context"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestWithRedis(t *testing.T) {
	req := testcontainers.ContainerRequest{
		Image:        "redis:latest",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections"),
	}
	redisC, _ := testcontainers.GenericContainer(context.Background(),
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})

	mapped, _ := redisC.MappedPort(context.Background(), "6379/tcp")

	t.Cleanup(func() { redisC.Terminate(context.Background()) })

	storage, err := NewRedis("127.0.0.1:"+mapped.Port(), "", 0)
	if err != nil {
		t.Fatal("failed to initialize the storage", err)
	}

	t.Run("test create and get", func(t *testing.T) {
		id := ID("foo")
		data := []byte("bar")

		err := storage.Create(id, data)
		if err != nil {
			t.Fatal("got error while creating", err)
		}
		retrieved, err := storage.Load(id)
		if err != nil {
			t.Fatal("got error while loading", err)
		}
		if string(retrieved) != "bar" {
			t.Fatalf("expecting bar, got %s", retrieved)
		}
	})

	t.Run("test duplicateid fails", func(t *testing.T) {
		id := ID("toDupe")
		data := []byte("bar")

		err := storage.Create(id, data)
		if err != nil {
			t.Fatal("got error while creating", err)
		}

		id1 := ID("toDupe")
		data1 := []byte("foo")

		err = storage.Create(id1, data1)
		if err == nil {
			t.Fatal("expecting error got nil")
		}
	})
}
