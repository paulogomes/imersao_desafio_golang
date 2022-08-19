package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type Account struct {
	Number string  `json:"account_number"`
	Amount float32 `json:"amount"`
}

type Transfer struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float32 `json:"amount"`
}

var DB *sql.DB

func initDB(filepath string) {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		fmt.Println(err)
	}
	if db == nil {
		fmt.Println("db nil")
	}
	DB = db
}

func CreateTable() {
	// create table if not exists
	sql_table := "CREATE TABLE IF NOT EXISTS account (id INTEGER PRIMARY KEY AUTOINCREMENT, account_number TEXT, amount FLOAT)"
	_, err := DB.Exec(sql_table)
	if err != nil {
		fmt.Println(err)
	}
}

// Show
func getAccount(c *gin.Context) {
	key := c.Param("account_number")
	var account Account
	query := "SELECT account_number, amount from account WHERE account_number = ?"
	statement, _ := DB.Prepare(query)
	statement.QueryRow(key).Scan(&account.Number, &account.Amount)
	c.JSON(http.StatusOK, gin.H{"account": account})
}

// Create
func createAccount(c *gin.Context) {
	var id int

	var newAcc Account
	json.NewDecoder(c.Request.Body).Decode(&newAcc)

	tx, err := DB.Begin()

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}

	insert := "INSERT INTO account (account_number, amount) VALUES (?, ?) RETURNING id"
	statement, err := tx.Prepare(insert)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}

	statement.QueryRow(newAcc.Number, newAcc.Amount).Scan(&id)
	defer statement.Close()

	e := tx.Commit()

	if e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": e})
	} else {
		c.JSON(http.StatusCreated, gin.H{"ID": id, "account_number": newAcc.Number})
	}
}

// Update
func transfer(c *gin.Context) {
	var FromNewAmount float32
	var ToNewAmount float32

	var newTransfer Transfer
	json.NewDecoder(c.Request.Body).Decode(&newTransfer)

	tx, err := DB.Begin()

	query := "SELECT amount FROM account WHERE account_number = ?"
	update := "UPDATE account SET amount = ? WHERE account_number = ?"
	statement1, _ := tx.Prepare(update)
	defer statement1.Close()
	statement1.Exec(newTransfer.Amount*-1, newTransfer.From)
	// Obtem novo Amount
	stm, err := tx.Prepare(query)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
	stm.QueryRow(newTransfer.From).Scan(&FromNewAmount)

	statement2, _ := tx.Prepare(update)
	defer statement2.Close()
	statement2.Exec(newTransfer.Amount, newTransfer.To)
	// Obtem novo Amount
	stm2, _ := tx.Prepare(query)
	stm2.QueryRow(newTransfer.From).Scan(&ToNewAmount)

	tx.Commit()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	} else {
		c.JSON(http.StatusOK, gin.H{"FromNewAmount": FromNewAmount, "ToNewAmount": ToNewAmount})
	}
}

func main() {
	initDB("./data/Accounts.db")
	CreateTable()

	// Instancia o Router
	r := gin.Default()

	// endpoints
	v1 := r.Group("/bank-accounts")
	{
		v1.POST("", createAccount)
		v1.GET("/:account_number", getAccount)
		v1.POST("/transfer", transfer)
	}

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	r.Run()

}
