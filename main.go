package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type employee struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	BAL  string `json:"bal"`
}

type JsonResponse struct {
	Type    string     `json:"type"`
	Data    []employee `json:"data"`
	Message string     `json:"message"`
}

func setupDB() *sql.DB {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/test")

	if err != nil {
		log.Fatal(err)
	}

	return db
}

func main() {
	rout := mux.NewRouter()
	rout.HandleFunc("/", show).Methods("GET")
	rout.HandleFunc("/", insert).Methods("POST")
	rout.HandleFunc("/delete/{id}", delete).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", rout))
}

func show(w http.ResponseWriter, r *http.Request) {
	db := setupDB()

	rows, err := db.Query("SELECT * FROM transaction")

	if err != nil {
		panic(err)
	}

	var emp []employee

	for rows.Next() {
		var Id string
		var Nam string
		var bl string

		err = rows.Scan(&Id, &Nam, &bl)

		if err != nil {
			panic(err)
		}

		emp = append(emp, employee{ID: Id, Name: Nam, BAL: bl})
	}

	var response = JsonResponse{Type: "success", Data: emp}

	json.NewEncoder(w).Encode(response)
}

func insert(w http.ResponseWriter, r *http.Request) {
	db := setupDB()
	var response = JsonResponse{}

	w.Header().Set("Content-Type", "application/json")

	stmt, err := db.Prepare("INSERT INTO transaction(id,name,balance) VALUES(?,?,?)")

	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		panic(err.Error())
	}

	keyVal := make(map[string]string)

	json.Unmarshal(body, &keyVal)

	_, err = stmt.Exec(keyVal["id"], keyVal["name"], keyVal["bal"])

	if err != nil {
		panic(err.Error())
	}
	response = JsonResponse{Type: "success", Message: "Record inserted successfully!"}

	json.NewEncoder(w).Encode(response)
}

func delete(w http.ResponseWriter, r *http.Request) {
	db := setupDB()
	var response = JsonResponse{}
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	stmt, err := db.Prepare("DELETE FROM transaction WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}
	_, err = stmt.Exec(params["id"])
	if err != nil {
		panic(err.Error())
	}
	response = JsonResponse{Type: "success", Message: "Record deleted successfully!"}
	json.NewEncoder(w).Encode(response)
}
