package postgres

import (
	"cmd/obyavigo/main.go/internal/config"
	"cmd/obyavigo/main.go/internal/models"
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Postgres struct {
	psql *sql.DB
	q    map[string]string
}

func Connect(cfg *config.Config, q map[string]string) (*Postgres, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Postgres.Host,
		cfg.Database.Postgres.Port,
		cfg.Database.Postgres.User,
		cfg.Database.Postgres.Password,
		cfg.Database.Postgres.DBName,
	)
	psql, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("postgres connection failed: %w", err)
	}
	return &Postgres{
		psql: psql,
		q:    q,
	}, nil
}

func (p *Postgres) Ping() error {
	if err := p.psql.Ping(); err != nil {
		return fmt.Errorf("postgres connection error: %w", err)
	}
	return nil
}

func (p *Postgres) CreateNewUser(u *models.User) (*uuid.UUID, error) {
	query, ok := p.q["CreateNewUser"]
	if !ok {
		return nil, fmt.Errorf("request 'CreateNewUser' not found")
	}
	var id uuid.UUID
	err := p.psql.QueryRow(query,
		u.Username,
		u.Email,
		u.PasswordHash,
		u.Role,
		u.Status,
		u.PhoneNumber,
		u.RegistrationDate,
		u.LastLogin,
		u.ProfilePictureID,
		u.Bio,
		u.Settings,
	).Scan(&id)

	if err != nil {
		return nil, fmt.Errorf("error while executing create user request: %w", err)
	}
	return &id, nil
}

func (p *Postgres) CheckEmail(email string) (bool, error) {
	query, ok := p.q["CheckEmail"]
	if !ok {
		return false, fmt.Errorf("request 'CheckEmail' not found")
	}

	var exists bool
	err := p.psql.QueryRow(query, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error while executing check email request: %w", err)
	}

	return exists, nil
}

func (p *Postgres) CreateEmailConfirmation(id *uuid.UUID, token string) error {
	query, ok := p.q["CreateEmailConfirmation"]
	if !ok {
		return fmt.Errorf("request 'CheckEmail' not found")
	}

	err := p.psql.QueryRow(query, id, token).Err()
	if err != nil {
		return fmt.Errorf("error while executing create email confirmation request: %w", err)
	}

	return nil
}

func (p *Postgres) IsConfirmationExpired(token string) (bool, error) {
	query, ok := p.q["IsConfirmationExpired"]
	if !ok {
		return false, fmt.Errorf("request 'IsConfirmationExpired' not found")
	}
	var exTime time.Time
	var confirmed bool
	err := p.psql.QueryRow(query, token).Scan(&exTime, &confirmed)
	if err != nil {
		return false, fmt.Errorf("error while calling the request to get the expiration date of the registration confirmation: %w", err)
	}

	if confirmed {
		return true, nil
	}

	if time.Now().After(exTime) {
		return false, nil
	}
	return true, nil
}

func (p *Postgres) DeleteAccountByEmail(email string) error {
	query, ok := p.q["DeleteAccountByEmail"]
	if !ok {
		return fmt.Errorf("request 'DeleteAccountByEmail' not found")
	}
	err := p.psql.QueryRow(query, email).Err()
	if err != nil {
		return fmt.Errorf("error while trying to delete account: %w", err)
	}
	return nil
}

func (p *Postgres) ConfirmAccount(token string) error {
	query, ok := p.q["ConfirmAccount"]
	if !ok {
		return fmt.Errorf("request 'ConfirmAccount' not found")
	}
	res, err := p.psql.Exec(query, token)
	if err != nil {
		return fmt.Errorf("error while calling the request to confirm the account: %w", err)
	}
	rowsAff, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error while getting the number of affected rows: %w", err)
	}
	if rowsAff == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (p *Postgres) GetAuthInfoByEmail(email string) (*models.AuhtInfo, error) {
	query, ok := p.q["GetAuthInfoByEmail"]
	if !ok {
		return nil, fmt.Errorf("request 'GetAuhtInfoByEmail' not found")
	}

	var id uuid.UUID
	var passHash string
	var confirmed bool

	err := p.psql.QueryRow(query, email).Scan(&id, &passHash, &confirmed)

	if err != nil {
		return nil, fmt.Errorf("error while trying to get user auth data by email: %w", err)
	}

	return &models.AuhtInfo{
		Id:           id,
		PasswordHash: passHash,
		Confirmed:    confirmed,
	}, nil
}

