# Creating a Simple Webserver in Go

![Webserver Illustration](/static/img/go_mascot.png)

This tutorial will guide you through building a simple webserver using Go. You will learn how to create a basic HTTP server, set up request handlers, and run your server locally.

## Prerequisites

- **Go Installed:** Ensure that you have [Go installed](https://golang.org/dl/).
- **Basic Terminal Knowledge:** Familiarity with the terminal/command prompt.
- **Basic Go Knowledge:** Understanding of Go syntax and packages.

## Step 1: Setting Up Your Project

1. Create a new directory for your project:

   ```bash
   mkdir simple-webserver
   cd simple-webserver
   ```

2. Initialize a new Go module (replace `example.com/simple-webserver` with your module path if needed):

   ```bash
   go mod init example.com/simple-webserver
   ```

## Step 2: Writing the Webserver Code

Create a file named `main.go` and add the following code:

```go
package main

import (
	"fmt"
	"log"
	"net/http"
)

// helloHandler responds to HTTP requests with a simple "Hello, World!" message.
func helloHandler(w http.ResponseWriter, r *http.Request) {
	// Write a response to the client
	fmt.Fprintf(w, "Hello, World!")
}

func main() {
	// Associate the helloHandler function with the "/" route.
	http.HandleFunc("/", helloHandler)

	// Inform the user that the server is running.
	fmt.Println("Server is running on http://localhost:8080")

	// Start the server on port 8080 and log any errors.
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

### Code Explanation

- **Imports:**

    - `fmt` is used for formatted I/O.
    - `log` handles logging of errors.
    - `net/http` provides HTTP client and server implementations.

- **Handler Function:**\
  The `helloHandler` function is triggered whenever an HTTP request is made to the root URL (`/`). It sends a simple "Hello, World!" message back to the client.

- **Main Function:**

    - `http.HandleFunc("/", helloHandler)` maps the root URL to the handler function.
    - `http.ListenAndServe(":8080", nil)` starts the server on port `8080`. If an error occurs, it will be logged using `log.Fatal`.

## Step 3: Running the Webserver

1. Open your terminal in the project directory.

2. Run the following command:

   ```bash
   go run main.go
   ```

3. Open a web browser and navigate to [http://localhost:8080](http://localhost:8080). You should see the message:

   ```
   Hello, World!
   ```

## Additional Tips

- **Changing the Port:**\
  To run the server on a different port, modify the argument in `http.ListenAndServe(":8080", nil)` to your desired port number, e.g., `":9090"`.

- **Handling More Routes:**\
  You can define additional handler functions and map them to different routes using `http.HandleFunc("/route", yourHandler)`.

- **Error Handling:**\
  For more robust applications, consider more comprehensive error handling and logging practices.

## Conclusion

You've successfully created a simple webserver in Go! This server is the foundation for building more complex web applications. Explore the [Go net/http documentation](https://golang.org/pkg/net/http/) to learn more about building web servers and handling HTTP requests.

Happy coding!

