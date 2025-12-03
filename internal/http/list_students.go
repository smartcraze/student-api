package httphandler

import (
	"net/http"
	"strconv"

	"github.com/smartcraze/student-api/internal/storage"
	"github.com/smartcraze/student-api/utils/response"
)

type ListStudentsResponse struct {
	Students []*StudentResponse `json:"students"`
	Total    int                `json:"total"`
	Limit    int                `json:"limit"`
	Offset   int                `json:"offset"`
}

func ListStudentsHandler(store storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse query parameters for pagination
		limitStr := r.URL.Query().Get("limit")
		offsetStr := r.URL.Query().Get("offset")

		// Default values
		limit := 10
		offset := 0

		// Parse limit
		if limitStr != "" {
			parsedLimit, err := strconv.Atoi(limitStr)
			if err != nil || parsedLimit < 1 {
				response.Writejson(w, http.StatusBadRequest, response.Response{
					Status: response.StatusError,
					Error:  "invalid limit parameter",
				})
				return
			}
			limit = parsedLimit
		}

		// Parse offset
		if offsetStr != "" {
			parsedOffset, err := strconv.Atoi(offsetStr)
			if err != nil || parsedOffset < 0 {
				response.Writejson(w, http.StatusBadRequest, response.Response{
					Status: response.StatusError,
					Error:  "invalid offset parameter",
				})
				return
			}
			offset = parsedOffset
		}

		// Limit max results to prevent excessive queries
		if limit > 100 {
			limit = 100
		}

		// Get students from database
		students, err := store.ListStudents(r.Context(), limit, offset)
		if err != nil {
			response.Writejson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		// Convert to response format (without passwords)
		studentResponses := make([]*StudentResponse, 0, len(students))
		for _, student := range students {
			studentResponses = append(studentResponses, &StudentResponse{
				ID:             student.ID,
				FirstName:      student.FirstName,
				LastName:       student.LastName,
				RegistrationNo: student.RegistrationNo,
				PhoneNumber:    student.PhoneNumber,
				Email:          student.Email,
				CreatedAt:      student.CreatedAt.Format("2006-01-02 15:04:05"),
				UpdatedAt:      student.UpdatedAt.Format("2006-01-02 15:04:05"),
			})
		}

		resp := ListStudentsResponse{
			Students: studentResponses,
			Total:    len(studentResponses),
			Limit:    limit,
			Offset:   offset,
		}

		response.Writejson(w, http.StatusOK, resp)
	}
}
