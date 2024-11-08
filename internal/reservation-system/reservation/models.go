package reservation

import "time"

type reservation struct {
	ID             *int       `db:"id"`
	ReservationUid *string    `db:"reservation_uid"`
	UserName       *string    `db:"username"`
	BookUid        *string    `db:"book_uid"`
	LibraryUid     *string    `db:"library_uid"`
	Status         *string    `db:"status"`
	StartDate      *time.Time `db:"start_date"`
	TillDate       *time.Time `db:"till_date"`
}
