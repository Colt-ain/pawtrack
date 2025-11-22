package dto

// CreateDogRequest for creating a dog
type CreateDogRequest struct {
	Name      string `json:"name" binding:"required,min=2,max=255" example:"Rex"`
	Breed     string `json:"breed" example:"German Shepherd"`
	BirthDate string `json:"birth_date" example:"2020-01-01T00:00:00Z"`
}

// UpdateDogRequest for updating a dog
type UpdateDogRequest struct {
	Name      string `json:"name" binding:"required,min=2,max=255" example:"Rex"`
	Breed     string `json:"breed" example:"German Shepherd"`
	BirthDate string `json:"birth_date" example:"2020-01-01T00:00:00Z"`
}