func (p *Postgres) CreateAd(m *models.AdTemplate) (*uuid.UUID, error) {
	query, ok := p.q["CreateAd"]
	if !ok {
		return nil, fmt.Errorf("request 'CreateAd' not found")
	}
	var adId uuid.UUID
	if err := p.psql.QueryRow(query,
		m.UserId,
		m.CategoryId,
		m.Title,
		m.Description,
		m.Price,
		m.LocationId,
		m.Condition,
		m.ContactPhone,
	).Scan(&adId); err != nil {
		return nil, fmt.Errorf("error while executing create ad request: %w", err)
	}

	return &adId, nil
}

func (p *Postgres) GetCreateAdDependencies(m *models.AdTemplate) error {
	query, ok := p.q["GetCreateAdDependencies"]
	if !ok {
		return fmt.Errorf("request 'GetCreateAdDependencies' not found")
	}

	var subcategoryName *string
	if m.SubcategoryName != "" {
		subcategoryName = &m.SubcategoryName
	}

	err := p.psql.QueryRow(query, subcategoryName, m.CategoryName, m.LocationName).Scan(&m.CategoryId, &m.LocationId)
	if err != nil {
		return fmt.Errorf("error while trying to get create ad dependencies: %w", err)
	}
	return nil
}

func (p *Postgres) InsertImages(userId *uuid.UUID, adId *uuid.UUID, ids []string) error {
	query, ok := p.q["InsertImage"]
	if !ok {
		return fmt.Errorf("request 'InsertImage' not found")
	}
	var valuesStr []string
	var args []interface{}
	for i, id := range ids {
		valuesStr = append(valuesStr, fmt.Sprintf("($%d, $%d, $%d)", i*3+1, i*3+2, i*3+3))
		args = append(args, adId, id, i)
	}
	query = fmt.Sprintf(query, strings.Join(valuesStr, ","))

	_, err := p.psql.ExecContext(context.Background(), query, args...)
	if err != nil {
		return fmt.Errorf("error while trying insert images: %w", err)
	}
	return nil
}

func (p *Postgres) GetAdInfo(adId *uuid.UUID) (*models.AdPage, error) {
	query, ok := p.q["GetAdInfo"]
	if !ok {
		return nil, fmt.Errorf("request 'GetAdInfo' not found")
	}
	var data models.AdPage
	loc := make([]string, 4)
	err := p.psql.QueryRow(query, adId).Scan(
		&data.Title,
		&data.Price,
		&data.CreatedAt,
		&data.ViewsCount,
		&data.Condition,
		&data.Description,
		&data.AdStatus,
		&data.SellerName,
		&data.Online,
		&data.UserID,
		&loc[0],
		&loc[1],
		&loc[2],
		&loc[3],
	)
	data.SellerCity = strings.Join(loc, " ")
	data.Price /= 100
	if err != nil {
		return nil, fmt.Errorf("error while executing get ad info request: %w", err)
	}

	return &data, nil
}

func (p *Postgres) GetAdImageIDs(adId *uuid.UUID) ([]string, error) {
	query, ok := p.q["GetAdImages"]
	if !ok {
		return nil, fmt.Errorf("request 'GetAdImages' not found")
	}

	imageRows, err := p.psql.Query(query, adId)
	if err != nil {
		return nil, fmt.Errorf("error while querying ad images: %w", err)
	}
	defer imageRows.Close()

	var imageIDs []string
	for imageRows.Next() {
		var imageID string
		if err := imageRows.Scan(&imageID); err != nil {
			return nil, fmt.Errorf("error while scanning image id: %w", err)
		}
		imageIDs = append(imageIDs, imageID)
	}

	return imageIDs, nil
}

