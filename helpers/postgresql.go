package helpers

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"math/rand"
	"net/url"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	postgresqlImage    = "docker.io/library/postgres:14.2"
	postgresqlPort     = "5432"
	postgresqlUser     = "eelbot"
	postgresqlPassword = "testpassword"
	postgresqlTimeout  = 10 * time.Second
)

// StartPostgreSQL starts and returns the URL to connect to a PostgreSQL instance. It first checks if the
// "TEST_POSTGRESQL_URL" environment variable is set, and if it is, it just returns that (useful for connecting to a
// locally running instance for debugging purposes). Otherwise, it spins up a new PostgreSQL instance via docker and
// returns the URL to connect to it.
func StartPostgreSQL() (uri string, close CloseFunc, err error) {
	// Check the "TEST_POSTGRESQL_URL" environment variable first.
	if u, ok := os.LookupEnv("TEST_POSTGRESQL_URL"); ok {
		uri, close = u, func() {}
		return
	}

	// Start a PostgreSQL instance via docker.
	var cli *client.Client
	if cli, err = client.NewClientWithOpts(client.FromEnv); err != nil {
		return
	}

	ctx := context.Background()
	var reader io.ReadCloser
	if reader, err = cli.ImagePull(ctx, postgresqlImage, types.ImagePullOptions{}); err != nil {
		return
	}
	if _, err = io.Copy(os.Stderr, reader); err != nil {
		return
	}
	_ = reader.Close()

	var instance container.ContainerCreateCreatedBody
	instance, err = cli.ContainerCreate(
		ctx,
		&container.Config{
			Image: postgresqlImage,
			Env: []string{
				fmt.Sprintf("POSTGRES_USER=%s", postgresqlUser),
				fmt.Sprintf("POSTGRES_PASSWORD=%s", postgresqlPassword),
			},
			ExposedPorts: nat.PortSet{nat.Port(fmt.Sprintf("%s/tcp", postgresqlPort)): struct{}{}},
		},
		&container.HostConfig{
			AutoRemove: true,
		},
		nil,
		nil,
		fmt.Sprintf("eelbot-test-postgresql-%d", rand.Int()),
	)
	if err != nil {
		return
	}

	if err = cli.ContainerStart(ctx, instance.ID, types.ContainerStartOptions{}); err != nil {
		_ = cli.ContainerRemove(ctx, instance.ID, types.ContainerRemoveOptions{Force: true})
		return
	}

	var details types.ContainerJSON
	if details, err = cli.ContainerInspect(ctx, instance.ID); err != nil {
		_ = cli.ContainerStop(ctx, instance.ID, nil)
		return
	}

	uri = (&url.URL{
		Scheme: "postgresql",
		Host:   fmt.Sprintf("%s:%s", details.NetworkSettings.IPAddress, postgresqlPort),
		Path:   postgresqlUser,
		User:   url.UserPassword(postgresqlUser, postgresqlPassword),
	}).String()

	if err = waitForPostgresql(uri); err != nil {
		_ = cli.ContainerStop(ctx, instance.ID, nil)
		return
	}

	close = func() {
		_ = cli.ContainerStop(ctx, instance.ID, nil)
	}
	return
}

func waitForPostgresql(uri string) error {
	db, err := sql.Open("pgx", uri)
	if err != nil {
		return err
	}

	giveup := time.Now().Add(postgresqlTimeout)
	for {
		if err = pingWithTimeout(db, 100*time.Millisecond); err == nil {
			return nil
		}
		if time.Now().After(giveup) {
			return err
		}
	}
}
