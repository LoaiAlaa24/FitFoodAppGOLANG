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

	dbHost  = envOrDefault("MYAPP_DATABASE_HOST", "localhost")
	dbPort  = envOrDefault("MYAPP_DATABASE_PORT", "27017")
	webPort = envOrDefault("MYAPP_WEB_PORT", "8080")
	dbName = envOrDefault("dbName", "" )
	dbUsername = envOrDefault("dbUsername","")
	dbPassword = envOrDefault("dbPassword","")

)

type mealResponse struct {
	meals []meal
}

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

type meal struct {
	Mealname     string `json:"mealname" bson:"mealname"`
	Username     string `json:"username" bson:"username"`
	NumberOfcals string `json:"numberOfcals" bson:"numberOfcals"`
}
type Config struct {
    DbName string `json:"dbName"`
    DbUsername string `json:"dbUsername"`
    DbPassword string `json:"dbPassword"`
}
type workout struct {
	Workoutname  string `json:"Workoutname" bson:"Workoutname"`
	Username     string `json:"username" bson:"username"`
	NumberOfcals string `json:"numberOfcals" bson:"numberOfcals"`
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

type HistoryResponse struct {
	Response        string
	Success         bool
	FoodResponse    string
	WorkoutResponse string
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

func createNewDb(url *mgo.DialInfo) (*mongoDbdatastore, error) {

	session, err := mgo.DialWithInfo(url)
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

func (m *mongoDbdatastore) CreateFoodItem(meal meal) error {

	session := m.Copy()

	defer session.Close()
	mealCollection := session.DB("FitFood").C("MealsHistory")
	err := mealCollection.Insert(&meal)
	if err != nil {
		return err
	}

	return nil
}

func (m *mongoDbdatastore) CreateExerciseItem(workout workout) error {

	session := m.Copy()

	defer session.Close()
	workoutCollection := session.DB("FitFood").C("WorkoutsHistory")
	err := workoutCollection.Insert(&workout)
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

func (m *mongoDbdatastore) getMealsByUsername(username string) ([]meal, error) {

	session := m.Copy()
	defer session.Close()
	mealsCollection := session.DB("FitFood").C("MealsHistory")
	//mea := meal{}
	meals := []meal{}
	err := mealsCollection.Find(bson.M{"username": username}).All(&meals)
	if err != nil {
		fmt.Println("meals array is empty")

		return nil, err
	}
	fmt.Println("meals array is not empty")
	fmt.Println(meals)
	return meals, nil

}

func (m *mongoDbdatastore) getMealsByname(name string) (meal, error) {

	session := m.Copy()
	defer session.Close()
	mealsCollection := session.DB("FitFood").C("MealsHistory")
	//mea := meal{}
	mea := meal{}
	err := mealsCollection.Find(bson.M{"username": name}).One(&mea)
	if err != nil {
		//fmt.Println("meals array is empty")

		return meal{}, err
	}
	//fmt.Println("meals array is not empty")
	//fmt.Println(meals)
	return mea, nil

}

func (m *mongoDbdatastore) getWorkoutsByUsername(username string) ([]workout, error) {

	session := m.Copy()
	defer session.Close()
	workoutsCollection := session.DB("FitFood").C("WorkoutsHistory")
	//wo := workout{}
	wo := []workout{}
	err := workoutsCollection.Find(bson.M{"username": username}).All(&wo)
	if err != nil {
		return nil, err
	}
	return wo, nil

}
func (m *mongoDbdatastore) getWorkoutsByname(name string) (workout, error) {

	session := m.Copy()
	defer session.Close()
	workoutsCollection := session.DB("FitFood").C("WorkoutsHistory")
	//wo := workout{}
	wo := workout{}
	err := workoutsCollection.Find(bson.M{"username": name}).One(&wo)
	if err != nil {
		return workout{}, err
	}
	return wo, nil

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
func myHandler2(e *mongoDbdatastore) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(uUsername) == 0 {

			http.Redirect(w, r, "/login", http.StatusSeeOther)

		} else {
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
					Workoutname := item["name"].(string)
					duration := item["duration_min"].(float64)
					numberOfcalories := item["nf_calories"].(float64)

					message := "number of calories burned during " + strconv.Itoa(int(duration)) + " minutes is " + strconv.Itoa(int(numberOfcalories)) + " Kcal"
					workoutTobeAdded := workout{Workoutname, uUsername, strconv.Itoa(int(numberOfcalories))}

					_, err = e.getWorkoutsByname(Workoutname)
					if err != nil {
						e.CreateExerciseItem(workoutTobeAdded)
					}
					z := Responser{string(data), true, message}
					tmpl.Execute(w, z)

				} else {
					z := Responser{"", true, "workout not found"}
					tmpl.Execute(w, z)

				}

			}
		}
	})
}
func myHandler(e *mongoDbdatastore) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if len(uUsername) == 0 {

			http.Redirect(w, r, "/login", http.StatusSeeOther)

		} else {
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
				//log.Fatal("The HTTP request failed with error %s\n", err)
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

					message := "food Name: " + foodName + "\n" + "Serving Quantity: " + strconv.Itoa(int(servingQty)) + "\n" +

						"Serving unit: " + servingUnit + "\n" +
						"serving grams: " + strconv.Itoa(int(servingGrams)) + "\n" +
						"number of calories: " + strconv.Itoa(int(numberOfcalories)) + "\n" +
						"carbs: " + strconv.Itoa(int(carbs)) + "\n" +
						"protein: " + strconv.Itoa(int(protein)) + "\n"
					mealTobeAdded := meal{foodName, uUsername, strconv.Itoa(int(numberOfcalories))}
					//e.CreateFoodItem(mealTobeAdded)
					_, err = e.getMealsByname(foodName)
					if err != nil {
						e.CreateFoodItem(mealTobeAdded)
					}
					z := Responser{string(data), true, message}
					tmpl.Execute(w, z)

				} else {
					z := Responser{"", true, "meal not found"}
					tmpl.Execute(w, z)

				}

			}
		}

	})
}

