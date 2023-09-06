package utils

import (
	"fmt"

	"github.com/google/uuid"
)

func GenerateId() string {
	return fmt.Sprintf("%s", uuid.New())
}
