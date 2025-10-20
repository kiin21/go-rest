package starter

type DeleteStarterRequest struct {
	Domain string `uri:"domain" validate:"required,min=2,max=100"`
}
