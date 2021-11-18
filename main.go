package main

import (
	"database/sql"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/squishd/usersession"
)

const (
	DatabaseLocation = "christmasList.db"
	WebsiteLogFile   = "christmasList.log"
	SalesTaxRate     = 0.06 // michigan sales tax
)

var (
	l *ListOrganizer
)

func doAddGift(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	s := usersession.New()
	s.Ip = r.RemoteAddr
	log.Println(s.GetID(), s.Ip, r.URL.Path)
	if f, err := ioutil.ReadFile("addgift.html"); err != nil {
		log.Fatalln("Error - can't serve addgift.html", err)
	} else {
		w.Write(f)
	}
}

func doInsert(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	s := usersession.New()
	s.Ip = r.RemoteAddr
	log.Println(s.GetID(), s.Ip, r.URL.Path)
	err := r.ParseForm()
	if err != nil {
		log.Fatalln(err)
	}
	rcpt := r.Form.Get("rcptname")
	giftname := r.Form.Get("giftname")
	priceString := r.Form.Get("price")
	priceFloat, _ := strconv.ParseFloat(priceString, 64)
	priceInt := int(priceFloat * 100)
	url := r.Form.Get("url")
	l.AddGift(rcpt, giftname, priceInt, url)
	http.Redirect(w, r, "/addgift", http.StatusFound)
}

func doGiftList(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	s := usersession.New()
	s.Ip = r.RemoteAddr
	log.Println(s.GetID(), s.Ip, r.URL.Path)
	var dataToServe []Recipient
	if len(r.URL.Query()) == 0 {
		dataToServe = l.GetAllUsersAndGifts()
	} else {
		username := r.URL.Query()["username"][0]
		if username == "" {
			dataToServe = l.GetAllUsersAndGifts()
		} else {
			recipient, err := l.GetRecipient(username)
			if err != nil {
				log.Println("got query for user", username, "but user not found")
				dataToServe = l.GetAllUsersAndGifts()
			} else {
				dataToServe = []Recipient{recipient}
			}
		}
	}
	if t, err := template.ParseFiles("giftlist.html"); err != nil {
		log.Fatalln("Error - can't serve giftlist.html", err)
	} else {
		t.Execute(w, dataToServe)
	}
}

func doPurchased(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	s := usersession.New()
	s.Ip = r.RemoteAddr
	log.Println(s.GetID(), s.Ip, r.URL.Path)
	itemIDString := r.URL.Query()["giftid"][0]
	username := r.URL.Query()["username"][0]
	id, err := strconv.Atoi(itemIDString)
	if err != nil {
		log.Fatalln("Error converting ID to int", err)
	}
	l.SetPurchased(id)
	http.Redirect(w, r, "/edit?username="+username, http.StatusFound)
}

func doDeleteUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	s := usersession.New()
	s.Ip = r.RemoteAddr
	log.Println(s.GetID(), s.Ip, r.URL.Path)
	username := r.URL.Query()["username"][0]
	userIDString := r.URL.Query()["userid"][0]
	id, err := strconv.Atoi(userIDString)
	if err != nil {
		log.Fatalln(err)
	}
	l.DropRecipient(id, username)
	http.Redirect(w, r, "/", http.StatusFound)
}

func doDeleteGift(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	s := usersession.New()
	s.Ip = r.RemoteAddr
	log.Println(s.GetID(), s.Ip, r.URL.Path)
	itemIDString := r.URL.Query()["giftid"][0]
	username := r.URL.Query()["username"][0]
	id, err := strconv.Atoi(itemIDString)
	if err != nil {
		log.Fatalln("Error:", err)
	}
	l.DropGift(id)
	http.Redirect(w, r, "/edit?username="+username, http.StatusFound)
}

func doFavicon(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	http.ServeFile(w, r, "favicon.ico")
}

func doIndex(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	s := usersession.New()
	s.Ip = r.RemoteAddr
	log.Println(s.GetID(), s.Ip, r.URL.Path)
	if t, err := template.ParseFiles("index.html"); err != nil {
		log.Fatalln(err)
	} else {
		databaseDump := l.GetAllUsersAndGifts()
		t.Execute(w, databaseDump)
	}
}

func setupHandlers() {
	mux.HandleFunc("/", doIndex)
	mux.HandleFunc("/addgift", doAddGift)
	mux.HandleFunc("/giftlist", doGiftList)
	mux.HandleFunc("/insert", doInsert)
	mux.HandleFunc("/purchased", doPurchased)
	mux.HandleFunc("/deletegift", doDeleteGift)
	mux.HandleFunc("/deleteuser", doDeleteUser)
	mux.HandleFunc("/favicon.ico", doFavicon)
}

var mux = http.NewServeMux()
var server = &http.Server{
	Handler: http.TimeoutHandler(mux, time.Second*60, "Request Timed Out!"),
	Addr:    ":8080",
}

func main() {
	// prep DB
	file, err := os.OpenFile(WebsiteLogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		log.Fatalln("Error creating log file:", err)
	}
	defer file.Close()
	log.SetOutput(file)
	setupHandlers()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Fatalln(server.ListenAndServe())
	}()
	theDatabase, err := sql.Open("sqlite3", DatabaseLocation)
	if err != nil {
		log.Fatalln("Error opening database:", err)
	}
	defer theDatabase.Close()
	createGiftsTable(theDatabase)
	createRecipientTable(theDatabase)
	l = NewListOrganizer(theDatabase, "ChristmasList")
	wg.Wait()
}