func (p *Postgres) GetUserRole(userId *uuid.UUID) (string, error) {
	query, ok := p.q["GetUserRole"]
	if !ok {
		return "", fmt.Errorf("request 'GetUserRole' not found")
	}

	var role string
	err := p.psql.QueryRow(query, userId).Scan(&role)
	if err != nil {
		return "", fmt.Errorf("error while executing get user role request: %w", err)
	}
	return role, nil
}

func (p *Postgres) GetUserData(userId *uuid.UUID) (*models.UserData, error) {
	query, ok := p.q["GetUserData"]
	if !ok {
		return nil, fmt.Errorf("request 'GetUserData' not found")
	}

	var data models.UserData
	p.psql.QueryRow(query, userId).Scan(
		&data.Username,
		&data.Email,
		&data.PhoneNumber,
		&data.Settings,
	)
	return &data, nil
}

func (p *Postgres) GetAdsList(limit, offset int) ([]models.AdTemplate, error) {
	query, ok := p.q["GetAdsList"]
	if !ok {
		return nil, fmt.Errorf("request 'GetAdsList' not found")
	}

	rows, err := p.psql.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error while executing get ads list request: %w", err)
	}
	defer rows.Close()

	var ads []models.AdTemplate
	for rows.Next() {
		var ad models.AdTemplate
		var updatedAt, expirationDate *time.Time
		err := rows.Scan(
			&ad.AdId,
			&ad.UserId,
			&ad.CategoryId,
			&ad.CategoryName,
			&ad.LocationId,
			&ad.LocationName,
			&ad.Title,
			&ad.Description,
			&ad.Price,
			&ad.Condition,
			&ad.ContactPhone,
			&ad.CreatedAt,
			&updatedAt,
			&expirationDate,
			&ad.Status,
			&ad.ViewsCount,
		)
		if err != nil {
			return nil, fmt.Errorf("error while scanning ads list: %w", err)
		}
		ad.UpdatedAt = updatedAt
		ad.ExpirationDate = expirationDate
		ad.Price = ad.Price / 100
		ads = append(ads, ad)
	}

	return ads, nil
}

func (p *Postgres) GetAdsCount() (int, error) {
	query, ok := p.q["GetAdsCount"]
	if !ok {
		return 0, fmt.Errorf("request 'GetAdsCount' not found")
	}

	var count int
	err := p.psql.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error while executing get ads count request: %w", err)
	}

	return count, nil
}

func (p *Postgres) GetUserChats(userId *uuid.UUID) ([]models.ChatPreview, error) {
	query, ok := p.q["GetUserChats"]
	if !ok {
		return nil, fmt.Errorf("request 'GetUserChats' not found")
	}

	rows, err := p.psql.Query(query, userId)
	if err != nil {
		return nil, fmt.Errorf("error while executing get user chats request: %w", err)
	}
	defer rows.Close()

	var chats []models.ChatPreview
	for rows.Next() {
		var chat models.ChatPreview
		err := rows.Scan(
			&chat.UserID,
			&chat.Username,
			&chat.LastMessage,
			&chat.LastMessageTime,
			&chat.UnreadCount,
		)
		if err != nil {
			return nil, fmt.Errorf("error while scanning user chats: %w", err)
		}
		chat.ID = *userId     // This would need to be adjusted based on actual chat ID logic
		chat.IsOnline = false // This would need to be determined based on last activity
		chats = append(chats, chat)
	}

	return chats, nil
}

func (p *Postgres) GetLikedAds(userId *uuid.UUID) ([]models.AdTemplate, error) {
	query, ok := p.q["GetLikedAds"]
	if !ok {
		return nil, fmt.Errorf("request 'GetLikedAds' not found")
	}

	rows, err := p.psql.Query(query, userId)
	if err != nil {
		return nil, fmt.Errorf("error while executing get liked ads request: %w", err)
	}
	defer rows.Close()

	var ads []models.AdTemplate
	for rows.Next() {
		var ad models.AdTemplate
		var updatedAt, expirationDate *time.Time
		err := rows.Scan(
			&ad.AdId,
			&ad.UserId,
			&ad.CategoryId,
			&ad.CategoryName,
			&ad.LocationId,
			&ad.LocationName,
			&ad.Title,
			&ad.Description,
			&ad.Price,
			&ad.Condition,
			&ad.ContactPhone,
			&ad.CreatedAt,
			&updatedAt,
			&expirationDate,
			&ad.Status,
			&ad.ViewsCount,
			&ad.ImageID,
		)
		if err != nil {
			return nil, fmt.Errorf("error while scanning liked ads: %w", err)
		}
		ad.UpdatedAt = updatedAt
		ad.ExpirationDate = expirationDate
		ad.Price = ad.Price / 100
		ads = append(ads, ad)
	}

	return ads, nil
}

