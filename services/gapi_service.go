package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"go.uber.org/zap"

	"library/constant"
	"library/entities"
)

type gapiResponse struct {
	TotalItems int64 `json:"totalItems,omitempty"`
	Items      []struct {
		VolumeInfo struct {
			Title           string   `json:"title,omitempty"`
			Authors         []string `json:"authors,omitempty"`
			PublicationYear string   `json:"publishedDate,omitempty"`
			AverageRating   float64  `json:"averageRating,omitempty"`
			Identifier      []struct {
				Type       string `json:"type,omitempty"`
				Identifier string `json:"identifier,omitempty"`
			} `json:"industryIdentifiers,omitempty"`
			Image struct {
				ImageURL      string `json:"thumbnail,omitempty"`
				SmallImageURL string `json:"smallThumbnail,omitempty"`
			} `json:"imageLinks,omitempty"`
		} `json:"volumeInfo,omitempty"`
	} `json:"items,omitempty"`
}

type GapiService struct{}

func NewGapiService() *GapiService {
	return &GapiService{}
}

func (svc *GapiService) GetFromGoogleAPI(ctx context.Context, isbn string) (*entities.Book, error) {
	resp, err := http.Get(fmt.Sprintf("https://www.googleapis.com/books/v1/volumes?q=isbn:%s", isbn))
	if err != nil {
		zap.L().Error(constant.ErrGapiError.Error(), zap.Error(err))
		return &entities.Book{}, constant.ErrGapiError
	}
	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		zap.L().Error(constant.ErrGapiError.Error(), zap.Error(err))
		return &entities.Book{}, constant.ErrGapiError
	}

	val := &gapiResponse{}
	r := bytes.NewReader([]byte(string(body)))
	decoder := json.NewDecoder(r)
	err = decoder.Decode(val)

	if err != nil {
		zap.L().Error(constant.ErrGapiError.Error(), zap.Error(err))
		return &entities.Book{}, constant.ErrGapiError
	}

	if val.TotalItems < 1 {
		zap.L().Error(constant.ErrBookNotFound.Error(), zap.Error(constant.ErrBookNotFound))
		return &entities.Book{}, constant.ErrBookNotFound
	}

	year, err := strconv.ParseInt(val.Items[0].VolumeInfo.PublicationYear, 10, 64)
	if err != nil {
		zap.L().Error(constant.ErrGapiError.Error(), zap.Error(err))
		return &entities.Book{}, constant.ErrGapiError
	}

	book := &entities.Book{
		ISBN:            isbn,
		Title:           val.Items[0].VolumeInfo.Title,
		Author:          val.Items[0].VolumeInfo.Authors[0],
		ImageURL:        val.Items[0].VolumeInfo.Image.ImageURL,
		SmallImageURL:   val.Items[0].VolumeInfo.Image.SmallImageURL,
		PublicationYear: year,
		AverageRating:   val.Items[0].VolumeInfo.AverageRating,
	}

	return book, nil
}
