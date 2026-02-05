package users_http

type registerDTO struct {
	FirstName      string `json:"first_name" validate:"required,min=2,max=30"`
	LastName       string `json:"last_name" validate:"required,min=2,max=30"`
	Email          string `json:"email" validate:"required,email"`
	BirthDate      string `json:"birth_date" validate:"required,datetime=2006-01-02"`
	PhoneNumber    string `json:"phone_number" validate:"omitempty,e164"`
	Gender         string `json:"gender" validate:"required,oneof=male female"`
	Password       string `json:"password" validate:"required,min=4,max=72"` // bcrypt лимит 72 байта
	RepeatPassword string `json:"repeat_password" validate:"required,min=4,max=72"`
}

type loginDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=4,max=72"`
}
