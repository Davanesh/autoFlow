# ⚙️ AutoFlow.AI

AutoFlow.AI is a **smart workflow orchestration platform** designed to automate, optimize, and visualize backend processes — powered by **Go microservices**, **AWS Cloud**, and **AI-driven automation**.

---

## 🚀 Features

- 🧩 **Drag-and-Drop Workflow Builder**
  - Build and connect tasks visually on an interactive canvas (React-based).
- ⚙️ **Go Microservices**
  - Backend engine written in Go for performance, scalability, and clean concurrency.
- ☁️ **AWS Integration**
  - Simulate and deploy workflows using AWS Lambda, ECS, and Step Functions.
- 🧠 **AI Optimization**
  - Intelligent suggestions for workflow efficiency and resource optimization.
- 🔐 **Secure Backend**
  - JWT-based authentication and role-based access management.
- 📊 **Real-Time Logs**
  - Monitor workflow executions and view live logs with MongoDB and WebSocket updates.

---

## 🏗️ Tech Stack

### **Frontend**
- React.js + Tailwind CSS  
- Redux Toolkit for State Management  
- Canvas-based workflow builder  
- Axios for API communication  

### **Backend**
- Go (Golang)  
- Gin / Fiber Framework  
- MongoDB Atlas  
- AWS SDK for Go  
- JSON Web Tokens (JWT)  
- REST API Architecture  

### **Cloud & DevOps**
- AWS Lambda, ECS, Step Functions  
- Terraform for IaC  
- Docker for containerization  
- CloudWatch for monitoring  

---

## 🧩 Architecture Overview

```text
Frontend (React + Redux)
        ↓
Gateway API (Go)
        ↓
Workflow Engine (Go Microservice)
        ↓
Task Executors (Lambda / Local Simulated)
        ↓
Database (MongoDB Atlas)
```

---

## 🧠 AI Automation Concept

The system analyzes workflows and:
- Suggests **optimized task ordering**
- Automates **retry logic and scaling**
- Can be extended to handle **auto-email or message responses** based on triggers

---

## 🧪 Local Development Setup

### 1️⃣ Clone the Repository
```bash
git clone https://github.com/Davanesh/autoflow.git
cd autoflow
```

### 2️⃣ Backend Setup
```bash
cd backend
cd orchestrator
go run main.go
```

### 3️⃣ Frontend Setup
```bash
cd frontend
npm install
npm run dev
```

### 4️⃣ Environment Variables
Create `.env` files for both backend and frontend.

**Backend .env**
```
PORT=8080
MONGO_URI=your_mongo_atlas_uri
JWT_SECRET=your_secret_key
AWS_REGION=ap-south-1
```

**Frontend .env**
```
VITE_API_URL=http://localhost:8080
```

---

## 🛠️ Roadmap

| Phase | Goal | Status |
|-------|------|--------|
| 1 | Backend core (Go microservices + MongoDB) | ✅ Done |
| 2 | Frontend Canvas Builder (React) | ✅ Done |
| 3 | Lambda / Step Functions Simulation | 🔄 In progress |
| 4 | AI Workflow Optimizer | ⏳ Planned |
| 5 | AWS Deployment + Terraform Setup | ⏳ Planned |

---

## 🤝 Contributing
Pull requests are welcome! For major changes, please open an issue first to discuss what you’d like to change.

---

## 📜 License
This project is licensed under the **MIT License**.

---

## 👨‍💻 Author
**Davanesh S**  
🚀 Full Stack Developer | Cloud & AI Enthusiast  
🌐 [Portfolio](https://davanesh.vercel.app/)  
💼 [LinkedIn](https://www.linkedin.com/in/davanesh-saminathan/)  