func myHandlerMenu(w http.ResponseWriter, r *http.Request) {
	if len(uUsername) == 0 {

		http.Redirect(w, r, "/login", http.StatusSeeOther)

	} else {
		tmpl := template.Must(template.ParseFiles("menu.html"))

		if r.Method != http.MethodPost {
			tmpl.Execute(w, nil)
			return
		}

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
			//log.Fatal(err)
			fmt.Println("email not found")

			z := Responser{"", true, "email not found"}
			tmpl.Execute(w, z)

		} else {
			if user.Password == password {

				// tmpl := template.Must(template.ParseFiles("menu.html"))
				// tmpl.Execute(w, nil)
				http.Redirect(w, r, "/menu", http.StatusSeeOther)
				uHeight = user.Height
				uWeight = user.Weight
				uAge = user.Age
				uGender = user.Gender
				uUsername = user.Username
				fmt.Println("Login in successfully !!")
				return

			} else {
				fmt.Println("incorrect password")
				z := Responser{"", true, "incorrect password"}
				tmpl.Execute(w, z)
			}
		}

	})
}

func myHandlerHistory(e *mongoDbdatastore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(uUsername) == 0 {

			http.Redirect(w, r, "/login", http.StatusSeeOther)

		} else {

			mealsResp, _ := e.getMealsByUsername(uUsername)
			fmt.Println(mealsResp)
			fmt.Println(uUsername)
			mealsArray := mealResponse{mealsResp}
			messageFood := ""
			for _, r := range mealsArray.meals {
				//fmt.Println(strings.Join(r.Mealname, ","))
				messageFood = messageFood + r.Mealname + " : " + r.NumberOfcals + "kcal " + "     |     "

			}
			fmt.Println(messageFood)

			workoutsReponse, _ := e.getWorkoutsByUsername(uUsername)
			fmt.Println(workoutsReponse)
			workoutsString := ""
			for _, r := range workoutsReponse {
				//fmt.Println(strings.Join(r.Mealname, ","))
				workoutsString = workoutsString + r.Workoutname + " : " + r.NumberOfcals + "kcal " + "      |      "

			}
			fmt.Println(workoutsString)

			tmpl := template.Must(template.ParseFiles("history.html"))
			//fmt.Println(test.Tests[0])
			z := HistoryResponse{"", true, messageFood, workoutsString}
			if r.Method != http.MethodPost {
				tmpl.Execute(w, z)
				return
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

		if age < 0 || weight < 0 || height < 0 {
			z := Responser{"", true, "Please enter a valid numbers for either the age or the weight or the height"}
			tmpl.Execute(w, z)
		} else {

			u := user{r.FormValue("username"), r.FormValue("email"), r.FormValue("password"),
				height, weight, age, r.FormValue("gender")}

			user, err := e.getUserEmail(u.Email)

			if err != nil {
				_, err := e.getUserUsername(u.Username)
				if err != nil {
					e.CreateUser(u)

					// tmpl := template.Must(template.ParseFiles("menu.html"))
					// tmpl.Execute(w, nil)
					http.Redirect(w, r, "/menu", http.StatusSeeOther)

					fmt.Println("Login in successfully !!")
					fmt.Println("user added")
					uHeight = u.Height
					uWeight = u.Weight
					uAge = u.Age
					uGender = u.Gender
					uUsername = u.Username
					return

				} else {
					fmt.Println("username already exists")
					z := Responser{"", true, "username already exists"}
					tmpl.Execute(w, z)
				}

			} else {
				if user.Email == u.Email {
					fmt.Println("email already exists")
					z := Responser{"", true, "email already exists"}
					tmpl.Execute(w, z)
				}

			}
		}

	})

}
func LoadConfiguration(file string) Config {
    var config Config
    configFile, err := os.Open(file)
    defer configFile.Close()
    if err != nil {
        fmt.Println(err.Error())
    }
    jsonParser := json.NewDecoder(configFile)
    jsonParser.Decode(&config)
    return config
}
func main() {

	//db, err := createNewDb(dbHost + ":27017")
	config := LoadConfiguration("config.json")
	log.Println(config)
	stringArray := []string{dbHost}
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    stringArray,
		Database: config.DbName,
		Username: config.DbUsername,
		Password: config.DbPassword,
	}

	db, _ := createNewDb(mongoDBDialInfo)

	login := myHandlerLogin(db)
	signUp := myHandlerSigning(db)
	mealSearch := myHandler(db)
	workoutsSearch := myHandler2(db)
	history := myHandlerHistory(db)

	http.HandleFunc("/menu", myHandlerMenu)
	http.Handle("/", login)
	http.Handle("/signUp", signUp)
	http.Handle("/meal", mealSearch)
	http.Handle("/exercise", workoutsSearch)
	http.Handle("/history", history)
	fmt.Println("Listening on :" + webPort + "...")
	http.ListenAndServe(":"+webPort, nil)
	//.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
}
