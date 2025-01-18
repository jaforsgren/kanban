# 🛠️ Multi-Backend Kanban App with Reverse Proxy 🚀

Experience the differences in how **Golang** and **Python** handle servers by building a Kanban application with a random reverse proxy routing to multiple backends. The frontend is kept simple with **HTML + jQuery** for dynamic interactivity.

---

## 🛠️ Setup

1. **Clone the Repository**:

   ```bash
   git clone https://github.com/your-repo/multi-backend-kanban.git
   cd multi-backend-kanbani

   ```

2. Start Services (Docker Compose FTW 🚢):

docker-compose up --build

3. Access the App:
   Open your browser and go to http://localhost:8080.

---

## 🔧 Backend Details

- **Golang Backend**:  
  Manages task and board creation with blazing speed ⚡.
- **Python Backend**:  
  The exact same thing in python.
- **Reverse Proxy**:  
  Randomly forwards requests to either backend for a fair comparison.

---

## 🖥️ Frontend

- Classic **HTML** + **jQuery** for simple yet effective UI.
- Dynamically loads boards and tasks based on routing.

---

## 🛠️ Database Setup

- Uses **PostgreSQL** for relational data storage.
- Migrations are auto-applied on container startup:
  - **Tasks Table**: Tracks tasks by `id`, `title`, `status`, and `board_id`.
  - **Boards Table**: Tracks boards by `id` and `title`.

---

## 📚 Learning Goals

- 🧵 **Explore Multithreading**: See how each language handles concurrent requests.
- 📊 **Database Agnosticism**: Observe how Golang and Python interact with PostgreSQL.
- 🧪 **Dynamic Frontend**: Use jQuery to make pages dynamic with minimal effort bc its boring.

---

## 🎯 Next Steps

1. actuall users
2. Add **metrics tracking** with a distributed counter service (e.g., Redis).
3. a php backend
4. A FE that is actually pleasant (stretch goal)
5. Integration with the unsplash api
6. Implement **real-time updates** using WebSockets.
7. Explore language-specific frameworks for optimizations:
   - **Golang**: Add gRPC support.
   - **Python**: Test with FastAPI for async performance.

```

```
