
```
package main

import (
	"fmt"
	"net/http"

	"github.com/k4s/httprouter"
)

func Whoami(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	firstName := params.Get(":first")
	lastName := params.Get(":last")
	fmt.Fprintf(w, "you are %s ,%s", firstName, lastName)
}
func Who(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	Name := params.Get(":id")
	fmt.Fprintf(w, "Iam %s", Name)
}

func main() {
	mux := httprouter.New()
	mux.Get("/:first(kas)/:last", Whoami)
	mux.Put("/user(user)/:id([0-9]+)", Who)

	//	http.Handle("/", mux)
	http.ListenAndServe(":8088", mux)
}
```