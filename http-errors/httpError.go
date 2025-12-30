package httperr

import (
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HttpExceptionJSONImpl struct {
	Status  int
	Message gin.H
}

func (e *HttpExceptionJSONImpl) Error() string {
	return e.Message["message"].(string)
}

func NewHttpExceptionJSON(status int, message gin.H) *HttpExceptionJSONImpl {
	return &HttpExceptionJSONImpl{
		Status:  status,
		Message: message,
	}
}

var (
	NotFound            = NewHttpExceptionJSON(404, gin.H{"message": "Not found"})
	BadRequest          = NewHttpExceptionJSON(400, gin.H{"message": "Bad request"})
	Success             = NewHttpExceptionJSON(200, gin.H{"message": "Success"})
	Created             = NewHttpExceptionJSON(201, gin.H{"message": "Success"})
	Unauthorized        = NewHttpExceptionJSON(401, gin.H{"message": "Unauthorized"})
	Forbidden           = NewHttpExceptionJSON(403, gin.H{"message": "Forbidden"})
	Conflict            = NewHttpExceptionJSON(409, gin.H{"message": "Conflict"})
	UnprocessableEntity = NewHttpExceptionJSON(422, gin.H{"message": "Unprocessable Entity"})
	TooManyRequests     = NewHttpExceptionJSON(429, gin.H{"message": "Too many requests"})
	InternalServerError = NewHttpExceptionJSON(500, gin.H{"message": "Internal server error"})
)

func HandlerError(err error) {
	if err == nil {
		return
	}

	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		err = NotFound
	case errors.Is(err, gorm.ErrInvalidData):
		err = BadRequest
	case errors.Is(err, gorm.ErrInvalidField):
		err = BadRequest
	case errors.Is(err, gorm.ErrDuplicatedKey):
		err = Conflict
	case errors.Is(err, gorm.ErrForeignKeyViolated):
		err = Conflict
	case errors.Is(err, gorm.ErrCheckConstraintViolated):
		err = Conflict
	case errors.Is(err, gorm.ErrInvalidValueOfLength):
		err = BadRequest
	case errors.Is(err, gorm.ErrEmptySlice):
		err = BadRequest
	case errors.Is(err, gorm.ErrInvalidDB):
		err = InternalServerError
	}

	panic(err)
}
