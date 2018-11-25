package main

import (
	"database/sql"
	"html/template"
	// "fmt"
	"log"
	// "math/rand"
	"net/http"
	"os"
	// "strconv"
	// "time"
	"encoding/json"
	"bytes"
	"io/ioutil"
	"github.com/go-redis/redis"
	_ "github.com/lib/pq"
)

var (
	dbHost     = envOrDefault("MYAPP_DATABASE_HOST", "localhost")
	dbPort     = envOrDefault("MYAPP_DATABASE_PORT", "5432")
	dbUser     = envOrDefault("MYAPP_DATABASE_USER", "root")
	dbPassword = envOrDefault("MYAPP_DATABASE_PASSWORD", "secret")
	dbName     = envOrDefault("MYAPP_DATABASE_NAME", "myapp")

	cacheHost = envOrDefault("MYAPP_CACHE_HOST", "localhost")
	cachePort = envOrDefault("MYAPP_CACHE_PORT", "6379")

	webHost = envOrDefault("MYAPP_WEB_HOST", "")
	webPort = envOrDefault("MYAPP_WEB_PORT", "8080")

	db    *sql.DB
	cache *redis.Client
)
type Gopher struct{
	Name string
}
type MealDetails struct {
	Name  string
	Subject string
	Message string
}
type Responser struct {
	Response string
	Success bool

}

func envOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func myHandler(w http.ResponseWriter, r *http.Request) {

	tmpl := template.Must(template.ParseFiles("forms.html"))

	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}

	details := MealDetails{
		Name:   r.FormValue("name"),
		Subject: r.FormValue("subject"),
		Message: r.FormValue("message"),
	}

	// do something with details
	_ = details
	log.Println(details.Name)
		jsonData := map[string]string{"query": details.Name}
	jsonValue, _ := json.Marshal(jsonData)
	request, _ := http.NewRequest("POST", "https://trackapi.nutritionix.com/v2/natural/nutrients", bytes.NewBuffer(jsonValue))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("x-app-id", "40523543")
	request.Header.Set("x-app-key", "44d9799d0bf08ca4a633dff233675a3d")
	request.Header.Set("x-remote-user-id", "0")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal("The HTTP request failed with error %s\n", err)
		z:= Responser{err.Error(),false}
		tmpl.Execute(w, z)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		log.Println(string(data))
		z:= Responser{string(data),true}
		tmpl.Execute(w, z)
	}

	
	// t, _:= template.ParseFiles("forms.html")
	// t.Execute(w, "Hello World!")

}



// func myCachedHandler(w http.ResponseWriter, r *http.Request) {
// 	n, err := cache.Get("n").Result()

// 	if err == redis.Nil {
// 		n = strconv.Itoa(rand.Intn(100))
// 		cache.Set("n", n, 5*time.Second)
// 	}

// 	fmt.Fprintf(w, "n = %s\n", n)
// }

func main() {
	// dbInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
	// 	dbHost, dbPort, dbUser, dbPassword, dbName)
	// var err error
	// db, err = sql.Open("postgres", dbInfo)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// if err = db.Ping(); err != nil {
	// 	log.Fatal(err)
	// }

	// cache = redis.NewClient(&redis.Options{
	// 	Addr: cacheHost + ":" + cachePort,
	// })
	// if _, err := cache.Ping().Result(); err != nil {
	// 	log.Fatal(err)
	// }

	http.HandleFunc("/", myHandler)
//	http.HandleFunc("/cache", myCachedHandler)
	log.Print("Listening on " + webHost + ":" + webPort + "...")
	http.ListenAndServe(webHost+":"+webPort, nil)
}
