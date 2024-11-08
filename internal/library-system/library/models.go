package library

type library struct {
	ID         int    `db:"id"`
	LibraryUid string `db:"libraryUid"`
	Name       string `db:"name"`
	Address    string `db:"address"`
	City       string `db:"city"`
}

type book struct {
	ID             int    `db:"id"`
	BookUid        string `db:"bookUid"`
	Name           string `db:"name"`
	Author         string `db:"author"`
	Genre          string `db:"genre"`
	Condition      string `db:"condition"`
	AvailableCount int    `db:"availableCount"`
}
