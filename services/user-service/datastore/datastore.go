package datastore

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/ankushChatterjee/jetpen/user-service/models"
	"github.com/ankushChatterjee/jetpen/user-service/pkg/utils"
	_ "github.com/lib/pq"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Printf("ERROR : %s: %s", msg, err)
	}
}

func fatalOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("ERROR : %s: %s", msg, err)
	}
}

func GetConnection() *sql.DB {
	dbHost := utils.GetEnvVar("DB_HOST")
	dbPassword := utils.GetEnvVar("DB_PASSWORD")
	dbPort := utils.GetEnvVar("DB_PORT")
	dbName := utils.GetEnvVar("DB_NAME")
	dbUserName := utils.GetEnvVar("DB_USERNAME")

	dbSQLInfo := fmt.Sprintf("host=%s port=%s user=%s "+"password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUserName, dbPassword, dbName)
	fmt.Println(dbSQLInfo)
	db, err := sql.Open("postgres", dbSQLInfo)
	fatalOnError(err, "Error Opening Database")
	return db
}

func ExecuteSQL(sqlString string, db *sql.DB) error {
	_, err := db.Exec(sqlString)
	failOnError(err, "Error Executing SQL")
	return err
}

func GetUser(username string, db *sql.DB) (*models.User, error) {
	user := models.User{}
	sqlString := "SELECT username, email, password, name, \"CreatedAt\" from jetpen.users WHERE username=$1"
	row := db.QueryRow(sqlString, username)
	switch err := row.Scan(&user.Username, &user.Email, &user.Password, &user.Name, &user.CreatedAt); err {
	case sql.ErrNoRows:
		failOnError(err, "No Such User Exist")
		return nil, err
	case nil:
		return &user, nil
	default:
		failOnError(err, "Error Reading users")
		return nil, err
	}
}

func InsertUser(user *models.User, db *sql.DB) error {
	sqlString := "INSERT INTO jetpen.users(username, email, name, password) VALUES ($1, $2, $3, $4)"
	_, err := db.Exec(sqlString, user.Username, user.Email, user.Name, user.Password)
	failOnError(err, "Error inserting user")
	return err
}

func UpdateNameOfUser(username string, name string, db *sql.DB) error {
	sqlString := "UPDATE jetpen.users SET name=$1 WHERE username=$2"
	_, err := db.Exec(sqlString, name, username)
	failOnError(err, "Error inserting user")
	return err
}

func UpdateUserPassword(username string, password string, db *sql.DB) error {
	sqlString := "UPDATE jetpen.users SET password=$1 WHERE username=$2"
	hashPassword, _ := utils.HashPassword(password)
	_, err := db.Exec(sqlString, hashPassword, username)
	failOnError(err, "Error inserting user")
	return err
}

func UpdateUserEmail(username string, email string, db *sql.DB) error {
	sqlString := "UPDATE jetpen.users SET email=$1 WHERE username=$2"
	_, err := db.Exec(sqlString, email, username)
	failOnError(err, "Error inserting user")
	return err
}

func DeleteUser(username string, db *sql.DB) error {
	sqlString := "DELETE FROM jetpen.users WHERE username=$1"
	_, err := db.Exec(sqlString, username)
	failOnError(err, "Error Deleting user")
	return err
}

func InsertTempUser(user *models.TempUser, db *sql.DB) error {
	sqlString := "INSERT INTO jetpen.temp_users(username, email, name, password, token) VALUES ($1, $2, $3, $4, $5)"
	_, err := db.Exec(sqlString, user.Username, user.Email, user.Name, user.Password, user.Token)
	failOnError(err, "Error inserting temp user")
	return err
}

func DeleteTempUser(username string, db *sql.DB) error {
	sqlString := "DELETE FROM jetpen.temp_users WHERE username=$1"
	_, err := db.Exec(sqlString, username)
	failOnError(err, "Error Deleting user")
	return err
}

func GetTempUser(username string, db *sql.DB) (*models.TempUser, error) {
	user := models.TempUser{}
	sqlString := "SELECT username, email, password, name, \"CreatedAt\", token from jetpen.temp_users WHERE username=$1"
	row := db.QueryRow(sqlString, username)
	switch err := row.Scan(&user.Username, &user.Email, &user.Password, &user.Name, &user.CreatedAt, &user.Token); err {
	case sql.ErrNoRows:
		failOnError(err, "No Such Temp User Exist")
		return nil, err
	case nil:
		return &user, nil
	default:
		failOnError(err, "Error Reading temp users")
		return nil, err
	}
}
