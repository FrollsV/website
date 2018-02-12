package pages

import (
	"testing"
)

func TestLoadArticle(t *testing.T) {
	sample := `{
    "title": "Hello coffee",
    "paragraphs": [
      {
        "content":"What a nice boat. Don't you want it back?",
        "image": "https://images.moviepilot.com/image/upload/c_fill,h_470,q_auto:good,w_620/w5vmzajmqceitr7vkrnc.jpg"
      },
      {
        "title": "don't accept stuff from strangers",
        "content": "You might fall pray of a psychopatic clown trapped in a storm drain!"
      },
      {
        "image": "data/cat.jpg"
      }
    ]
  }
  `

	p, err := LoadArticle(sample)

	if err != nil {
		t.Error(err)
	}

	if p.Title != "Hello coffee" {
		t.Log("title not as expected")
		t.Fail()
	}

	if len(p.Paragraphs) != 3 {
		t.Log("count of paragraphs is not 3")
		t.Fail()
	}

	if p.Paragraphs[0].Title != "" {
		t.Log("1st paragraph has a title")
		t.Fail()
	}
}
