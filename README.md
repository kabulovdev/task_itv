# Task ITV Project


## Running the Project

1. **Build the Docker images:**
    ```bash
    docker compose build
    ```

2. **Start the services:**
    ```bash
    docker compose up
    ```

3. **Access Swagger UI:**

    Open your browser and go to:  
    [http://localhost:8080/swagger](http://localhost:8080/swagger)

---

Feel free to check the Swagger documentation for API details.

## Data Model Structure

The project uses two main data tables: **User** and **Movie**.

### User Table

- **ID**: Primary key, auto-incremented.
- **Username**: Unique, required.
- **Email**: Unique, required, must be a valid email format.
- **Password**: Required.
- **Role**: String with a default value of 'user', required.

### Movie Table

- **ID**: Primary key, auto-incremented.
- **Title**: Required.
- **Director**: Required.
- **Year**: Required, integer.
- **Plot**: Optional description of the movie.

These structures are managed using GORM, ensuring data integrity and validation at the database level.