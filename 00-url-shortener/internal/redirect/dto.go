package redirect

type RedirectDto struct {
	Alias string `uri:"alias" binding:"required,len=6,hexadecimal"`
}
