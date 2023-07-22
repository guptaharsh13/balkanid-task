[![Review Assignment Due Date](https://classroom.github.com/assets/deadline-readme-button-24ddc0f5d75046c5622901739e7c5dd533143b0c8e959d652212380cedb1ea36.svg)](https://classroom.github.com/a/YCCXVJKc)
[![Open in Visual Studio Code](https://classroom.github.com/assets/open-in-vscode-718a45dd9cf7e7f842a935f5ebbe5719a5e09af4491e668f4dbf3b35d5cca122.svg)](https://classroom.github.com/online_ide?assignment_repo_id=11471296&assignment_repo_type=AssignmentRepo)

# BalkanID Task

Backend Assignment | BalkanID

![ERD](erd.png)

## About the Backend

A robust containerized task management system that prioritizes security and user access management. Users can register and authenticate securely, while also having the option to deactivate or delete their accounts, adhering to data retention policies. It offers role-based and group-based access management for tasks, ensuring that sensitive information remains protected. It ensures data integrity and user safety. Additionally, it provides seamless support for bulk upload using CSV files, guaranteeing that relationships between users and tasks are accurately preserved.

## Key Features
- Authentication and Authorization
- Custom Role-based and Group-based access management
- Secure Bulk Uploads from CSV files

## Tech Stack
- Go
- Gin
- Gorm
- PostgreSQL
- Docker

## Getting Started

### Prerequisites

- Install [Go](https://go.dev/doc/install)
- Install [PostgreSQL](https://www.postgresql.org/download/)

### Setup

- Clone the repo.

    ```shell
    git clone https://github.com/BalkanID-University/balkanid-fte-hiring-task-vit-vellore-2023-guptaharsh13
    ```

- Change into the directory.

    ```shell
    cd balkanid-fte-hiring-task-vit-vellore-2023-guptaharsh13
    ```

- **Environment Variables**

  ```shell
  touch .env
  ```

  **For running this project successfully you'll need to create a `.env` file and store your PostgreSQL there like [`.env.sample`](https://github.com/BalkanID-University/balkanid-fte-hiring-task-vit-vellore-2023-guptaharsh13/blob/master/.env.sample).**

## Usage

### Local Run

```shell
go mod tidy
go run cmd/*.go 
```

### Docker Run

```shell
docker-compose up --build -d
```

### Create Admin User
```shell
go run cmd/*.go create admin <username> <email> <password>
```


Now, may access the app on http://localhost:PORT

PORT: You may have optionally specified in the .env file

## Contact

Harsh Gupta - hg242322@gmail.com

GitHub Link: [https://github.com/guptaharsh13](https://github.com/guptaharsh13)

<p align="center">Made with ❤ by Harsh Gupta</p>
