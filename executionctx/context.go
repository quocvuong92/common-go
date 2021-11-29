package executionctx

import (
	"os"
	"strings"

	"github.com/cjlapao/common-go/helper"
	"github.com/google/uuid"
)

var contextService *Context

// Context entity
type Context struct {
	Configuration *Configuration
	CorrelationId string
	Environment   string
	IsDevelopment bool
	Debug         bool
	Init          func() error
}

func NewContext() (*Context, error) {
	if contextService != nil {
		contextService = nil
	}

	return InitNewContext(nil)
}

func InitNewContext(init func() error) (*Context, error) {
	contextService = &Context{
		IsDevelopment: false,
		Debug:         false,
		Init:          init,
	}

	contextService.CorrelationId = uuid.NewString()
	environment := os.Getenv("CJ_ENVIRONMENT")
	debug := os.Getenv("CJ_ENABLE_DEBUG")

	if !helper.IsNilOrEmpty(environment) {
		if strings.ToLower(environment) == "development" {
			contextService.IsDevelopment = true
			contextService.Environment = "Development"
		} else {
			contextService.IsDevelopment = false
			switch strings.ToLower(environment) {
			case "production":
				contextService.Environment = "Production"
			case "release":
				contextService.Environment = "Release"
			case "ci":
				contextService.IsDevelopment = true
				contextService.Environment = "CI"
			case "devprod":
				contextService.Environment = "DevProd"
			default:
				contextService.Environment = "Production"
			}
		}
	} else {
		contextService.IsDevelopment = false
		contextService.Environment = "Production"
	}

	if !helper.IsNilOrEmpty(debug) && strings.ToLower(debug) == "true" {
		contextService.Debug = true
	} else {
		contextService.Debug = false
	}

	contextService.Configuration = GetConfigService()

	if contextService.Init != nil {
		err := contextService.Init()
		if err != nil {
			contextService = nil
			return contextService, err
		}
	}

	return contextService, nil
}

func GetContext() *Context {
	if contextService != nil {
		return contextService
	}

	NewContext()

	return contextService
}