func (p *Postgres) GetAdminStats() (*models.AdminStats, error) {
	query, ok := p.q["GetAdminStats"]
	if !ok {
		return nil, fmt.Errorf("request 'GetAdminStats' not found")
	}

	var stats models.AdminStats
	err := p.psql.QueryRow(query).Scan(
		&stats.TotalAds,
		&stats.TotalUsers,
		&stats.PendingReports,
		&stats.PendingModeration,
	)
	if err != nil {
		return nil, fmt.Errorf("error while executing get admin stats request: %w", err)
	}

	return &stats, nil
}

func (p *Postgres) GetAllAds(limit, offset int) ([]models.AdTemplate, error) {
	query, ok := p.q["GetAllAds"]
	if !ok {
		return nil, fmt.Errorf("request 'GetAllAds' not found")
	}

	rows, err := p.psql.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error while executing get all ads request: %w", err)
	}
	defer rows.Close()

	var ads []models.AdTemplate
	for rows.Next() {
		var ad models.AdTemplate
		var updatedAt, expirationDate *time.Time
		err := rows.Scan(
			&ad.AdId,
			&ad.UserId,
			&ad.CategoryId,
			&ad.CategoryName,
			&ad.LocationId,
			&ad.LocationName,
			&ad.Title,
			&ad.Description,
			&ad.Price,
			&ad.Condition,
			&ad.ContactPhone,
			&ad.CreatedAt,
			&updatedAt,
			&expirationDate,
			&ad.Status,
			&ad.ViewsCount,
		)
		if err != nil {
			return nil, fmt.Errorf("error while scanning all ads: %w", err)
		}
		ad.UpdatedAt = updatedAt
		ad.ExpirationDate = expirationDate
		ad.Price = ad.Price / 100
		ads = append(ads, ad)
	}

	return ads, nil
}

func (p *Postgres) GetAllUsers(limit, offset int) ([]models.UserTemplate, error) {
	query, ok := p.q["GetAllUsers"]
	if !ok {
		return nil, fmt.Errorf("request 'GetAllUsers' not found")
	}

	rows, err := p.psql.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error while executing get all users request: %w", err)
	}
	defer rows.Close()

	var users []models.UserTemplate
	for rows.Next() {
		var user models.UserTemplate
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.PhoneNumber,
			&user.RegistrationDate,
			&user.Status,
			&user.Role,
		)
		if err != nil {
			return nil, fmt.Errorf("error while scanning all users: %w", err)
		}
		users = append(users, user)
	}

	return users, nil
}

func (p *Postgres) GetComplaints(limit, offset int) ([]models.ComplaintTemplate, error) {
	query, ok := p.q["GetComplaints"]
	if !ok {
		return nil, fmt.Errorf("request 'GetComplaints' not found")
	}

	rows, err := p.psql.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error while executing get complaints request: %w", err)
	}
	defer rows.Close()

	var complaints []models.ComplaintTemplate
	for rows.Next() {
		var complaint models.ComplaintTemplate
		err := rows.Scan(
			&complaint.ID,
			&complaint.ListingID,
			&complaint.TargetUserID,
			&complaint.ComplainantID,
			&complaint.ComplaintType,
			&complaint.Description,
			&complaint.Status,
			&complaint.CreatedAt,
			&complaint.UpdatedAt,
			&complaint.AdminID,
			&complaint.ResolutionComment,
		)
		if err != nil {
			return nil, fmt.Errorf("error while scanning complaints: %w", err)
		}
		complaints = append(complaints, complaint)
	}

	return complaints, nil
}

