# Task Queue System

A simple task queue system built with Go that processes background jobs using worker pools and channels.

## What does this do?

Imagine you have a web application and users can do things that take time - like sending emails, processing uploaded images, or generating reports. You don't want users to wait for these tasks to finish. Instead, you queue them up and process them in the background.

That's what this project does. You submit tasks via an API, they go into a queue, and workers pick them up and process them one by one.

## Main Features

- Submit tasks through REST API
- 5 workers process tasks concurrently
- Tasks can be: sending emails, processing images, or generating reports (all simulated)
- If a task fails, it automatically retries (up to 3 times)
- Check task status anytime
- See system statistics (how many tasks completed, failed, etc.)

## How it works

```
User → API → Queue → Workers → Process Task → Update Status
```

1. You submit a task through the API
2. Task gets added to a queue (like a waiting line)
3. One of the 5 workers picks it up
4. Worker processes the task (sends email, processes image, etc.)
5. Worker updates the task status (completed or failed)
6. You can check the status anytime

## Project Structure

The project follows Clean Architecture, which basically means:
- **domain/** - Core business logic (what a Task is, what statuses it can have)
- **usecase/** - Application logic (how to submit a task, get a task, etc.)
- **infrastructure/** - Implementation details (where tasks are stored, how workers work)
- **delivery/http/** - REST API (how you interact with the system)
- **cmd/server/** - Main application that starts everything

Think of it like layers of an onion - the core business logic doesn't know about HTTP or databases, making it easy to test and change.

## Running the project

1. Make sure you have Go installed (1.21 or higher)

2. Download dependencies:
   ```bash
   go mod download
   ```

3. Start the server:
   ```bash
   go run cmd/server/main.go
   ```

4. The server starts on `http://localhost:8080`

You'll see logs like:
```
🚀 Starting Task Queue System...
✅ Repository initialized
✅ Queue initialized (capacity: 100)
✅ Worker pool started (5 workers)
🌐 Server starting on http://localhost:8080
✨ Ready to accept requests!
```


## Notes

- Tasks are stored in memory, so they're lost when you restart the server
- The actual task processing is simulated (just sleeps and logs)
- In a real system, you'd use Redis or a database for the queue
- You could add authentication, rate limiting, etc.

