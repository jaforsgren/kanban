from .database import db
import uuid
from sqlalchemy.dialects.postgresql import UUID
from sqlalchemy.orm import relationship

class User(db.Model):
    __tablename__ = 'users'
    id = db.Column(UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
    username = db.Column(db.String(80), nullable=False)
    email = db.Column(db.String(120), nullable=False)

class Board(db.Model):
    __tablename__ = 'boards'
    id = db.Column(UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
    user_id = db.Column(UUID(as_uuid=True), db.ForeignKey('users.id'), nullable=False)
    title = db.Column(db.String(100), nullable=False)
    tasks = relationship('Task', back_populates='board')

class Task(db.Model):
    __tablename__ = 'tasks'
    id = db.Column(UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
    board_id = db.Column(UUID(as_uuid=True), db.ForeignKey('boards.id'), nullable=False)
    title = db.Column(db.String(100), nullable=False)
    description = db.Column(db.String(250), nullable=True)
    status = db.Column(db.String(50), nullable=False)
    board = relationship('Board', back_populates='tasks')
