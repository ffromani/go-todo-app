package store

import (
	"context"
	"fmt"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// exercise

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
		t.Fatalf("failed to initialize the storage", err)
	}

	t.Run("test create and get", func(t *testing.T) {
		fmt.Println("TODO add a test where we create an item and we verify it", storage)
	})

	t.Run("test duplicateid fails", func(t *testing.T) {
		fmt.Println("TODO add a test where we create a duplicate id fails", storage)
	})
}
