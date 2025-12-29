package response

import "github.com/gogf/gf/v2/errors/gcode"

var (
	CodeOK       = gcode.New(0, "OK", nil)
	CodeBadReq   = gcode.New(40001, "Bad Request", nil)
	CodeNotFound = gcode.New(40401, "Not Found", nil)
	CodeConflict = gcode.New(40901, "Conflict", nil)
	CodeInternal = gcode.New(50000, "Internal Error", nil)
)
