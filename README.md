---

# 📘 BACKEND README (Go + Gin)

```md
# 🚀 Sootika Backend

Backend of the Sootika e-commerce application built using **Go (Gin Framework)**.

---

## 🚀 Tech Stack

- Go (Golang)
- Gin Framework
- JWT Authentication
- PostgreSQL / MySQL
- Cloudinary (Image Upload)
- Razorpay (Payments)

---

## 📁 Project Structure


src/
├── controllers/ # Request handlers
├── middleware/ # Auth & Admin middleware
├── repository/ # Database layer
├── routes/ # API routes
├── services/ # Business logic
└── main.go


---

## ⚙️ Setup Instructions

### 1. Clone repo

```bash
git clone <your-repo-url>
cd backend
2. Install dependencies
go mod tidy
3. Run server
go run main.go

Server runs on:

http://localhost:8080

🌐 Environment Variables

Create .env:

PORT=8080
JWT_SECRET=xxxxxx

DB_URL=xxxxxx

RAZORPAY_KEY_ID=xxxxxxxx
RAZORPAY_SECRET=xxxxxxxx

CLOUDINARY_URL=xxxxxxxx

🔐 Authentication
JWT-based auth
Access + Refresh tokens
Middleware:
AuthMiddleware
AdminMiddleware

📡 API Endpoints

Auth
POST /auth/signup
POST /auth/login
POST /auth/refresh
POST /auth/logout

User
GET /user/dashboard
Products
GET /products/
GET /products/:id

Cart
GET /cart/
POST /cart/add
DELETE /cart/remove/:id

Wishlist
GET /wishlist/
POST /wishlist/add

Admin
GET /admin/dashboard
POST /admin/products
PUT /admin/products/:id
DELETE /admin/products/:id

💳 Razorpay Integration
Order created from backend
Payment verified using signature

Test Mode supported

🛡️ Security
Authorization header required:
Authorization: Bearer <token>
Tokens validated in middleware

⚠️ Notes
Never expose Razorpay secret key
Always validate payment signature
Ensure database is connected before running