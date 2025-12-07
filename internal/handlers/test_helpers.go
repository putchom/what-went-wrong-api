package handlers

import (
	"context"
	"fmt"
	"testing"
	"time"
	"what-went-wrong-api/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupTestDB(t *testing.T) (*gorm.DB, func()) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:15-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "test",
			"POSTGRES_PASSWORD": "test",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForAll(
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
			wait.ForListeningPort("5432/tcp"),
		).WithDeadline(60 * time.Second),
	}

	postgresContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	assert.NoError(t, err, "PostgreSQLコンテナの起動に失敗しました")

	host, err := postgresContainer.Host(ctx)
	assert.NoError(t, err)

	port, err := postgresContainer.MappedPort(ctx, "5432")
	assert.NoError(t, err)

	dsn := fmt.Sprintf("host=%s user=test password=test dbname=testdb port=%s sslmode=disable TimeZone=Asia/Tokyo", host, port.Port())

	// Retry connection up to 10 times with 1 second delay
	var db *gorm.DB
	for i := 0; i < 10; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			// Test the connection
			sqlDB, err := db.DB()
			if err == nil {
				err = sqlDB.Ping()
				if err == nil {
					break
				}
			}
		}
		t.Logf("Connection attempt %d failed, retrying...", i+1)
		time.Sleep(1 * time.Second)
	}
	assert.NoError(t, err, "テストデータベースへの接続に失敗しました")

	err = db.AutoMigrate(
		&models.Goal{},
		&models.ExcuseEntry{},
		&models.UserPlan{},
	)
	assert.NoError(t, err, "マイグレーションに失敗しました")

	cleanup := func() {
		postgresContainer.Terminate(ctx)
	}

	return db, cleanup
}
