# Web Development Assignment: README

## Overview
This assignment focuses on implementing several key web development features, including filtering, sorting, pagination, structured logging, error handling, rate limiting, graceful shutdown, sending emails, and a user profile page. These functionalities improve application usability, maintainability, and performance while ensuring robust error handling and user engagement.

---

## Table of Contents
1. [Technologies Used](#technologies-used)
2. [Features](#features)
   - Filtering, Sorting, and Pagination
   - Structured Logging, Error Handling, Rate Limiting, and Graceful Shutdown
   - Sending Emails
   - User Profile Page
3. [Installation and Setup](#installation-and-setup)
4. [Usage Instructions](#usage-instructions)
5. [Testing](#testing)
6. [Future Enhancements](#future-enhancements)
7. [Contact](#contact)

---

## Technologies Used
- Backend: Node.js with Express
- Database: PostgreSQL
- Frontend: HTML, CSS, JavaScript
- Email Service: Nodemailer (or equivalent)
- Logging: Winston (for structured logging)

---

## Features

### 1. Filtering, Sorting, and Pagination
#### Description:
- **Filtering:** Users can filter data (e.g., user database, products) based on criteria like name, age, or email.
- **Sorting:** Data can be ordered in ascending or descending order by fields such as registration date or username.
- **Pagination:** Data is displayed across multiple pages for better performance and user experience.

#### Implementation:
- **Backend:**
  - Filtering: Query parameters are used to fetch filtered data.
  - Sorting: SQL `ORDER BY` is implemented dynamically based on user input.
  - Pagination: SQL `LIMIT` and `OFFSET` are used to divide data.
- **Frontend:**
  - Dynamic forms allow users to input filtering and sorting criteria.
  - Pagination controls (e.g., previous, next, page numbers) are provided.

---

### 2. Structured Logging, Error Handling, Rate Limiting, and Graceful Shutdown
#### Description:
- **Structured Logging:** Standardized JSON logging for easier analysis.
- **Error Handling:** Comprehensive error detection and response system.
- **Rate Limiting:** Restricts requests per user or IP within a specific timeframe.
- **Graceful Shutdown:** Ensures safe termination of the application.

#### Implementation:
- **Structured Logging:** Winston library for JSON-formatted logs.
- **Error Handling:** Centralized middleware to handle and log errors.
- **Rate Limiting:** Middleware using `express-rate-limit`.
- **Graceful Shutdown:** Signals (e.g., SIGINT) handle shutdown tasks like database disconnection.

---

### 3. Sending Emails
#### Description:
Users can send support messages to the helpdesk with text, attachments, and images.

#### Implementation:
- **Backend:**
  - Nodemailer configured to send emails via SMTP.
  - Attachment handling through file uploads.
- **Frontend:**
  - Form to compose messages and upload attachments.

#### Key Features:
- Plain text and rich text messages.
- File validation for attachments.

---

### 4. User Profile Page
#### Description:
A dedicated page for users to manage their personal information and view historical data.

#### Features:
- **Personal Information Management:**
  - Update password and email address.
- **Order Management:**
  - View order history and interaction logs.
- **Support Integration:**
  - Send messages to the support team.

#### Implementation:
- **Frontend:**
  - Profile dashboard with interactive forms.
- **Backend:**
  - API endpoints for CRUD operations on user data.

---

## Installation and Setup

### Prerequisites:
- Node.js and npm installed.
- PostgreSQL installed and running.

### Steps:
1. Clone the repository:
   ```bash
   git clone []
   ```
2. Install dependencies:
   ```bash
   npm install
   ```
3. Set up the PostgreSQL database:
   - Create a database and configure connection settings in `.env`.
   - Run migrations:
     ```bash
     npx sequelize-cli db:migrate
     ```
4. Start the server:
   ```bash
   npm start
   ```

---

## Usage Instructions
1. Open the application in your browser:
   ```
   http://localhost:3000
   ```
2. Navigate to:
   - Filtering, Sorting, and Pagination: `/data` page.
   - User Profile: `/profile` page.
   - Support: `/support` page.
3. Test error handling by entering invalid inputs.
4. Test rate limiting by sending multiple rapid requests.

---

## Testing
- Unit tests: Written using Mocha and Chai.
- Testing structured logging and error handling.
- Mock testing email functionality with a local SMTP server.

Run tests:
```bash
npm test
```

---

## Future Enhancements
- Implement advanced filtering with multiple criteria.
- Add front-end validations for file uploads.
- Extend support functionality for multi-agent ticketing systems.
- Enhance logging with external log management services.
![image](https://github.com/user-attachments/assets/6afdd3d8-44aa-4215-92df-075d28c9d320)




