package main

import (
	"fmt"
	"io"
	"net/http"

	"strings"

	"golang.org/x/net/html"
)

// Word defines a structure with the word itself and possible translations of that word
type Word struct {
	word         string
	translations []string
}

// Method for transforming Word struct to string
func (w Word) String() string {
	return fmt.Sprintf("*%s*\n%s", w.word, strings.Join(w.translations, ", "))
}

// GetTranslations from slovnik.seznam.cz for specified word
func GetTranslations(word string) (Word, error) {
	url := fmt.Sprintf("https://slovnik.seznam.cz/cz-ru/?q=%s", word)
	resp, err := http.Get(url)
	if err != nil {
		return Word{}, err
	}
	return parsePage(resp.Body), nil
}

func parsePage(pageBody io.Reader) Word {
	z := html.NewTokenizer(pageBody)

	inWord := false
	inTranslations := false
	inTranslationLink := false

	w := Word{}
	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			return w

		case tt == html.StartTagToken:
			t := z.Token()

			if t.Data == "h3" {
				for _, attr := range t.Attr {
					if attr.Key == "lang" && attr.Val == "cs" {
						inWord = true
					}
				}
			}

			if t.Data == "div" {
				for _, attr := range t.Attr {
					if attr.Key == "id" && attr.Val == "fastMeanings" {
						inTranslations = true
					}
				}
			}

			if t.Data == "a" && inTranslations {
				inTranslationLink = true
			}

			break

		case tt == html.EndTagToken:
			t := z.Token()
			if t.Data == "div" {
				inTranslations = false
			}

			if t.Data == "h3" {
				inWord = false
			}
			if t.Data == "a" && inTranslationLink {
				inTranslationLink = false
			}

			break

		case tt == html.TextToken:
			t := z.Token()
			if inWord {
				w.word = t.Data
			}
			if inTranslationLink {
				w.translations = append(w.translations, t.Data)
			}
			break
		}
	}
}
