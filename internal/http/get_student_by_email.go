package httphandler

import (
	"net/http"

	"github.com/smartcraze/student-api/internal/storage"
	"github.com/smartcraze/student-api/utils/response"
)

func GetStudentByEmailHandler(store storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get email from query parameter
		email := r.URL.Query().Get("email")
		if email == "" {
			response.Writejson(w, http.StatusBadRequest, response.Response{
				Status: response.StatusError,
				Error:  "email parameter is required",
			})
			return
		}

		// Get student from database by email
		student, err := store.GetStudentByEmail(r.Context(), email)
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