func (p *Postgres) GetModerationAds(limit, offset int) ([]models.AdTemplate, error) {
	query, ok := p.q["GetModerationAds"]
	if !ok {
		return nil, fmt.Errorf("request 'GetModerationAds' not found")
	}

	rows, err := p.psql.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error while executing get moderation ads request: %w", err)
	}
	defer rows.Close()

	var ads []models.AdTemplate
	for rows.Next() {
		var ad models.AdTemplate
		var updatedAt, expirationDate *time.Time
		err := rows.Scan(
			&ad.AdId,
			&ad.UserId,
			&ad.CategoryId,
			&ad.CategoryName,
			&ad.LocationId,
			&ad.LocationName,
			&ad.Title,
			&ad.Description,
			&ad.Price,
			&ad.Condition,
			&ad.ContactPhone,
			&ad.CreatedAt,
			&updatedAt,
			&expirationDate,
			&ad.Status,
			&ad.ViewsCount,
		)
		if err != nil {
			return nil, fmt.Errorf("error while scanning moderation ads: %w", err)
		}
		ad.UpdatedAt = updatedAt
		ad.ExpirationDate = expirationDate
		ad.Price = ad.Price / 100
		ads = append(ads, ad)
	}

	return ads, nil
}

func (p *Postgres) UpdatePassword(userId *uuid.UUID, newPasswordHash string) error {
	query, ok := p.q["UpdateUserPassword"]
	if !ok {
		return fmt.Errorf("request 'UpdateUserPassword' not found")
	}

	_, err := p.psql.Exec(query, newPasswordHash, userId)
	if err != nil {
		return fmt.Errorf("error while executing update password request: %w", err)
	}

	return nil
}

func (p *Postgres) DeleteAccount(userId *uuid.UUID) error {
	query, ok := p.q["DeleteUserAccount"]
	if !ok {
		return fmt.Errorf("request 'DeleteUserAccount' not found")
	}

	_, err := p.psql.Exec(query, userId)
	if err != nil {
		return fmt.Errorf("error while executing delete account request: %w", err)
	}

	return nil
}

func (p *Postgres) GetAdsListByCategory(category string, limit, offset int) ([]models.AdTemplate, error) {
	query, ok := p.q["GetAdsListByCategory"]
	if !ok {
		return nil, fmt.Errorf("request 'GetAdsListByCategory' not found")
	}

	rows, err := p.psql.Query(query, category, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error while executing get ads list by category request: %w", err)
	}
	defer rows.Close()

	var ads []models.AdTemplate
	for rows.Next() {
		var ad models.AdTemplate
		var updatedAt, expirationDate *time.Time
		err := rows.Scan(
			&ad.AdId,
			&ad.UserId,
			&ad.CategoryId,
			&ad.CategoryName,
			&ad.LocationId,
			&ad.LocationName,
			&ad.Title,
			&ad.Description,
			&ad.Price,
			&ad.Condition,
			&ad.ContactPhone,
			&ad.CreatedAt,
			&updatedAt,
			&expirationDate,
			&ad.Status,
			&ad.ViewsCount,
		)
		if err != nil {
			return nil, fmt.Errorf("error while scanning ads list by category: %w", err)
		}
		ad.UpdatedAt = updatedAt
		ad.ExpirationDate = expirationDate
		ad.Price = ad.Price / 100
		ads = append(ads, ad)
	}

	return ads, nil
}

func (p *Postgres) GetAdsCountByCategory(category string) (int, error) {
	query, ok := p.q["GetAdsCountByCategory"]
	if !ok {
		return 0, fmt.Errorf("request 'GetAdsCountByCategory' not found")
	}

	var count int
	err := p.psql.QueryRow(query, category).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error while executing get ads count by category request: %w", err)
	}

	return count, nil
}

