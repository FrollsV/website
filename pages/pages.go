package pages

import (
	"encoding/json"
)

type Paragraph struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Image   string `json:"image"`
}

type Page struct {
	Title      string      `json:"title"`
	Paragraphs []Paragraph `json:"paragraphs"`
}

func LoadArticle(content string) (Page, error) {
	page := Page{}
	err := json.Unmarshal([]byte(content), &page)
	return page, err

}
