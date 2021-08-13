package services

import (
	"context"
	"fmt"
	"library/entities"
	"library/repo/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpsert(t *testing.T) {
	type testCase struct {
		name   string
		desc   string
		expRes *entities.Book
		expErr error
	}

	testCases := []testCase{
		{
			name: "Happy Case",
			desc: "all ok",
			expRes: &entities.Book{
				BookID:          0,
				ISBN:            "isbn",
				Title:           "title",
				Authors:         "authors",
				ImageURL:        "imageURL",
				SmallImageURL:   "smallImageURL",
				PublicationYear: 2021,
				Publisher:       "publisher",
				UserID:          "userId",
				Status:          1,
				Description:     "description",
				PageCount:       999,
				Categories:      "categories",
				Language:        "language",
				Source:          "source"},
		},
		{
			name:   "Sad Case",
			desc:   "repo returns error",
			expErr: fmt.Errorf("mock error"),
		},
	}

	for _, v := range testCases {
		repo := mocks.IdbRepo{}
		dbSvc := NewDbService(&repo)
		repo.On("Upsert", context.Background(), &entities.Book{BookID: 0, ISBN: "isbn", Title: "title", Authors: "authors", ImageURL: "imageURL", SmallImageURL: "smallImageURL", PublicationYear: 2021, Publisher: "publisher", UserID: "userId", Status: 1, Description: "description", PageCount: 999, Categories: "categories", Language: "language", Source: "source"}).Return(v.expRes, v.expErr)
		actRes, actErr := dbSvc.Upsert(context.Background(), "isbn", "title", "authors", "imageURL", "smallImageURL", "publisher", "userId", "description", "categories", "language", "source", 2021, 1, 999)
		assert.Equal(t, v.expRes, actRes)
		assert.Equal(t, v.expErr, actErr)
	}
}

func TestDBGet(t *testing.T) {
	type testCase struct {
		name   string
		desc   string
		expRes *entities.Book
		expErr error
	}

	testCases := []testCase{
		{
			name: "Happy Case",
			desc: "all ok",
			expRes: &entities.Book{
				BookID:          0,
				ISBN:            "isbn",
				Title:           "title",
				Authors:         "authors",
				ImageURL:        "imageURL",
				SmallImageURL:   "smallImageURL",
				PublicationYear: 2021,
				Publisher:       "publisher",
				UserID:          "userId",
				Status:          1,
				Description:     "description",
				PageCount:       999,
				Categories:      "categories",
				Language:        "language",
				Source:          "source"},
		},
		{
			name:   "Sad Case",
			desc:   "repo returns error",
			expErr: fmt.Errorf("mock error"),
		},
	}

	for _, v := range testCases {
		repo := mocks.IdbRepo{}
		dbSvc := NewDbService(&repo)
		repo.On("Get", context.Background(), &entities.Book{BookID: 0, ISBN: "isbn", Title: "", Authors: "", ImageURL: "", SmallImageURL: "", PublicationYear: 0, Publisher: "", UserID: "userid", Status: 0, Description: "", PageCount: 0, Categories: "", Language: "", Source: ""}).Return(v.expRes, v.expErr)
		actRes, actErr := dbSvc.Get(context.Background(), "isbn", "userid")
		assert.Equal(t, v.expRes, actRes)
		assert.Equal(t, v.expErr, actErr)
	}
}

func TestList(t *testing.T) {
	type testCase struct {
		name   string
		desc   string
		expRes []*entities.Book
		expErr error
	}
	testCases := []testCase{
		{
			name: "Happy Case",
			desc: "all ok",
			expRes: []*entities.Book{
				{
					BookID:          0,
					ISBN:            "isbn",
					Title:           "title",
					Authors:         "authors",
					ImageURL:        "imageURL",
					SmallImageURL:   "smallImageURL",
					PublicationYear: 2021,
					Publisher:       "publisher",
					UserID:          "userId",
					Status:          1,
					Description:     "description",
					PageCount:       999,
					Categories:      "categories",
					Language:        "language",
					Source:          "source",
				},
			},
		},
		{
			name:   "Sad Case",
			desc:   "repo returns error",
			expErr: fmt.Errorf("mock error"),
		},
	}
	for _, v := range testCases {
		repo := mocks.IdbRepo{}
		dbSvc := NewDbService(&repo)
		repo.On("List", context.Background(), int64(10), int64(0), "userid").Return(v.expRes, v.expErr)
		actRes, actErr := dbSvc.List(context.Background(), 10, 0, "userid")
		assert.Equal(t, v.expRes, actRes)
		assert.Equal(t, v.expErr, actErr)
	}
}
