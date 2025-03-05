# Creating a Simple Webserver in Go

![Webserver Illustration](/static/img/go_mascot.png)

This tutorial will guide you through building a simple web server using Go. You will learn how to create a basic HTTP server, set up request handlers, and run your server locally.

## Prerequisites

- **Go Installed:** Ensure you have [Go installed](https://golang.org/dl/).
- **Basic Terminal Knowledge:** Familiarity with the terminal/command prompt.
- **Basic Go Knowledge:** Understanding of Go syntax and packages.

## Step 1: Setting Up Your Project

1. Create a new directory for your project:

   ```bash
   mkdir simple-webserver
   cd simple-webserver
   ```

2. Initialize a new Go module:

   ```bash
   go mod init example.com/simple-webserver
   ```

## Step 2: Writing the Web Server Code

Create a file named `main.go` and add the following code:

```go
package main

import (
	"fmt"
	"log"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

func main() {
	http.HandleFunc("/", helloHandler)
	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

### Code Explanation

- **Imports:**
    - `fmt` for formatted I/O.
    - `log` for logging errors.
    - `net/http` for HTTP functionality.
- **Handler Function:**
    - `helloHandler` responds with "Hello, World!" when accessed at `/`.
- **Main Function:**
    - Maps `/` to `helloHandler`.
    - Starts the server on port `8080`.

## Step 3: Running the Web Server

1. Open your terminal in the project directory.
2. Run:

   ```bash
   go run main.go
   ```

3. Open [http://localhost:8080](http://localhost:8080) in a browser to see:

   ```
   Hello, World!
   ```

## Additional Tips

- **Change the Port:** Modify `http.ListenAndServe(":8080", nil)` to use a different port.
- **Handle More Routes:** Add more handlers using `http.HandleFunc("/route", yourHandler)`.
- **Error Handling:** Implement better logging for robustness.

## Conclusion

You've built a simple web server in Go! Explore the [Go net/http documentation](https://golang.org/pkg/net/http/) to expand your knowledge.

Happy coding!
