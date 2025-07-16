package main

import (
	"context"
	"fmt"
	"github/michaellimmm/turakkingu/internal/adapter"
	"github/michaellimmm/turakkingu/internal/core"
	"github/michaellimmm/turakkingu/internal/repository"
	"github/michaellimmm/turakkingu/internal/usecase"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	migrateData()

	config, err := core.NewConfig()
	if err != nil {
		slog.Error("failed to get config", slog.String("error", err.Error()))
		os.Exit(1)
	}

	repo, err := repository.NewRepo(config)
	if err != nil {
		slog.Error("failed to initialize repository", slog.String("error", err.Error()))
		os.Exit(1)
	}

	uc := usecase.NewUseCase(config, repo)
	server := adapter.NewAdapter(config, uc)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	serverErrChan := make(chan error, 1)

	go func() {
		err = server.Run()
		if err != nil {
			slog.Error("error", slog.String("error", err.Error()))
			serverErrChan <- err
		}
	}()

	slog.Info("server is running....")

	select {
	case sig := <-sigChan:
		slog.Info("shutdown signal received, starting graceful shutdown", slog.String("signal", sig.String()))
		gracefulShutdown(server, repo)
	case err := <-serverErrChan:
		slog.Error("server failed to start", slog.String("error", err.Error()))
		os.Exit(1)
	}

	slog.Info("shutting down services")
}

func gracefulShutdown(server adapter.AdapterCloser, repo repository.RepoCloser) {
	slog.Info("stopping server...")
	if err := server.Close(context.Background()); err != nil {
		slog.Error("server shutdown error", slog.String("error", err.Error()))
	}

	slog.Info("closing database connections...")
	if err := repo.Close(context.Background()); err != nil {
		slog.Error("repository close error", slog.String("error", err.Error()))
	}
}

func migrateData() {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	// Insert Tracking Settings
	InsertTrackingSettings(client, ctx)
	InsertThankYouPage(client, ctx)
}

const (
	Tenant1 string = "tenant1"
	Katene  string = "katene"
	IRobot  string = "irobot"
	Takami  string = "takami"
)

type ObjectType string

const (
	TrackingSetting ObjectType = "tracking_setting"
)

func GetObjectIDs() map[ObjectType]map[string]bson.ObjectID {
	ids := make(map[ObjectType]map[string]bson.ObjectID)

	trackingSettingTenant1ID, err := bson.ObjectIDFromHex("68759be6ae3bd247fc9b4fd4")
	if err != nil {
		log.Fatal("Invalid ObjectID:", err)
	}

	trackingSettingKateneID, err := bson.ObjectIDFromHex("68759be6ae3bd247fc9b4fd5")
	if err != nil {
		log.Fatal("Invalid ObjectID:", err)
	}

	trackingSettingIRobotID, err := bson.ObjectIDFromHex("68759be6ae3bd247fc9b4fd6")
	if err != nil {
		log.Fatal("Invalid ObjectID:", err)
	}

	trackingSettingTakamiID, err := bson.ObjectIDFromHex("68759be6ae3bd247fc9b4fd7")
	if err != nil {
		log.Fatal("Invalid ObjectID:", err)
	}

	ids[TrackingSetting] = map[string]bson.ObjectID{
		Tenant1: trackingSettingTenant1ID,
		Katene:  trackingSettingKateneID,
		IRobot:  trackingSettingIRobotID,
		Takami:  trackingSettingTakamiID,
	}

	return ids
}

func InsertTrackingSettings(client *mongo.Client, ctx context.Context) {
	// Select database and collection
	db := client.Database("conversionTracking")
	collectionName := "tracking_setting"

	collections, err := db.ListCollectionNames(ctx, bson.M{"name": collectionName})
	if err != nil {
		log.Fatal("Error listing collections:", err)
	}

	// If it exists, drop it
	if len(collections) > 0 {
		fmt.Println("Collection exists. Dropping...")
		if err := db.Collection(collectionName).Drop(ctx); err != nil {
			log.Fatal("Error dropping collection:", err)
		}
		fmt.Println("Collection dropped.")
	} else {
		fmt.Println("Collection does not exist. No need to drop.")
	}

	// Re-create collection reference
	collection := db.Collection(collectionName)

	objectIDs := GetObjectIDs()

	timestamp, err := time.Parse(time.RFC3339, "2025-07-15T00:08:06.666Z")
	if err != nil {
		log.Fatal("Invalid timestamp:", err)
	}

	// Create the document
	docs := []interface{}{
		bson.M{
			"_id":        objectIDs[TrackingSetting][Tenant1],
			"tenant_id":  "tenant1",
			"created_at": timestamp,
			"updated_at": timestamp,
		},
		bson.M{
			"_id":        objectIDs[TrackingSetting][Katene],
			"tenant_id":  "ulototeseyalixun",
			"created_at": timestamp,
			"updated_at": timestamp,
		},
		bson.M{
			"_id":        objectIDs[TrackingSetting][IRobot],
			"tenant_id":  "irobotonlinestore",
			"created_at": timestamp,
			"updated_at": timestamp,
		},
		bson.M{
			"_id":        objectIDs[TrackingSetting][Takami],
			"tenant_id":  "esufajicopajaqoz",
			"created_at": timestamp,
			"updated_at": timestamp,
		},
	}

	// Insert all documents
	res, err := collection.InsertMany(ctx, docs)
	if err != nil {
		log.Fatal("InsertMany failed:", err)
	}

	fmt.Println("Inserted document ID:", res.InsertedIDs)
}

