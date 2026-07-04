package handlers

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"url_shortner/database"
	"url_shortner/models"
)

func generateShortURL(OriginalURL string) string {
	hasher := md5.New()
	hasher.Write([]byte(OriginalURL))
	fmt.Println("converted into bits\n", hasher)
	data := hasher.Sum(nil)
	fmt.Println(data)
	return hex.EncodeToString(data[:8])

}
func CreateURL(originalURL string, userID int) (string, error) {
	shortURL := generateShortURL(originalURL)
	query := `
	INSERT INTO urls (
		original_url,
		short_url,
		user_id
	)
	VALUES ($1,$2,$3)
	`
	_, err := database.DB.Exec(
		context.Background(),
		query,
		originalURL,
		shortURL,
		userID,
	)
	if err != nil {
		return "", err
	}

	// Store in Redis for 24 hours
	err = database.RedisClient.Set(
		database.Ctx,
		shortURL,
		originalURL,
		24*time.Hour,
	).Err()

	if err != nil {
		fmt.Println("Redis SET error:", err)
	}

	return shortURL, nil
}
func getURL(shortURL string) (models.URL, error) {
	fmt.Println("📦 Fetching from PostgreSQL")

	var url models.URL

	query := `
	SELECT
		id,
		original_url,
		short_url,
		creation_date
	FROM urls
	WHERE short_url = $1
	`

	err := database.DB.QueryRow(
		context.Background(),
		query,
		shortURL,
	).Scan(
		&url.ID,
		&url.OriginalURL,
		&url.ShortenURL,
		&url.CreationDate,
	)

	if err != nil {
		return models.URL{}, err
	}

	return url, nil
}
func HandleRedirect(w http.ResponseWriter, r *http.Request) {
	fmt.Println("➡️ HandleRedirect called:", r.Method, r.URL.Path)

	id := r.URL.Path[1:]

	originalURL, err := database.RedisClient.Get(
		database.Ctx,
		id,
	).Result()

	if err == nil {
		fmt.Println("✅ Cache Hit")
		http.Redirect(w, r, originalURL, http.StatusFound)
		return
	}

	fmt.Println("❌ Cache Miss")

	urlData, err := getURL(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	err = database.RedisClient.Set(
		database.Ctx,
		id,
		urlData.OriginalURL,
		24*time.Hour,
	).Err()

	if err != nil {
		fmt.Println("Redis SET error:", err)
	}
	http.Redirect(w, r, urlData.OriginalURL, http.StatusFound)
}
func ShortenURLHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		URL string `json:"url"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	userID := r.Context().
		Value("userID").(int)

	fmt.Println("Logged in User ID:", userID)
	shortURL, err :=
		CreateURL(
			data.URL,
			userID,
		)
	if err != nil {
		fmt.Println("DATABASE ERROR:", err)

		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	fmt.Fprintln(w, shortURL)

}
func MyURLsHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)

	query := `
	SELECT
		id,
		original_url,
		short_url,
		creation_date
	FROM urls
	WHERE user_id=$1
	`

	rows, err := database.DB.Query(context.Background(), query, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var urls []models.URL
	for rows.Next() {
		var url models.URL
		if err := rows.Scan(&url.ID, &url.OriginalURL, &url.ShortenURL, &url.CreationDate); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		urls = append(urls, url)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(urls)
}
