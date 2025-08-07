package v1

import (
	"github.com/go-playground/validator/v10"
	"github.com/xiabin827/task-machinery/internal/usecase"
	"github.com/xiabin827/task-machinery/pkg/logger"
)

// V1 -.
type V1 struct {
	t usecase.Translation
	l logger.Interface
	v *validator.Validate
}
