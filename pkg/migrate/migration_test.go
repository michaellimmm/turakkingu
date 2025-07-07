package mongodb

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	dt "github.com/golang-migrate/migrate/v4/database/testing"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func waitForMongoReady(uri string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}
	defer func() {
		_ = client.Disconnect(ctx)
	}()

	return client.Ping(ctx, nil)
}

func TestRunMongoDB(t *testing.T) {
	container, host, port, err := setupMongoContainer("mongo:8.0")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = container.Terminate(context.TODO())
	}()

	uri := fmt.Sprintf("mongodb://%s:%s/testMigration", host, port)
	if err := waitForMongoReady(uri); err != nil {
		t.Fatal(err)
	}

	p := &Mongo{}
	d, err := p.Open(uri)
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()

	dt.TestNilVersion(t, d)
	dt.TestLockAndUnlock(t, d)
	dt.TestRun(t, d, bytes.NewReader([]byte(`[{"insert":"hello","documents":[{"wild":"world"}]}]`)))
	dt.TestSetVersion(t, d)
	dt.TestDrop(t, d)
}

func TestMigrate(t *testing.T) {
	container, host, port, err := setupMongoContainer("mongo:8.0")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = container.Terminate(context.TODO())
	}()

	uri := fmt.Sprintf("mongodb://%s:%s/testMigration", host, port)
	if err := waitForMongoReady(uri); err != nil {
		t.Fatal(err)
	}

	p := &Mongo{}
	d, err := p.Open(uri)
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()

	m, err := migrate.NewWithDatabaseInstance("file://./examples/migrations", "", d)
	if err != nil {
		t.Fatal(err)
	}
	dt.TestMigrate(t, m)
}

