package entity

import "time"

type Transaction struct {
	ID         int `gorm:"primary_key:auto_increment"`
	CampaignID int
	UserID     int
	Amount     int
	Status     string
	Code       string
	User       User
	PaymentURL string
	Campaign   Campaign
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
