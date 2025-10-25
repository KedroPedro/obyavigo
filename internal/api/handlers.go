package api

import (
	"cmd/obyavigo/main.go/internal/config"
	"cmd/obyavigo/main.go/internal/database"
	"cmd/obyavigo/main.go/internal/mail"
	"cmd/obyavigo/main.go/internal/models"
	"cmd/obyavigo/main.go/internal/secure"
	"html/template"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type Handlers struct {
	db   *database.DB
	tmpl *template.Template
	mail *mail.Mail
	jwt  *secure.JWT
}

func New(db *database.DB, cfg *config.Config, mail *mail.Mail) *Handlers {
	tmpl := template.Must(template.ParseGlob(cfg.Environment.TemplatesDir + "*.html"))

	return &Handlers{
		db:   db,
		tmpl: tmpl,
		mail: mail,
		jwt:  secure.NewJWT(cfg),
	}
}

func (h *Handlers) GetMainPage() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/" {
				h.sendNotFound(w)
				return
			}

			h.tmpl.ExecuteTemplate(w, "index.html", nil)

		},
	)
}

func (h *Handlers) GetProfilePage() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/profile" && r.URL.Path != "/profile/" {
				h.sendNotFound(w)
				return
			}

			userID, err := userIDFromCtx(r)
			if err != nil {
				http.Redirect(w, r, "/auth/", http.StatusPermanentRedirect)
				return
			}

			userData, err := h.db.Psql.GetUserData(userID)
			if handleError(w, err, http.StatusInternalServerError, "error while executing get user preview data") {
				h.tmpl.ExecuteTemplate(w, "profile.html", nil)
				return
			}

			err = h.tmpl.ExecuteTemplate(w, "profile.html", userData)
			if handleError(w, err, http.StatusInternalServerError, "error while trying to send profile page") {
				h.sendNotFound(w)
			}
		},
	)
}

func (h *Handlers) GetMessagesPage() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/messages" && r.URL.Path != "/messages/" {
				h.sendNotFound(w)
				return
			}

			userID, err := userIDFromCtx(r)
			if err != nil {
				http.Redirect(w, r, "/auth/", http.StatusPermanentRedirect)
				return
			}

			// Get user chats
			chats, err := h.db.Psql.GetUserChats(userID)
			if handleError(w, err, http.StatusInternalServerError, "error while getting user chats") {
				return
			}

			data := models.MessagesPageData{
				Chats: chats,
			}

			err = h.tmpl.ExecuteTemplate(w, "messages.html", data)
			if handleError(w, err, http.StatusInternalServerError, "error while executing messages template") {
				h.sendNotFound(w)
			}
		},
	)
}

func (h *Handlers) GetCreateAdPage() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/create-ad" && r.URL.Path != "/create-ad/" {
				h.sendNotFound(w)
				return
			}
			h.tmpl.ExecuteTemplate(w, "create-ad.html", nil)
		},
	)
}

func (h *Handlers) GetAuthPage() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/auth" && r.URL.Path != "/auth/" {
				h.sendNotFound(w)
				return
			}
			if id, _ := userIDFromCtx(r); id != nil {
				http.Redirect(w, r, "/", http.StatusFound)
				return
			}
			h.tmpl.ExecuteTemplate(w, "auth.html", nil)
		},
	)
}

func (h *Handlers) GetAdsPage() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/ads" && r.URL.Path != "/ads/" {
				h.sendNotFound(w)
				return
			}

			// Get query parameters for pagination
			page := 1
			limit := 20
			if p := r.URL.Query().Get("page"); p != "" {
				if parsedPage, err := strconv.Atoi(p); err == nil && parsedPage > 0 {
					page = parsedPage
				}
			}

			offset := (page - 1) * limit

			// Get ads data
			ads, err := h.db.Psql.GetAdsList(limit, offset)
			if handleError(w, err, http.StatusInternalServerError, "error while getting ads list") {
				return
			}

			totalCount, err := h.db.Psql.GetAdsCount()
			if handleError(w, err, http.StatusInternalServerError, "error while getting ads count") {
				return
			}

			data := models.AdsPageData{
				Ads:        ads,
				TotalCount: totalCount,
				Page:       page,
				Limit:      limit,
			}

			err = h.tmpl.ExecuteTemplate(w, "ads.html", data)
			if handleError(w, err, http.StatusInternalServerError, "error while executing ads template") {
				h.sendNotFound(w)
			}
		},
	)
}

