package services

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"library/constant"
	"library/entities"

	"github.com/PuerkitoBio/goquery"
)

type BookService struct {
}

func NewBookService() *BookService {
	return &BookService{}
}

func (svc *BookService) Get(ctx context.Context, isbn string) (*entities.Book, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", fmt.Sprintf("https://isbndb.com/book/%s", isbn), nil)
	req.Header.Set("cookie", "_ga=GA1.2.885462694.1626779278; SESSab6de86aea7caa3f48ba6097cf7cdcf6=EEgG0nrbk7rMaChfagD5rU6GRDSUF4ugoT5iePIMMkk; __stripe_mid=6fbb6b27-b7fc-4fdb-b2c1-7bb5781d032a841978; _gid=GA1.2.1646186259.1626935305; AWSALB=0gdmlLUlv6jOXKTEbbfAx2OQWsho065Xg+dbDxFh2nHgWaZ0bazyJ2+swZKgYOK4/QTRaBM17ITAXLVxWCG6h6JdNVuKIWPxN1tZXo7wdTqixu3akEgRQukgj6CQ; AWSALBCORS=0gdmlLUlv6jOXKTEbbfAx2OQWsho065Xg+dbDxFh2nHgWaZ0bazyJ2+swZKgYOK4/QTRaBM17ITAXLVxWCG6h6JdNVuKIWPxN1tZXo7wdTqixu3akEgRQukgj6CQ")
	res, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	book := &entities.Book{}
	doc.Find("body div table tr").Each(func(i int, s *goquery.Selection) {
		if s.Find("th").Text() == "Full Title" {
			book.Title = s.Find("td").Text()
		}
		if s.Find("th").Text() == "ISBN13" {
			book.ISBN = s.Find("td").Text()
		}
		if s.Find("th").Text() == "Publisher" {
			book.Publisher = s.Find("td").Text()
		}
		if s.Find("th").Text() == "Authors" {
			book.Author = strings.TrimSpace(s.Find("td").Text())
		}
	})

	doc.Find("body div .container .col-md-3 object").Each(func(i int, s *goquery.Selection) {
		if img, ok := s.Attr("data"); ok {
			book.ImageURL = img
		}
	})

	if book.Title == "" {
		return nil, constant.ErrBookNotFound
	}

	return book, nil
}