func (p *Postgres) SearchAds(filters *models.AdSearchFilters) ([]models.AdTemplate, error) {
	query, ok := p.q["SearchAdsBase"]
	if !ok {
		return nil, fmt.Errorf("request 'SearchAdsBase' not found")
	}

	var args []interface{}
	argIndex := 1

	if filters.Category != "" && filters.Subcategory == "" {
		query += fmt.Sprintf(" AND (c.name = $%d OR c.parent_id = (SELECT id FROM site.categories WHERE name = $%d))", argIndex, argIndex)
		args = append(args, filters.Category)
		argIndex++
	}

	if filters.Region != "" {
		query += fmt.Sprintf(" AND loc.region = $%d", argIndex)
		args = append(args, filters.Region)
		argIndex++
	}

	if filters.Location != "" {
		query += fmt.Sprintf(" AND loc.city ILIKE $%d", argIndex)
		args = append(args, "%"+filters.Location+"%")
		argIndex++
	}

	if filters.Condition != "" {
		query += fmt.Sprintf(" AND l.condition = $%d", argIndex)
		args = append(args, filters.Condition)
		argIndex++
	}

	if filters.MinPrice != nil {
		query += fmt.Sprintf(" AND l.price >= $%d", argIndex)
		args = append(args, *filters.MinPrice*100)
		argIndex++
	}

	if filters.MaxPrice != nil {
		query += fmt.Sprintf(" AND l.price <= $%d", argIndex)
		args = append(args, *filters.MaxPrice*100)
		argIndex++
	}

	if filters.Subcategory != "" {
		query += fmt.Sprintf(" AND c.name = $%d", argIndex)
		args = append(args, filters.Subcategory)
		argIndex++
	}

	if filters.SearchQuery != "" {
		query += fmt.Sprintf(" AND (l.title ILIKE $%d OR l.description ILIKE $%d)", argIndex, argIndex)
		args = append(args, "%"+filters.SearchQuery+"%")
		argIndex++
	}

	sortOrder := "l.created_at DESC"
	switch filters.SortBy {
	case "oldest":
		sortOrder = "l.created_at ASC"
	case "price-low":
		sortOrder = "l.price ASC"
	case "price-high":
		sortOrder = "l.price DESC"
	case "popular":
		sortOrder = "l.views_count DESC"
	}
	query += " ORDER BY " + sortOrder
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, filters.Limit, (filters.Page-1)*filters.Limit)

	rows, err := p.psql.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error while executing search ads request: %w", err)
	}
	defer rows.Close()

	var ads []models.AdTemplate
	for rows.Next() {
		var ad models.AdTemplate
		var updatedAt, expirationDate *time.Time
		err := rows.Scan(
			&ad.AdId,
			&ad.UserId,
			&ad.CategoryId,
			&ad.CategoryName,
			&ad.LocationId,
			&ad.LocationName,
			&ad.Title,
			&ad.Description,
			&ad.Price,
			&ad.Condition,
			&ad.ContactPhone,
			&ad.CreatedAt,
			&updatedAt,
			&expirationDate,
			&ad.Status,
			&ad.ViewsCount,
			&ad.ImageID,
		)
		if err != nil {
			return nil, fmt.Errorf("error while scanning search ads: %w", err)
		}
		ad.UpdatedAt = updatedAt
		ad.ExpirationDate = expirationDate
		ad.Price = ad.Price / 100
		ads = append(ads, ad)
	}

	return ads, nil
}

func (p *Postgres) SearchAdsCount(filters *models.AdSearchFilters) (int, error) {
	query, ok := p.q["SearchAdsCountBase"]
	if !ok {
		return 0, fmt.Errorf("request 'SearchAdsCountBase' not found")
	}

	var args []interface{}
	argIndex := 1

	if filters.Category != "" && filters.Subcategory == "" {
		query += fmt.Sprintf(" AND (c.name = $%d OR c.parent_id = (SELECT id FROM site.categories WHERE name = $%d))", argIndex, argIndex)
		args = append(args, filters.Category)
		argIndex++
	}

	if filters.Region != "" {
		query += fmt.Sprintf(" AND loc.region = $%d", argIndex)
		args = append(args, filters.Region)
		argIndex++
	}

	if filters.Location != "" {
		query += fmt.Sprintf(" AND loc.city ILIKE $%d", argIndex)
		args = append(args, "%"+filters.Location+"%")
		argIndex++
	}

	if filters.Condition != "" {
		query += fmt.Sprintf(" AND l.condition = $%d", argIndex)
		args = append(args, filters.Condition)
		argIndex++
	}

	if filters.MinPrice != nil {
		query += fmt.Sprintf(" AND l.price >= $%d", argIndex)
		args = append(args, *filters.MinPrice*100)
		argIndex++
	}

	if filters.MaxPrice != nil {
		query += fmt.Sprintf(" AND l.price <= $%d", argIndex)
		args = append(args, *filters.MaxPrice*100)
		argIndex++
	}

	if filters.Subcategory != "" {
		query += fmt.Sprintf(" AND c.name = $%d", argIndex)
		args = append(args, filters.Subcategory)
		argIndex++
	}

	if filters.SearchQuery != "" {
		query += fmt.Sprintf(" AND (l.title ILIKE $%d OR l.description ILIKE $%d)", argIndex, argIndex)
		args = append(args, "%"+filters.SearchQuery+"%")
		argIndex++
	}

	var count int
	err := p.psql.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error while executing search ads count request: %w", err)
	}

	return count, nil
}

