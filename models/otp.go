package models

import "time"

type OTP struct {
	RollNo               string    `json:"roll_no" db:"roll_no"`
	OTPLastGenerated     int       `json:"otp_last_generated" db:"otp_last_generated"` 
	OTPLastGeneratedTime time.Time `json:"otp_last_generated_time" db:"otp_last_generated_time"` 
}
