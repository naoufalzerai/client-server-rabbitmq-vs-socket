package traitement

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Encaissement struct {
	NumPiece string
	NumFac   string
	Mnt      string
	Caissier string
}

func TraitementSuperComplique(dataJson []byte) error {

	data := new(Encaissement)

	err := json.Unmarshal(dataJson, &data)
	if err != nil {
		return err
	}

	log.Println("Traitement de : ", data)

	//init database
	username := "postgres"
	password := "postgres"
	dbName := "postgres"
	dbHost := "localhost"

	dbUri := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, username, dbName, password) //Build connection string

	conn, err := gorm.Open("postgres", dbUri)
	if err != nil {
		return err
	}

	conn.LogMode(false)
	gorm.DefaultTableNameHandler = func(dbVeiculosGorm *gorm.DB, defaultTableName string) string {
		return "somei." + defaultTableName
	}

	db = conn
	db.AutoMigrate(&Encaissement{}) //Database migration

	db.Create(&data)

	return nil
}

var db *gorm.DB //database

//returns a handle to the DB object
func GetDB() *gorm.DB {
	return db
}
