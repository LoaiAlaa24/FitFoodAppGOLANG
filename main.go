package main

import (
	"fmt"
	"html/template"
	"strconv"

	// "fmt"
	"log"
	// "math/rand"
	"net/http"
	"os"

	// "strconv"
	// "time"
	"bytes"
	"encoding/json"
	"io/ioutil"

	_ "github.com/lib/pq"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	uHeight   float64
	uWeight   float64
	uAge      float64
	uGender   string
	uUsername string

	dbHost     = envOrDefault("MYAPP_DATABASE_HOST", "localhost")
	dbPort     = envOrDefault("MYAPP_DATABASE_PORT", "5432")
	dbUser     = envOrDefault("MYAPP_DATABASE_USER", "root")
	dbPassword = envOrDefault("MYAPP_DATABASE_PASSWORD", "secret")
	dbName     = envOrDefault("MYAPP_DATABASE_NAME", "myapp")

	cacheHost = envOrDefault("MYAPP_CACHE_HOST", "localhost")
	cachePort = envOrDefault("MYAPP_CACHE_PORT", "6379")

	webHost = envOrDefault("MYAPP_WEB_HOST", "")
	webPort = envOrDefault("MYAPP_WEB_PORT", "8080")

	// db    *sql.DB
	// cache *redis.Client
)

// type Result struct {
//     Common[]   string        `json:"common"`
//     Branded[] struct{
// 		ItemID	`json:"nix_item_id"`
// 		ServingUnit `json:"serving_unit"`
// 		Photo `json:"photo"`
// 		Calories `json:"nf_calories"`
// 	}`json:"branded"`

// }

type mongoDbdatastore struct {
	*mgo.Session
}

func FloatToString(input_num float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(input_num, 'f', 6, 64)
}

type user struct {
	//UUID     string  `json:"uuid" bson:"uuid"`
	Username string  `json:"username" bson:"username"`
	Email    string  `json:"email" bson:"email"`
	Password string  `json:"password" bson:"password"`
	Height   float64 `json:"height" bson:"height"`
	Weight   float64 `json:"weight" bson:"weight"`
	Age      float64 `json:"age" bson:"age"`
	Gender   string  `json:"gender" bson:"gender"`
}

type MealDetails struct {
	Name string
}
type ExDetails struct {
	Name string
	Age  string
}
type Responser struct {
	Response string
	Success  bool
	Message  string
}

type signUp struct {
	Email    string
	Password string
	Username string
	Height   float64
	Age      float64
	Gender   string
}

type logIn struct {
	Email    string
	Password string
}

func createNewDb(url string) (*mongoDbdatastore, error) {

	session, err := mgo.Dial(url)
	if err != nil {
		return nil, err
	}
	return &mongoDbdatastore{
		Session: session,
	}, nil
}
func (m *mongoDbdatastore) CreateUser(user user) error {

	session := m.Copy()

	defer session.Close()
	userCollection := session.DB("FitFood").C("Users")
	err := userCollection.Insert(&user)
	if err != nil {
		return err
	}

	return nil
}

func (m *mongoDbdatastore) getUserEmail(email string) (user, error) {

	session := m.Copy()
	defer session.Close()
	userCollection := session.DB("FitFood").C("Users")
	u := user{}
	err := userCollection.Find(bson.M{"email": email}).One(&u)
	if err != nil {

		return user{}, err
	}
	log.Println(u)
	return u, nil

}

func (m *mongoDbdatastore) getUserUsername(username string) (user, error) {

	session := m.Copy()
	defer session.Close()
	userCollection := session.DB("FitFood").C("Users")
	u := user{}
	err := userCollection.Find(bson.M{"username": username}).One(&u)
	if err != nil {
		return user{}, err
	}
	return u, nil

}

func (m *mongoDbdatastore) Close() {
	m.Close()
}

func envOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func myHandler2(w http.ResponseWriter, r *http.Request) {

	tmpl := template.Must(template.ParseFiles("forms2.html"))

	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}

	details := MealDetails{
		Name: r.FormValue("name"),
	}
	exdetails := ExDetails{
		Name: r.FormValue("name"),
	}
	// do something with details
	_ = details
	log.Println(details.Name)
	jsonData := map[string]string{"query": exdetails.Name, "age": FloatToString(uAge), "height_cm": FloatToString(uHeight), "weight_kg": FloatToString(uWeight), "gender": uGender}
	jsonValue, _ := json.Marshal(jsonData)
	request, _ := http.NewRequest("POST", "https://trackapi.nutritionix.com/v2/natural/exercise", bytes.NewBuffer(jsonValue))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("x-app-id", "40523543")
	request.Header.Set("x-app-key", "44d9799d0bf08ca4a633dff233675a3d")
	request.Header.Set("x-remote-user-id", "0")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		//log.Fatal("The HTTP request failed with error %s\n", err)
		//fmt.Println("weselna lel error")
		// z := Responser{err.Error(), false, "workout not found"}
		// tmpl.Execute(w, z)
	} else {
		data, _ := ioutil.ReadAll(response.Body)

		log.Println(string(data))

		var result map[string]interface{}
		json.Unmarshal([]byte(data), &result)
		exe := result["exercises"].([]interface{})
		if len(exe) != 0 {
			item := exe[0].(map[string]interface{})

			duration := item["duration_min"].(float64)
			numberOfcalories := item["nf_calories"].(float64)

			message := "number of calories burned during " + FloatToString(duration) + " minutes is " + FloatToString(numberOfcalories) + " Kcal"

			z := Responser{string(data), true, message}
			tmpl.Execute(w, z)

		} else {
			z := Responser{"", true, "workout not found"}
			tmpl.Execute(w, z)

		}

	}

}

