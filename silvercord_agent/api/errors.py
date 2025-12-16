from datetime import datetime
from flask import jsonify
from marshmallow import ValidationError


def register_error_handlers(app):
    """Register global error handlers for the Flask app."""

    @app.errorhandler(ValidationError)
    def handle_validation_error(error):
        """Handle Marshmallow validation errors."""
        return jsonify({
            "error": {
                "code": "INVALID_INPUT",
                "message": "Request validation failed",
                "details": error.messages,
                "timestamp": datetime.utcnow().isoformat()
            }
        }), 400

    @app.errorhandler(404)
    def handle_not_found(error):
        """Handle 404 errors."""
        return jsonify({
            "error": {
                "code": "NOT_FOUND",
                "message": "The requested endpoint does not exist",
                "timestamp": datetime.utcnow().isoformat()
            }
        }), 404

    @app.errorhandler(500)
    def handle_internal_error(error):
        """Handle internal server errors."""
        return jsonify({
            "error": {
                "code": "INTERNAL_ERROR",
                "message": "An internal server error occurred",
                "timestamp": datetime.utcnow().isoformat()
            }
        }), 500

    @app.errorhandler(Exception)
    def handle_generic_exception(error):
        """Handle all uncaught exceptions."""
        # Log the error for debugging
        app.logger.error(f"Unhandled exception: {str(error)}", exc_info=True)

        # Determine error type
        error_message = str(error)
        if "professor" in error_message.lower() and "not found" in error_message.lower():
            code = "PROFESSOR_NOT_FOUND"
            message = "Could not find information about this professor"
        elif "chromadb" in error_message.lower() or "database" in error_message.lower():
            code = "DATABASE_ERROR"
            message = "Database connection error"
        elif "cohere" in error_message.lower() or "openai" in error_message.lower():
            code = "EXTERNAL_API_ERROR"
            message = "External API service error"
        elif "scraping" in error_message.lower():
            code = "SCRAPING_ERROR"
            message = "Failed to scrape professor data"
        else:
            code = "INTERNAL_ERROR"
            message = "An unexpected error occurred"

        return jsonify({
            "error": {
                "code": code,
                "message": message,
                "timestamp": datetime.utcnow().isoformat()
            }
        }), 500
