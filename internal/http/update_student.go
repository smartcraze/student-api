package httphandler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/smartcraze/student-api/internal/storage"
	"github.com/smartcraze/student-api/utils/response"
)

type UpdateStudentRequest struct {
	FirstName      string `json:"first_name" validate:"required"`
	LastName       string `json:"last_name" validate:"required"`
	RegistrationNo int    `json:"reg_no" validate:"required"`
	PhoneNumber    int64  `json:"phone_number" validate:"required"`
	Email          string `json:"email" validate:"required,email"`
}

func UpdateStudentHandler(store storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get student ID from URL path parameter
		idStr := r.PathValue("id")
		if idStr == "" {
			response.Writejson(w, http.StatusBadRequest, response.Response{
				Status: response.StatusError,
				Error:  "student ID is required",
			})
			return
		}

		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			response.Writejson(w, http.StatusBadRequest, response.Response{
				Status: response.StatusError,
				Error:  "invalid student ID",
			})
			return
		}

		// Parse request body
		var req UpdateStudentRequest
		err = json.NewDecoder(r.Body).Decode(&req)
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

		// Check if student exists
		existingStudent, err := store.GetStudentByID(r.Context(), id)
		if err != nil {
			response.Writejson(w, http.StatusNotFound, response.Response{
				Status: response.StatusError,
				Error:  "student not found",
			})
			return
		}

		// Update student data
		existingStudent.FirstName = req.FirstName
		existingStudent.LastName = req.LastName
		existingStudent.RegistrationNo = req.RegistrationNo
		existingStudent.PhoneNumber = req.PhoneNumber
		existingStudent.Email = req.Email

		// Save updated student
		err = store.UpdateStudent(r.Context(), existingStudent)
		if err != nil {
			response.Writejson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		// Return updated student
		resp := StudentResponse{
			ID:             existingStudent.ID,
			FirstName:      existingStudent.FirstName,
			LastName:       existingStudent.LastName,
			RegistrationNo: existingStudent.RegistrationNo,
			PhoneNumber:    existingStudent.PhoneNumber,
			Email:          existingStudent.Email,
			CreatedAt:      existingStudent.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:      existingStudent.UpdatedAt.Format("2006-01-02 15:04:05"),
		}

		response.Writejson(w, http.StatusOK, resp)
	}
}
