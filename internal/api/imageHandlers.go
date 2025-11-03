package api

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *Handlers) GetImageByID() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			imageIDStr := r.PathValue("id")
			if imageIDStr == "" {
				http.Error(w, "Image ID is required", http.StatusBadRequest)
				return
			}

			imageID, err := uuid.Parse(imageIDStr)
			if err != nil {
				http.Error(w, "Invalid image ID", http.StatusBadRequest)
				return
			}

			bucket, err := h.db.Mongo.GetBucket()
			if err != nil {
				http.Error(w, "Failed to access image storage", http.StatusInternalServerError)
				return
			}

			binUUID := primitive.Binary{Subtype: 4, Data: imageID[:]}

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			downloadStream, err := bucket.OpenDownloadStream(binUUID)
			if err != nil {
				http.Error(w, "Image not found", http.StatusNotFound)
				return
			}
			defer downloadStream.Close()

			buf := new(bytes.Buffer)
			if _, err := buf.ReadFrom(downloadStream); err != nil {
				http.Error(w, "Failed to read image", http.StatusInternalServerError)
				return
			}

			var fileDoc bson.M
			filter := bson.M{"_id": binUUID}
			err = bucket.GetFilesCollection().FindOne(ctx, filter).Decode(&fileDoc)
			if err == nil {
				if metadata, ok := fileDoc["metadata"].(bson.M); ok {
					if contentType, ok := metadata["contentType"].(string); ok && contentType != "" {
						w.Header().Set("Content-Type", contentType)
					}
			}
		}

		w.Header().Set("Cache-Control", "public, max-age=31536000")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", buf.Len()))

			if w.Header().Get("Content-Type") == "" {
				w.Header().Set("Content-Type", "image/jpeg")
			}

			w.WriteHeader(http.StatusOK)
			w.Write(buf.Bytes())
		},
	)
}

func (h *Handlers) GetAvatarByID() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			avatarIDStr := r.PathValue("id")
			if avatarIDStr == "" {
				http.Error(w, "Avatar ID is required", http.StatusBadRequest)
				return
			}

			avatarID, err := uuid.Parse(avatarIDStr)
			if err != nil {
				http.Error(w, "Invalid avatar ID", http.StatusBadRequest)
				return
			}

			
			bucket, err := h.db.Mongo.GetAvatarBucket()
			if err != nil {
				http.Error(w, "Failed to access avatar storage", http.StatusInternalServerError)
				return
			}

			binUUID := primitive.Binary{Subtype: 4, Data: avatarID[:]}

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			downloadStream, err := bucket.OpenDownloadStream(binUUID)
			if err != nil {
				http.Error(w, "Avatar not found", http.StatusNotFound)
				return
			}
			defer downloadStream.Close()

			buf := new(bytes.Buffer)
			if _, err := buf.ReadFrom(downloadStream); err != nil {
				http.Error(w, "Failed to read avatar", http.StatusInternalServerError)
				return
			}

			var fileDoc bson.M
			filter := bson.M{"_id": binUUID}
			err = bucket.GetFilesCollection().FindOne(ctx, filter).Decode(&fileDoc)
			if err == nil {
				if metadata, ok := fileDoc["metadata"].(bson.M); ok {
					if contentType, ok := metadata["contentType"].(string); ok && contentType != "" {
						w.Header().Set("Content-Type", contentType)
					}
				}
			}

			w.Header().Set("Cache-Control", "public, max-age=31536000")
			w.Header().Set("Content-Length", fmt.Sprintf("%d", buf.Len()))

			if w.Header().Get("Content-Type") == "" {
				w.Header().Set("Content-Type", "image/jpeg")
			}

			w.WriteHeader(http.StatusOK)
			w.Write(buf.Bytes())
		},
	)
}
