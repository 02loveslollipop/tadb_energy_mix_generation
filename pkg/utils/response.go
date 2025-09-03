package utils

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

// Response represents a standard API response
type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// SuccessResponse sends a successful response
func SuccessResponse(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, Response{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

// ErrorResponse sends an error response
func ErrorResponse(c *gin.Context, code int, message string) {
	c.JSON(code, Response{
		Status: "error",
		Error:  message,
	})
}

// LogError logs an error with context
func LogError(context string, err error) {
	log.Printf("ERROR [%s]: %v", context, err)
}

// LogInfo logs an info message
func LogInfo(message string) {
	log.Printf("INFO: %s", message)
}

// PrettyPrint prints a struct in a formatted JSON
func PrettyPrint(v interface{}) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println(string(b))
}

// ValidateJSON validates if a string is valid JSON
func ValidateJSON(data string) error {
	var js json.RawMessage
	return json.Unmarshal([]byte(data), &js)
}
