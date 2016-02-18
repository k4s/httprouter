
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
	Name := params.Get(":name")
	fmt.Fprintf(w, "Iam %s", Name)
}

func main() {
	mux := httprouter.New()
	mux.Get("/:first/:last", Whoami)
	mux.Put("/:name", Who)

	//	http.Handle("/", mux)
	http.ListenAndServe(":8088", mux)
}
```