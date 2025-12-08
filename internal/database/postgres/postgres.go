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

func uuidPointer(u uuid.UUID) *uuid.UUID {
	val := u
	return &val
}

func timePointer(t time.Time) *time.Time {
	val := t
	return &val
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
	var status string

	err := p.psql.QueryRow(query, email).Scan(&id, &passHash, &confirmed, &status)

	if err != nil {
		return nil, fmt.Errorf("error while trying to get user auth data by email: %w", err)
	}

	return &models.AuhtInfo{
		Id:           id,
		PasswordHash: passHash,
		Confirmed:    confirmed,
		Status:       status,
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
	data.AdId = *adId
	var sellerPhone sql.NullString
	var country, region, city, postalCode sql.NullString

	err := p.psql.QueryRow(query, adId).Scan(
		&data.Title,
		&data.Price,
		&data.CreatedAt,
		&data.ViewsCount,
		&data.Condition,
		&data.Description,
		&data.ContactPhone,
		&data.AdStatus,
		&data.SellerName,
		&data.Online,
		&data.UserID,
		&country,
		&region,
		&city,
		&postalCode,
		&data.SellerAvatarID,
		&sellerPhone,
	)

	// Формируем SellerCity только из непустых значений
	var cityParts []string
	if country.Valid && country.String != "" {
		cityParts = append(cityParts, country.String)
	}
	if region.Valid && region.String != "" {
		cityParts = append(cityParts, region.String)
	}
	if city.Valid && city.String != "" {
		cityParts = append(cityParts, city.String)
	}
	if postalCode.Valid && postalCode.String != "" {
		cityParts = append(cityParts, postalCode.String)
	}
	data.SellerCity = strings.Join(cityParts, " ")

	data.Price /= 100
	if sellerPhone.Valid && sellerPhone.String != "" {
		data.SellerPhone = sellerPhone.String
	}

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

func (p *Postgres) GetUserStatus(userId *uuid.UUID) (string, error) {
	query, ok := p.q["GetUserStatus"]
	if !ok {
		return "", fmt.Errorf("request 'GetUserStatus' not found")
	}
	var status string
	err := p.psql.QueryRow(query, userId).Scan(&status)
	if err != nil {
		return "", fmt.Errorf("error while executing get user status request: %w", err)
	}
	return status, nil
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

func (p *Postgres) UpdateUserAvatar(userId *uuid.UUID, avatarID string) error {
	_, err := p.psql.Exec(`UPDATE site.users SET profile_picture_id = $1 WHERE id = $2`, avatarID, userId)
	if err != nil {
		return fmt.Errorf("error while updating user avatar: %w", err)
	}
	return nil
}

func (p *Postgres) GetUserAvatar(userId *uuid.UUID) (*string, error) {
	var avatarID *string
	err := p.psql.QueryRow(`SELECT profile_picture_id FROM site.users WHERE id = $1`, userId).Scan(&avatarID)
	if err != nil {
		return nil, fmt.Errorf("error while getting user avatar: %w", err)
	}
	return avatarID, nil
}

func (p *Postgres) CreateChat(chat *models.Chat) (*models.Chat, error) {
	query, ok := p.q["CreateChat"]
	if !ok {
		return nil, fmt.Errorf("request 'CreateChat' not found")
	}

	var chatID uuid.UUID
	err := p.psql.QueryRow(
		query,
		chat.CustomerId,
		chat.ListingId,
	).Scan(&chatID)
	if err != nil {
		return nil, fmt.Errorf("error while trying to create new chat: %w", err)
	}

	return p.GetChatById(chatID)
}

func (p *Postgres) GetChatById(chatID uuid.UUID) (*models.Chat, error) {
	query, ok := p.q["GetChatById"]
	if !ok {
		return nil, fmt.Errorf("request 'GetChatId' not found")
	}

	var (
		id         uuid.UUID
		sellerID   uuid.UUID
		customerID uuid.UUID
		listingID  uuid.UUID
		createdAt  time.Time
	)
	err := p.psql.QueryRow(query, chatID).Scan(
		&id,
		&sellerID,
		&customerID,
		&listingID,
		&createdAt,
	)

	if err != nil {
		return nil, fmt.Errorf("error while trying to get chat data: %w", err)
	}

	chat := models.Chat{
		ChatId:     uuidPointer(id),
		SellerId:   uuidPointer(sellerID),
		CustomerId: uuidPointer(customerID),
		ListingId:  uuidPointer(listingID),
		CreatedAt:  timePointer(createdAt),
	}

	return &chat, nil
}

func (p *Postgres) CreateMessage(msg *models.Message) error {
	query, ok := p.q["CreateMessage"]
	if !ok {
		return fmt.Errorf("request 'CreateMessage' not found")
	}

	if msg.ChatId == nil {
		return fmt.Errorf("chat id is required for message")
	}
	if msg.SenderId == nil {
		return fmt.Errorf("sender id is required for message")
	}

	var messageID uuid.UUID
	var createdAt time.Time
	err := p.psql.QueryRow(
		query,
		*msg.ChatId,
		*msg.SenderId,
		msg.Text,
	).Scan(&messageID, &createdAt)
	if err != nil {
		return fmt.Errorf("error while trying to create new message: %w", err)
	}
	msg.Id = uuidPointer(messageID)
	msg.CreatedAt = timePointer(createdAt)
	return nil
}

func (p *Postgres) GetMessages(chatId uuid.UUID, limit, offset int) ([]models.Message, error) {
	query, ok := p.q["GetMessages"]
	if !ok {
		return nil, fmt.Errorf("request 'GetMessages' not found")
	}

	rows, err := p.psql.Query(query, chatId, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error while trying to execute query: %w", err)
	}
	defer rows.Close()

	var messages []models.Message

	for rows.Next() {
		var m models.Message
		var (
			id        uuid.UUID
			senderID  uuid.UUID
			createdAt time.Time
			chatUUID  uuid.UUID
		)

		if err := rows.Scan(&id, &senderID, &m.Text, &createdAt, &chatUUID); err != nil {
			return nil, err
		}
		m.Id = uuidPointer(id)
		m.SenderId = uuidPointer(senderID)
		m.CreatedAt = timePointer(createdAt)
		m.ChatId = uuidPointer(chatUUID)
		messages = append(messages, m)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return messages, nil
}

func (p *Postgres) GetUserChats(userId *uuid.UUID) ([]models.ChatPreview, error) {
	query, ok := p.q["GetUserChats"]
	if !ok {
		return nil, fmt.Errorf("request 'GetUserChats' not found")
	}

	rows, err := p.psql.Query(query, userId)
	if err != nil {
		return nil, fmt.Errorf("error while trying to execute query: %w", err)
	}
	defer rows.Close()

	var chats []models.ChatPreview

	for rows.Next() {
		var chat models.ChatPreview
		var avatarID sql.NullString
		var lastMessageTime sql.NullTime

		if err := rows.Scan(
			&chat.ChatId,
			&chat.ListingId,
			&chat.LastMessage,
			&lastMessageTime,
			&chat.CompanionId,
			&chat.CompanionName,
			&avatarID,
			&chat.ListingTitle,
		); err != nil {
			return nil, err
		}

		if avatarID.Valid {
			value := avatarID.String
			chat.CompanionAvatarId = &value
		}

		if lastMessageTime.Valid {
			t := lastMessageTime.Time
			chat.LastMessageTime = &t
		}

		chats = append(chats, chat)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return chats, nil
}

func (p *Postgres) GetAdminStats() (*models.AdminStats, error) {
	var stats models.AdminStats

	if err := p.psql.QueryRow(`
		SELECT
			(SELECT COUNT(*) FROM site.listings),
			(SELECT COUNT(*) FROM site.users WHERE status = 'active'),
			(SELECT COUNT(*) FROM site.complaints WHERE status = 'pending'),
			(SELECT COUNT(*) FROM site.listings WHERE status = 'moderation')
	`).Scan(&stats.TotalAds, &stats.TotalUsers, &stats.PendingReports, &stats.PendingModeration); err != nil {
		return nil, fmt.Errorf("error while getting admin stats: %w", err)
	}

	return &stats, nil
}

func (p *Postgres) UpdateAdStatus(adId *uuid.UUID, status string) error {
	_, err := p.psql.Exec(`UPDATE site.listings SET status = $1, updated_at = NOW() WHERE id = $2`, status, adId)
	if err != nil {
		return fmt.Errorf("error while updating ad status: %w", err)
	}
	return nil
}

func (p *Postgres) UpdateUserStatus(userId *uuid.UUID, status string) error {
	_, err := p.psql.Exec(`UPDATE site.users SET status = $1 WHERE id = $2`, status, userId)
	if err != nil {
		return fmt.Errorf("error while updating user status: %w", err)
	}
	return nil
}

func (p *Postgres) UpdateUserRole(userId *uuid.UUID, role string) error {
	query, ok := p.q["UpdateUserRole"]
	if !ok {
		return fmt.Errorf("request 'UpdateUserRole' not found")
	}
	_, err := p.psql.Exec(query, role, userId)
	if err != nil {
		return fmt.Errorf("error while updating user role: %w", err)
	}
	return nil
}

func (p *Postgres) UpdateLastLogin(userId *uuid.UUID) error {
	query, ok := p.q["UpdateLastLogin"]
	if !ok {
		return fmt.Errorf("request 'UpdateLastLogin' not found")
	}
	_, err := p.psql.Exec(query, userId)
	if err != nil {
		return fmt.Errorf("error while updating last login: %w", err)
	}
	return nil
}

func (p *Postgres) GetUsers(page, limit int, roleFilter, search string) ([]models.UserTemplate, error) {
	query := `
		SELECT id, username, email, phone_number, registration_date, status, role
		FROM site.users
		WHERE 1=1
	`
	var args []interface{}
	argIndex := 1

	if roleFilter != "" && roleFilter != "all" {
		query += fmt.Sprintf(" AND role = $%d", argIndex)
		args = append(args, roleFilter)
		argIndex++
	}

	if search != "" {
		query += fmt.Sprintf(" AND (email ILIKE $%d OR username ILIKE $%d)", argIndex, argIndex)
		args = append(args, "%"+search+"%")
		argIndex++
	}

	query += " ORDER BY registration_date DESC"
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, (page-1)*limit)

	rows, err := p.psql.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error while getting users: %w", err)
	}
	defer rows.Close()

	var users []models.UserTemplate
	for rows.Next() {
		var user models.UserTemplate
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.PhoneNumber, &user.RegistrationDate, &user.Status, &user.Role); err != nil {
			return nil, fmt.Errorf("error while scanning user: %w", err)
		}
		users = append(users, user)
	}

	return users, nil
}

func (p *Postgres) GetReports(page, limit int, statusFilter string) ([]models.ComplaintTemplate, error) {
	queryBase, ok := p.q["GetReportsBase"]
	if !ok {
		return nil, fmt.Errorf("request 'GetReportsBase' not found")
	}

	query := queryBase
	var args []interface{}
	argIndex := 1

	if statusFilter != "" && statusFilter != "all" {
		query += fmt.Sprintf(" AND c.status = $%d", argIndex)
		args = append(args, statusFilter)
		argIndex++
	} else {
		query += " AND c.status = 'pending'"
	}

	query += " ORDER BY c.created_at DESC"
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, (page-1)*limit)

	rows, err := p.psql.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error while getting reports: %w", err)
	}
	defer rows.Close()

	var reports []models.ComplaintTemplate
	for rows.Next() {
		var report models.ComplaintTemplate
		if err := rows.Scan(
			&report.ID,
			&report.ListingID,
			&report.TargetUserID,
			&report.ComplainantID,
			&report.ComplaintType,
			&report.Description,
			&report.Status,
			&report.CreatedAt,
			&report.UpdatedAt,
			&report.AdminID,
			&report.ResolutionComment,
			&report.ComplainantEmail,
			&report.AdOwnerID,
		); err != nil {
			return nil, fmt.Errorf("error while scanning report: %w", err)
		}
		reports = append(reports, report)
	}

	return reports, nil
}