func (p *Postgres) AddToFavorites(userId, adId *uuid.UUID) error {
	query, ok := p.q["AddToFavorites"]
	if !ok {
		return fmt.Errorf("request 'AddToFavorites' not found")
	}

	_, err := p.psql.Exec(query, userId, adId)
	if err != nil {
		return fmt.Errorf("error while adding to favorites: %w", err)
	}

	return nil
}

func (p *Postgres) RemoveFromFavorites(userId, adId *uuid.UUID) error {
	query, ok := p.q["RemoveFromFavorites"]
	if !ok {
		return fmt.Errorf("request 'RemoveFromFavorites' not found")
	}

	res, err := p.psql.Exec(query, userId, adId)
	if err != nil {
		return fmt.Errorf("error while removing from favorites: %w", err)
	}

	rowsAff, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error while getting affected rows: %w", err)
	}

	if rowsAff == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (p *Postgres) CheckIfFavorite(userId, adId *uuid.UUID) (bool, error) {
	query, ok := p.q["CheckIfFavorite"]
	if !ok {
		return false, fmt.Errorf("request 'CheckIfFavorite' not found")
	}

	var isFavorite bool
	err := p.psql.QueryRow(query, userId, adId).Scan(&isFavorite)
	if err != nil {
		return false, fmt.Errorf("error while checking favorite status: %w", err)
	}

	return isFavorite, nil
}

func (p *Postgres) GetUserAds(userId *uuid.UUID) ([]models.AdTemplate, error) {
	query, ok := p.q["GetUserAds"]
	if !ok {
		return nil, fmt.Errorf("request 'GetUserAds' not found")
	}

	rows, err := p.psql.Query(query, userId)
	if err != nil {
		return nil, fmt.Errorf("error while executing get user ads request: %w", err)
	}
	defer rows.Close()

	var ads []models.AdTemplate
	for rows.Next() {
		var ad models.AdTemplate
		var updatedAt, expirationDate *time.Time
		err := rows.Scan(
			&ad.AdId,
			&ad.UserId,
			&ad.CategoryId,
			&ad.CategoryName,
			&ad.LocationId,
			&ad.LocationName,
			&ad.Title,
			&ad.Description,
			&ad.Price,
			&ad.Condition,
			&ad.ContactPhone,
			&ad.CreatedAt,
			&updatedAt,
			&expirationDate,
			&ad.Status,
			&ad.ViewsCount,
			&ad.ImageID,
		)
		if err != nil {
			return nil, fmt.Errorf("error while scanning user ads: %w", err)
		}
		ad.UpdatedAt = updatedAt
		ad.ExpirationDate = expirationDate
		ad.Price = ad.Price / 100
		ads = append(ads, ad)
	}

	return ads, nil
}

func (p *Postgres) UpdateUserProfile(userId *uuid.UUID, username string, phoneNumber *string) error {
	query, ok := p.q["UpdateUserProfile"]
	if !ok {
		return fmt.Errorf("request 'UpdateUserProfile' not found")
	}

	_, err := p.psql.Exec(query, username, phoneNumber, userId)
	if err != nil {
		return fmt.Errorf("error while updating user profile: %w", err)
	}

	return nil
}
