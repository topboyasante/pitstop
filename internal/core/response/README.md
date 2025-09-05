# Structured API Responses

This package provides a standardized way to return structured API responses, ensuring consistency across all endpoints.

## Response Structure

All API responses follow this structure:

```json
{
  "success": true,
  "message": "Operation completed successfully",
  "data": { ... },
  "error": null,
  "meta": { ... },
  "timestamp": "2023-12-01T10:30:00Z"
}
```

### Fields

- **success**: Boolean indicating if the request was successful
- **message**: Human-readable message describing the result
- **data**: The actual response data (only present on success)
- **error**: Error details (only present on failure)
- **meta**: Metadata like pagination info (optional)
- **timestamp**: ISO timestamp of when the response was generated

## Usage Examples

### Success Response

```go
// Simple success
return response.SuccessJSON(c, user, "User retrieved successfully")

// Success with pagination metadata
meta := response.NewPaginationMeta(page, limit, total, hasNext)
return response.SuccessJSONWithMeta(c, users, "Users retrieved successfully", meta)

// Created response
return response.CreatedJSON(c, newUser, "User created successfully")
```

### Error Responses

```go
// Validation error
return response.ValidationErrorJSON(c, "Invalid email format", err.Error())

// Not found error
return response.NotFoundJSON(c, "User")

// Internal server error
return response.InternalErrorJSON(c, "Database connection failed")

// Unauthorized
return response.UnauthorizedJSON(c)

// Custom error
return response.ErrorJSON(c, fiber.StatusTeapot, "CUSTOM_ERROR", "I'm a teapot", "")
```

## Response Examples

### Successful User Retrieval

```json
{
  "success": true,
  "message": "User retrieved successfully",
  "data": {
    "id": "123",
    "email": "user@example.com",
    "name": "John Doe"
  },
  "timestamp": "2023-12-01T10:30:00Z"
}
```

### Paginated Users List

```json
{
  "success": true,
  "message": "Users retrieved successfully",
  "data": [
    { "id": "1", "name": "User 1" },
    { "id": "2", "name": "User 2" }
  ],
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 25,
    "total_pages": 3,
    "has_next": true,
    "has_prev": false
  },
  "timestamp": "2023-12-01T10:30:00Z"
}
```

### Validation Error

```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid email format",
    "details": "email must be a valid email address"
  },
  "timestamp": "2023-12-01T10:30:00Z"
}
```

### Not Found Error

```json
{
  "success": false,
  "error": {
    "code": "NOT_FOUND",
    "message": "User not found"
  },
  "timestamp": "2023-12-01T10:30:00Z"
}
```

## Error Codes

Standard error codes used throughout the API:

- `VALIDATION_ERROR`: Request validation failed
- `NOT_FOUND`: Requested resource not found
- `UNAUTHORIZED`: Authentication required
- `FORBIDDEN`: Access denied
- `INTERNAL_ERROR`: Internal server error

## Migration from Old Responses

**Before:**
```go
return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
    "error": "Invalid request",
})
```

**After:**
```go
return response.ValidationErrorJSON(c, "Invalid request", err.Error())
```

This ensures all responses have consistent structure and error handling.