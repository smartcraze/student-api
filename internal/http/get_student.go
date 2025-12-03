package httphandler

import (
	"net/http"
	"strconv"

	"github.com/smartcraze/student-api/internal/storage"
	"github.com/smartcraze/student-api/utils/response"
)

type StudentResponse struct {
	ID             int64  `json:"id"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	RegistrationNo int    `json:"reg_no"`
	PhoneNumber    int64  `json:"phone_number"`
	Email          string `json:"email"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

func GetStudentHandler(store storage.Storage) http.HandlerFunc {
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

		// Get student from database
		student, err := store.GetStudentByID(r.Context(), id)
		if err != nil {
			response.Writejson(w, http.StatusNotFound, response.Response{
				Status: response.StatusError,
				Error:  err.Error(),
			})
			return
		}

		// Return student without password
		resp := StudentResponse{
			ID:             student.ID,
			FirstName:      student.FirstName,
			LastName:       student.LastName,
			RegistrationNo: student.RegistrationNo,
			PhoneNumber:    student.PhoneNumber,
			Email:          student.Email,
			CreatedAt:      student.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:      student.UpdatedAt.Format("2006-01-02 15:04:05"),
		}

		response.Writejson(w, http.StatusOK, resp)
	}
}
