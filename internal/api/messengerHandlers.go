package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"cmd/obyavigo/main.go/internal/models"

	"golang.org/x/net/websocket"

	"github.com/google/uuid"
)

type CreateChatRequest struct {
	ListingId uuid.UUID `json:"listing_id"`
}

type MessageJSON struct {
	ChatID    string     `json:"chat_id"`
	Text      string     `json:"text"`
	SenderID  string     `json:"sender_id,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}

func (h *Handlers) CreateChat() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			userID, err := userIDFromCtx(r)
			if err != nil {
				sendJSONError(w, http.StatusUnauthorized, "unauthorized")
				return
			}

			var req CreateChatRequest
			if handleError(w, json.NewDecoder(r.Body).Decode(&req), http.StatusBadRequest, "invalid request body") {
				return
			}

			if req.ListingId == uuid.Nil {
				sendJSONError(w, http.StatusBadRequest, "listing_id is required")
				return
			}

			p, err := h.db.Psql.GetAdInfo(&req.ListingId)
			if handleError(w, err, http.StatusNotFound, "объявление не найдено") {
				sendJSONError(w, http.StatusForbidden, "нельзя создать чат с самим собой")
				return
			}

			if p.UserID.String() == userID.String() {

			}

			chat := &models.Chat{
				ListingId:  &req.ListingId,
				CustomerId: userID,
			}

			createdChat, err := h.db.Psql.CreateChat(chat)
			if handleError(w, err, http.StatusInternalServerError, "error while creating chat") {
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]string{
				"chat_id": createdChat.ChatId.String(),
			})
		},
	)
}

func (h *Handlers) GetUserChats() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			userID, err := userIDFromCtx(r)
			if err != nil {
				sendJSONError(w, http.StatusUnauthorized, "unauthorized")
				return
			}

			chats, err := h.db.Psql.GetUserChats(userID)
			if handleError(w, err, http.StatusInternalServerError, "error while getting chats") {
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"chats": chats,
			})
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

				msgJSON.Text = strings.TrimSpace(msgJSON.Text)
				if msgJSON.Text == "" {
					continue
				}

				chatID, err := uuid.Parse(msgJSON.ChatID)
				if err != nil {
					slog.Error("invalid chat id", slog.String("chat_id", msgJSON.ChatID))
					break
				}

				chat, err := h.db.Psql.GetChatById(chatID)
				if err != nil {
					slog.Error("error while trying to get chat by id", slog.String("error", err.Error()))
					break
				}

				if chat.ChatId == nil {
					idCopy := chatID
					chat.ChatId = &idCopy
				}

				if chat.CustomerId == nil || chat.SellerId == nil {
					slog.Error("chat participants not found", slog.String("chat_id", msgJSON.ChatID))
					break
				}

				if userID.String() != chat.CustomerId.String() && userID.String() != chat.SellerId.String() {
					slog.Warn("user tried to send message to foreign chat", slog.String("user_id", userID.String()), slog.String("chat_id", msgJSON.ChatID))
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

				msgJSON.SenderID = userID.String()
				msgJSON.CreatedAt = msg.CreatedAt

				val, ok := h.ws.Load(toID.String())
				if ok {
					toConn, ok := val.(*websocket.Conn)
					if !ok {
						slog.Error("stored connection is not a *websocket.Conn")
						h.ws.Delete(toID.String())
						continue
					}
					err = websocket.JSON.Send(toConn, msgJSON)
					if err != nil {
						slog.Error("error sending message", slog.String("error", err.Error()))
						h.ws.Delete(toID.String())
						toConn.Close()
					}
				}

				sendErr := websocket.JSON.Send(ws, msgJSON)
				if sendErr != nil {
					slog.Error("error sending message back to sender", slog.String("error", sendErr.Error()))
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
				sendJSONError(w, http.StatusUnauthorized, "unauthorized")
				return
			}

			chatIdParam := r.PathValue("chatId")
			chatUUID, err := uuid.Parse(chatIdParam)
			if err != nil {
				sendJSONError(w, http.StatusBadRequest, "invalid chat id")
				return
			}

			chat, err := h.db.Psql.GetChatById(chatUUID)
			if handleError(w, err, http.StatusBadRequest, "error while trying to get chat by id") {
				return
			}

			role, err := h.db.Psql.GetUserRole(userID)
			if handleError(w, err, http.StatusForbidden, "error while trying to get user role") {
				return
			}

			if role != "admin" && role != "moderator" &&
				(chat.CustomerId == nil || chat.SellerId == nil ||
					(chat.CustomerId.String() != userID.String() && chat.SellerId.String() != userID.String())) {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			limit := parseIntParam(r, "limit", 50, 100)
			offset := 0
			if rawOffset := r.URL.Query().Get("offset"); rawOffset != "" {
				if parsed, convErr := strconv.Atoi(rawOffset); convErr == nil && parsed >= 0 {
					offset = parsed
				}
			}

			if chat.ChatId == nil {
				sendJSONError(w, http.StatusInternalServerError, "chat id is missing")
				return
			}

			msgs, err := h.db.Psql.GetMessages(*chat.ChatId, limit, offset)
			if handleError(w, err, http.StatusBadRequest, "get messages error") {
				return
			}

			type MessagesResponse struct {
				Messages   []models.Message `json:"messages"`
				NextOffset int              `json:"next_offset"`
				HasMore    bool             `json:"has_more"`
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(&MessagesResponse{
				Messages:   msgs,
				NextOffset: offset + len(msgs),
				HasMore:    len(msgs) == limit,
			})
		},
	)
}
