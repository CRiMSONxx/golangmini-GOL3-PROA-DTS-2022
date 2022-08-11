package main

import (
	"net/http"
    "strconv"

	"github.com/gin-gonic/gin"
)

var db = make(map[string]string)
// Task ...
type Task struct {
	Id_task int
    Isi  string
    Dikerjakan_oleh string
	Deadline string
	Status int
}

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()
	r.LoadHTMLGlob("html/*")

	//dummy
	tasks := make([]Task, 0)
    tasks = append(tasks, Task{
		Id_task: 0,
        Isi:  "Fitur Menambahkan Task (Create)",
        Dikerjakan_oleh: "P1",
        Deadline: "10 Agustus 2022",
        Status: 0,
    })
    tasks = append(tasks, Task{
		Id_task: 1,
        Isi:  "Fitur Membaca Task (Read)",
        Dikerjakan_oleh: "P2",
        Deadline: "11 Desember 2022",
        Status: 1,
    })
    tasks = append(tasks, Task{
		Id_task: 2,
        Isi:  "Fitur Mengedit Task (Update)",
        Dikerjakan_oleh: "P1",
        Deadline: "12 Oktober 2022",
        Status: 5,
    })
	

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "tasks_list.html", gin.H{
            "tasks": tasks,
        })
	})
	

	r.GET("/edit/:id", func(c *gin.Context) {
		var idtask int
		idtask, err := strconv.Atoi(c.Params.ByName("id"))
		
		if err == nil {
			tugas := tasks[idtask]
			c.HTML(http.StatusOK, "tasks_edit.html", gin.H{
				"tugas_single": tugas,
			})
		} else {
			c.HTML(http.StatusOK, "error.html", gin.H{
				
			})
		}
	})

	// ======================================
	// default start
	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// Get user value
	r.GET("/user/:name", func(c *gin.Context) {
		user := c.Params.ByName("name")
		value, ok := db[user]
		if ok {
			c.JSON(http.StatusOK, gin.H{"user": user, "value": value})
		} else {
			c.JSON(http.StatusOK, gin.H{"user": user, "status": "no value"})
		}
	})

	// Authorized group (uses gin.BasicAuth() middleware)
	// Same than:
	// authorized := r.Group("/")
	// authorized.Use(gin.BasicAuth(gin.Credentials{
	//	  "foo":  "bar",
	//	  "manu": "123",
	//}))
	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
		"foo":  "bar", // user:foo password:bar
		"manu": "123", // user:manu password:123
	}))

	/* example curl for /admin with basicauth header
	   Zm9vOmJhcg== is base64("foo:bar")

		curl -X POST \
	  	http://localhost:8080/admin \
	  	-H 'authorization: Basic Zm9vOmJhcg==' \
	  	-H 'content-type: application/json' \
	  	-d '{"value":"bar"}'
	*/
	authorized.POST("admin", func(c *gin.Context) {
		user := c.MustGet(gin.AuthUserKey).(string)

		// Parse JSON
		var json struct {
			Value string `json:"value" binding:"required"`
		}

		if c.Bind(&json) == nil {
			db[user] = json.Value
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		}
	})
	//default end
	return r
}

func main() {
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
