package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
	"github.com/justinas/nosurf"
	"github.com/olahol/go-imageupload"
)

type Tribute struct {
	gorm.Model
	Name         string
	Relationship string
	Profile      string
	Tribute      string `gorm:"type:text"`
}

var (
	db *gorm.DB
)

func initDb() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var db_err error

	db, db_err = gorm.Open("mysql", fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("username"),
		os.Getenv("password"), os.Getenv("host"), os.Getenv("db_name")))
	if db_err != nil {
		fmt.Println(db_err.Error())
		panic("failed to connect database")
	}

}

func initialMigration() {

	// Migrate the schema
	db.AutoMigrate(&Tribute{})
}

func mainRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler)
	r.HandleFunc("/biography", biographHandler)
	r.HandleFunc("/tributes", tributeHandler)

	initialMigration()

	staticFileDirectory := http.Dir("./static/")

	staticFileHandler := http.StripPrefix("/static/", http.FileServer(staticFileDirectory))

	r.PathPrefix("/static/").Handler(staticFileHandler).Methods("GET")

	return r
}

func main() {

	initDb()

	r := mainRouter()

	defer db.Close()

	http.ListenAndServe(":8080", nosurf.New(r))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {

	var Tribute Tribute

	db.Last(&Tribute)

	tmpl := template.Must(template.ParseFiles("templates/home.html"))
	tmpl.Execute(w, Tribute)

}

func biographHandler(w http.ResponseWriter, r *http.Request) {

	tmpl := template.Must(template.ParseFiles("templates/biography.html"))
	tmpl.Execute(w, nil)
}

func tributeHandler(w http.ResponseWriter, r *http.Request) {

	var image_path string

	if r.Method == "POST" {

		r.ParseForm()

		img, err := imageupload.Process(r, "profile")

		if err != nil {

			image_path = ""

		} else {

			thumb, err := imageupload.ThumbnailPNG(img, 60, 60)

			if err != nil {
				panic(err)
			}

			image_path = fmt.Sprintf("%d.png", time.Now().Unix())

			thumb.Save("static/uploads/" + image_path)

		}

		name := r.Form["name"][0]
		relationship := r.Form["relationship"][0]
		tribute := r.Form["tribute"][0]
		profile := image_path

		var tributeRecord = Tribute{Name: name, Relationship: relationship, Tribute: tribute, Profile: profile}

		db.Create(&tributeRecord)

		http.Redirect(w, r, "/tributes", http.StatusSeeOther)

	} else {

		var Tributes []Tribute
		db.Order("id desc").Find(&Tributes)

		context := make(map[string]interface{})
		context["token"] = nosurf.Token(r)
		context["Tributes"] = Tributes

		tmpl := template.Must(template.ParseFiles("templates/tribute.html"))
		tmpl.Execute(w, context)
	}

}
