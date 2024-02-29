package main

import (
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/Unleash/unleash-client-go/v3"
)

type metricsInterface struct{}

func init() {
	unleash.Initialize(
		unleash.WithUrl("https://gitlab.com/api/v4/feature_flags/unleash/40951967"),
		unleash.WithInstanceId("8bZJ99faxtsW3anLf2ak"),
		unleash.WithAppName("production"), // Set to the running environment of your application
		unleash.WithListener(&metricsInterface{}),
	)
}

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", CatHandler)
	log.Printf("Listening on port 5000!")
	http.ListenAndServe(":5000", nil)
}

func Random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

func CatHandler(w http.ResponseWriter, r *http.Request) {
	// Fetch hostname of container
	name, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	// Enable Grumpy Cat Feature Flag
	var catpic int
	var message string
	if unleash.IsEnabled("grumpy-cat") {
		catpic = Random(11, 15)
		message = "Grumpy Cat Feature Flag Enabled"
	} else {
		catpic = Random(1, 10)
		message = "Grumpy Cat is Off - Have Fun :)"
	}

	// Introduce XSS Vulnerability
	userInput := r.URL.Query().Get("userInput") // Get user input from query parameter

	// Parse index.html template
	t, err := template.ParseFiles("index.html")
	if err != nil {
		fmt.Println(err)
	}

	items := struct {
		Url       int
		Hostname  string
		Message   string
		UserInput string // Add UserInput to the data passed to the template
	}{
		Url:       catpic,
		Hostname:  name,
		Message:   message,
		UserInput: userInput, // Pass the unescaped user input to the template
	}

	t.Execute(w, items)
	log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
}
