package starter

type FindStarterRequest struct {
	Domain string `uri:"domain" binding:"required,min=2,max=100"`
}