func (p *Postgres) UpdateReportStatus(reportId *uuid.UUID, status, resolutionComment string, adminId *uuid.UUID) error {
	query, ok := p.q["UpdateReportStatus"]
	if !ok {
		return fmt.Errorf("request 'UpdateReportStatus' not found")
	}
	_, err := p.psql.Exec(query, status, resolutionComment, adminId, reportId)
	if err != nil {
		return fmt.Errorf("error while updating report status: %w", err)
	}
	return nil
}

func (p *Postgres) CreateReport(userId, adId *uuid.UUID, reportType, description string) error {
	query, ok := p.q["CreateReport"]
	if !ok {
		return fmt.Errorf("request 'CreateReport' not found")
	}
	_, err := p.psql.Exec(query, adId, userId, reportType, description)
	if err != nil {
		return fmt.Errorf("error while creating report: %w", err)
	}
	return nil
}

func (p *Postgres) GetAdminAds(page, limit int, status, search string) ([]map[string]interface{}, int, error) {
	queryBase, ok := p.q["GetAdminAdsBase"]
	if !ok {
		return nil, 0, fmt.Errorf("request 'GetAdminAdsBase' not found")
	}
	countQueryBase, ok := p.q["GetAdminAdsCountBase"]
	if !ok {
		return nil, 0, fmt.Errorf("request 'GetAdminAdsCountBase' not found")
	}

	var args []interface{}
	argIndex := 1
	query := queryBase
	countQuery := countQueryBase

	if status != "" && status != "all" {
		query += fmt.Sprintf(" AND l.status = $%d", argIndex)
		countQuery += fmt.Sprintf(" AND l.status = $%d", argIndex)
		args = append(args, status)
		argIndex++
	}

	if search != "" {
		query += fmt.Sprintf(" AND (l.title ILIKE $%d OR l.description ILIKE $%d)", argIndex, argIndex)
		countQuery += fmt.Sprintf(" AND (l.title ILIKE $%d OR l.description ILIKE $%d)", argIndex, argIndex)
		args = append(args, "%"+search+"%")
		argIndex++
	}

	query += " ORDER BY l.created_at DESC"
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	queryArgs := append(args, limit, (page-1)*limit)

	rows, err := p.psql.Query(query, queryArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("error while getting admin ads: %w", err)
	}
	defer rows.Close()

	var ads []map[string]interface{}
	for rows.Next() {
		var id, userId, title, userEmail, adStatus string
		var createdAt time.Time
		if err := rows.Scan(&id, &userId, &title, &createdAt, &userEmail, &adStatus); err != nil {
			return nil, 0, fmt.Errorf("error while scanning admin ad: %w", err)
		}
		ads = append(ads, map[string]interface{}{
			"id":         id,
			"user_id":    userId,
			"title":      title,
			"created_at": createdAt.Format("2006-01-02T15:04:05Z"),
			"user_email": userEmail,
			"ad_status":  adStatus,
		})
	}

	var totalCount int
	err = p.psql.QueryRow(countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("error while getting admin ads count: %w", err)
	}

	return ads, totalCount, nil
}

