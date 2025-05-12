package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Result Service is running...")

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	http.ListenAndServe(":8083", nil)
}
