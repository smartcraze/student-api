package httphandler

import (
	"net/http"
	"strconv"

	"github.com/smartcraze/student-api/internal/storage"
	"github.com/smartcraze/student-api/utils/response"
)

func DeleteStudentHandler(store storage.Storage) http.HandlerFunc {
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

		// Delete student from database
		err = store.DeleteStudent(r.Context(), id)
		if err != nil {
			response.Writejson(w, http.StatusNotFound, response.Response{
				Status: response.StatusError,
				Error:  err.Error(),
			})
			return
		}

		// Return success response
		response.Writejson(w, http.StatusOK, response.Response{
			Status: response.StatusOK,
			Error:  "student deleted successfully",
		})
	}
}
