# HTMX Todo App

## Running

Ensure you have Go installed. If you don't, you should be able to `brew install go`.

Install the dependencies with `go mod download`.

Once Go is installed, you can use `go run .` to run the app. You can then navigate to `http://localhost:8080` to see the app running.

Note that as the app is running locally and using an in-memory array for todos, latency is basically non-existent. However, HTMX has the ability to swap loading indicators in/out of the DOM when requests are in-flight. If you apply arbitrary latency through your browsers dev tools, you should be able to see it in action.