func TestWithAuth(t *testing.T) {
	container, host, port, err := setupMongoContainer("mongo:8.0")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = container.Terminate(context.TODO())
	}()

	uri := fmt.Sprintf("mongodb://%s:%s/testMigration", host, port)
	if err := waitForMongoReady(uri); err != nil {
		t.Fatal(err)
	}

	p := &Mongo{}
	d, err := p.Open(uri)
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()

	createUserCMD := []byte(`[{"createUser":"deminem","pwd":"gogo","roles":[{"role":"readWrite","db":"testMigration"}]}]`)
	err = d.Run(bytes.NewReader(createUserCMD))
	if err != nil {
		t.Fatal(err)
	}
	testcases := []struct {
		name            string
		connectUri      string
		isErrorExpected bool
	}{
		{"right auth data", "mongodb://deminem:gogo@%s:%v/testMigration", false},
		{"wrong auth data", "mongodb://wrong:auth@%s:%v/testMigration", true},
	}

	for _, tcase := range testcases {
		t.Run(tcase.name, func(t *testing.T) {
			mc := &Mongo{}
			d, err := mc.Open(fmt.Sprintf(tcase.connectUri, host, strings.ReplaceAll(port, "/tcp", "")))
			if err == nil {
				defer func() {
					if err := d.Close(); err != nil {
						t.Error(err)
					}
				}()
			}

			switch {
			case tcase.isErrorExpected && err == nil:
				t.Fatalf("no error when expected")
			case !tcase.isErrorExpected && err != nil:
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestLockWorks(t *testing.T) {
	container, host, port, err := setupMongoContainer("mongo:8.0")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = container.Terminate(context.TODO())
	}()

	uri := fmt.Sprintf("mongodb://%s:%s/testMigration", host, port)
	if err := waitForMongoReady(uri); err != nil {
		t.Fatal(err)
	}

	p := &Mongo{}
	d, err := p.Open(uri)
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()

	dt.TestRun(t, d, bytes.NewReader([]byte(`[{"insert":"hello","documents":[{"wild":"world"}]}]`)))

	mc := d.(*Mongo)

	err = mc.Lock()
	if err != nil {
		t.Fatal(err)
	}
	err = mc.Unlock()
	if err != nil {
		t.Fatal(err)
	}

	err = mc.Lock()
	if err != nil {
		t.Fatal(err)
	}
	err = mc.Unlock()
	if err != nil {
		t.Fatal(err)
	}

	// enable locking,
	//try to hit a lock conflict
	mc.config.Locking.Enabled = true
	mc.config.Locking.Timeout = 1
	err = mc.Lock()
	if err != nil {
		t.Fatal(err)
	}
	err = mc.Lock()
	if err == nil {
		t.Fatal("should have failed, mongo should be locked already")
	}
}

func TestTransaction(t *testing.T) {
	cmd := []string{"mongod", "--bind_ip_all", "--replSet", "rs0"}
	container, host, port, err := setupMongoContainer("mongo:4.0", cmd...)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = container.Terminate(context.TODO())
	}()

	uri := fmt.Sprintf("mongodb://%s:%s/testMigration?connect=direct", host, port)
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		t.Fatal(err)
	}

	// Run replSetInitiate
	err = client.Database("admin").RunCommand(context.TODO(), bson.D{{Key: "replSetInitiate", Value: bson.D{}}}).Err()
	if err != nil {
		t.Fatal(err)
	}

	// Wait for primary
	err = waitForPrimary(client)
	if err != nil {
		t.Fatal(err)
	}

	d, err := WithInstance(client, &Config{
		DatabaseName:    "testMigration",
		TransactionMode: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()

	insertCMD := []byte(`[
		{"create":"hello"},
		{"createIndexes": "hello",
			"indexes": [{
				"key": {"wild": 1},
				"name": "unique_wild",
				"unique": true,
				"background": true
			}]
		}]`)
	if err := d.Run(bytes.NewReader(insertCMD)); err != nil {
		t.Fatal(err)
	}

	testcases := []struct {
		name            string
		cmds            []byte
		documentsCount  int64
		isErrorExpected bool
	}{
		{
			name: "success transaction",
			cmds: []byte(`[{"insert":"hello","documents":[
				{"wild":"world"}, {"wild":"west"}, {"wild":"natural"}
			]}]`),
			isErrorExpected: false,
		},
		{
			name: "failure transaction",
			cmds: []byte(`[{"insert":"hello","documents":[
				{"wild":"flower"}, {"wild":"cat"}, {"wild":"west"}
			]}]`),
			isErrorExpected: true,
		},
	}

	for _, tcase := range testcases {
		t.Run(tcase.name, func(t *testing.T) {
			client, err := mongo.Connect(options.Client().ApplyURI(uri))
			if err != nil {
				t.Fatal(err)
			}
			if err := client.Ping(context.TODO(), nil); err != nil {
				t.Fatal(err)
			}

			d, err := WithInstance(client, &Config{
				DatabaseName:    "testMigration",
				TransactionMode: true,
			})
			if err != nil {
				t.Fatal(err)
			}
			defer d.Close()

			runErr := d.Run(bytes.NewReader(tcase.cmds))
			if runErr != nil && !tcase.isErrorExpected {
				t.Fatal(runErr)
			}
		})
	}
}

func setupMongoContainer(image string, cmd ...string) (testcontainers.Container, string, string, error) {
	req := testcontainers.ContainerRequest{
		Image:        image,
		ExposedPorts: []string{"27017/tcp"},
		WaitingFor:   wait.ForListeningPort("27017/tcp").WithStartupTimeout(30 * time.Second),
		Cmd:          cmd,
	}
	container, err := testcontainers.GenericContainer(context.TODO(), testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, "", "", err
	}

	port, err := container.MappedPort(context.TODO(), "27017")
	if err != nil {
		return nil, "", "", err
	}
	host, err := container.Host(context.TODO())
	if err != nil {
		return nil, "", "", err
	}

	return container, host, port.Port(), nil
}

func waitForPrimary(client *mongo.Client) error {
	timeout := time.After(15 * time.Second)
	tick := time.Tick(500 * time.Millisecond)

	for {
		select {
		case <-timeout:
			return fmt.Errorf("timed out waiting for replica set to become primary")
		case <-tick:
			var result bson.M
			err := client.Database("admin").RunCommand(context.TODO(), bson.D{{Key: "isMaster", Value: 1}}).Decode(&result)
			if err == nil && result["ismaster"] == true {
				return nil
			}
		}
	}
}
