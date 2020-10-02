package main

import (
	"assignment/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"text/template"

	_ "github.com/mattn/go-sqlite3"
)

func main() {

	http.HandleFunc("/home", home)
	http.HandleFunc("/register", registrasi)
	http.HandleFunc("/delete", deleteUser)
	http.HandleFunc("/update", update)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("Error running service: ", err)
	}
}

func ConnectDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./project.db")
	if err != nil {
		panic(err)
		return nil
	}

	return db
}

func registrasi(w http.ResponseWriter, r *http.Request) {
	db := ConnectDB()

	fmt.Println("Method: ", r.Method)
	if r.Method == "GET" {
		t, _ := template.ParseFiles("./templates/register.gtpl")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		if len(r.Form["nama"][0]) < 5 {
			fmt.Println("Nama harus lebih dari 5 karakter: ", r.Form["username"])
		}

		if m, _ := regexp.MatchString("^[0-9]+$", r.Form.Get("umur")); !m {
			fmt.Println("Umur yang di input bukan bilangan positif: ", r.Form["age"])
		}

		if len(r.Form["alamat"][0]) < 5 {
			fmt.Println("Alamat harus lebih dari 5 karakter: ", r.Form["username"])
		}

		if m, _ := regexp.MatchString("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$", r.Form.Get("email")); !m {
			fmt.Println("Harus sesuai dengan penulisan email", r.Form["email"])
		}

		if len(r.Form["role"][0]) < 3 {
			fmt.Println("role harus lebih dari 3 karakter: ", r.Form["role"])
		}

		nama := r.FormValue("nama")
		alamat := r.FormValue("alamat")
		umur := r.FormValue("umur")
		newUmur, _ := strconv.Atoi(umur)
		email := r.FormValue("email")
		role := r.FormValue("role")

		newData := &models.Users{nama, newUmur, alamat, email, role}

		dataJson, _ := json.Marshal(newData)
		err := ioutil.WriteFile("user.json", dataJson, 0644)
		fmt.Println(err)

		// insert
		stmt, err := db.Prepare("INSERT INTO users (nama, umur, alamat, email, role) values (?, ?, ?, ?, ?)")
		checkErr(err)

		res, err := stmt.Exec(nama, umur, alamat, email, role)
		checkErr(err)

		id, err := res.LastInsertId()
		checkErr(err)
		fmt.Println(id)

	}
}

func update(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Method: ", r.Method)
	if r.Method == "GET" {
		t, _ := template.ParseFiles("./templates/update.gtpl")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		if len(r.Form["nama"][0]) < 5 {
			fmt.Println("Nama harus lebih dari 5 karakter: ", r.Form["username"])
		}

		if m, _ := regexp.MatchString("^[0-9]+$", r.Form.Get("umur")); !m {
			fmt.Println("Umur yang di input bukan bilangan positif: ", r.Form["age"])
		}

		if len(r.Form["alamat"][0]) < 5 {
			fmt.Println("Alamat harus lebih dari 5 karakter: ", r.Form["username"])
		}

		if m, _ := regexp.MatchString("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$", r.Form.Get("email")); !m {
			fmt.Println("Harus sesuai dengan penulisan email", r.Form["email"])
		}

		if m, _ := regexp.MatchString("^[0-9]+$", r.Form.Get("id")); !m {
			fmt.Println("ID yang di input bukan bilangan positif: ", r.Form["id"])
		}

		nama := r.FormValue("nama")
		alamat := r.FormValue("alamat")
		umur := r.FormValue("umur")
		newUmur, _ := strconv.Atoi(umur)
		email := r.FormValue("email")
		role := r.FormValue("role")
		id := r.FormValue("id")

		newData := &models.Users{nama, newUmur, alamat, email, role}

		dataJson, _ := json.Marshal(newData)
		err := ioutil.WriteFile("user.json", dataJson, 0644)
		fmt.Println(err)

		db := ConnectDB()

		// update
		stmt, err := db.Prepare("UPDATE users SET nama=?, umur=?, alamat=?, email=? where user_id=?")
		checkErr(err)

		res, err := stmt.Exec(nama, umur, alamat, email, id)
		checkErr(err)

		affect, err := res.RowsAffected()
		checkErr(err)
		fmt.Println(affect)

	}
}

func home(w http.ResponseWriter, r *http.Request) {
	db := ConnectDB()

	rows, err := db.Query("SELECT * FROM users")
	checkErr(err)

	for rows.Next() {
		var userID int
		var nama string
		var umur int
		var alamat string
		var email string
		var role string
		err = rows.Scan(&userID, &nama, &umur, &alamat, &email, &role)
		checkErr(err)

		t, err := template.ParseFiles("./templates/table.gtpl")
		checkErr(err)

		users := struct {
			UserID int
			Nama   string
			Umur   int
			Alamat string
			Email  string
			Role   string
		}{userID, nama, umur, alamat, email, role}

		err = t.Execute(w, users)
		checkErr(err)
	}

}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	db := ConnectDB()

	stmt, err := db.Prepare("DELETE FROM users WHERE user_id=?")
	checkErr(err)

	res, err := stmt.Exec(1)
	checkErr(err)

	_, err = res.RowsAffected()
	checkErr(err)

	db.Close()

}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
