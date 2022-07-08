package forms

type LoginForm struct {
	State string `json:"state" binding:"required"`
}
