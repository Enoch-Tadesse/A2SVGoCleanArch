
package infrastructure_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/A2SVTask7/Infrastructure"
	"github.com/stretchr/testify/require"
)

func TestInitMongo_Success(t *testing.T) {
	// Set environment variable for DB_NAME if your function still uses it
	os.Setenv("DB_NAME", "testdb")
	defer os.Unsetenv("DB_NAME")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := infrastructure.InitMongo(ctx, "mongodb://localhost:27017", "testdb")
	require.NoError(t, err)
	require.NotNil(t, db)
	require.Equal(t, "testdb", db.Name())
}

func TestInitMongo_Failure(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Invalid URI should return error
	db, err := infrastructure.InitMongo(ctx, "mongodb://invalid:27017", "testdb")
	require.Error(t, err)
	require.Nil(t, db)
}
