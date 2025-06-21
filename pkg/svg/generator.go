package svg

import (
	"bytes"
	"html/template"
	"log"
)

const svgTemplate = `
<svg xmlns="http://www.w3.org/2000/svg" width="380" height="60" viewBox="0 0 380 60">
    <style>
        .container {
            font-family: 'Segoe UI', 'Roboto', 'Helvetica Neue', 'Arial', sans-serif;
            font-size: 14px;
            fill: #333;
        }
        .repo-name {
            font-weight: 600;
            font-size: 16px;
        }
        .label {
            font-weight: 400;
            fill: #555;
        }
        .lang-color {
            fill: {{ .Color }};
        }
    </style>
    <rect x="0" y="0" width="380" height="60" fill="transparent"/>
    <g class="container" transform="translate(10, 25)">
        <text class="label">Last commit:</text>
        <text x="85" y="0" class="repo-name">{{ .RepoName }}</text>
    </g>
    <g class="container" transform="translate(10, 48)">
        <text class="label">Main language:</text>
        <text x="105" y="0" class="lang-color">{{ .Language }}</text>
    </g>
</svg>
`

type TemplateData struct {
	RepoName string
	Language string
	Color    string
}

func GenerateBadge(data TemplateData) (string, error) {
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
	errorData := TemplateData{
		RepoName: "Error",
		Language: message,
		Color:    "#de2f2f",
	}
	svg, err := GenerateBadge(errorData)
	if err != nil {
		return "<svg>Error generating badge</svg>"
	}
	return svg
}
