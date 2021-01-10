package main

import (
	"bytes"
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/saintfish/chardet"
	"golang.org/x/net/html/charset"
	"io/ioutil"
	"os"
	"regexp"
)

var (
	removeNumTitle = regexp.MustCompile("(?m)^\\d+$")
	detector = chardet.NewTextDetector()
)

func convertToUTF8(strBytes []byte, origEncoding string) ([]byte, error) {
	byteReader := bytes.NewReader(strBytes)
	reader, err := charset.NewReaderLabel(origEncoding, byteReader)
	if err != nil {
		return nil, err
	}
	strBytes, err = ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return strBytes, nil
}

func WriteFile(filename string, data string) error {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	_, err = f.Write([]byte(data))
	if err1 := f.Close(); err == nil {
		err = err1
	}

	return err
}

func getTextEncoding(dataBytes []byte) (string, error) {
	result, err := detector.DetectBest(dataBytes)
	if err != nil {
		return "", err
	}
	return result.Charset, nil
}

func main() {
	c := colly.NewCollector()

	c.OnHTML("html", func(e *colly.HTMLElement) {
		text, err  := processText(e.ChildText("p") + "\n")
		if err != nil {
			panic(err)
		}

		err = WriteFile("data/nos.txt", text)
		if err != nil {
			panic(err)
		}
		e.ForEach("a", func(_ int, element *colly.HTMLElement) {
			if element.Text == "Next" {
				e.Request.Visit(element.Attr("href"))
			}
		})
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit("https://www.sacred-texts.com/nos/preface.htm")
}

func processText(text string) (string, error) {
	dataBytes := []byte(text)
	textEncoding, err := getTextEncoding(dataBytes)
	if err != nil {
		return "", err
	}
	if textEncoding != "UTF-8" {
		dataBytes, err = convertToUTF8(dataBytes, textEncoding)
		if err != nil {
			return "", err
		}
	}
	text = removeNumTitle.ReplaceAllString(string(dataBytes), "")

	return text, nil
}
