from flask import Blueprint, request, jsonify
from .models import db, User, Board, Task
import uuid
import json

bp = Blueprint('api', __name__)

@bp.route('/health', methods=['GET'])
def health_check():
    return 'Flask API is running!', 200

@bp.route('/users', methods=['POST'])
def create_user():
    data = request.get_json()
    user = User(id=uuid.uuid4(), username=data['username'], email=data['email'])
    db.session.add(user)
    db.session.commit()
    return jsonify({'id': str(user.id), 'username': user.username, 'email': user.email}), 201

@bp.route('/boards', methods=['POST'])
def create_board():
    data = request.get_json()
    board = Board(id=uuid.uuid4(), user_id=data['user_id'], title=data['title'])
    db.session.add(board)
    db.session.commit()
    return json.dumps({'id': str(board.id), 'title': board.title}), 201

@bp.route('/boards', methods=['GET'])
def get_all_boards():
    boards = Board.query.all()
    return jsonify([{'id': str(board.id), 'title': board.title} for board in boards]), 200

@bp.route('/boards/<uuid:board_id>', methods=['DELETE'])
def delete_board(board_id):
    board = Board.query.get_or_404(board_id)
    db.session.delete(board)
    db.session.commit()
    return jsonify({'message': 'Board deleted successfully'}), 200

@bp.route('/boards/<uuid:board_id>', methods=['GET'])
def get_board(board_id):
    tasks = Task.query.filter(Task.board_id == board_id).all()
    tasks_list = [
        {'id': task.id, 'title': task.title, 'board_id': str(task.board_id), 'description': task.description, 'status': task.status} for task in tasks
    ]

    return jsonify(tasks_list), 200

@bp.route('/tasks', methods=['PATCH'])
def update_task():
    data = request.get_json()

    data = request.get_json()
    if not data:
        return jsonify({'error': 'Invalid or missing JSON payload'}), 400

    board_id = data.get('board_id')
    description = data.get('description')
    status = data.get('status')
    title = data.get('title')

    task = Task.query.get_or_404(data.get('id'))

    if board_id:
        task.board_id = board_id
    if description:
        task.description = description
    if status:
        task.status = status
    if title:
        task.title = title

    db.session.commit()

    updated_task = {
        'id': task.id,
        'board_id': str(task.board_id),
        'description': task.description,
        'status': task.status,
        'title': task.title,
    }

    return jsonify(updated_task), 200

@bp.route('/tasks', methods=['POST'])
def create_task():
    data = request.get_json()
    task = Task(id=uuid.uuid4(), board_id=data['board_id'], title=data['title'], status=data['status'], description=data['description'])
    db.session.add(task)
    db.session.commit()
    
    updated_task = {
        'id': task.id,
        'board_id': str(task.board_id),
        'description': task.description,
        'status': task.status,
        'title': task.title,
    }

    return jsonify(updated_task), 200

@bp.route('/tasks/<uuid:task_id>', methods=['DELETE'])
def delete_task(task_id):
    task = Task.query.get_or_404(task_id)
    db.session.delete(task)
    db.session.commit()
    return jsonify({'message': 'Task deleted successfully'}), 200