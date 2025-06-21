package handler

import (
	"fmt"
	"gitpulse/pkg/github"
	"gitpulse/pkg/svg"
	"net/http"
	"regexp"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		w.Header().Set("Content-Type", "image/svg+xml")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, svg.GenerateErrorBadge("Username is required"))
		return
	}

	color := r.URL.Query().Get("color")
	if color == "" || !isValidHexColor(color) {
		color = "#007acc"
	}

	repo, err := github.GetLastUpdatedRepo(username)
	if err != nil {
		w.Header().Set("Content-Type", "image/svg+xml")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, svg.GenerateErrorBadge("User or repo not found"))
		return
	}

	language, err := github.GetRepoPrimaryLanguage(username, repo.Name)
	if err != nil {
		w.Header().Set("Content-Type", "image/svg+xml")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, svg.GenerateErrorBadge("Could not get language"))
		return
	}

	svgData := svg.TemplateData{
		RepoName: repo.Name,
		Language: language,
		Color:    color,
	}

	badge, err := svg.GenerateBadge(svgData)
	if err != nil {
		w.Header().Set("Content-Type", "image/svg+xml")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, svg.GenerateErrorBadge("Failed to generate SVG"))
		return
	}

	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Set("Cache-Control", "s-maxage=3600, stale-while-revalidate")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, badge)
}

func isValidHexColor(s string) bool {
	if s[0] == '#' {
		s = s[1:]
	}
	match, _ := regexp.MatchString(`^[0-9a-fA-F]{3,6}$`, s)
	return match
}
