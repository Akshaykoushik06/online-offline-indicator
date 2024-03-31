package main

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func init() {
	_db, err := sql.Open("mysql", "asdasd")
	if err != nil {
		panic(err)
	}
	DB = _db
}

func main() {
	ge := gin.Default()

	ge.POST("/heartbeats", func(ctx *gin.Context) {
		data := map[string]interface{}{}
		ctx.Bind(&data)
		if _, err := DB.Exec("REPLACE INTO oo_heartbeats (user_id, last_hb) values (?, ?);", data["user_id"], time.Now().Unix()); err != nil {
			panic(err)
		}

		ctx.JSON(200, map[string]interface{}{"message": "ok"})
	})

	ge.GET("/heartbeats/status/:user_id", func(ctx *gin.Context) {
		var lastHB int
		row := DB.QueryRow("SELECT last_hb from oo_heartbeats WHERE user_id = ?;", ctx.Param("user_id"))
		row.Scan(&lastHB)
		ctx.JSON(200, map[string]interface{}{"is_online": lastHB > int(time.Now().Unix()-30)})
	})

	ge.GET("/heartbeats/status", func(ctx *gin.Context) {
		rows, err := DB.Query("SELECT user_id, last_hb FROM oo_heartbeats WHERE user_id IN (?);", ctx.Query("user_ids"))
		if err != nil {
			panic(err)
		}

		var statusMap map[string]bool = make(map[string]bool)
		var userID, lastHB int
		for rows.Next() {
			if err := rows.Scan(&userID, &lastHB); err != nil {
				panic(err)
			}
			statusMap[fmt.Sprintf("%d", userID)] = lastHB > int(time.Now().Unix())-30
		}
		rows.Close()
		ctx.JSON(200, statusMap)
	})

	ge.GET("/heartbeats/status_nonop", func(ctx *gin.Context) {
		var statusMap map[string]bool = make(map[string]bool)
		for _, userID := range strings.Split(ctx.Query("user_ids"), ",") {
			var lastHB int
			row := DB.QueryRow("SELECT last_hb FROM oo_heartbeats WHERE user_id = ?;", userID)
			row.Scan(&lastHB)
			statusMap[userID] = lastHB > int(time.Now().Unix())-30
		}
		ctx.JSON(200, statusMap)
	})

	ge.Run(":9000")
}

//  go mod init connection-pooling
