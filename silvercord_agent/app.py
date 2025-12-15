#!/usr/bin/env python
"""
Silvercord Agent Flask API Server

Main entry point for the Flask application.
Run with: python app.py (development) or gunicorn app:app (production)
"""

from api import create_app

app = create_app()

if __name__ == '__main__':
    # Development server only
    # In production, use Gunicorn: gunicorn --bind 0.0.0.0:5000 app:app
    app.run(host='0.0.0.0', port=5000, debug=False)
