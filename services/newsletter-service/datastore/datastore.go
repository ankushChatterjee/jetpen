package datastore

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ankushChatterjee/jetpen/newsletter-service/models"
	"github.com/ankushChatterjee/jetpen/newsletter-service/pkg/utils"
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

func decodeCursor(encodedCursor string) (res time.Time, uuid string, err error) {
	byt, err := base64.StdEncoding.DecodeString(encodedCursor)
	if err != nil {
		return
	}

	arrStr := strings.Split(string(byt), ",")
	if len(arrStr) != 2 {
		err = errors.New("cursor is invalid")
		return
	}

	res, err = time.Parse(time.RFC3339Nano, arrStr[0])
	if err != nil {
		return
	}
	uuid = arrStr[1]
	return
}

func encodeCursor(t time.Time, uuid string) string {
	key := fmt.Sprintf("%s,%s", t.Format(time.RFC3339Nano), uuid)
	return base64.StdEncoding.EncodeToString([]byte(key))
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

func InsertNewsletter(newsletter *models.Newsletter, db *sql.DB) error {
	sqlString := "INSERT INTO jetpen.newsletter(name, description, owner) VALUES($1,$2,$3)"
	_, err := db.Exec(sqlString, newsletter.Name, newsletter.Description, newsletter.Owner)
	failOnError(err, "Error inserting newsletter")
	return err
}

func InsertLetter(letter *models.Letter, db *sql.DB) error {
	sqlString := "INSERT INTO jetpen.letter(subject, owner, nid, content,\"isPublished\",\"PublishedAt\") values($1, $2, $3, $4, $5, $6)"
	_, err := db.Exec(sqlString, letter.Subject, letter.Owner, letter.Nid, letter.Content, letter.IsPublished, letter.PublishedAt.Time)
	failOnError(err, "Error Inserting letter")
	return err
}

func GetNewsLettersForUser(owner string, cursor string, limit int, db *sql.DB) ([]*models.Newsletter, *string, error) {
	newsletters := make([]*models.Newsletter, 0)
	sqlString := ""
	var rows *sql.Rows
	if len(cursor) == 0 {
		sqlString = "SELECT name, description, \"CreatedAt\",id from jetpen.newsletter where owner=$1 ORDER BY \"CreatedAt\" DESC LIMIT $2"
		var err error
		rows, err = db.Query(sqlString, owner, limit)
		failOnError(err, "Error inserting newsletter")
		if err != nil {
			return nil, nil, err
		}
	} else {
		createdAt, _, err := decodeCursor(cursor)
		failOnError(err, "Error Inserting letter")
		if err != nil {
			return nil, nil, err
		}
		sqlString = "SELECT name, description, \"CreatedAt\",id from jetpen.newsletter WHERE \"CreatedAt\" <= $1 AND owner=$2 ORDER BY \"CreatedAt\" DESC LIMIT $3"
		rows, err = db.Query(sqlString, createdAt, owner, limit)
		failOnError(err, "Error inserting newsletter")
		if err != nil {
			return nil, nil, err
		}
	}
	defer rows.Close()

	retCursor := ""
	if rows == nil {
		log.Println("Rows is nil")
	}
	for rows.Next() {
		newsletter := new(models.Newsletter)
		err := rows.Scan(&newsletter.Name, &newsletter.Description, &newsletter.CreatedAt, &newsletter.Id)
		failOnError(err, "Error Scanning Rows")
		if err != nil {
			return nil, nil, err
		}
		newsletter.Owner = owner
		newsletters = append(newsletters, newsletter)
		retCursor = encodeCursor(newsletter.CreatedAt, newsletter.Id)
	}
	return newsletters, &retCursor, nil
}

func GetLettersForNewsletter(nid string, cursor string, limit int, db *sql.DB) ([]*models.Letter, *string, error) {
	letters := make([]*models.Letter, 0)
	sqlString := ""
	var rows *sql.Rows
	if len(cursor) == 0 {
		sqlString = "SELECT id, subject, \"CreatedAt\", \"PublishedAt\", \"isPublished\" FROM jetpen.letter WHERE nid=$1 ORDER BY \"CreatedAt\" DESC LIMIT $2"
		var err error
		rows, err = db.Query(sqlString, nid, limit)
		failOnError(err, "Error inserting newsletter")
		if err != nil {
			return nil, nil, err
		}
	} else {
		createdAt, _, err := decodeCursor(cursor)
		failOnError(err, "Error Inserting letter")
		if err != nil {
			return nil, nil, err
		}
		log.Println(createdAt)
		sqlString = "SELECT id, subject, \"CreatedAt\", \"PublishedAt\", \"isPublished\" FROM jetpen.letter WHERE \"CreatedAt\" <= $1 AND nid=$2 ORDER BY \"CreatedAt\" DESC LIMIT $3"
		rows, err = db.Query(sqlString, createdAt, nid, limit)
		failOnError(err, "Error inserting newsletter")
		if err != nil {
			return nil, nil, err
		}
	}
	defer rows.Close()
	retCursor := ""
	for rows.Next() {
		letter := new(models.Letter)
		err := rows.Scan(&letter.Id, &letter.Subject, &letter.CreatedAt, &letter.PublishedAt, &letter.IsPublished)
		failOnError(err, "Error Scanning Rows")
		if err != nil {
			return nil, nil, err
		}
		letter.Nid = nid
		letters = append(letters, letter)
		retCursor = encodeCursor(letter.CreatedAt, letter.Id)
	}
	return letters, &retCursor, nil

}

func GetLetterContent(id string, db *sql.DB) (string, error) {
	sqlString := "SELECT content from jetpen.letter WHERE id = $1"
	content := ""
	row := db.QueryRow(sqlString, id)
	err := row.Scan(&content)
	return content, err
}

func IsLetterPublished(id string, db *sql.DB) (bool, error) {
	sqlStringSelect := "SELECT \"isPublished\" FROM  jetpen.letter WHERE id=$1"
	row := db.QueryRow(sqlStringSelect, id)
	var isPublished bool
	err := row.Scan(&isPublished)
	failOnError(err, "Error getting if letter is published")
	if err != nil {
		return true, err
	}
	return isPublished, nil
}

func DeleteLetter(id string, db *sql.DB) error {
	sqlString := "DELETE FROM jetpen.letter WHERE id = $1"
	_, err := db.Exec(sqlString, id)
	failOnError(err, "Error deleting letter")
	return err
}

func DeleteNewsletter(id string, db *sql.DB) error {
	sqlString := "DELETE FROM jetpen.newsletter WHERE id = $1"
	_, err := db.Exec(sqlString, id)
	failOnError(err, "Error Delete Newsletter")
	return err
}

func UpdateLetterContentAndSubject(id string, subject string, content string, db *sql.DB) error {
	sqlString := "UPDATE jetpen.letter SET content=$1, subject=$2 WHERE id=$3"
	_, err := db.Exec(sqlString, content, subject, id)
	failOnError(err, "Error updating Letter subject")
	return err
}

func UpdateLetterContentSubjectAndPublish(id string, subject string, content string, db *sql.DB) error {
	sqlString := "UPDATE jetpen.letter SET content=$1, subject=$2, \"PublishedAt\"=$3, \"isPublished\"=true WHERE id=$4"
	_, err := db.Exec(sqlString, content, subject, time.Now().UTC(), id)
	failOnError(err, "Error updating Letter subject")
	return err
}

func UpdateNewsletterName(id string, name string, db *sql.DB) error {
	sqlStringUpdate := "UPDATE jetpen.newsletter SET name=$1 WHERE id=$2"
	_, err := db.Exec(sqlStringUpdate, name, id)
	failOnError(err, "Error updating newsletter name")
	return err
}

func UpdateNewsletterDescription(id string, description string, db *sql.DB) error {
	sqlStringUpdate := "UPDATE jetpen.newsletter SET description=$1 WHERE id=$2"
	_, err := db.Exec(sqlStringUpdate, description, id)
	failOnError(err, "Error updating newsletter description")
	return err
}

func GetSubEmails(nid string, db *sql.DB) ([]string, error) {
	sqlString := "SELECT email FROM jetpen.subscription WHERE nid=$1"
	rows, err := db.Query(sqlString, nid)
	failOnError(err, "Error On get subscriptions")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	emails := make([]string, 0)
	for rows.Next() {
		email := ""
		err := rows.Scan(&email)
		failOnError(err, "Error Reading email of subscriber")
		if err != nil {
			return nil, err
		}
		emails = append(emails, email)
	}
	return emails, nil
}

func GetNewsLetterOwner(nid string, db *sql.DB) (string, error) {
	sqlString := "SELECT owner from jetpen.newsletter WHERE id=$1"
	row := db.QueryRow(sqlString, nid)
	owner := ""
	switch err := row.Scan(&owner); err {
	case sql.ErrNoRows:
		failOnError(err, "No newsletter Exist")
		return "", err
	case nil:
		return owner, nil
	default:
		failOnError(err, "Error getting newsletter owner")
		return "", err
	}
}

func GetLetterOwner(nid string, db *sql.DB) (string, error) {
	sqlString := "SELECT owner from jetpen.letter WHERE id=$1"
	row := db.QueryRow(sqlString, nid)
	owner := ""
	switch err := row.Scan(&owner); err {
	case sql.ErrNoRows:
		failOnError(err, "Error  letter owner")
		return "", err
	case nil:
		return owner, nil
	default:
		failOnError(err, "Error getting letter owner")
		return "", err
	}
}
