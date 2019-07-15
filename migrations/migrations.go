package migrations

type Link struct {
	Hash	string	`gorm:"primary_key"`
	URL		string
	Hits	int 	`gorm:"default:0"`
}
