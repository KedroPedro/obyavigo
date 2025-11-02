package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"cmd/obyavigo/main.go/internal/models"

	"golang.org/x/net/websocket"

	"github.com/google/uuid"
)

type CreateChatRequest struct {
	ListingId uuid.UUID `json:"listing_id"`
}

type MessageJSON struct {
	ChatID string `json:"chat_id"`
	Text   string `json:"text"`
}

func (h *Handlers) CreateChat() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			userID, err := userIDFromCtx(r)
			if err != nil {
				http.Redirect(w, r, "/auth/", http.StatusPermanentRedirect)
				return
			}

			var req CreateChatRequest
			if handleError(w, json.NewDecoder(r.Body).Decode(&req), http.StatusBadRequest, "invalid request body") {
				return
			}

			chat := &models.Chat{
				ListingId:  &req.ListingId,
				CustomerId: userID,
			}

			err = h.db.Psql.CreateChat(chat)
			if handleError(w, err, http.StatusInternalServerError, "error while creating chat") {
				return
			}

			http.Redirect(w, r, "/messages/", http.StatusPermanentRedirect)
		},
	)
}

func (h *Handlers) ChatWebSocketHandler() websocket.Handler {
	return websocket.Handler(
		func(ws *websocket.Conn) {
			jwtToken := ws.Request().URL.Query().Get("token")

			userID, err := h.jwt.ParseJWTToken(jwtToken)
			if err != nil || userID == nil {
				ws.Close()
				return
			}
			h.ws.Store(userID.String(), ws)

			defer func() {
				h.ws.Delete(userID.String())
				ws.Close()
			}()

			for {
				var msgJSON MessageJSON
				err := websocket.JSON.Receive(ws, &msgJSON)
				if err != nil {
					slog.Error("error when receiving message", slog.String("error", err.Error()))
					break
				}

				chat, err := h.db.Psql.GetChatById(msgJSON.ChatID)
				if err != nil {
					slog.Error("error while trying to get chat by id", slog.String("error", err.Error()))
					break
				}
				toID := chat.CustomerId
				if userID.String() == chat.CustomerId.String() {
					toID = chat.SellerId
				}

				msg := models.Message{
					Text:     msgJSON.Text,
					ChatId:   chat.ChatId,
					SenderId: userID,
				}

				if err = h.db.Psql.CreateMessage(&msg); err != nil {
					slog.Error("error while trying to save message", slog.String("error", err.Error()))
					break
				}

				val, ok := h.ws.Load(toID.String())
				if ok {
					toConn, ok := val.(*websocket.Conn)
					if !ok {
						slog.Error("stored connection is not a *websocket.Conn")
						h.ws.Delete(toID.String())
						continue
					}
					err = websocket.Message.Send(toConn, msgJSON)
					if err != nil {
						slog.Error("error sending message", slog.String("error", err.Error()))
						h.ws.Delete(toID.String())
						toConn.Close()
					}
				}
			}
		},
	)
}

func (h *Handlers) GetMessages() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			userID, err := userIDFromCtx(r)
			if err != nil {
				http.Redirect(w, r, "/auth/", http.StatusPermanentRedirect)
				return
			}

			chatId := r.URL.Query().Get("chatId")
			offset := r.URL.Query().Get("offset")

			chat, err := h.db.Psql.GetChatById(chatId)

			role, err := h.db.Psql.GetUserRole(userID)
			if handleError(w, err, http.StatusForbidden, "error while trying to get user role") {
				return
			}

			if role != "admin" && role != "moderator" {
				if chat.CustomerId.String() != userID.String() || chat.SellerId.String() != userID.String() {
					w.WriteHeader(http.StatusForbidden)
					return
				}
			}

			if handleError(w, err, http.StatusBadRequest, "error while trying to get chat by id") {
				return
			}

			msgs, err := h.db.Psql.GetMessages(chatId, offset)
			if handleError(w, err, http.StatusBadRequest, "get messages error") {
				return
			}

			type MessagesResponse struct {
				Messages *[]models.Message `json:"messages"`
				Offset   string            `json:"offset"`
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(&MessagesResponse{
				Messages: msgs,
				Offset:   offset,
			})
		},
	)
}
