package controllers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Bioskop struct {
	ID     int
	Nama   string
	Lokasi string
	Rating float64
}

var BioskopDatas = []Bioskop{}

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "echo"
	dbname   = "formative13"
)

func ConnectDB() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to Database")
	return db, nil
}

func CreateBioskop(ctx *gin.Context) {
	db, err := ConnectDB()
	if err != nil {
		panic(err)
	}

	var newBioskop Bioskop

	if err := ctx.ShouldBindJSON(&newBioskop); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// newBioskop.ID, _ = strconv.Atoi(fmt.Sprintf("c%d", len(BioskopDatas)+1))
	// BioskopDatas = append(BioskopDatas, newBioskop)

	sqlStatement := `
	INSERT INTO Bioskop (Nama, Lokasi, Rating) VALUES
	($1, $2, $3)
	Returning *
	`

	err = db.QueryRow(
		sqlStatement,
		newBioskop.Nama,
		newBioskop.Lokasi,
		newBioskop.Rating,
	).Scan(
		&newBioskop.ID,
		&newBioskop.Nama,
		&newBioskop.Lokasi,
		&newBioskop.Rating,
	)

	if err != nil {
		panic(err)
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"bioskop": newBioskop,
	})

}