func InsertThankYouPage(client *mongo.Client, ctx context.Context) {
	// Select database and collection
	db := client.Database("conversionTracking")
	collectionName := "thank_you_page"

	collections, err := db.ListCollectionNames(ctx, bson.M{"name": collectionName})
	if err != nil {
		log.Fatal("Error listing collections:", err)
	}

	// If it exists, drop it
	if len(collections) > 0 {
		fmt.Println("Collection exists. Dropping...")
		if err := db.Collection(collectionName).Drop(ctx); err != nil {
			log.Fatal("Error dropping collection:", err)
		}
		fmt.Println("Collection dropped.")
	} else {
		fmt.Println("Collection does not exist. No need to drop.")
	}

	// Re-create collection reference
	collection := db.Collection(collectionName)

	// Prepare fixed ObjectID and timestamp
	objectID1, err := bson.ObjectIDFromHex("11159be6ae3bd247fc9b4fd4")
	if err != nil {
		log.Fatal("Invalid ObjectID:", err)
	}

	objectID2, err := bson.ObjectIDFromHex("11259be6ae3bd247fc9b4fd5")
	if err != nil {
		log.Fatal("Invalid ObjectID:", err)
	}

	objectID3, err := bson.ObjectIDFromHex("11359be6ae3bd247fc9b4fd6")
	if err != nil {
		log.Fatal("Invalid ObjectID:", err)
	}

	objectID4, err := bson.ObjectIDFromHex("11459be6ae3bd247fc9b4fd6")
	if err != nil {
		log.Fatal("Invalid ObjectID:", err)
	}

	timestamp, err := time.Parse(time.RFC3339, "2025-07-15T00:08:06.666Z")
	if err != nil {
		log.Fatal("Invalid timestamp:", err)
	}

	ids := GetObjectIDs()

	// Create the document
	docs := []interface{}{
		bson.M{
			"_id":                 objectID1,
			"tracking_setting_id": ids[TrackingSetting][Katene],
			"url":                 "https://katene.chuden.jp/clubkatene/id/idConfirm.do",
			"point":               0,
			"name":                "中部電力（カテエネ会員登録）Thank you page",
			"tracking_status":     1,
			"created_at":          timestamp,
			"updated_at":          timestamp,
		},
		bson.M{
			"_id":                 objectID2,
			"tracking_setting_id": ids[TrackingSetting][IRobot],
			"url":                 "https://store.irobot-jp.com/cart_complete.html",
			"point":               0,
			"name":                "アイロボット オンラインストア Thank you page",
			"tracking_status":     1,
			"created_at":          timestamp,
			"updated_at":          timestamp,
		},
		bson.M{
			"_id":                 objectID3,
			"tracking_setting_id": ids[TrackingSetting][Takami],
			"url":                 "https://www.takami-labo.com/cart/completion/completion.html",
			"point":               0,
			"name":                "タカミスキンピール Thank you page",
			"tracking_status":     1,
			"created_at":          timestamp,
			"updated_at":          timestamp,
		},
		bson.M{
			"_id":                 objectID4,
			"tracking_setting_id": ids[TrackingSetting][Tenant1],
			"url":                 "https://car-form.ngrok.app/thank-you",
			"point":               0,
			"name":                "Car form Thank you page",
			"tracking_status":     1,
			"created_at":          timestamp,
			"updated_at":          timestamp,
		},
	}

	// Insert all documents
	res, err := collection.InsertMany(ctx, docs)
	if err != nil {
		log.Fatal("InsertMany failed:", err)
	}

	fmt.Println("Inserted document ID:", res.InsertedIDs)
}
