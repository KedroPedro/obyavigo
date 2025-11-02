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
	"sync"
	"time"

	"github.com/google/uuid"
)

type Handlers struct {
	db   *database.DB
	tmpl *template.Template
	mail *mail.Mail
	jwt  *secure.JWT
	ws   sync.Map
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

			_ = userID

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

			_, err := userIDFromCtx(r)
			if err != nil {
				http.Redirect(w, r, "/auth/", http.StatusPermanentRedirect)
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

			filters := models.AdSearchFilters{
				Page:        parseIntParam(r, "page", 1, 0),
				Limit:       20,
				Category:    r.URL.Query().Get("category"),
				Subcategory: r.URL.Query().Get("subcategory"),
				Region:      r.URL.Query().Get("region"),
				Location:    r.URL.Query().Get("location"),
				Condition:   r.URL.Query().Get("condition"),
				SearchQuery: r.URL.Query().Get("q"),
				SortBy:      r.URL.Query().Get("sort"),
				MinPrice:    parsePriceParam(r, "min_price"),
				MaxPrice:    parsePriceParam(r, "max_price"),
			}

			ads, err := h.db.Psql.SearchAds(&filters)
			if handleError(w, err, http.StatusInternalServerError, "error while searching ads") {
				return
			}

			totalCount, err := h.db.Psql.SearchAdsCount(&filters)
			if handleError(w, err, http.StatusInternalServerError, "error while getting ads count") {
				return
			}

			data := models.AdsPageData{
				Ads:        ads,
				TotalCount: totalCount,
				Page:       filters.Page,
				Limit:      filters.Limit,
				Category:   filters.Category,
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

			imageIDs, err := h.db.Psql.GetAdImageIDs(&adID)
			if err == nil && len(imageIDs) > 0 {
				adData.Images = imageIDs
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

			type AdPageView struct {
				models.AdPage
				FormattedDate string
			}
			pageData := AdPageView{
				AdPage:        *adData,
				FormattedDate: adData.CreatedAt.Format("02.01.2006"),
			}

			err = h.tmpl.ExecuteTemplate(w, "ad.html", pageData)

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

		if role != "admin" && role != "moderator" {
			h.sendNotFound(w)
			return
		}

			err = h.tmpl.ExecuteTemplate(w, "admin-panel.html", nil)
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
