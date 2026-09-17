package link

type CreateUserDto struct {
	URL            string `json:"long_url" binding:"required,url"`
	ExpirationDate string `json:"expiration_date"`
}
