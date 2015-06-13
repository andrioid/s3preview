package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	//"os"
)

// S3 creds (delete before github push)
// AKIAIKE57LTDKBK5W6FA
// 6sFwzVUXFLmzhMewi0jQQXBL0OJsTUmKA3TSTRkz

var preview_bucket = "andridk-assets"
var preview_prefix = "s3preview"

func main() {
	r := mux.NewRouter()

	http.Handle("/", r)
	registerHandlers(r)
	fmt.Printf("I'm listening...\n")
	http.ListenAndServe(":80", nil)

}
