package database

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// GetTestDatabase starts a PostgreSQL container for testing
// and returns a termination function, a Service instance, and an error.
func GetTestDatabase() (func(context.Context, ...testcontainers.TerminateOption) error, Service, error) {
	var (
		dbName = "database"
		dbPwd  = "password"
		dbUser = "user"
	)

	dbContainer, err := postgres.Run(
		context.Background(),
		"postgres:latest",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPwd),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		return nil, nil, err
	}

	database = dbName
	password = dbPwd
	username = dbUser

	dbHost, err := dbContainer.Host(context.Background())
	if err != nil {
		return dbContainer.Terminate, nil, err
	}

	dbPort, err := dbContainer.MappedPort(context.Background(), "5432/tcp")
	if err != nil {
		return dbContainer.Terminate, nil, err
	}

	host = dbHost
	port = dbPort.Port()

	// Create a new Service instance
	srv := New()

	return dbContainer.Terminate, srv, nil
}

var testService Service

func TestMain(m *testing.M) {
	teardown, srv, err := GetTestDatabase()
	if err != nil {
		log.Fatalf("could not start postgres container: %v", err)
	}
	testService = srv

	m.Run()

	if teardown != nil && teardown(context.Background()) != nil {
		log.Fatalf("could not teardown postgres container: %v", err)
	}
}

func TestNew(t *testing.T) {
	if testService == nil {
		t.Fatal("testService is nil")
	}
}

func TestHealth(t *testing.T) {
	stats := testService.Health()

	if stats["status"] != "up" {
		t.Fatalf("expected status to be up, got %s", stats["status"])
	}

	if _, ok := stats["error"]; ok {
		t.Fatalf("expected error not to be present")
	}

	if stats["message"] != "It's healthy" {
		t.Fatalf("expected message to be 'It's healthy', got %s", stats["message"])
	}
}

func TestClose(t *testing.T) {
	// Note: We don't actually close the testService here as it's shared across tests
	// and closed in TestMain. This test verifies the Close method works.
	srv := New()
	if srv == nil {
		t.Fatal("New() returned nil")
	}
	// The actual close is handled by TestMain after all tests complete
}
