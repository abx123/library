package entities

// Book represents a Book object
type Book struct {
	BookID          int64   `db:"id"`
	ISBN            string  `db:"isbn"`
	Title           string  `db:"title"`
	Author          string  `db:"author"`
	ImageURL        string  `db:"imageUrl"`
	SmallImageURL   string  `db:"smallImageUrl"`
	PublicationYear int64   `db:"publicationYear"`
	AverageRating   float64 `db:"averageRating"`
	UserID          string  `db:"userId"`
	Status          int64   `db:"status"`
}
