package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type users struct {
	gorm.Model
	NAME  string `gorm:"not null" json:"name" binding:"required,min=2"`
	EMAIL string `gorm:"not null" json:"email" binding :"required,email"`
}

var db *gorm.DB
var err1 error

func main() {
	dsn := "host=localhost user=postgres password=1234 dbname=ecommerce_db port=5432 sslmode=disable"
	db, err1 = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err1 != nil {
		log.Println("error during connecting the db", err1)
	}
	if err := db.AutoMigrate(&users{}); err != nil {
		log.Println("auto migration fail")
		return
	}
	fmt.Println("auto migrate connected")

	r:=gin.Default()
	v1:=r.Group("/v1")
	{
		v1.GET("/",GetUsers)
		v1.POST("/",CreateUser)
		v1.PUT("/:id",updateuser)
		v1.DELETE("/:id",deleteUser)
	}

	r.Run(":8080")
}
func GetUsers(c *gin.Context){
	var user users
	result:=db.Find(&user)
	if result.Error != nil{
		c.JSON(400,gin.H{
			"message":result.Error.Error(),
		})
		c.JSON(200,user)
	}

}

func CreateUser (c *gin.Context){
	var user users
	if err:=c.ShouldBindJSON(&user); err !=nil{
		c.JSON(400,gin.H{"error":"something not right on the input"})
	}
	result:=db.Create(&user)
	if result.Error !=nil {
		c.JSON(400,gin.H{"message":"something not right"})
	}
	c.JSON(200,&user)
}

func updateuser (c *gin.Context){
	var user,input users
	id:=c.Param("id")
	result:=db.Find(&user,id)
	if result.Error != nil{
		c.JSON(400,gin.H{"message":"user not found"})
	}
	if err:=c.ShouldBindJSON(&input); err != nil{
		c.JSON(400,gin.H{"error":"error"})
	}
	results:=db.Model(&user).Updates(&input)
	if results.Error != nil{
		c.JSON(400,results.Error.Error())
	}
	c.JSON(200,input)
}

func deleteUser(c *gin.Context){
	id:=c.Param("id")
	result:=db.Delete(&users{},id)
	if result.Error != nil{
		c.JSON(400,gin.H{"message":"user not found"})
	}
	c.JSON(200,gin.H{"message":"user deleted"})
}




