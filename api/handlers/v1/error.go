package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golanguzb70/simple-post-app/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *handlerV1) HandleDatabaseLevelWithMessage(c *gin.Context, err error, message string, args ...interface{}) bool {
	status_err, _ := status.FromError(err)
	if err != nil {
		errorCode := InternalServerError
		statuscode := http.StatusInternalServerError
		message := status_err.Message()
		switch status_err.Code() {
		case codes.NotFound:
			errorCode = NotFound
			statuscode = http.StatusNotFound
		case codes.Unknown:
			errorCode = InternalServerError
			statuscode = http.StatusBadRequest
			message = "Ooops something went wrong"
		case codes.Aborted:
			errorCode = BadRequest
			statuscode = http.StatusBadRequest
		case codes.InvalidArgument:
			errorCode = BadRequest
			statuscode = http.StatusBadRequest
		}

		h.log.Error(message, err, args)
		c.AbortWithStatusJSON(statuscode, models.StandardResponse{
			StatusId: errorCode,
			Message:  message,
		})
		return true
	}
	return false
}

// Handles response according to err arguments. If err is nil it returns false otherwise true
func (h *handlerV1) HandleResponse(c *gin.Context, err error, httpStatusCode int, statusId, message string, data any, args ...any) bool {
	if err != nil {
		if statusId != InternalServerError {
			c.AbortWithStatusJSON(httpStatusCode, models.StandardResponse{
				StatusId: statusId,
				Message:  message,
				Data:     data,
			})
		} else {
			h.log.Error(message, err, args)
			c.AbortWithStatusJSON(httpStatusCode, models.StandardResponse{
				StatusId: statusId,
				Message:  "Internal server error",
				Data:     data,
			})
		}
		return true
	} else if statusId == "success" {
		c.JSON(httpStatusCode, models.StandardResponse{
			StatusId: statusId,
			Message:  message,
			Data:     data,
		})
	}

	return false
}
