# 🛍️ E-commerce Backend with Go and Fiber

This repository contains the source code for an **e-commerce backend**, developed with the goal of **learning and practicing the Go language** using the **Fiber** framework. The project includes essential functionalities for operating an e-commerce platform, following best practices for code organization, design principles (such as SOLID), and clean architecture.

---

## 🚀 Implemented Features

- **Authentication and Authorization**
  - User registration and login.
  - JWT tokens for secure authentication.
- **Product Management**
  - Creation, editing, listing, and deletion of products.
- **Orders**
  - Order creation from the cart.
  - Order history per user.
- **Categories**
  - Creation of categories.
  - Category creation history per user.
- **Users**
  - User creation.
  - Users can also be sellers.
- **Administration**
  - CRUD for products and categories.
  - User management.

---

## 🛠️ Technologies Used

- **Language**: Go
- **Web Framework**: Fiber
- **Database**: PostgreSQL
- **ORM**: GORM
- **Authentication**: JWT
- **Configuration Management**: Viper
- **Data Validation**: Validator
- **Logging**: Logrus
- **Testing**: Go Testing (with simple examples for learning)

---

## 📂 Project Structure

```plaintext
.
├── config/             # Initial application configuration
├── infra/              # Configuration files
├── internal/           # Domain logic
│   ├── handlers/       # Route handlers
│   ├── services/       # Business logic
│   ├── repositories/   # Database access
│   ├── dto/            # Data formatting files
│   ├── helper/         # Helper functions
│   ├── api/            # Route configuration
│   ├── domain/         # Domain objects
├── pkg/                # Reusable packages
│   ├── notification/   # Notification logic     
└── main.go             # Application entry point
```

---
## 🌱 Project Goals
#### This project was developed to:

1. **Learn Go**: Explore the syntax and fundamental concepts of the language.
2. **Practice using Fiber**: Understand how to create fast and efficient web applications.
3. **Implement Best Practices**: Apply design principles like SOLID and modular organization.
4. **Simulate a Real-World Scenario**: Work on a backend with common functionalities found in real applications.

---

## 🧰 How to Run the Project
### Prerequisites
- Go (version 1.19 or higher)
- PostgreSQL
- Go Fiber
- A tool for managing environment variables (e.g., dotenv)

---
## 📚 Swagger Documentation

### Generate the Documentation:
Run the following command at the root of your project (where go.mod is located) to generate the Swagger documentation files in the docs/ directory:
```Bash
  swag init  
```

Now, start your server and navigate to http://localhost:3000/swagger/index.html to view the Swagger UI.

---

## 📝 Next Steps

- Add unit and integration tests.
- Implement image upload for products.
- Integrate a simulated payment system.
- Improve API documentation using Swagger.

---

## 🤝 Contributions
If you have suggestions or improvements, feel free to open an issue or submit a pull request! 😊
