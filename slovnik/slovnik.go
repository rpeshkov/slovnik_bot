package slovnik

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"strings"

	"golang.org/x/net/html"
)

// Word defines a structure with the word itself and possible translations of that word
type Word struct {
	Word         string
	Translations []string
	WordType     string
	Synonyms     []string
	Antonyms     []string
}

// Method for transforming Word struct to string
func (w Word) String() string {
	out := fmt.Sprintf("*%s*\n", w.WordType)
	out += fmt.Sprintln(strings.Join(w.Translations, ", "))
	if len(w.Synonyms) > 0 {
		out += fmt.Sprintln("\n*Synonyms:*")
		out += fmt.Sprintln(strings.Join(w.Synonyms, ", "))
	}
	if len(w.Antonyms) > 0 {
		out += fmt.Sprintln("\n*Antonyms:*")
		out += fmt.Sprintln(strings.Join(w.Antonyms, ", "))
	}
	return out
}

// GetTranslations from slovnik.seznam.cz for specified word
func GetTranslations(word string, langcode string) (Word, error) {
	urls := map[string]string{
		"cs": "https://slovnik.seznam.cz/cz-ru/",
		"ru": "https://slovnik.seznam.cz/ru/",
	}

	query, _ := url.Parse(urls[langcode])

	p := url.Values{}
	p.Add("q", word)
	p.Add("shortView", "0")

	query.RawQuery = p.Encode()

	resp, err := http.Get(query.String())
	if err != nil {
		return Word{}, err
	}
	return parsePage(resp.Body), nil
}

// DetectLanguage used to find out which language is used for the input string
func DetectLanguage(input string) string {
	const ru = "абвгдеёжзийклмнопрстуфхцчшщьыъэюя"
	for _, ch := range input {
		if strings.Contains(ru, strings.ToLower(string(ch))) {
			return "ru"
		}
	}
	return "cs"
}

var f func(w *Word, data string)

func parsePage(pageBody io.Reader) Word {
	z := html.NewTokenizer(pageBody)

	inTranslations := false
	foundSynonymsHeader := false
	inSynonymsBlock := false
	foundAntonymsHeader := false
	inAntonymsBlock := false
	prevTag := ""

	w := Word{}
	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			return w

		case tt == html.StartTagToken:
			t := z.Token()

			if t.Data == "h3" {
				lang := getAttr(t.Attr, "lang")
				if lang == "cs" || lang == "ru" {
					f = addWord
				}
			}

			if t.Data == "div" {
				inTranslations = getAttr(t.Attr, "id") == "fastMeanings"
				inSynonymsBlock = getAttr(t.Attr, "class") == "other-meaning" && foundSynonymsHeader
				inAntonymsBlock = getAttr(t.Attr, "class") == "other-meaning" && foundAntonymsHeader
			}

			if t.Data == "a" && inTranslations {
				f = addTranslation
			}

			if t.Data == "a" && prevTag == "a" && inTranslations {
				f = updateLastTranslation
			}

			if t.Data == "span" && inTranslations {
				if getAttr(t.Attr, "class") != "comma" {
					f = addTranslation
				}
			}

			if t.Data == "a" && inSynonymsBlock {
				f = addSynonym
			}

			if t.Data == "a" && inAntonymsBlock {
				f = addAntonym
			}

			if t.Data == "span" && getAttr(t.Attr, "class") == "morf" {
				f = addWordType
			}

			prevTag = t.Data

			break

		case tt == html.SelfClosingTagToken:
			t := z.Token()
			prevTag = t.Data
			break

		case tt == html.EndTagToken:
			t := z.Token()
			if t.Data == "div" {
				inTranslations = false
			}

			if t.Data == "div" && inSynonymsBlock {
				inSynonymsBlock = false
				foundSynonymsHeader = false
			}

			if t.Data == "div" && inAntonymsBlock {
				inAntonymsBlock = false
				foundAntonymsHeader = false
			}

			break

		case tt == html.TextToken:
			t := z.Token()
			if f != nil {
				f(&w, t.Data)
				f = nil
			}

			if t.Data == "Synonyma" {
				foundSynonymsHeader = true
			}

			if t.Data == "Antonyma" {
				foundAntonymsHeader = true
			}

			break
		}
	}
}

func addWord(w *Word, data string) {
	w.Word = data
}
func addWordType(w *Word, data string) {
	w.WordType = data
}

func addTranslation(w *Word, data string) {
	w.Translations = append(w.Translations, data)
}

func updateLastTranslation(w *Word, data string) {
	if len(w.Translations) > 0 {
		lastTranslation := w.Translations[len(w.Translations)-1]
		lastTranslation = lastTranslation + " " + data
		w.Translations[len(w.Translations)-1] = lastTranslation
	}
}

func addSynonym(w *Word, data string) {
	w.Synonyms = append(w.Synonyms, data)
}

func addAntonym(w *Word, data string) {
	w.Antonyms = append(w.Antonyms, data)
}

func getAttr(attrs []html.Attribute, name string) string {
	for _, a := range attrs {
		if a.Key == name {
			return a.Val
		}
	}

	return ""
}
