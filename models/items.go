package models

type Item struct {
	ID          int    `json:"id" db:"id"`            
	Name        string `json:"name" db:"name"`         
	Cost        int    `json:"cost" db:"cost"`        
	ImageURL    string `json:"image_url" db:"image_url"` 
	Description string `json:"description" db:"description"` 
	VendorID    int    `json:"vendor_id" db:"vendor_id"`
	IsAvailable    bool    `json:"is_available " db:"is_available "` 
}
