package models

type Vendor struct {
	ID          int    `json:"id" db:"id"`                  
	Name        string `json:"name" db:"name"`               
	Description string `json:"description" db:"description"` 
	Password string `json:"password" db:"password"` 
	ImageURL    string `json:"image_url" db:"image_url"`    
}
