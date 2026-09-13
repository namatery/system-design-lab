package link

type CreateUserDto struct {
	URL        string `json:"URL" binding:"required,url"`
	Expiration string `json:"Expiration" binding:"required,datetime=2006-01-02"`
}
