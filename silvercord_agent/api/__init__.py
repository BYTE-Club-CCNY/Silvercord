import os
from flask import Flask
from flask_cors import CORS
from flasgger import Swagger


def create_app():
    """Flask application factory."""
    app = Flask(__name__)

    # Configuration
    app.config['JSON_SORT_KEYS'] = False
    app.config['SWAGGER'] = {
        'title': 'Silvercord Agent API',
        'version': '1.0.0',
        'description': 'REST API for Silvercord professor information agent',
        'uiversion': 3
    }

    # CORS configuration
    allowed_origins = [
        'http://localhost:*',
        'http://127.0.0.1:*',
        'http://discord-bot:*',
        'http://silvercord-discord:*'
    ]

    # Add custom origin if provided
    custom_origin = os.getenv('DISCORD_BOT_URL', '')
    if custom_origin:
        allowed_origins.append(custom_origin)

    CORS(app, origins=allowed_origins, supports_credentials=True)

    # Initialize Swagger documentation
    Swagger(app)

    # Register blueprints
    from api.routes import api_bp
    app.register_blueprint(api_bp, url_prefix='/api/v1')

    # Register error handlers
    from api.errors import register_error_handlers
    register_error_handlers(app)

    # Initialize database clients
    from services.db_service import init_clients
    with app.app_context():
        init_clients()
        app.logger.info("Database clients initialized successfully")

    return app
