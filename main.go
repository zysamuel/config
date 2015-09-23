package main

import (
	"net/http"
       )	

func main() {
    restRtr := createNewRestRouter()
    http.ListenAndServe(":8080", restRtr)
}
