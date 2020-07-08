package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	_ "github.com/mattn/go-sqlite3"
)

type Todo struct {
	gorm.Model  
	Text   string
	Status string
}

func dbInit() {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("データベース開ません!!(dbInit)")
	}
	db.AutoMigrate(&Todo{})
	defer db.Close()
}

func dbInsert(text string, status string) {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("データベース開ません!!!(dbInsert)")
  }
  db.Create(&Todo{Text: text, Status: status})
	defer db.Close()
}

func dbGetAll() []Todo {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("データベース開ません!!!(dbGetAll())")
	}
	var todos []Todo
	db.Order("created_at desc").Find(&todos)
	db.Close()
	return todos
}

func dbGetOne(id int) Todo {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("データベース開ません!!!(dbGetOne())")
	}
	var todo Todo
	db.First(&todo, id)
	db.Close()
	return todo
}

func dbUpdate(id int, text string, status string) {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("データベース開ません!!!(dbUpdate)")
	}
	var todo Todo
	db.First(&todo, id)
	todo.Text = text
	todo.Status = status
	db.Save(&todo)
	db.Close()
}

func dbDelete(id int) {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("データベース開ません!!!(dbDelete)")
	}
	var todo Todo
	db.First(&todo, id)
	db.Delete(&todo)
	db.Close()
}

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*.html")

	dbInit()

	router.GET("/", func(ctx *gin.Context) {
		todos := dbGetAll()
		ctx.HTML(200, "index.html", gin.H{
			"todos": todos,
		})
	})

	router.GET("/new", func(ctx *gin.Context) {
		text := ctx.PostForm("text")
		status := ctx.PostForm("status")
		dbInsert(text, status)
		ctx.Redirect(302, "/")
	})

	router.GET("/detail/:id", func(ctx *gin.Context) {
		n := ctx.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic(err)
		}
		todos := dbGetOne(id)
		ctx.HTML(200, "index.html", gin.H{
			"todos": todos,
		})
  })
  
  router.POST("/update/:id", func(ctx *gin.Context) {
    n := ctx.Param("id")
    id, err := strconv.Atoi(n)
    if err != nil {
      panic("ERROR")
    }
    text := ctx.PostForm("text")
    status := ctx.PostForm("status")
    dbUpdate(id, text, status)
    ctx.Redirect(302, "/")
  })

  router.GET("/delete_check/:id",func(ctx *gin.Context){
    n := ctx.Param("id")
    id, err := strconv.Atoi(n)
    if err != nil {
      panic("ERROR")
    }
    todo := dbGetOne(id)
    ctx.HTML(200, "delete.html", gin.H{
      "todo": todo,
    })
  })

  router.POST("/delete/:id", func(ctx *gin.Context) {
    n := ctx.Param("id")
    id, err := strconv.Atoi(n)
    if err != nil {
      panic("ERROR")
    }
    dbDelete(id)
    ctx.Redirect(302, "/")
  })
	router.Run()
}
