from flask import Flask
from .database import init_db
from flask_cors import CORS

def create_app():
    app = Flask(__name__)
    CORS(app, resources={r"/*": {"origins": "*"}})
    app.config["SQLALCHEMY_DATABASE_URI"] = "postgresql://kanban_user:kanban_password@db:5432/kanban_db"
    app.config["SQLALCHEMY_TRACK_MODIFICATIONS"] = False

    init_db(app)

    from .routes import bp
    app.register_blueprint(bp, url_prefix="/api")

    return app
