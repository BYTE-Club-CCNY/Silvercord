from marshmallow import Schema, fields, validate, ValidationError


class ProfessorRequestSchema(Schema):
    """Schema for professor endpoint request validation."""
    professor_name = fields.Str(
        required=True,
        validate=validate.Length(min=1, max=100),
        error_messages={
            "required": "Professor name is required.",
            "invalid": "Professor name must be a string.",
        }
    )
    question = fields.Str(
        required=False,
        missing="How is the following professor?",
        validate=validate.Length(max=500)
    )


class BreakRequestSchema(Schema):
    """Schema for break/calendar endpoint request validation."""
    year = fields.Str(
        required=True,
        validate=validate.Regexp(
            r"^\d{4}-\d{4}$",
            error="Year must be in format YYYY-YYYY (e.g., 2025-2026)"
        ),
        error_messages={
            "required": "Academic year is required.",
            "invalid": "Year must be a string.",
        }
    )


class ResponseSchema(Schema):
    """Schema for API responses."""
    name = fields.Str(required=True)
    link = fields.Str(allow_none=True)
    response = fields.Str(required=True)
    processed_at = fields.Str(required=True)


class ErrorSchema(Schema):
    """Schema for error responses."""
    error = fields.Dict(keys=fields.Str(), values=fields.Raw())


class HealthResponseSchema(Schema):
    """Schema for health check response."""
    status = fields.Str(required=True)
    services = fields.Dict(keys=fields.Str(), values=fields.Str())
    timestamp = fields.Str(required=True)
