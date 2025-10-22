package mongodb

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"time"

	"cmd/obyavigo/main.go/internal/config"

	"github.com/google/uuid"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongo struct {
	mng *mongo.Client
	q   map[string]string
}

func Connect(cfg *config.Config, q map[string]string) (*Mongo, error) {
	mongoURI := fmt.Sprintf("mongodb://%s:%s@%s:%d/%s",
		cfg.Database.Mongo.User,
		cfg.Database.Mongo.Password,
		cfg.Database.Mongo.Host,
		cfg.Database.Mongo.Port,
		cfg.Database.Mongo.DBName,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, fmt.Errorf("mongodb connection failed: %w", err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		_ = client.Disconnect(context.Background())
		return nil, fmt.Errorf("failed to ping mongodb: %w", err)
	}

	return &Mongo{
		mng: client,
		q:   q,
	}, nil
}

func (m *Mongo) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := m.mng.Ping(ctx, nil); err != nil {
		return fmt.Errorf("mongodb connection error: %w", err)
	}
	return nil
}

func (m *Mongo) UploadImages(ctx context.Context, files []*multipart.FileHeader, adID string) ([]string, error) {
	bucket, err := gridfs.NewBucket(
		m.mng.Database("obyavigopics"),
		options.GridFSBucket().SetName("fs"),
	)
	if err != nil {
		return nil, fmt.Errorf("gridfs bucket create error: %w", err)
	}

	var ids []string
	for _, fh := range files {
		file, err := fh.Open()
		if err != nil {
			return nil, fmt.Errorf("cannot open image file: %w", err)
		}
		defer file.Close()

		generatedUUID := uuid.New()
		binUUID := primitive.Binary{Subtype: 4, Data: generatedUUID[:]}

		uploadOpts := options.GridFSUpload().
			SetMetadata(bson.M{
				"ad_id":      adID,
				"length":     fh.Size,
				"uploadDate": time.Now(),
			})

		uploadStream, err := bucket.OpenUploadStreamWithID(binUUID, "", uploadOpts)
		if err != nil {
			return nil, fmt.Errorf("open upload stream: %w", err)
		}

		_, err = io.Copy(uploadStream, file)
		uploadStream.Close()
		if err != nil {
			return nil, fmt.Errorf("copy image data: %w", err)
		}

		ids = append(ids, generatedUUID.String())
	}
	return ids, nil
}

func (m *Mongo) DownloadImagesByID(ctx context.Context, adID string) ([][]byte, error) {
	bucket, err := gridfs.NewBucket(
		m.mng.Database("obyavigopics"),
		options.GridFSBucket().SetName("fs"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create GridFS bucket: %w", err)
	}

	filter := bson.M{"metadata.ad_id": adID}
	cursor, err := bucket.Find(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find files by adID: %w", err)
	}
	defer cursor.Close(ctx)

	var filesData [][]byte
	for cursor.Next(ctx) {
		var fileDoc bson.M
		if err := cursor.Decode(&fileDoc); err != nil {
			return nil, fmt.Errorf("failed to decode file info: %w", err)
		}

		fileID, ok := fileDoc["_id"].(primitive.Binary)
		if !ok {
			return nil, fmt.Errorf("file id not in expected format")
		}

		downloadStream, err := bucket.OpenDownloadStream(fileID)
		if err != nil {
			return nil, fmt.Errorf("failed to open download stream: %w", err)
		}
		defer downloadStream.Close()

		buf := new(bytes.Buffer)
		if _, err := io.Copy(buf, downloadStream); err != nil {
			return nil, fmt.Errorf("failed to read file data: %w", err)
		}

		filesData = append(filesData, buf.Bytes())
	}
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return filesData, nil
}
