package svg

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"strings"
)

const svgTemplate = `<svg xmlns="http://www.w3.org/2000/svg" width="{{ .TotalWidth }}" height="20">
  <linearGradient id="b" x2="0" y2="100%">
    <stop offset="0" stop-color="#bbb" stop-opacity=".1"/>
    <stop offset="1" stop-opacity=".1"/>
  </linearGradient>
  <clipPath id="a">
    <rect width="{{ .TotalWidth }}" height="20" rx="3" fill="#fff"/>
  </clipPath>
  <g clip-path="url(#a)">
    <path fill="{{ .LabelColor }}" d="M0 0h{{ .LabelWidth }}v20H0z"/>
    <path fill="{{ .MessageColor }}" d="M{{ .LabelWidth }} 0h{{ .MessageWidth }}v20H{{ .LabelWidth }}z"/>
    <path fill="url(#b)" d="M0 0h{{ .TotalWidth }}v20H0z"/>
  </g>
  <g fill="#fff" text-anchor="middle" font-family="Verdana,Geneva,DejaVu Sans,sans-serif" font-size="110">
    <text x="{{ .LabelX }}" y="150" fill="#010101" fill-opacity=".3" transform="scale(.1)"><tspan font-weight="bold">{{ .Username }}</tspan><tspan>/{{ .RepoName }}</tspan></text>
    <text x="{{ .LabelX }}" y="140" transform="scale(.1)"><tspan font-weight="bold">{{ .Username }}</tspan><tspan>/{{ .RepoName }}</tspan></text>
    <text x="{{ .MessageX }}" y="150" fill="#010101" fill-opacity=".3" transform="scale(.1)">{{ .MessageText }}</text>
    <text x="{{ .MessageX }}" y="140" transform="scale(.1)">{{ .MessageText }}</text>
  </g>
</svg>
`

type TemplateData struct {
	Username     string
	RepoName     string
	MessageText  string
	LabelColor   string
	MessageColor string
	TotalWidth   int
	LabelWidth   int
	MessageWidth int
	LabelX       int
	MessageX     int
}

func calculateTextWidth(text string) int {
	return len(text) * 7
}

func GenerateBadge(username, repoName, message, color string) (string, error) {
	fullLabel := fmt.Sprintf("%s/%s", username, repoName)
	labelWidth := calculateTextWidth(fullLabel) + 10
	messageWidth := calculateTextWidth(message) + 10

	data := TemplateData{
		Username:     username,
		RepoName:     repoName,
		MessageText:  message,
		LabelColor:   "#555",
		MessageColor: color,
		TotalWidth:   labelWidth + messageWidth,
		LabelWidth:   labelWidth,
		MessageWidth: messageWidth,
		LabelX:       labelWidth * 5,
		MessageX:     (labelWidth + messageWidth/2) * 10,
	}

	tmpl, err := template.New("svg").Parse(svgTemplate)
	if err != nil {
		log.Printf("Error parsing SVG template: %v", err)
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		log.Printf("Error executing SVG template: %v", err)
		return "", err
	}

	return buf.String(), nil
}

func GenerateErrorBadge(message string) string {
	cleanedMessage := strings.ReplaceAll(message, "-", "--")
	svg, err := GenerateBadge("gitpulse", "error", cleanedMessage, "#e05d44")
	if err != nil {
		return "<svg>Error</svg>"
	}
	return svg
}
