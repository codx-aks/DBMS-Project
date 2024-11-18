package models

type Vendor struct {
	ID          int    `json:"id" db:"id"`                  
	Name        string `json:"name" db:"name"`               
	Description string `json:"description" db:"description"` 
	Password string `json:"password" db:"password"` 
	ImageURL    string `json:"image_url" db:"image_url"`
	IsActive    bool    `json:"is_active" db:"is_active"`    
}

type RespVendor struct {
	ID          int    `json:"id" db:"id"`                  
	Name        string `json:"name" db:"name"`               
	Description string `json:"description" db:"description"` 
	ImageURL    string `json:"image_url" db:"image_url"`    
}
