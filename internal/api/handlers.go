package api

import (
	"cmd/obyavigo/main.go/internal/config"
	"cmd/obyavigo/main.go/internal/database"
	"cmd/obyavigo/main.go/internal/mail"
	"cmd/obyavigo/main.go/internal/models"
	"cmd/obyavigo/main.go/internal/secure"
	"html/template"
	"net/http"

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

			userID, err := userIDFromCtx(r)
			if handleError(w, err, http.StatusInternalServerError, "error while trying to get user id from context") {
				h.tmpl.ExecuteTemplate(w, "index.html", nil)
				return
			}
			username, err := h.db.Psql.GetUserPreviewData(userID)
			if handleError(w, err, http.StatusInternalServerError, "error while executing get user preview data") {
				h.tmpl.ExecuteTemplate(w, "index.html", nil)
				return
			}

			data := models.UserPreviewPageData{
				UserName: username,
			}

			h.tmpl.ExecuteTemplate(w, "index.html", data)

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
			if handleError(w, err, http.StatusInternalServerError, "error while trying to get user id from context") {
				h.tmpl.ExecuteTemplate(w, "profile.html", nil)
				return
			}
			username, err := h.db.Psql.GetUserPreviewData(userID)
			if handleError(w, err, http.StatusInternalServerError, "error while executing get user preview data") {
				h.tmpl.ExecuteTemplate(w, "profile.html", nil)
				return
			}

			data := models.UserPreviewPageData{
				UserName: username,
			}

			err = h.tmpl.ExecuteTemplate(w, "profile.html", data)
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
			h.tmpl.ExecuteTemplate(w, "messages.html", nil)
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
			h.tmpl.ExecuteTemplate(w, "ads.html", nil)
		},
	)
}

func (h *Handlers) GetAdPage() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			userID, err := userIDFromCtx(r)
			if handleError(w, err, http.StatusInternalServerError, "error while trying to get user id from context") {
				return
			}
			token := r.PathValue("token")
			adID, err := uuid.Parse(token)
			if err != nil {
				h.sendNotFound(w)
				return
			}

			data, err := h.db.Psql.GetAdInfo(&adID)
			if handleError(w, err, http.StatusInternalServerError, "error while trying to get ad data") {
				return
			}
			role, err := h.db.Psql.GetUserRole(userID)
			if handleError(w, err, http.StatusInternalServerError, "error while trying to get user role") {
				return
			}
			if (*userID != data.UserId && role == "user") && data.Status != "public" {
				h.sendNotFound(w)
				return
			}
			h.tmpl.ExecuteTemplate(w, "ad.html", data)
		},
	)
}
