#!/bin/bash

# Module generator script for Pitstop modular monolith
# Usage: ./scripts/generate-module.sh <module-name> [description]

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Check if module name is provided
if [ -z "$1" ]; then
    echo -e "${RED}Error: Module name is required${NC}"
    echo "Usage: $0 <module-name> [description]"
    echo "Example: $0 garage 'Car garage management'"
    exit 1
fi

MODULE_NAME="$1"
MODULE_DESCRIPTION="${2:-$MODULE_NAME module}"

# Convert module name to capitalized version for types
MODULE_NAME_CAP=$(echo "$MODULE_NAME" | awk '{print toupper(substr($0,1,1)) tolower(substr($0,2))}')

# Validate module name (lowercase, alphanumeric, no spaces)
if [[ ! "$MODULE_NAME" =~ ^[a-z][a-z0-9]*$ ]]; then
    echo -e "${RED}Error: Module name must be lowercase alphanumeric, starting with a letter${NC}"
    echo "Valid examples: garage, notification, search"
    exit 1
fi

# Check if module already exists
MODULE_PATH="internal/modules/$MODULE_NAME"
if [ -d "$MODULE_PATH" ]; then
    echo -e "${YELLOW}Warning: Module '$MODULE_NAME' already exists at $MODULE_PATH${NC}"
    read -p "Do you want to overwrite it? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "Aborted."
        exit 1
    fi
    rm -rf "$MODULE_PATH"
fi

echo -e "${BLUE}Generating module: $MODULE_NAME${NC}"
echo -e "${BLUE}Description: $MODULE_DESCRIPTION${NC}"
echo

# Create module directory structure
echo -e "${GREEN}Creating directory structure...${NC}"
mkdir -p "$MODULE_PATH"/{domain,repository,service,handler,dto}

# Generate domain file
echo -e "${GREEN}Creating domain model...${NC}"
cat > "$MODULE_PATH/domain/${MODULE_NAME}.go" << EOF
package domain

import (
	"time"

	"gorm.io/gorm"
)

