package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type ToDo struct {
	Id        string `json:"id"`
	Todo      string `json:"todo"`
	Completed bool   `json:"completed"`
}

func start() *sql.DB {
	db, err := sql.Open("mysql", MYSQL_URI)

	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func main() {
	db := start()
	defer db.Close()
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello world!",
		})
	})
	router.GET("/todos", func(c *gin.Context) {
		getTodos(c, db)
	})
	router.POST("/todos", func(c *gin.Context) {
		addTodo(c, db)
	})
	router.PATCH("/todos/:id",func(c *gin.Context) {
		updateTodo(c, db)
	} )
	router.DELETE("/todos/:id", func(c *gin.Context) {
		deleteTodo(c, db)
	})

	router.Run("localhost:8080")
}

func getTodos(c *gin.Context, db *sql.DB) {
	var (
		id        string
		todo      string
		completed bool
	)
	rows, err := db.Query("SELECT * FROM todos")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	list := []ToDo{}
	for rows.Next() {
		if err := rows.Scan(&id, &todo, &completed); err != nil {
			log.Fatal(err)
		}
		list = append(list, ToDo{id, todo, completed})
	}
	c.JSON(http.StatusOK, gin.H{
		"todos": list,
	})
}
func addTodo(c *gin.Context, db *sql.DB) {
	id := c.PostForm("id")
	todo := c.PostForm("todo")
	stmt, err := db.Prepare("INSERT INTO todos VALUES(?,?,?)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(id,todo,false)
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusCreated, gin.H{
		"msg": "Added Todo",
	})
}
func updateTodo(c *gin.Context,db *sql.DB) {
	id := c.Param("id")
	value, _ := strconv.ParseBool(c.Query("value"))
	stmt,err := db.Prepare("UPDATE todos SET completed =? WHERE id=?")
	if err!=nil{
		log.Fatal(err)
	}
	if _,err = stmt.Exec(value,id); err!=nil{
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "Updated todo " + c.Param("id"),
	})
}
func deleteTodo(c *gin.Context,db *sql.DB) {
	id := c.Param("id")
	stmt,err := db.Prepare("DELETE FROM todos WHERE id=?")
	if err!=nil{
		log.Fatal(err)
	}
	if _,err = stmt.Exec(id); err!=nil{
		log.Fatal(err)
	}	
	c.JSON(http.StatusOK, gin.H{
		"msg": "Updated todo " + c.Param("id"),
	})
}
