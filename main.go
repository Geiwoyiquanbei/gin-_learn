package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
)

type Todo struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Status bool   `json:"status"`
}

func main() {

	dsn := "root:suxianjin@tcp(127.0.0.1:3306)/teett?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	db.AutoMigrate(&Todo{})

	r := gin.Default()
	r.Static("/static", "static") //访问静态css文件
	r.LoadHTMLGlob("templates/*") //加载模板html
	r.GET("index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	v1Group := r.Group("v1")
	{
		// 待办事项
		// 添加
		v1Group.POST("/todo", func(c *gin.Context) {
			// 前端页面填写待办事项 点击提交 会发请求到这里
			// 1. 从请求中把数据拿出来
			var todo Todo
			c.BindJSON(&todo)
			// 2. 存入数据库
			//err = DB.Create(&todo).Error
			//if err!= nil {
			//}
			// 3. 反回响应
			if err := db.Create(&todo).Error; err != nil {
				c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusOK, todo)
				//c.JSON(http.StatusOK, gin.H{
				//	"code": 2000,
				//	"msg": "success",
				//	"data": todo,
				//})
			}
		})
		v1Group.GET("/todo", func(c *gin.Context) {
			// 查询todo这个表里的所有数据
			var todoList []Todo
			if err := db.Find(&todoList).Error; err != nil {
				c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusOK, todoList)
			}
		})
		v1Group.GET("/todo/id", func(c *gin.Context) {

		})
		// 修改某一个待办事项
		v1Group.PUT("/todo/id", func(c *gin.Context) {
			id, ok := c.Params.Get("id")
			if !ok {
				c.JSON(http.StatusOK, gin.H{"error": "无效的id"})
				return
			}
			var todo Todo
			if err := db.Where("id=?", id).First(&todo).Error; err != nil {
				c.JSON(http.StatusOK, gin.H{"error": err.Error()})
				return
			}
			c.BindJSON(&todo)
			if err := db.Save(&todo).Error; err != nil {
				c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusOK, todo)
			}

		})
		// 删除某一个待办事项
		v1Group.DELETE("/todo/id", func(c *gin.Context) {
			id, ok := c.Params.Get("id")
			if !ok {
				c.JSON(http.StatusOK, gin.H{"error": "无效的id"})
				return
			}
			if err := db.Where("id=?", id).Delete(Todo{}).Error; err != nil {
				c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusOK, gin.H{id: "deleted"})
			}
		})

	}
	r.Run()
}