func (h *Handlers) GetAdPage() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			token := r.PathValue("token")
			adID, err := uuid.Parse(token)
			if err != nil {
				h.sendNotFound(w)
				return
			}

			userID, err := userIDFromCtx(r)
			if err != nil {
				userID = nil
			}

			if userID == nil {
				slog.Info("the guest visited the ad page", slog.String("page_id", adID.String()))
			} else {
				slog.Info("the user visited the ad page",
					slog.String("user_id", userID.String()),
					slog.String("page_id", adID.String()))
			}
			adData, err := h.db.Psql.GetAdInfo(&adID)

			if adData == nil || err != nil {
				h.sendNotFound(w)
				return
			}
			role := "user"
			if userID != nil {
				role, err = h.db.Psql.GetUserRole(userID)
				if handleError(w, err, http.StatusInternalServerError, "error while trying to get user role") {
					return
				}
			}
			if userID == nil || (*userID != adData.UserID && role == "user") {
				if adData.AdStatus != "public" {
					h.sendNotFound(w)
					return
				}
			}
			err = h.tmpl.ExecuteTemplate(w, "ad.html", adData)

			if handleError(w, err, http.StatusInternalServerError, "template execution error") {
				return
			}
		},
	)
}

func (h *Handlers) UserLogout() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			http.SetCookie(w, &http.Cookie{
				Name:     "auth_token",
				Value:    "",
				Path:     "/",
				Expires:  time.Unix(0, 0),
				MaxAge:   -1,
				HttpOnly: true,
				Secure:   true,
			})
			w.WriteHeader(http.StatusOK)
		},
	)
}

func (h *Handlers) GetAdminPanelPage() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/admin-panel" && r.URL.Path != "/admin-panel/" {
				h.sendNotFound(w)
				return
			}

			userID, err := userIDFromCtx(r)
			if err != nil {
				http.Redirect(w, r, "/auth/", http.StatusPermanentRedirect)
				return
			}

			role, err := h.db.Psql.GetUserRole(userID)
			if handleError(w, err, http.StatusInternalServerError, "error while trying to get user role") {
				return
			}

			if role != "admin" {
				h.sendNotFound(w)
				return
			}

			// Get admin panel data
			stats, err := h.db.Psql.GetAdminStats()
			if handleError(w, err, http.StatusInternalServerError, "error while getting admin stats") {
				return
			}

			// Get limited data for each section
			ads, err := h.db.Psql.GetAllAds(10, 0)
			if handleError(w, err, http.StatusInternalServerError, "error while getting all ads") {
				return
			}

			users, err := h.db.Psql.GetAllUsers(10, 0)
			if handleError(w, err, http.StatusInternalServerError, "error while getting all users") {
				return
			}

			complaints, err := h.db.Psql.GetComplaints(10, 0)
			if handleError(w, err, http.StatusInternalServerError, "error while getting complaints") {
				return
			}

			moderation, err := h.db.Psql.GetModerationAds(10, 0)
			if handleError(w, err, http.StatusInternalServerError, "error while getting moderation ads") {
				return
			}

			data := models.AdminPanelData{
				Stats:      *stats,
				Ads:        ads,
				Users:      users,
				Complaints: complaints,
				Moderation: moderation,
			}

			err = h.tmpl.ExecuteTemplate(w, "admin-panel.html", data)
			if handleError(w, err, http.StatusInternalServerError, "error while executing admin panel template") {
				h.sendNotFound(w)
			}
		},
	)
}

func (h *Handlers) GetLikedAdsPage() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/liked-ads" && r.URL.Path != "/liked-ads/" {
				h.sendNotFound(w)
				return
			}

			userID, err := userIDFromCtx(r)
			if err != nil {
				http.Redirect(w, r, "/auth/", http.StatusPermanentRedirect)
				return
			}

			ads, err := h.db.Psql.GetLikedAds(userID)
			if handleError(w, err, http.StatusInternalServerError, "error while getting liked ads") {
				return
			}

			data := models.LikedAdsPageData{
				Ads: ads,
			}

			err = h.tmpl.ExecuteTemplate(w, "liked-ads.html", data)
			if handleError(w, err, http.StatusInternalServerError, "error while executing liked ads template") {
				h.sendNotFound(w)
			}
		},
	)
}
