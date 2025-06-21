package handler

import (
	"fmt"
	"gitpulse/pkg/github"
	"gitpulse/pkg/svg"
	"net/http"
	"regexp"
	"sync"
	"time"
)

type CacheEntry struct {
	SVGData   string
	Timestamp time.Time
}

var (
	cache      = make(map[string]CacheEntry)
	cacheMutex = &sync.Mutex{}
	cacheTTL   = 8 * time.Hour
)

func Handler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		sendErrorResponse(w, "Username is required", http.StatusBadRequest)
		return
	}

	color := r.URL.Query().Get("color")
	if color == "" || !isValidHexColor(color) {
		color = "#007acc"
	}
	if color[0] != '#' {
		color = "#" + color
	}

	cacheKey := fmt.Sprintf("%s:%s", username, color)

	cacheMutex.Lock()
	entry, found := cache[cacheKey]
	if found && time.Since(entry.Timestamp) < cacheTTL {
		cacheMutex.Unlock()
		sendSuccessResponse(w, entry.SVGData)
		return
	}
	cacheMutex.Unlock()

	repo, err := github.GetLastUpdatedRepo(username)
	if err != nil {
		sendErrorResponse(w, "User or repo not found", http.StatusNotFound)
		return
	}

	language, err := github.GetRepoPrimaryLanguage(username, repo.Name)
	if err != nil {
		sendErrorResponse(w, "Could not get language", http.StatusInternalServerError)
		return
	}

	badge, err := svg.GenerateBadge(repo.Name, language, color)
	if err != nil {
		sendErrorResponse(w, "Failed to generate SVG", http.StatusInternalServerError)
		return
	}

	cacheMutex.Lock()
	cache[cacheKey] = CacheEntry{
		SVGData:   badge,
		Timestamp: time.Now(),
	}
	cacheMutex.Unlock()

	sendSuccessResponse(w, badge)
}

func sendSuccessResponse(w http.ResponseWriter, svgData string) {
	w.Header().Set("Content-Type", "image/svg+xml; charset=utf-8")
	w.Header().Set("Cache-Control", "s-maxage=3600, stale-while-revalidate")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, svgData)
}

func sendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "image/svg+xml; charset=utf-8")
	w.WriteHeader(statusCode)
	fmt.Fprint(w, svg.GenerateErrorBadge(message))
}

func isValidHexColor(s string) bool {
	if s[0] == '#' {
		s = s[1:]
	}
	match, _ := regexp.MatchString(`^([0-9a-fA-F]{3}|[0-9a-fA-F]{6})$`, s)
	return match
}
