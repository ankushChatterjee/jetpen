package datastore

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/ankushChatterjee/jetpen/sub-service/models"
	"github.com/ankushChatterjee/jetpen/sub-service/pkg/utils"
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

func InsertSubscription(sub *models.Subscription, db *sql.DB) error {
	sqlString := "INSERT INTO jetpen.subscription(id,email, nid, subToken) VALUES($1, $2, $3, $4)"
	_, err := db.Exec(sqlString, sub.Id, sub.Email, sub.Nid, sub.SubToken)
	failOnError(err, "Error Inserting Subscription")
	return err
}

func DeleteSubscription(id string, db *sql.DB) error {
	sqlString := "DELETE FROM jetpen.subscription WHERE id = $1"
	_, err := db.Exec(sqlString, id)
	failOnError(err, "Error Deleting subscription")
	return err
}

func GetSubscription(id string, db *sql.DB) (*models.Subscription, error) {
	sqlString := "SELECT email, nid, \"CreatedAt\", subToken FROM jetpen.subscription WHERE id = $1"
	row := db.QueryRow(sqlString, id)
	sub := new(models.Subscription)
	sub.Id = id
	err := row.Scan(&sub.Email, &sub.Nid, &sub.CreatedAt, &sub.SubToken)
	failOnError(err, "Getting Subscription")
	return sub, err
}

func InsertTempSubscription(sub *models.Subscription, db *sql.DB) error {
	sqlString := "INSERT INTO jetpen.temp_subscription(id, email, nid, subToken) VALUES($1, $2, $3, $4)"
	_, err := db.Exec(sqlString, sub.Id, sub.Email, sub.Nid, sub.SubToken)
	failOnError(err, "Error Inserting Subscription")
	return err
}

func DeleteTempSubscription(id string, db *sql.DB) error {
	sqlString := "DELETE FROM jetpen.temp_subscription WHERE id = $1"
	_, err := db.Exec(sqlString, id)
	failOnError(err, "Error Deleting subscription")
	return err
}

func GetTempSubscription(id string, db *sql.DB) (*models.Subscription, error) {
	sqlString := "SELECT email, nid, \"CreatedAt\", subToken FROM jetpen.temp_subscription WHERE id = $1"
	row := db.QueryRow(sqlString, id)
	sub := new(models.Subscription)
	sub.Id = id
	err := row.Scan(&sub.Email, &sub.Nid, &sub.CreatedAt, &sub.SubToken)
	failOnError(err, "Getting Subscription")
	return sub, err
}

func GetTempSubscriptionWithEmail(nid string, email string, db *sql.DB) (*models.Subscription, error) {
	sqlString := "SELECT id, \"CreatedAt\", subToken FROM jetpen.temp_subscription WHERE nid=$1 AND email=$2"
	row := db.QueryRow(sqlString, nid, email)
	sub := new(models.Subscription)
	sub.Email = email
	sub.Nid = nid
	err := row.Scan(&sub.Id, &sub.CreatedAt, &sub.SubToken)
	failOnError(err, "Getting Subscription")
	return sub, err
}

func GetSubscriptionWithEmail(nid string, email string, db *sql.DB) (*models.Subscription, error) {
	sqlString := "SELECT id, \"CreatedAt\", subToken FROM jetpen.subscription WHERE nid=$1 AND email=$2"
	row := db.QueryRow(sqlString, nid, email)
	sub := new(models.Subscription)
	sub.Email = email
	sub.Nid = nid
	err := row.Scan(&sub.Id, &sub.CreatedAt, &sub.SubToken)
	failOnError(err, "Getting Subscription")
	return sub, err
}
