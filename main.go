package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/olahol/go-imageupload"
)

type Tribute struct {
	gorm.Model
	Name         string
	Relationship string
	Profile      string
	Tribute      string
}

func initialMigration() {
	db, err := gorm.Open("mysql", "root@/ripobed?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect database")
	}
	defer db.Close()

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

	r := mainRouter()

	http.ListenAndServe(":8080", r)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {

	tmpl := template.Must(template.ParseFiles("templates/home.html"))
	tmpl.Execute(w, nil)

}

func biographHandler(w http.ResponseWriter, r *http.Request) {

	tmpl := template.Must(template.ParseFiles("templates/biography.html"))
	tmpl.Execute(w, nil)
}

func tributeHandler(w http.ResponseWriter, r *http.Request) {

	db, err := gorm.Open("mysql", "root@/ripobed?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect database")
	}
	defer db.Close()

	if r.Method == "POST" {

		r.ParseForm()
		img, err := imageupload.Process(r, "profile")

		if err != nil {
			panic(err)
		}

		thumb, err := imageupload.ThumbnailPNG(img, 60, 60)

		if err != nil {
			panic(err)
		}

		image_path := fmt.Sprintf("%d.png", time.Now().Unix())

		thumb.Save("static/uploads/" + image_path)

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

		tmpl := template.Must(template.ParseFiles("templates/tribute.html"))
		tmpl.Execute(w, Tributes)
	}

}
