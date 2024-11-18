package models


type User struct {
	RollNo        string  `json:"roll_no" db:"roll_no"`
	Name          string  `json:"name" db:"name"`
	Email         string  `json:"email" db:"email"`
	Pin      string  `json:"pin" db:"pin"`
	IsApproved    bool    `json:"is_approved" db:"is_approved"`
	IsVerified    bool    `json:"is_verified" db:"is_verified"`
	WalletBalance int `json:"wallet_balance" db:"wallet_balance"`
}

type UserResp struct {
	RollNo        string  `json:"roll_no" db:"roll_no"`
	Name          string  `json:"name" db:"name"`
	Email         string  `json:"email" db:"email"`
	IsApproved    bool    `json:"is_approved" db:"is_approved"`
	IsVerified    bool    `json:"is_verified" db:"is_verified"`
	WalletBalance int `json:"wallet_balance" db:"wallet_balance"`
}





