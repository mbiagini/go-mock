package apierrors

import "fmt"

// error codes
const (
	// default
	ERR_NOT_DEFINED 	  ErrorCode = -1
	// business errors
	ENDPOINT_NOT_FOUND    ErrorCode = 1001
	ERR_SAVING_FILE       ErrorCode = 1002
	// client errors
	OPERATION_NOT_DEFINED ErrorCode = 2001
	INVALID_ARGUMENT 	  ErrorCode = 2002
	FILES_SIZE_EXCEEDED   ErrorCode = 2003
	// server errors
	INTERNAL_SERVER_ERROR ErrorCode = 5001
	IO_FILE_ERROR 		  ErrorCode = 5002
	JSON_PARSING_ERROR 	  ErrorCode = 5003
	READER_ERROR          ErrorCode = 5004
)

var errors = map[ErrorCode]ErrorMessage {
	// default
	ERR_NOT_DEFINED: {
		Label: "ERR_NOT_DEFINED", 
		Message: "Unrecognized error encountered. Contact support team",
	},
	// business errors
	ENDPOINT_NOT_FOUND: {
		Label: "ENDPOINT_NOT_FOUND", 
		Message: "Endpoint could not be found",
	},
	ERR_SAVING_FILE: {
		Label: "ERR_SAVING_FILE",
		Message: "Uploaded file could not be saved",
	},
	// client errors
	OPERATION_NOT_DEFINED: {
		Label: "OPERATION_NOT_DEFINED",
		Message: "Invoked path/operation is not defined",
	},
	INVALID_ARGUMENT: {
		Label: "INVALID_ARGUMENT",
		Message: "Client specified invalid request parameter",
	},
	FILES_SIZE_EXCEEDED: {
		Label: "FILES_SIZE_EXCEEDED",
		Message: "Size of the uploaded files exceed the maximum allowed",
	},
	// server errors
	INTERNAL_SERVER_ERROR: {
		Label: "INTERNAL_SERVER_ERROR",
		Message: "Internal error found. Please contact support team",
	},
	IO_FILE_ERROR: {
		Label: "IO_FILE_ERROR",
		Message: "Input/Output error while reading a file",
	},
	JSON_PARSING_ERROR: {
		Label: "JSON_PARSING_ERROR",
		Message: "Error while parsing a JSON content",
	},
	READER_ERROR: {
		Label: "READER_ERROR",
		Message: "Error while reading Body from Request",
	},
}

type ErrorCode int

type ErrorMessage struct {
	Label 	string
	Message string
}

// the serializable error structure
type Error struct {
	Code 	ErrorCode `json:"code"`
	Label 	string 	  `json:"label"`
	Message string    `json:"message"`
}

func (e *Error) Error() string {
	return e.ToString()
}

func (e *Error) ToString() string {
	return fmt.Sprintf("[%d] %s: %s", e.Code, e.Label, e.Message)
}

func New(ec ErrorCode) *Error {
	errMsg := errors[ERR_NOT_DEFINED]
	if em, ok := errors[ec]; ok {
		errMsg = em
	}
	return &Error{
		Code: ec,
		Label: errMsg.Label,
		Message: errMsg.Message,
	}
}

func NewWithMsg(ec ErrorCode, m string) *Error {
	errMsg := errors[ERR_NOT_DEFINED]
	if em, ok := errors[ec]; ok {
		errMsg = em
	}
	return &Error{
		Code: ec,
		Label: errMsg.Label,
		Message: m,
	}
}