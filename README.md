# Portone

### This project is a payment processing API built using Go (Golang) and Stripe. It provides a set of RESTful endpoints for handling various aspects of payment processing, such as creating payment intents, confirming payment intents, capturing payments, creating refunds, and retrieving a list of payment intents.

## Table of Contents

-  [Getting Started](#getting-started)
   -  [Prerequisites](#prerequisites)
   -  [Installation](#installation)
   -  [Configuration](#configuration)
-  [Usage](#usage)
-  [Endpoints](#endpoints)
-  [Testing](#testing)
-  [Contributing](#contributing)
-  [License](#license)

## Getting Started

Welcome to Portone! This section will guide you through the process of setting up and using the project.

### Prerequisites

Before you begin, ensure you have the following prerequisites installed:

-  [Go](https://golang.org/doc/install)
-  [Docker](https://www.docker.com/get-started)
-  [Stripe Account](https://dashboard.stripe.com/register) - Create a sandbox account to obtain API keys.

### Installation

Follow these steps to install and run the project:

1. Clone the repository:

   ```bash
   git clone https://github.com/AshiqurRahaman02/portone.git
   ```

2. Navigate to the project directory:

   ```bash
   cd portone
   ```

3. Set up your environment variables by creating a `.env` file in the root directory. Example content:

   ```env
   STRIPE_API_KEY=your_stripe_api_key
   ```

4. Build and run the Docker containers:

   ```bash
   docker-compose up -d
   ```

### Configuration

Before running the application, make sure to configure the necessary environment variables in the `.env` file.

## Usage

To use the project, follow these steps:

1. Create a payment intent:

   ```bash
   curl -X POST -d '{"amount": 2000, "currency": "usd"}' http://localhost:8080/api/v1/create_intent
   ```

2. Confirm a payment intent:

   ```bash
   curl -X POST http://localhost:8080/api/v1/confirm_payment_intent/{payment_intent_id}
   ```

3. Capture the created intent:

   ```bash
   curl -X POST http://localhost:8080/api/v1/capture_intent/{payment_intent_id}
   ```

4. Create a refund for the created intent:

   ```bash
   curl -X POST http://localhost:8080/api/v1/create_refund/{payment_intent_id}
   ```

5. Get a list of all payment intents:

   ```bash
   curl http://localhost:8080/api/v1/get_intents
   ```

For more details on each endpoint and additional options, refer to the [Endpoints](#endpoints) section.

## Endpoints

Document the available API endpoints and their functionality.

-  **Create Intent**
   -  **Endpoint:** `POST /api/v1/create_intent`
   -  **Description:** Creates a payment intent.
   -  **Request:**
      -  **Example:**
         ```json
         {
         	"amount": 2000,
         	"currency": "usd"
         }
         ```
      -  **Parameters:**
         -  `amount` (required): The amount of the payment in the smallest currency unit (e.g., cents).
         -  `currency` (required): The currency in which the payment is made (e.g., "usd").
   -  **Response:**
      -  **Example:**
         ```json
         {
         	"success": true,
         	"payment_intent_id": "pi_3MtwBwLkdIwHu7ix28a3tqPa",
         	"payment": {
         		"id": "pi_3MtwBwLkdIwHu7ix28a3tqPa",
         		"object": "payment_intent",
         		"amount": 2000,
         		"currency": "usd",
         		"status": "requires_confirmation",
         		"created": 1680800504,
         		"client_secret": "pi_3MtwBwLkdIwHu7ix28a3tqPa_secret_YrKJUKribcBjcG8HVhfZluoGH"
         		// ... other payment intent properties
         	}
         }
         ```
      -  **Response Properties:**
         -  `success`: A boolean indicating the success of the operation.
         -  `payment_intent_id`: The ID of the created payment intent.
         -  `payment`: Details of the created payment intent.
   -  **Error Responses:**
      -  `400 Bad Request`: If the request is malformed or missing required parameters.
      -  `500 Internal Server Error`: If there is an internal error during payment intent creation.
   -  **Usage:**
      ```bash
      curl -X POST -H "Content-Type: application/json" -d '{"amount": 2000, "currency": "usd"}' http://localhost:8080/api/v1/create_intent
      ```
   -  **Notes:**
      -  Ensure that the `STRIPE_API_KEY` environment variable is set before making the request.

---

-  **Confirm Intent**
   -  **Endpoint:** `POST /api/v1/confirm_payment_intent/{id}`
   -  **Description:** Confirm a payment intent.
   -  **Request:**
      -  **Example:**
         ```json
         // No request body required
         ```
      -  **Path Parameters:**
         -  `id` (required): The ID of the payment intent to confirm.
   -  **Response:**
      -  **Example:**
         ```json
         {
         	"success": true,
         	"payment_intent_id": "pi_3MtwBwLkdIwHu7ix28a3tqPa",
         	"status": "requires_action",
         	"client_secret": "pi_3MtwBwLkdIwHu7ix28a3tqPa_secret_YrKJUKribcBjcG8HVhfZluoGH"
         }
         ```
      -  **Response Properties:**
         -  `success`: A boolean indicating the success of the operation.
         -  `payment_intent_id`: The ID of the confirmed payment intent.
         -  `status`: The updated status of the payment intent (e.g., "requires_capture").
         -  `client_secret`: The updated client secret of the payment intent.
   -  **Error Responses:**
      -  `400 Bad Request`: If the request is malformed or missing required parameters.
      -  `500 Internal Server Error`: If there is an internal error during payment intent confirmation.
   -  **Usage:**
      ```bash
      curl -X POST http://localhost:8080/api/v1/confirm_payment_intent/pi_3MtwBwLkdIwHu7ix28a3tqPa
      ```
   -  **Notes:**
      -  Ensure that the `STRIPE_API_KEY` environment variable is set before making the request.

---

-  **Capture Intent**
   -  **Endpoint:** `POST /api/v1/capture_intent/{id}`
   -  **Description:** Captures a payment intent.
   -  **Request:**
      -  **Example:**
         ```json
         // No request body required
         ```
      -  **Path Parameters:**
         -  `id` (required): The ID of the payment intent to capture.
   -  **Response:**
      -  **Example:**
         ```json
         {
         	"success": true,
         	"payment_intent_id": "pi_3MtwBwLkdIwHu7ix28a3tqPa",
         	"status": "succeeded"
         }
         ```
      -  **Response Properties:**
         -  `success`: A boolean indicating the success of the operation.
         -  `payment_intent_id`: The ID of the captured payment intent.
         -  `status`: The updated status of the payment intent after capture (e.g., "succeeded").
   -  **Error Responses:**
      -  `400 Bad Request`: If the request is malformed or missing required parameters.
      -  `500 Internal Server Error`: If there is an internal error during payment intent capture.
   -  **Usage:**
      ```bash
      curl -X POST http://localhost:8080/api/v1/capture_intent/pi_3MtwBwLkdIwHu7ix28a3tqPa
      ```
   -  **Notes:**
      -  Ensure that the `STRIPE_API_KEY` environment variable is set before making the request.

---

-  **Create Refund**
   -  **Endpoint:** `POST /api/v1/create_refund/{id}`
   -  **Description:** Creates a refund for a payment intent.
   -  **Request:**
      -  **Example:**
         ```json
         // No request body required
         ```
      -  **Path Parameters:**
         -  `id` (required): The ID of the payment intent for which to create a refund.
   -  **Response:**
      -  **Example:**
         ```json
         {
         	"success": true,
         	"refund_id": "re_1234567890abcdef",
         	"status": "refunded",
         	"amount_refunded": 2000
         }
         ```
      -  **Response Properties:**
         -  `success`: A boolean indicating the success of the refund creation.
         -  `refund_id`: The ID of the created refund.
         -  `status`: The status of the refund (e.g., "succeeded").
         -  `amount_refunded`: The amount refunded in cents.
   -  **Error Responses:**
      -  `400 Bad Request`: If the request is malformed or missing required parameters.
      -  `500 Internal Server Error`: If there is an internal error during refund creation.
   -  **Usage:**
      ```bash
      curl -X POST http://localhost:8080/api/v1/create_refund/pi_3MtwBwLkdIwHu7ix28a3tqPa
      ```
   -  **Notes:**
      -  Ensure that the `STRIPE_API_KEY` environment variable is set before making the request.

---

-  **Get Intents**
   -  **Endpoint:** `GET /api/v1/get_intents`
   -  **Description:** Gets a list of all payment intents.
   -  **Request:**
      -  **Example:**
         ```json
         // No request body required
         ```
   -  **Response:**
      -  **Example:**
         ```json
         {
         	"success": true,
         	"paymentIntents": [
         		{
         			"id": "pi_1234567890abcdef",
         			"amount": 2000,
         			"currency": "usd",
         			"status": "requires_confirmation",
         			"created": 1680800504,
         			"clientSecret": "pi_1234567890abcdef_secret_YrKJUKribcBjcG8HVhfZluoGH"
         		}
         		// ... (more payment intents)
         	]
         }
         ```
      -  **Response Properties:**
         -  `success`: A boolean indicating the success of fetching payment intents.
         -  `paymentIntents`: An array of payment intents, each containing relevant information.
            -  `id`: The ID of the payment intent.
            -  `amount`: The amount in cents.
            -  `currency`: The currency code (e.g., "usd").
            -  `status`: The current status of the payment intent.
            -  `created`: The timestamp of when the payment intent was created.
            -  `clientSecret`: The client secret associated with the payment intent.
   -  **Error Responses:**
      -  `500 Internal Server Error`: If there is an internal error during payment intent retrieval.
   -  **Usage:**
      ```bash
      curl http://localhost:8080/api/v1/get_intents
      ```
   -  **Notes:**
      -  Ensure that the `STRIPE_API_KEY` environment variable is set before making the request.

---

## Testing

To run the tests for this project, follow these steps:

1. Make sure you have Go installed on your machine.

2. Navigate to the project's root directory in your terminal.

3. Run the following command to execute the tests:

   ```bash
   go test ./tests
   ```

   This command will run all the tests in the project.

4. Review the test results to ensure that all tests pass successfully.

Note: Ensure that you have set up the necessary environment variables, such as `STRIPE_API_KEY`, before running the tests. You can use a `.env` file to manage your environment variables.

## Contributing

Thank you for considering contributing to our project! Whether you're reporting a bug, proposing a feature, or submitting code changes, your contributions are highly appreciated.

## Issues

If you find a bug, have a question, or want to propose a new feature, check our issue tracker for existing topics. If not found, feel free to open a new issue and provide details such as a clear title, steps to reproduce, and your environment.

## Feature Requests

Have a feature in mind? We welcome new ideas and enhancements. Open an issue on our GitHub repository to discuss and share your thoughts with the community.

## Pull Requests

Contributions through pull requests are welcome. To contribute:

1. Fork the repository.

2. Create a new branch for your changes: git checkout -b feature/your-feature.

3. Make changes following our coding standards.

4. Write tests and run existing tests.

5. Push changes to your fork: git push origin feature/your-feature.

6. Open a pull request on GitHub with a clear description of your changes.

## Coding Standards

-  **Indentation and Formatting:**

   1. Use tabs for indentation.
   2. Follow the standard Go formatting guidelines. You can use the gofmt tool to automatically format your code.
   3. Variable Naming:

-  **Variable Naming**

   1. Use meaningful and descriptive names for variables.
   2. Follow the camelCase naming convention for variables.

-  **Function Naming:**

   1. Use camelCase for function names.
   2. Choose function names that clearly indicate their purpose.

-  **Comments:**

   1. Include comments to explain complex sections of code or to provide context.
   2. Write clear and concise comments.

-  **Error Handling:**

   1. Properly handle errors using the if err != nil pattern.
   2. Avoid generic error messages; provide specific details when handling errors.

-  **Testing:**

   1. Write comprehensive unit tests for your code.
   2. Ensure that tests cover different scenarios and edge cases.

-  **Documentation:**

   1. Provide documentation for public functions and packages.
   2. Use GoDoc-style comments for documenting functions and packages.

-  **Imports:**

   1. Group imports into standard library packages, third-party packages, and local packages.
   2. Avoid unused imports.

-  **Concurrency and Goroutines:**

   1. Use goroutines and channels responsibly.
   2. Ensure proper synchronization to avoid race conditions.

-  **Code Modularity:**

   1. Encapsulate functionality into modular functions and packages.
   2. Aim for a clear separation of concerns.

-  **Security:**

   1. Follow security best practices, especially when dealing with user input.
   2. Be mindful of potential vulnerabilities and address them promptly.

-  **Version Control:**

   1. Make small, meaningful commits with clear commit messages.
   2. Avoid committing large binary files or sensitive information.

## Getting Help

For questions or assistance, open an issue or join community discussions.

##

```
Thank you for contributing! Feel free to customize it based on your project's specifics.
```
