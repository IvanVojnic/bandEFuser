package repository

import (
	"context"
	"fmt"
	"github.com/IvanVojnic/bandEFuser/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ory/dockertest/v3"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"os"
	"os/exec"
	"testing"
)

var db *pgxpool.Pool

var testUserValidData = models.User{
	Name:     `User`,
	Email:    `user@gmail.com`,
	Password: `user123`,
	ID:       uuid.New(),
}

var testRequestValidData = models.Request{
	SenderID:   uuid.New(),
	ReceiverID: uuid.New(),
}

var testNotValidData = models.User{
	Name:     ``,
	Email:    `aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa`,
	Password: ``,
}

// TestMain used to test main func
func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("unix:///home/ivanvoynich/.docker/desktop/docker.sock")
	if err != nil {
		logrus.Fatalf("Could not construct pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		logrus.Fatalf("Could not connect to Docker: %s", err)
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "latest",
		Env: []string{
			"POSTGRES_DB=postgres",
			"POSTGRES_USER=postgres",
			"POSTGRES_PASSWORD=postgres",
		},
	})
	if err != nil {
		logrus.Fatalf("Could not start resource: %s", err)
	}

	hostAndPort := resource.GetHostPort("5432/tcp")
	databaseURL := fmt.Sprintf("postgresql://postgres:postgres@%s/postgres", hostAndPort)

	if errConnWrap := pool.Retry(func() error {
		var errConn error
		db, errConn = pgxpool.New(context.Background(), databaseURL)
		if errConn != nil {
			return errConn
		}
		return db.Ping(context.Background())
	}); errConnWrap != nil {
		logrus.Fatalf("Could not connect to database: %s", errConnWrap)
	}

	cmd := exec.Command("flyway", "-user=postgres", "-password=postgres",
		"-locations=filesystem:../../migrations/sql",
		fmt.Sprintf("-url=jdbc:postgresql://%v/postgres", hostAndPort), "migrate")
	err = cmd.Run()
	if err != nil {
		logrus.Fatalf("can't run flyway: %s", err)
	}
	code := m.Run()

	if errR1 := pool.Purge(resource); errR1 != nil {
		logrus.Fatalf("Could not purge resource: %s", errR1)
	}

	os.Exit(code)
}

// Test_GetFriends used to get friends
func Test_GetFriends(t *testing.T) {
	ctx := context.Background()
	repos := NewUserCommPostgres(db)
	_, err := repos.GetFriends(ctx, testUserValidData.ID)
	require.NoError(t, err, "get friends error")
}

// Test_SendFriendsRequest used to send request
func Test_SendFriendsRequest(t *testing.T) {
	ctx := context.Background()
	userSenderID := uuid.New()
	userReceiverID := uuid.New()
	sender := models.User{Name: "u1", Email: "u1@g.c", Password: "u1u1", ID: userSenderID}
	receiver := models.User{Name: "u2", Email: "u2@g.c", Password: "u2u2", ID: userReceiverID}
	repos := NewUserPostgres(db)
	errSender := repos.SignUp(ctx, &sender)
	errReceiver := repos.SignUp(ctx, &receiver)
	require.NoError(t, errSender, "create user error")
	require.NoError(t, errReceiver, "create user error")
	reposSend := NewUserCommPostgres(db)
	err := reposSend.SendFriendsRequest(ctx, userSenderID, userReceiverID)
	require.NoError(t, err, "send request error")
}

// Test_AcceptFriendsRequest used to accept request
func Test_AcceptFriendsRequest(t *testing.T) {
	ctx := context.Background()
	repos := NewUserCommPostgres(db)
	err := repos.AcceptFriendsRequest(ctx, testRequestValidData.SenderID, testRequestValidData.ReceiverID)
	require.NoError(t, err, "accept request error")
}

// Test_DeclineFriendsRequest used to decline request
func Test_DeclineFriendsRequest(t *testing.T) {
	ctx := context.Background()
	repos := NewUserCommPostgres(db)
	err := repos.DeclineFriendsRequest(ctx, testRequestValidData.SenderID, testRequestValidData.ReceiverID)
	require.NoError(t, err, "decline request error")
}

// Test_FindUser used to find user
func Test_FindUser(t *testing.T) {
	ctx := context.Background()
	repos := NewUserPostgres(db)
	errSU := repos.SignUp(ctx, &testUserValidData)
	require.NoError(t, errSU, "create user error")
	reposComm := NewUserCommPostgres(db)
	_, err := reposComm.FindUser(ctx, testUserValidData.Email)
	require.NoError(t, err, "find user error")
}

// Test_GetRequest used to get requests
func Test_GetRequest(t *testing.T) {
	ctx := context.Background()
	repos := NewUserCommPostgres(db)
	_, err := repos.GetRequest(ctx, testUserValidData.ID)
	require.NoError(t, err, "get friends error")
}
