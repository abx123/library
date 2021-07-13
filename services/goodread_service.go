package services

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"

	"go.uber.org/zap"

	"library/constant"
	"library/entities"
)

type grResponse struct {
	Search struct {
		Results struct {
			Work struct {
				PublicationYear int64   `xml:"original_publication_year"`
				AverageRating   float64 `xml:"average_rating"`
				Book            struct {
					ID     int64  `xml:"id"`
					Title  string `xml:"title"`
					Author struct {
						ID   int64  `xml:"id"`
						Name string `xml:"name"`
					} `xml:"author"`
					ImageURL      string `xml:"image_url"`
					SmallImageURL string `xml:"small_image_url"`
				} `xml:"best_book"`
			} `xml:"work"`
		} `xml:"results"`
	} `xml:"search"`
}

type GoodreadService struct{}

func NewGoodreadService() *GoodreadService {
	return &GoodreadService{}
}

func (svc *GoodreadService) GetFromGoodread(ctx context.Context, isbn string) (*entities.Book, error) {

	url := fmt.Sprintf("https://www.goodreads.com/search/index.xml?q=%s&key=6qVbqOjnzhHws97M5gYYA", isbn)
	resp, err := http.Get(url)
	if err != nil {
		zap.L().Error(constant.ErrGoodreadError.Error(), zap.Error(err))
		return &entities.Book{}, constant.ErrGoodreadError
	}
	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		zap.L().Error(constant.ErrGoodreadError.Error(), zap.Error(err))
		return &entities.Book{}, constant.ErrGoodreadError
	}
	xmlReader := bytes.NewReader([]byte(string(body)))
	xmlBook := new(grResponse)
	if err := xml.NewDecoder(xmlReader).Decode(xmlBook); err != nil {
		zap.L().Error(constant.ErrGoodreadError.Error(), zap.Error(err))
		return &entities.Book{}, constant.ErrGoodreadError
	}

	if xmlBook.Search.Results.Work.Book.Title == "" {
		zap.L().Error(constant.ErrBookNotFound.Error(), zap.Error(constant.ErrBookNotFound))
		return &entities.Book{}, constant.ErrBookNotFound
	}

	book := &entities.Book{
		ISBN:            isbn,
		Title:           xmlBook.Search.Results.Work.Book.Title,
		Author:          xmlBook.Search.Results.Work.Book.Author.Name,
		ImageURL:        xmlBook.Search.Results.Work.Book.ImageURL,
		SmallImageURL:   xmlBook.Search.Results.Work.Book.SmallImageURL,
		PublicationYear: xmlBook.Search.Results.Work.PublicationYear,
		AverageRating:   xmlBook.Search.Results.Work.AverageRating,
	}

	return book, nil
}
