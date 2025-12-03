package httphandler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/smartcraze/student-api/internal/storage"
	"github.com/smartcraze/student-api/utils/response"
	"golang.org/x/crypto/bcrypt"
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

func CreateStudentHandler(store storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateStudent

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			response.Writejson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		// Validate request
		validate := validator.New()
		if err := validate.Struct(req); err != nil {
			response.Writejson(w, http.StatusBadRequest, response.ValidationError(err.(validator.ValidationErrors)))
			return
		}

		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			response.Writejson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		// Create student in database
		student := &storage.Student{
			FirstName:      req.FirstName,
			LastName:       req.LastName,
			RegistrationNo: req.RegistrationNo,
			PhoneNumber:    req.PhoneNumber,
			Email:          req.Email,
			Password:       string(hashedPassword),
		}

		err = store.CreateStudent(r.Context(), student)
		if err != nil {
			response.Writejson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		// Return created student (without password)
		req.Id = student.ID
		req.CreatedAt = student.CreatedAt
		req.Password = "" // Don't send password back

		response.Writejson(w, http.StatusCreated, req)
	}
}
