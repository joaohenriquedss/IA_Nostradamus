package main

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"os"
)
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

func main() {
	c := colly.NewCollector()

	c.OnHTML("html", func(e *colly.HTMLElement) {
		err := WriteFile("data/nos.txt", e.ChildText("p") + "\n")
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
