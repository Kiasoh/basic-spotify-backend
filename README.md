# Basic Spotify Backend

This project is a backend application developed in Go, simulating core functionalities of a music streaming service like Spotify. It provides user authentication, track listing and searching, playlist management, and user interaction tracking. A key feature is the ability to show a user's interaction state (liked, disliked, or neutral) with tracks, conditional on user authentication.

## Features

-   **User Authentication**: Register and log in users securely using JWT (JSON Web Tokens).
-   **Track Management**:
    -   Retrieve single track details by ID.
    -   List all available tracks with pagination and sorting options.
    -   Search tracks by various criteria.
    -   **Conditional Interaction State**: Tracks returned by listing, searching, or single-track retrieval APIs now include an `interaction_state` field (liked, disliked, neutral) if the request is made by an authenticated user. This state reflects the user's latest interaction with that specific track.
-   **Playlist Management**:
    -   Create and manage personal playlists.
    -   Add or remove tracks from playlists.
    -   View tracks within a specific playlist.
    -   Update playlist details (name, description).
-   **User Interaction Tracking**: Record user interactions with tracks (e.g., likes, dislikes, skips, plays, additions/removals from playlists).
-   **User Interest Modeling**: Implicitly models user preferences based on interactions to potentially drive future recommendation features (via `AvgInterest` in the user profile).
-   **Asynchronous Processing**: Uses Apache Kafka for asynchronous event processing (e.g., user interactions).

## Technologies Used

-   **Go**: The primary programming language.
-   **PostgreSQL**: Relational database for storing application data.
-   **Apache Kafka**: Distributed streaming platform for handling real-time data feeds and asynchronous operations.
-   **Chi**: A lightweight, idiomatic, and composable router for building HTTP services in Go.
-   **JWT (JSON Web Tokens)**: For secure user authentication.
-   **Docker & Docker Compose**: For containerization and orchestrating the application's services (Go backend, PostgreSQL, Kafka).

## Project Structure

The project follows a standard layered architecture to ensure separation of concerns and maintainability:

-   `main.go`: Application entry point, handles dependency injection and route initialization.
-   `handlers/`: Contains HTTP handler functions that process requests and return responses.
-   `services/`: Implements business logic and interacts with repositories.
-   `repository/`: Manages database interactions for different data models.
-   `models/`: Defines the data structures (structs) for the application.
-   `middleware/`: Contains HTTP middleware for concerns like authentication and request processing.
-   `init.sql`: SQL script for initializing the PostgreSQL database schema.
-   `Dockerfile`: Defines the Docker image for the Go backend application.
-   `docker-compose.yaml`: Orchestrates the multi-container application (backend, PostgreSQL, Kafka).

## Setup Instructions

### Prerequisites

-   [Docker](https://www.docker.com/get-started) and [Docker Compose](https://docs.docker.com/compose/install/) (recommended for easy setup)
-   [Go (1.21 or later)](https://golang.org/doc/install) (if running manually)

### Using Docker Compose (Recommended)

1.  **Clone the repository**:
    ```bash
    git clone https://github.com/kiasoh/basic-spotify-backend.git
    cd basic-spotify-backend
    ```
2.  **Start the services**:
    This command will build the Docker images (if not already built), start the PostgreSQL database, Kafka broker, and the Go backend application.
    ```bash
    docker-compose up --build
    ```
3.  **Initialize the database**:
    Once PostgreSQL is running, you need to initialize the database schema. The `init.sql` script creates the necessary tables. You can execute it manually, or if your `docker-compose.yaml` is set up to run it, it might be automatic. If manual:
    ```bash
    docker exec -it <postgres_container_id> psql -U niflheim -d dsdb -f /docker-entrypoint-initdb.d/init.sql
    # Replace <postgres_container_id> with the actual ID or name of your PostgreSQL container (e.g., 'project-postgres_ds-1')
    ```
    *(Note: Ensure `init.sql` has necessary `CREATE DATABASE` and `CREATE TABLE` statements if not handled by `docker-entrypoint-initdb.d` in your setup.)*

4.  **Access the application**:
    The backend API will be available at `http://localhost:8081`.

### Manual Setup (Without Docker)

1.  **Install PostgreSQL and Kafka**: Ensure you have a running PostgreSQL instance and Kafka broker accessible.
2.  **Configure Database Connection**: Update the DSN in `main.go` (or via environment variables if implemented) to point to your PostgreSQL instance.
3.  **Configure Kafka Connection**: Update the `kafkaURL` in `main.go` (or via environment variables) to point to your Kafka broker.
4.  **Install Go dependencies**:
    ```bash
    go mod tidy
    ```
5.  **Run the application**:
    ```bash
    go run main.go
    ```
    The application will run on `http://localhost:8081`.

## API Endpoints

All API endpoints are prefixed with `http://localhost:8081`.

### Authentication

-   `POST /register`: Register a new user.
-   `POST /login`: Authenticate a user and receive a JWT.

### Tracks

-   `GET /tracks/{trackID}`: Retrieve details for a single track.
    -   **Optional Authentication**: If a valid JWT is provided, the response includes `interaction_state`.
-   `GET /tracks`: List tracks with pagination, sorting, and optional authentication.
    -   **Query Parameters**: `limit`, `offset`, `sort_by`, `order`.
    -   **Optional Authentication**: If a valid JWT is provided, each track in the response includes `interaction_state`.
-   `GET /tracks/search`: Search tracks by query and field, with pagination and optional authentication.
    -   **Query Parameters**: `q` (query string), `field` (e.g., `track_name`, `artist`), `limit`, `offset`.
    -   **Optional Authentication**: If a valid JWT is provided, each track in the response includes `interaction_state`.

### Playlists

-   `GET /playlists/{playlistID}/tracks`: Retrieve tracks within a specific playlist.
    -   **Optional Authentication**: If a valid JWT is provided, each track in the response includes `interaction_state`.
-   `GET /playlists` (Protected): List all playlists owned by the authenticated user.
-   `POST /playlists` (Protected): Create a new playlist.
-   `PUT /playlists/{playlistID}` (Protected): Update details of an existing playlist.
-   `POST /playlists/{playlistID}/tracks/{trackID}` (Protected): Add a track to a playlist.
-   `DELETE /playlists/{playlistID}/tracks/{trackID}` (Protected): Remove a track from a playlist.

### User Interactions

-   `POST /tracks/{trackID}/interact` (Protected): Record a user interaction with a track.
    -   **Body**: `{ "type": "like" | "dislike" | "skip" | "play" | "add_to_playlist" | "remove_from_playlist" }`
-   `GET /tracks/{trackID}/interactions` (Protected): Get all interactions for a specific track (likely for administrative/debugging purposes).

## Authentication

This application uses JWTs for authentication.

-   Upon successful login (`POST /login`), a JWT is returned.
-   For **protected routes**, this JWT must be included in the `Authorization` header of subsequent requests in the format: `Authorization: Bearer <your_jwt_token>`.
-   For **optional authentication routes** (`/tracks`, `/tracks/search`, `/tracks/{trackID}`, `/playlists/{playlistID}/tracks`), providing a valid JWT will enrich the response with the user's `interaction_state`. If no JWT is provided or it's invalid, the request proceeds, but without the `interaction_state`.

## Contributing

Feel free to fork the repository and contribute!

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.