func (p *Postgres) DeleteAd(adId *uuid.UUID) error {
	_, err := p.psql.Exec(`DELETE FROM site.listings WHERE id = $1`, adId)
	if err != nil {
		return fmt.Errorf("error while deleting ad: %w", err)
	}
	return nil
}

func (p *Postgres) UpdateAd(adId *uuid.UUID, title, description string, price int, condition, contactPhone string) error {
	_, err := p.psql.Exec(`
		UPDATE site.listings 
		SET title = $1, description = $2, price = $3, condition = $4, contact_phone = $5, status = 'moderation', updated_at = NOW()
		WHERE id = $6
	`, title, description, price*100, condition, contactPhone, adId)
	if err != nil {
		return fmt.Errorf("error while updating ad: %w", err)
	}
	return nil
}

func (p *Postgres) CheckAdOwnership(userId, adId *uuid.UUID) (bool, error) {
	var isOwner bool
	err := p.psql.QueryRow(`SELECT user_id = $1 FROM site.listings WHERE id = $2`, userId, adId).Scan(&isOwner)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("error while checking ad ownership: %w", err)
	}
	return isOwner, nil
}

func (p *Postgres) DeleteAdImage(imageId string) error {
	_, err := p.psql.Exec(`DELETE FROM site.listing_images WHERE image_id = $1`, imageId)
	if err != nil {
		return fmt.Errorf("error while deleting image record: %w", err)
	}
	return nil
}