func myHandler(w http.ResponseWriter, r *http.Request) {

	tmpl := template.Must(template.ParseFiles("forms.html"))

	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}

	details := MealDetails{
		Name: r.FormValue("name"),
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
		z := Responser{err.Error(), false, ""}
		tmpl.Execute(w, z)
	} else {

		data, _ := ioutil.ReadAll(response.Body)

		log.Println(string(data))

		var result map[string]interface{}
		json.Unmarshal([]byte(data), &result)
		mea := result["foods"].([]interface{})
		if len(mea) != 0 {
			item := mea[0].(map[string]interface{})

			foodName := item["food_name"].(string)
			servingQty := item["serving_qty"].(float64)
			servingUnit := item["serving_unit"].(string)
			servingGrams := item["serving_weight_grams"].(float64)
			numberOfcalories := item["nf_calories"].(float64)
			carbs := item["nf_total_carbohydrate"].(float64)
			protein := item["nf_protein"].(float64)

			message := "food Name: " + foodName + "\n" + "Serving Quantity: " + FloatToString(servingQty) + "\n" +
				"Serving unit: " + servingUnit + "\n" +
				"serving grams: " + FloatToString(servingGrams) + "\n" +
				"number of calories: " + FloatToString(numberOfcalories) + "\n" +
				"carbs: " + FloatToString(carbs) + "\n" +
				"protein: " + FloatToString(protein) + "\n"

			z := Responser{string(data), true, message}
			tmpl.Execute(w, z)

		} else {
			z := Responser{"", true, "meal not found"}
			tmpl.Execute(w, z)

		}

	}

}

func myHandlerMenu(w http.ResponseWriter, r *http.Request) {

	tmpl := template.Must(template.ParseFiles("menu.html"))

	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}

}

func myHandlerLogin(e *mongoDbdatastore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tmpl := template.Must(template.ParseFiles("login.html"))

		if r.Method != http.MethodPost {
			tmpl.Execute(w, nil)
			return
		}

		email := r.FormValue("email")
		password := r.FormValue("password")

		user, err := e.getUserEmail(email)
		if err != nil {
			log.Fatal(err)
			fmt.Println("user was not found")

		} else {
			if user.Password == password {

				tmpl := template.Must(template.ParseFiles("menu.html"))
				tmpl.Execute(w, nil)
				uHeight = user.Height
				uWeight = user.Weight
				uAge = user.Age
				uGender = user.Gender
				uUsername = user.Username
				fmt.Println("Login in successfully !!")
				return

			} else {
				fmt.Println("incorrect password")
			}
		}

	})
}

func myHandlerSigning(e *mongoDbdatastore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tmpl := template.Must(template.ParseFiles("signUp.html"))

		if r.Method != http.MethodPost {
			tmpl.Execute(w, nil)
			return
		}

		weight, _ := strconv.ParseFloat(r.FormValue("weight"), 64)
		height, _ := strconv.ParseFloat(r.FormValue("height"), 64)
		age, _ := strconv.ParseFloat(r.FormValue("age"), 64)

		u := user{r.FormValue("username"), r.FormValue("email"), r.FormValue("password"),
			height, weight, age, r.FormValue("gender")}

		user, err := e.getUserEmail(u.Email)

		if err != nil {
			_, err := e.getUserUsername(u.Username)
			if err != nil {
				e.CreateUser(u)

				tmpl := template.Must(template.ParseFiles("menu.html"))
				tmpl.Execute(w, nil)
				fmt.Println("Login in successfully !!")
				fmt.Println("user added")
				return

			} else {
				fmt.Println("username already exists")
			}

		} else {
			if user.Email == u.Email {
				fmt.Println("email already exists")
			}

		}

	})

}

func main() {

	db, err := createNewDb(dbHost + ":27017")
	//db, err := createNewDb("localhost:2000")
	if err != nil {
		fmt.Println("error")
		log.Fatal(err)
	}
	//else {
	// 	u := user{"marwanihab", "marwanihab95@gmail.com", "123456", 70, 50, 160, "male"}

	// 	db.CreateUser(u)
	// 	// user, err := db.GetUser("marwanihab")

	// 	// if err != nil {
	// 	// 	log.Fatal(err)
	// 	// 	fmt.Println("user was not found")
	// 	// }

	// 	//fmt.Println("success")
	// }

	login := myHandlerLogin(db)
	signUp := myHandlerSigning(db)

	//http.HandleFunc("/", myHandlerMenu)
	http.HandleFunc("/menu", myHandlerMenu)
	http.Handle("/", login)
	http.Handle("/signUp", signUp)
	http.HandleFunc("/meal", myHandler)
	http.HandleFunc("/exercise", myHandler2)
	log.Print("Listening on " + webHost + ":" + webPort + "...")
	http.ListenAndServe(webHost+":"+webPort, nil)
}
