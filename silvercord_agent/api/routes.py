from datetime import datetime

from api.schemas import BreakRequestSchema, ProfessorRequestSchema
from flasgger import swag_from
from flask import Blueprint, jsonify, request
from marshmallow import ValidationError
from services.agent_service import get_break_info, get_professor_info
from services.db_service import check_health

api_bp = Blueprint('api', __name__)


@api_bp.route('/professor', methods=['POST'])
def professor_endpoint():
    """
    Retrieve professor information using RAG
    ---
    tags:
      - Professor
    parameters:
      - in: body
        name: body
        required: true
        schema:
          type: object
          required:
            - professor_name
          properties:
            professor_name:
              type: string
              description: Name of the professor (1-100 characters)
              example: "Douglas Troeger"
            question:
              type: string
              description: Optional question about the professor
              example: "How is the following professor?"
              default: "How is the following professor?"
    responses:
      200:
        description: Successfully retrieved professor information
        schema:
          type: object
          properties:
            name:
              type: string
              example: "Douglas Troeger"
            link:
              type: string
              nullable: true
              example: "https://www.ratemyprofessors.com/professor/..."
            response:
              type: string
              example: "Professor Troeger is known for his engaging teaching style..."
            processed_at:
              type: string
              format: date-time
              example: "2025-12-15T10:30:00"
      400:
        description: Invalid request parameters
      500:
        description: Internal server error
      503:
        description: Service unavailable
    """
    schema = ProfessorRequestSchema()

    try:
        data = schema.load(request.json or {})
    except ValidationError as e:
        raise e

    result = get_professor_info(
        professor_name=data['professor_name'],
        question=data.get('question', 'How is the following professor?')
    )

    return jsonify(result), 200


@api_bp.route('/break', methods=['POST'])
def break_endpoint():
    """
    Retrieve academic calendar/break information
    ---
    tags:
      - Calendar
    parameters:
      - in: body
        name: body
        required: true
        schema:
          type: object
          required:
            - year
          properties:
            year:
              type: string
              description: Academic year in format YYYY-YYYY
              example: "2025-2026"
              pattern: r"^\\d{4}-\\d{4}$"
    responses:
      200:
        description: Successfully retrieved calendar information
        schema:
          type: object
          properties:
            name:
              type: string
              example: "Calendar"
            link:
              type: string
              nullable: true
              example: null
            response:
              type: string
              example: "This command is under maintenance. Stay tuned!"
            processed_at:
              type: string
              format: date-time
              example: "2025-12-15T10:30:00"
      400:
        description: Invalid request parameters
      500:
        description: Internal server error
      501:
        description: Feature not implemented
    """
    schema = BreakRequestSchema()

    try:
        data = schema.load(request.json or {})
    except ValidationError as e:
        raise e

    result = get_break_info(year=data['year'])

    return jsonify(result), 200


@api_bp.route('/health', methods=['GET'])
def health():
    """
    Health check endpoint
    ---
    tags:
      - Health
    responses:
      200:
        description: Service is healthy
        schema:
          type: object
          properties:
            status:
              type: string
              enum: [healthy, degraded]
              example: "healthy"
            services:
              type: object
              properties:
                chromadb:
                  type: string
                  enum: [connected, disconnected]
                  example: "connected"
                cohere:
                  type: string
                  enum: [configured, missing]
                  example: "configured"
                openai:
                  type: string
                  enum: [configured, missing]
                  example: "configured"
            timestamp:
              type: string
              format: date-time
              example: "2025-12-15T10:30:00"
      503:
        description: Service is degraded or unhealthy
    """
    health_status = check_health()

    # Determine overall status
    all_healthy = (
        health_status['chromadb'] == 'connected' and
        health_status['cohere'] == 'configured' and
        health_status['openai'] == 'configured'
    )

    status = 'healthy' if all_healthy else 'degraded'
    code = 200 if status == 'healthy' else 503

    return jsonify({
        'status': status,
        'services': health_status,
        'timestamp': datetime.utcnow().isoformat()
    }), code
