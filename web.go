package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

type Configuration struct {
	KENSA_TEMPLATE_MYADDON_URL string
}

type AddonResource struct {
	Id     string        `json:"id"`
	Config Configuration `json:"config"`
}
type AddonResources []AddonResource

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/", Index)
	router.HandleFunc("/heroku/resources", use(Resources, basicAuth))
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), router))
}

// "/" Handler
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Kensa Tempalte for Golang!, %q", html.EscapeString(r.URL.Path))
}

func Resources(w http.ResponseWriter, r *http.Request) {
	c := Configuration{KENSA_TEMPLATE_MYADDON_URL: "https://example.com"}
	addonResource := AddonResource{Id: "kensa_template", Config: c}
	json.NewEncoder(w).Encode(addonResource)
}

func use(h http.HandlerFunc, middleware ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, m := range middleware {
		h = m(h)
	}

	return h
}

// Leverages nemo's answer in http://stackoverflow.com/a/21937924/556573
func basicAuth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)

		s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
		if len(s) != 2 {
			http.Error(w, "Not authorized", 401)
			return
		}

		b, err := base64.StdEncoding.DecodeString(s[1])
		if err != nil {
			http.Error(w, err.Error(), 401)
			return
		}

		pair := strings.SplitN(string(b), ":", 2)
		if len(pair) != 2 {
			http.Error(w, "Not authorized", 401)
			return
		}

		if pair[0] != "kensa_template" && pair[1] != "tYpx1jt652dRGIcK" {
			http.Error(w, "Not authorized", 401)
			return
		}

		h.ServeHTTP(w, r)
	}
}
