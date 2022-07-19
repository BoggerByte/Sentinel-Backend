package forms

type GenerateURLForm struct{}

type Oauth2RedirectForm struct {
	Code  string `form:"code" binding:"required"`
	State string `form:"state" binding:"required"`
}
