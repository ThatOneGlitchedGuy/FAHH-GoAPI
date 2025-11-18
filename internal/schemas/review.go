package schemas

import (
	"time"
)

type ReviewCreate struct {
	Rating  int    `json:"rating" validate:"required,min=1,max=5"`
	Comment string `json:"comment" validate:"max=4000"`
}

type ReviewOut struct {
	ID        uint      `json:"id"`
	ProductID uint      `json:"product_id"`
	UserID    uint      `json:"user_id"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
	User      *UserReviewerOut `json:"user,omitempty"`
}

type UserReviewerOut struct {
	ID       uint   `json:"id"`
	FullName string `json:"full_name"`
}
