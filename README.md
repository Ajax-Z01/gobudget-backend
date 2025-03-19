### 📌 **GoBudget-Backend**  

GoBudget-Backend is a **RESTful API** built with **Golang** and **PostgreSQL** to power **GoBudget**, a personal finance management application. This backend handles user authentication, expense tracking, budget management, and financial reports.  

## 🚀 **Features**  
✔️ **User Authentication** – Secure login and registration with JWT  
✔️ **Expense Tracking** – CRUD operations for transactions  
✔️ **Budget Management** – Set spending limits per category  
✔️ **Financial Reports** – Generate summaries and analytics  
✔️ **Localization Ready** – API supports multiple languages  
✔️ **Soft Delete** – Uses `deleted_at` for reversible deletion  

## 🏗️ **Tech Stack**  
- **Golang (Gin Framework)** – Fast and lightweight web framework  
- **PostgreSQL** – Reliable relational database  
- **GORM** – ORM for seamless database interactions  
- **JWT Authentication** – Token-based user authentication  
- **Swagger** – API documentation for easy integration  

## 📂 **Project Structure**  
```
/gobudget-backend
│── main.go          # Entry point of the app
│── router.go        # Defines API routes
│── controllers/     # Handles request logic
│── models/          # Database schema & ORM models
│── middleware/      # JWT auth & request validation
│── database.go      # PostgreSQL connection setup
│── config/          # App configuration settings
│── seeder.go        # Initial database seed data
│── api_test.go      # API testing
```

## 🛠️ **Setup & Installation**  

1️⃣ **Clone the repository**  
```bash
git clone https://github.com/Ajax-Z01/gobudget-backend.git
cd gobudget-backend
```

2️⃣ **Set up the environment**  
Create a `.env` file and configure database credentials:  
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_user
DB_PASSWORD=your_password
DB_NAME=gobudget
JWT_SECRET=your_secret_key
```

3️⃣ **Install dependencies**  
```bash
go mod tidy
```

4️⃣ **Run the application**  
```bash
go run main.go
```

5️⃣ **API Documentation** (Swagger)  
After running the server, access API docs at:  
```
http://localhost:8080/swagger/index.html
```

## 📌 **Contributing**  
Feel free to fork this repo and submit pull requests! Any contributions to improve the project are welcome.  

🔗 **Frontend Repo**: [GoBudget-Frontend](https://github.com/Ajax-Z01/gobudget-frontend)  

---

Let me know if you need any changes! 🚀
