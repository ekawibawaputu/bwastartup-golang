package transaction

import (
	"bwastartup/user"
	"time"
)

type Transaction struct {
	ID        int
	Campaign  int
	UserID    int
	Amount    int
	Status    string
	Code      string
	User user.User
	CreatedAt time.Time
	UpdatedAt time.Time
}