func (p *Postgres) CreatePasswordResetToken(email, token string) error {
	query, ok := p.q["CreatePasswordResetToken"]
	if !ok {
		return fmt.Errorf("request 'CreatePasswordResetToken' not found")
	}

	expiresAt := time.Now().Add(1 * time.Hour)
	_, err := p.psql.Exec(query, token, expiresAt, email)
	if err != nil {
		return fmt.Errorf("error while creating password reset token: %w", err)
	}
	return nil
}

func (p *Postgres) GetPasswordResetToken(token string) (*models.PasswordResetToken, error) {
	query, ok := p.q["GetPasswordResetToken"]
	if !ok {
		return nil, fmt.Errorf("request 'GetPasswordResetToken' not found")
	}

	var resetToken models.PasswordResetToken
	err := p.psql.QueryRow(query, token).Scan(&resetToken.Token, &resetToken.UserID, &resetToken.ExpiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error while getting password reset token: %w", err)
	}
	return &resetToken, nil
}

func (p *Postgres) DeletePasswordResetToken(token string) error {
	query, ok := p.q["DeletePasswordResetToken"]
	if !ok {
		return fmt.Errorf("request 'DeletePasswordResetToken' not found")
	}

	_, err := p.psql.Exec(query, token)
	if err != nil {
		return fmt.Errorf("error while deleting password reset token: %w", err)
	}
	return nil
}

func (p *Postgres) GetUserIdByEmail(email string) (*uuid.UUID, error) {
	query, ok := p.q["GetUserIdByEmail"]
	if !ok {
		return nil, fmt.Errorf("request 'GetUserIdByEmail' not found")
	}

	var userId uuid.UUID
	err := p.psql.QueryRow(query, email).Scan(&userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error while getting user id by email: %w", err)
	}
	return &userId, nil
}
