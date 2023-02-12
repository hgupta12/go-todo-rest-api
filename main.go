package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)


type ToDo struct{
	Id string `gorm:"primaryKey"`
	Todo string
	Completed bool
}

func start() *gorm.DB {
	db, err := gorm.Open(mysql.Open(MYSQL_URI), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func main() {
	db := start()

	db.AutoMigrate(&ToDo{})

	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello world!",
		})
	})
	router.GET("/todos", func(c *gin.Context) {
		getTodos(c,db)
	})
	router.POST("/todos", func(c *gin.Context) {
		addTodo(c,db)
	})
	router.PATCH("/todos/:id",func(c *gin.Context) {
		updateTodo(c,db)
	} )
	router.DELETE("/todos/:id", func(c *gin.Context) {
		deleteTodo(c,db)
	})

	router.Run("localhost:8080")
}

func getTodos(c *gin.Context, db *gorm.DB) {
	var todos []ToDo
	db.Find(&todos)
	c.JSON(http.StatusOK, gin.H{
		"todos": todos,
	})
}
func addTodo(c *gin.Context, db *gorm.DB) {
	id := c.PostForm("id")
	todo := c.PostForm("todo")
	db.Create(&ToDo{id,todo,false})
	c.JSON(http.StatusCreated, gin.H{
		"msg": "Added Todo",
	})
}
func updateTodo(c *gin.Context,db *gorm.DB) {
	id := c.Param("id")
	value, _ := strconv.ParseBool(c.Query("value"))
	var todo ToDo

	db.First(&todo,"id=?",id)
	db.Model(&todo).Update("completed",value)

	c.JSON(http.StatusOK, gin.H{
		"msg": "Updated todo " + c.Param("id"),
	})
}
func deleteTodo(c *gin.Context,db *gorm.DB) {
	id := c.Param("id")
	var todo ToDo
	db.Delete(&todo,"id=?",id)
	c.JSON(http.StatusOK, gin.H{
		"msg": "Deleted todo " + c.Param("id"),
	})
}
