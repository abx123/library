package presenter

type Book struct {
	ISBN            string  `json:"isbn,omitempty"`
	Title           string  `json:"title,omitempty"`
	Author          string  `json:"author,omitempty"`
	ImageURL        string  `json:"imageURL,omitempty"`
	SmallImageURL   string  `json:"smallImageURL,omitempty"`
	PublicationYear int64   `json:"publicationYear,omitempty"`
	AverageRating   float64 `json:"averageRating,omitempty"`
	UserID          string  `json:"userId,omitempty"`
	Status          int64   `json:"status,omitempty"`
}
