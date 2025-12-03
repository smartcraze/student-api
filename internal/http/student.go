package http

import (
	"encoding/json"
	"net/http"
	"time"
)

type CreateStudent struct {
	Id             int64     `json:"id,omitempty"`
	FirstName      string    `json:"first_name" validate:"required"`
	LastName       string    `json:"last_name"  validate:"required"`
	RegistrationNo int       `json:"reg_no" validate:"required"`
	PhoneNumber    int64     `json:"phone_number" validate:"required"`
	Email          string    `json:"email" validate:"required,email"`
	Password       string    `json:"password,omitempty" validate:"required"`
	CreatedAt      time.Time `json:"created_at,omitempty"`
}

func CreateStudentHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateStudent

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "Invalide json Fields", http.StatusBadRequest)
		return
	}
	validate := 

}
