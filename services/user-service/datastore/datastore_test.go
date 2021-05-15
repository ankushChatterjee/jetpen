package datastore

import (
	"fmt"
	"testing"

	"github.com/ankushChatterjee/jetpen/user-service/models"
)

func TestConnectAndExecute(t *testing.T) {
	db := GetConnection()
	defer db.Close()
	fmt.Println("Connected to DB")
	user := models.User{}
	user.Username = "ankush"
	user.Name = "Ankush Chatterjee"
	user.Email = "ac.ankush15@gmail.com"
	user.Password = "xxas"

	usr, _ := GetUser("ankush", db)
	fmt.Println(usr)
}
