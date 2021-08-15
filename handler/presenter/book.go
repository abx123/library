package presenter

// Book defines book object
type Book struct {
	ISBN            string `json:"isbn,omitempty"`
	Title           string `json:"title,omitempty"`
	Author          string `json:"author,omitempty"`
	ImageURL        string `json:"imageURL,omitempty"`
	SmallImageURL   string `json:"smallImageURL,omitempty"`
	Publisher       string `json:"publisher,omitempty"`
	Description     string `json:"description,omitempty"`
	PageCount       int64  `json:"pageCount,omitempty"`
	Categories      string `json:"categories,omitempty"`
	Language        string `json:"language,omitempty"`
	PublicationYear int64  `json:"publicationYear,omitempty"`
	UserID          string `json:"userId,omitempty"`
	Status          int64  `json:"status,omitempty"`
	Source          string `json:"source,omitempty"`
}