// ${MODULE_NAME_CAP} represents a ${MODULE_NAME} entity
type ${MODULE_NAME_CAP} struct {
	ID        uint           \`gorm:"primarykey" json:"id"\`
	Name      string         \`gorm:"not null;size:255" json:"name" validate:"required,min=1,max=255"\`
	CreatedAt time.Time      \`json:"created_at"\`
	UpdatedAt time.Time      \`json:"updated_at"\`
	DeletedAt gorm.DeletedAt \`gorm:"index" json:"-"\`
}

// TableName specifies the table name for the ${MODULE_NAME_CAP} model
func (${MODULE_NAME_CAP}) TableName() string {
	return "${MODULE_NAME}s"
}
EOF

# Generate repository file
echo -e "${GREEN}Creating repository...${NC}"
cat > "$MODULE_PATH/repository/${MODULE_NAME}_repository.go" << EOF
package repository

import (
	"github.com/topboyasante/pitstop/internal/modules/${MODULE_NAME}/domain"
	"gorm.io/gorm"
)

// ${MODULE_NAME_CAP}Repository handles ${MODULE_NAME} data operations
type ${MODULE_NAME_CAP}Repository struct {
	db *gorm.DB
}

// New${MODULE_NAME_CAP}Repository creates a new ${MODULE_NAME} repository instance
func New${MODULE_NAME_CAP}Repository(db *gorm.DB) *${MODULE_NAME_CAP}Repository {
	return &${MODULE_NAME_CAP}Repository{db: db}
}

// Create creates a new ${MODULE_NAME}
func (r *${MODULE_NAME_CAP}Repository) Create(${MODULE_NAME} *domain.${MODULE_NAME_CAP}) error {
	return r.db.Create(${MODULE_NAME}).Error
}

// GetByID retrieves a ${MODULE_NAME} by ID
func (r *${MODULE_NAME_CAP}Repository) GetByID(id uint) (*domain.${MODULE_NAME_CAP}, error) {
	var ${MODULE_NAME} domain.${MODULE_NAME_CAP}
	err := r.db.First(&${MODULE_NAME}, id).Error
	if err != nil {
		return nil, err
	}
	return &${MODULE_NAME}, nil
}

// GetAll retrieves all ${MODULE_NAME}s with pagination
func (r *${MODULE_NAME_CAP}Repository) GetAll(page, limit int) ([]domain.${MODULE_NAME_CAP}, int64, error) {
	var ${MODULE_NAME}s []domain.${MODULE_NAME_CAP}
	var totalCount int64

	offset := (page - 1) * limit

	// Get total count
	if err := r.db.Model(&domain.${MODULE_NAME_CAP}{}).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Get ${MODULE_NAME}s
	if err := r.db.Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&${MODULE_NAME}s).Error; err != nil {
		return nil, 0, err
	}

	return ${MODULE_NAME}s, totalCount, nil
}

// Update updates a ${MODULE_NAME}
func (r *${MODULE_NAME_CAP}Repository) Update(${MODULE_NAME} *domain.${MODULE_NAME_CAP}) error {
	return r.db.Save(${MODULE_NAME}).Error
}

// Delete soft deletes a ${MODULE_NAME}
func (r *${MODULE_NAME_CAP}Repository) Delete(id uint) error {
	return r.db.Delete(&domain.${MODULE_NAME_CAP}{}, id).Error
}
EOF

# Generate service file
echo -e "${GREEN}Creating service...${NC}"
cat > "$MODULE_PATH/service/${MODULE_NAME}_service.go" << EOF
package service

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/topboyasante/pitstop/internal/modules/${MODULE_NAME}/domain"
	"github.com/topboyasante/pitstop/internal/modules/${MODULE_NAME}/dto"
	"github.com/topboyasante/pitstop/internal/modules/${MODULE_NAME}/repository"
	"github.com/topboyasante/pitstop/internal/shared/events"
	"github.com/topboyasante/pitstop/internal/core/logger"
)

// ${MODULE_NAME_CAP}Service handles ${MODULE_NAME} business logic
type ${MODULE_NAME_CAP}Service struct {
	${MODULE_NAME}Repo *repository.${MODULE_NAME_CAP}Repository
	validator *validator.Validate
	eventBus  *events.EventBus
}

// New${MODULE_NAME_CAP}Service creates a new ${MODULE_NAME} service instance
func New${MODULE_NAME_CAP}Service(${MODULE_NAME}Repo *repository.${MODULE_NAME_CAP}Repository, validator *validator.Validate, eventBus *events.EventBus) *${MODULE_NAME_CAP}Service {
	return &${MODULE_NAME_CAP}Service{
		${MODULE_NAME}Repo: ${MODULE_NAME}Repo,
		validator: validator,
		eventBus:  eventBus,
	}
}

// Create${MODULE_NAME_CAP} creates a new ${MODULE_NAME}
func (s *${MODULE_NAME_CAP}Service) Create${MODULE_NAME_CAP}(req dto.Create${MODULE_NAME_CAP}Request) (*dto.${MODULE_NAME_CAP}Response, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	${MODULE_NAME} := &domain.${MODULE_NAME_CAP}{
		Name: req.Name,
	}

	if err := s.${MODULE_NAME}Repo.Create(${MODULE_NAME}); err != nil {
		logger.Error("Failed to create ${MODULE_NAME}", "error", err)
		return nil, fmt.Errorf("failed to create ${MODULE_NAME}: %w", err)
	}

	logger.Info("${MODULE_NAME_CAP} created successfully", "${MODULE_NAME}_id", ${MODULE_NAME}.ID)

	return &dto.${MODULE_NAME_CAP}Response{
		ID:        ${MODULE_NAME}.ID,
		Name:      ${MODULE_NAME}.Name,
		CreatedAt: ${MODULE_NAME}.CreatedAt,
	}, nil
}

// Get${MODULE_NAME_CAP}ByID retrieves a ${MODULE_NAME} by ID
func (s *${MODULE_NAME_CAP}Service) Get${MODULE_NAME_CAP}ByID(id uint) (*dto.${MODULE_NAME_CAP}Response, error) {
	${MODULE_NAME}, err := s.${MODULE_NAME}Repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("${MODULE_NAME} not found: %w", err)
	}

	return &dto.${MODULE_NAME_CAP}Response{
		ID:        ${MODULE_NAME}.ID,
		Name:      ${MODULE_NAME}.Name,
		CreatedAt: ${MODULE_NAME}.CreatedAt,
	}, nil
}

// GetAll${MODULE_NAME_CAP}s retrieves all ${MODULE_NAME}s with pagination
func (s *${MODULE_NAME_CAP}Service) GetAll${MODULE_NAME_CAP}s(page, limit int) (*dto.${MODULE_NAME_CAP}sResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	${MODULE_NAME}s, totalCount, err := s.${MODULE_NAME}Repo.GetAll(page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve ${MODULE_NAME}s: %w", err)
	}

	${MODULE_NAME}Responses := make([]dto.${MODULE_NAME_CAP}Response, len(${MODULE_NAME}s))
	for i, ${MODULE_NAME} := range ${MODULE_NAME}s {
		${MODULE_NAME}Responses[i] = dto.${MODULE_NAME_CAP}Response{
			ID:        ${MODULE_NAME}.ID,
			Name:      ${MODULE_NAME}.Name,
			CreatedAt: ${MODULE_NAME}.CreatedAt,
		}
	}

	hasNext := int64((page-1)*limit+len(${MODULE_NAME}s)) < totalCount

	return &dto.${MODULE_NAME_CAP}sResponse{
		${MODULE_NAME_CAP}s:    ${MODULE_NAME}Responses,
		TotalCount: totalCount,
		Page:       page,
		Limit:      limit,
		HasNext:    hasNext,
	}, nil
}
EOF

# Generate DTOs
echo -e "${GREEN}Creating DTOs...${NC}"
cat > "$MODULE_PATH/dto/${MODULE_NAME}_dto.go" << EOF
package dto

import "time"

// Create${MODULE_NAME_CAP}Request represents a request to create a new ${MODULE_NAME}
type Create${MODULE_NAME_CAP}Request struct {
	Name string \`json:"name" validate:"required,min=1,max=255"\`
}

// Update${MODULE_NAME_CAP}Request represents a request to update a ${MODULE_NAME}
type Update${MODULE_NAME_CAP}Request struct {
	Name string \`json:"name,omitempty" validate:"omitempty,min=1,max=255"\`
}

// ${MODULE_NAME_CAP}Response represents a ${MODULE_NAME} in API responses
type ${MODULE_NAME_CAP}Response struct {
	ID        uint      \`json:"id"\`
	Name      string    \`json:"name"\`
	CreatedAt time.Time \`json:"created_at"\`
}

// ${MODULE_NAME_CAP}sResponse represents a paginated list of ${MODULE_NAME}s
type ${MODULE_NAME_CAP}sResponse struct {
	${MODULE_NAME_CAP}s    []${MODULE_NAME_CAP}Response \`json:"${MODULE_NAME}s"\`
	TotalCount int64                    \`json:"total_count"\`
	Page       int                      \`json:"page"\`
	Limit      int                      \`json:"limit"\`
	HasNext    bool                     \`json:"has_next"\`
}
EOF

# Generate handler
echo -e "${GREEN}Creating handler...${NC}"
cat > "$MODULE_PATH/handler/${MODULE_NAME}_handler.go" << EOF
package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/topboyasante/pitstop/internal/modules/${MODULE_NAME}/dto"
	"github.com/topboyasante/pitstop/internal/modules/${MODULE_NAME}/service"
	"github.com/topboyasante/pitstop/internal/core/logger"
)

// ${MODULE_NAME_CAP}Handler handles HTTP requests for ${MODULE_NAME}s
type ${MODULE_NAME_CAP}Handler struct {
	${MODULE_NAME}Service *service.${MODULE_NAME_CAP}Service
}

// New${MODULE_NAME_CAP}Handler creates a new ${MODULE_NAME} handler instance
func New${MODULE_NAME_CAP}Handler(${MODULE_NAME}Service *service.${MODULE_NAME_CAP}Service) *${MODULE_NAME_CAP}Handler {
	return &${MODULE_NAME_CAP}Handler{
		${MODULE_NAME}Service: ${MODULE_NAME}Service,
	}
}

// GetAll${MODULE_NAME_CAP}s retrieves all ${MODULE_NAME}s
// @Summary Get all ${MODULE_NAME}s
// @Description Retrieve a paginated list of ${MODULE_NAME}s
// @Tags ${MODULE_NAME}s
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "${MODULE_NAME_CAP}s per page" default(20)
// @Success 200 {object} dto.${MODULE_NAME_CAP}sResponse
// @Router /${MODULE_NAME}s [get]
func (h *${MODULE_NAME_CAP}Handler) GetAll${MODULE_NAME_CAP}s(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	${MODULE_NAME}s, err := h.${MODULE_NAME}Service.GetAll${MODULE_NAME_CAP}s(page, limit)
	if err != nil {
		logger.Error("Failed to retrieve ${MODULE_NAME}s", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve ${MODULE_NAME}s",
		})
	}

	return c.JSON(${MODULE_NAME}s)
}

// Create${MODULE_NAME_CAP} creates a new ${MODULE_NAME}
// @Summary Create a new ${MODULE_NAME}
// @Description Create a new ${MODULE_NAME}
// @Tags ${MODULE_NAME}s
// @Accept json
// @Produce json
// @Param request body dto.Create${MODULE_NAME_CAP}Request true "${MODULE_NAME_CAP} details"
// @Success 201 {object} dto.${MODULE_NAME_CAP}Response
// @Failure 400 {object} map[string]string
// @Router /${MODULE_NAME}s [post]
func (h *${MODULE_NAME_CAP}Handler) Create${MODULE_NAME_CAP}(c *fiber.Ctx) error {
	var req dto.Create${MODULE_NAME_CAP}Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	${MODULE_NAME}, err := h.${MODULE_NAME}Service.Create${MODULE_NAME_CAP}(req)
	if err != nil {
		logger.Error("Failed to create ${MODULE_NAME}", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	logger.Info("${MODULE_NAME_CAP} created successfully", "${MODULE_NAME}_id", ${MODULE_NAME}.ID)
	return c.Status(fiber.StatusCreated).JSON(${MODULE_NAME})
}

// Get${MODULE_NAME_CAP} retrieves a specific ${MODULE_NAME} by ID
// @Summary Get a ${MODULE_NAME} by ID
// @Description Retrieve a specific ${MODULE_NAME}
// @Tags ${MODULE_NAME}s
// @Accept json
// @Produce json
// @Param id path int true "${MODULE_NAME_CAP} ID"
// @Success 200 {object} dto.${MODULE_NAME_CAP}Response
// @Failure 404 {object} map[string]string
// @Router /${MODULE_NAME}s/{id} [get]
func (h *${MODULE_NAME_CAP}Handler) Get${MODULE_NAME_CAP}(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ${MODULE_NAME} ID",
		})
	}

	${MODULE_NAME}, err := h.${MODULE_NAME}Service.Get${MODULE_NAME_CAP}ByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "${MODULE_NAME_CAP} not found",
		})
	}

	return c.JSON(${MODULE_NAME})
}
EOF

# Generate routes file
echo -e "${GREEN}Creating routes...${NC}"
cat > "$MODULE_PATH/routes.go" << EOF
package ${MODULE_NAME}

import (
	"github.com/gofiber/fiber/v2"
	"github.com/topboyasante/pitstop/internal/modules/${MODULE_NAME}/handler"
)

// RegisterRoutes registers all ${MODULE_NAME}-related routes
func RegisterRoutes(router fiber.Router, ${MODULE_NAME}Handler *handler.${MODULE_NAME_CAP}Handler) {
	${MODULE_NAME}s := router.Group("/${MODULE_NAME}s")
	
	${MODULE_NAME}s.Get("/", ${MODULE_NAME}Handler.GetAll${MODULE_NAME_CAP}s)
	${MODULE_NAME}s.Post("/", ${MODULE_NAME}Handler.Create${MODULE_NAME_CAP})        // Requires auth
	${MODULE_NAME}s.Get("/:id", ${MODULE_NAME}Handler.Get${MODULE_NAME_CAP})
}
EOF

# Generate events (optional - add to shared/events/types.go)
echo -e "${GREEN}Creating event types...${NC}"
cat > "/tmp/${MODULE_NAME}_events.go" << EOF
// Add these events to internal/shared/events/types.go

// ${MODULE_NAME_CAP} Events
type ${MODULE_NAME_CAP}Created struct {
	BaseEvent
	${MODULE_NAME_CAP}ID uint   \`json:"${MODULE_NAME}_id"\`
	Name      string \`json:"name"\`
}

func New${MODULE_NAME_CAP}Created(${MODULE_NAME}ID uint, name string) *${MODULE_NAME_CAP}Created {
	return &${MODULE_NAME_CAP}Created{
		BaseEvent: BaseEvent{
			Name:      "${MODULE_NAME}.created",
			Timestamp: time.Now(),
		},
		${MODULE_NAME_CAP}ID: ${MODULE_NAME}ID,
		Name:      name,
	}
}
EOF

# Auto-register the module in provider and main.go
echo -e "${GREEN}Auto-registering module...${NC}"

# Add import to provider
PROVIDER_FILE="internal/provider/provider.go"
if ! grep -q "internal/modules/$MODULE_NAME/handler" "$PROVIDER_FILE"; then
    # Find the last import line and add after it
    sed -i '' "/contentService \"github.com\/topboyasante\/pitstop\/internal\/modules\/content\/service\"/a\\
	${MODULE_NAME}Handler \"github.com/topboyasante/pitstop/internal/modules/$MODULE_NAME/handler\"\\
	${MODULE_NAME}Repository \"github.com/topboyasante/pitstop/internal/modules/$MODULE_NAME/repository\"\\
	${MODULE_NAME}Service \"github.com/topboyasante/pitstop/internal/modules/$MODULE_NAME/service\"
" "$PROVIDER_FILE"
    
    # Add handler to Provider struct
    sed -i '' "/PostHandler \*contentHandler\.PostHandler/a\\
	\\
	// ${MODULE_NAME_CAP} module\\
	${MODULE_NAME_CAP}Handler *${MODULE_NAME}Handler.${MODULE_NAME_CAP}Handler
" "$PROVIDER_FILE"
    
    # Add initialization in NewProvider function
    sed -i '' "/postHandler := contentHandler\.NewPostHandler(postService)/a\\
	\\
	// Initialize ${MODULE_NAME_CAP} module\\
	${MODULE_NAME}Repo := ${MODULE_NAME}Repository.New${MODULE_NAME_CAP}Repository(db)\\
	${MODULE_NAME}Service := ${MODULE_NAME}Service.New${MODULE_NAME_CAP}Service(${MODULE_NAME}Repo, validator, eventBus)\\
	${MODULE_NAME}Handler := ${MODULE_NAME}Handler.New${MODULE_NAME_CAP}Handler(${MODULE_NAME}Service)
" "$PROVIDER_FILE"
    
    # Add to return statement
    sed -i '' "/PostHandler: postHandler,/a\\
	\\
		${MODULE_NAME_CAP}Handler: ${MODULE_NAME}Handler,
" "$PROVIDER_FILE"
    
    echo "   ✅ Added ${MODULE_NAME_CAP} module to provider"
else
    echo "   ⚠️  ${MODULE_NAME_CAP} module already in provider"
fi

# Add import and route registration to main.go
MAIN_FILE="cmd/server/main.go"
if ! grep -q "internal/modules/$MODULE_NAME" "$MAIN_FILE"; then
    # Add import
    sed -i '' "/\"github.com\/topboyasante\/pitstop\/internal\/modules\/content\"/a\\
	\"github.com/topboyasante/pitstop/internal/modules/$MODULE_NAME\"
" "$MAIN_FILE"
    
    # Add route registration
    sed -i '' "/content\.RegisterRoutes(v1, provider\.PostHandler)/a\\
	${MODULE_NAME}.RegisterRoutes(v1, provider.${MODULE_NAME_CAP}Handler)
" "$MAIN_FILE"
    
    echo "   ✅ Added ${MODULE_NAME_CAP} routes to main.go"
else
    echo "   ⚠️  ${MODULE_NAME_CAP} routes already in main.go"
fi

echo
echo -e "${GREEN}✅ Module '$MODULE_NAME' generated and registered successfully!${NC}"
echo
echo -e "${YELLOW}📁 Generated files:${NC}"
echo "   $MODULE_PATH/domain/${MODULE_NAME}.go"
echo "   $MODULE_PATH/repository/${MODULE_NAME}_repository.go"
echo "   $MODULE_PATH/service/${MODULE_NAME}_service.go"
echo "   $MODULE_PATH/handler/${MODULE_NAME}_handler.go"
echo "   $MODULE_PATH/dto/${MODULE_NAME}_dto.go"
echo "   $MODULE_PATH/routes.go"
echo
echo -e "${YELLOW}🔧 Auto-registered:${NC}"
echo "   ✅ Added imports to internal/provider/provider.go"
echo "   ✅ Added ${MODULE_NAME_CAP}Handler to Provider struct"
echo "   ✅ Added module initialization in NewProvider()"
echo "   ✅ Added route registration in cmd/server/main.go"
echo
echo -e "${YELLOW}📋 Next steps:${NC}"
echo "1. Add event types from /tmp/${MODULE_NAME}_events.go to internal/shared/events/types.go"
echo "2. Run database migrations to create the '${MODULE_NAME}s' table"
echo "3. Run 'make docs -B' to regenerate Swagger documentation"
echo "4. Test the endpoints:"
echo "   GET    /api/v1/${MODULE_NAME}s"
echo "   POST   /api/v1/${MODULE_NAME}s"
echo "   GET    /api/v1/${MODULE_NAME}s/:id"
echo
echo -e "${BLUE}🚀 Module '$MODULE_NAME' is ready for development and fully integrated!${NC}"