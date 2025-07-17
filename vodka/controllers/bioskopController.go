package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

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
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Something went wrong, we are still figuring out what's the problem",
		})
	}

	var newBioskop Bioskop

	if err := ctx.ShouldBindJSON(&newBioskop); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

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

func GetBioskop(ctx *gin.Context) {
	db, err := ConnectDB()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Something went wrong, we are still figuring out what's the problem",
		})
		return
	}

	var results = []Bioskop{}

	sqlStatement := `SELECT * FROM Bioskop`

	rows, err := db.Query(sqlStatement)

	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var bioskop = Bioskop{}

		err = rows.Scan(&bioskop.ID, &bioskop.Nama, &bioskop.Lokasi, &bioskop.Rating)

		if err != nil {
			panic(err)
		}

		results = append(results, bioskop)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"bioskop": results,
	})
}

func GetBioskopById(ctx *gin.Context) {
	id := ctx.Param("id")

	bioskopID, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	db, err := ConnectDB()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Something went wrong, we are still figuring out what's the problem",
		})
		return
	}

	var result = Bioskop{}

	sqlStatement := `SELECT * FROM Bioskop WHERE Id = $1`

	row := db.QueryRow(sqlStatement, bioskopID)

	err = row.Scan(&result.ID, &result.Nama, &result.Lokasi, &result.Rating)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "Data not found",
			})
		} else {
			ctx.AbortWithError(http.StatusBadRequest, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"bioskop": result,
	})
}

func UpdateBioskop(ctx *gin.Context) {
	id := ctx.Param("id")

	bioskopID, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	db, err := ConnectDB()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Something went wrong, we are still figuring out what's the problem",
		})
		return
	}

	var updateBioskop Bioskop

	if err := ctx.ShouldBindJSON(&updateBioskop); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	sqlStatement := `SELECT Id FROM Bioskop WHERE Id = $1`

	err = db.QueryRow(sqlStatement, bioskopID).Scan(&updateBioskop.ID)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "Data not found",
			})
		} else {
			ctx.AbortWithError(http.StatusBadRequest, err)
		}
		return
	}

	sqlStatement = `
	UPDATE Bioskop 
	SET Nama = $1,
		Lokasi = $2,
		Rating = $3
	WHERE Id = $4
	Returning *
	`

	err = db.QueryRow(
		sqlStatement,
		updateBioskop.Nama,
		updateBioskop.Lokasi,
		updateBioskop.Rating,
		updateBioskop.ID,
	).Scan(
		&updateBioskop.ID,
		&updateBioskop.Nama,
		&updateBioskop.Lokasi,
		&updateBioskop.Rating,
	)

	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"message": "Updated successfully",
	})
}

func DeleteBioskop(ctx *gin.Context) {
	id := ctx.Param("id")

	bioskopID, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	db, err := ConnectDB()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Something went wrong, we are still figuring out what's the problem",
		})
		return
	}

	var deleteBioskop Bioskop

	sqlStatement := `SELECT Id FROM Bioskop WHERE Id = $1`

	err = db.QueryRow(sqlStatement, bioskopID).Scan(&deleteBioskop.ID)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "Data not found",
			})
		} else {
			ctx.AbortWithError(http.StatusBadRequest, err)
		}
		return
	}

	sqlStatement = `
	DELETE FROM Bioskop 
	WHERE Id = $1
	Returning *
	`

	res, err := db.Exec(sqlStatement, deleteBioskop.ID)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
	}

	_, err = res.RowsAffected()
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"message": "Deleted successfully",
	})
}
