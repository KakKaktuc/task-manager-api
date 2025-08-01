package main 

import (
	"fmt"
	"net/http"

	"github.com/KakKaktuc/task-manager-api/internal/handler"
	"github.com/KakKaktuc/task-manager-api/internal/repository"
	"github.com/KakKaktuc/task-manager-api/pkg/utils"
)

func pingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong\n")
}

func headers(w http.ResponseWriter, req *http.Request) {

	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func main() {
	userRepo := repository.NewUserRepository()
	userHandler := handler.NewUserHandler(userRepo)

	http.Handle("/users", utils.RecoverMiddleware((userHandler)))
	http.Handle("/users/", utils.RecoverMiddleware((userHandler)))

	http.HandleFunc("/headers", headers)
	http.HandleFunc("/ping", pingHandler)

	fmt.Println("Server is running on port 8090")
	http.ListenAndServe(":8090", nil)
}
