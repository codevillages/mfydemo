package response

import (
	"net/http"

	"github.com/gogf/gf/v2/errors/gcode"
)

// HTTPStatus maps business code to HTTP status.
func HTTPStatus(code int) int {
	switch code {
	case CodeBadReq.Code(), gcode.CodeValidationFailed.Code():
		return http.StatusBadRequest
	case CodeNotFound.Code():
		return http.StatusNotFound
	case CodeConflict.Code():
		return http.StatusConflict
	case CodeOK.Code():
		return http.StatusOK
	default:
		return http.StatusInternalServerError
	}
}
