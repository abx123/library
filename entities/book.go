package entities

// Book represents a Book object
type Book struct {
	BookID          int64  `db:"id"`
	ISBN            string `db:"isbn"`
	Title           string `db:"title"`
	Authors         string `db:"authors"`
	ImageURL        string `db:"imageUrl"`
	SmallImageURL   string `db:"smallImageUrl"`
	PublicationYear int64  `db:"publicationYear"`
	Publisher       string `db:"publisher"`
	UserID          string `db:"userId"`
	Status          int64  `db:"status"`

	Description string `db:"description"`
	PageCount   int64  `db:"pageCount"`
	Categories  string `db:"categories"`
	Language    string `db:"language"`
	Source      string `db:"source"`
}
