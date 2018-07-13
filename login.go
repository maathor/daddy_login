package main

import (
	"fmt"
	"time"
	"github.com/jinzhu/gorm"
	"github.com/gin-gonic/gin"
	 _ "github.com/go-sql-driver/mysql"
)

//struct for database

type User struct {
	ID        			int        `sql:"AUTO_INCREMENT" gorm:"primary_key"`
	Firstname      		string     `sql:"size:255"`
	Lastname			string	   `sql:"size:255"`
	Password			string	   `sql:"size:255;index"`
	CreationDate		time.Time
	LastModification	time.Time 	`sql:"DEFAULT:current_timestamp"`
	PasswordReset		[]PasswordReset
	Address				Address
	Role				Role
	FacebookID			string	   `sql:"size:255;index"`
	GoogleID			string	   `sql:"size:255;index"`
	Deleted          	bool       `sql:"DEFAULT:false"`
}

type PasswordReset struct {
	Token			string	   `sql:"size:255;index"`
	ExpirationDate	time.Time `sql:"DEFAULT:current_timestamp"`
}

type Address struct {
	Country  string `gorm:"primary_key"`
	City     string `gorm:"primary_key"`
	PostCode string `gorm:"primary_key"`
	Line1    string `sql:"size:255"`
	Line2    string `sql:"size:255"`
}


type Role struct {
	Name			string	   `sql:"size:255;index"`
	CreationDate 	time.Time  `sql:"DEFAULT:current_timestamp"`
}


// globals vars
var db *gorm.DB
var err error

func main() {
	db, err = gorm.Open("mysql", "root:bjorn@/login?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	db.AutoMigrate(&User{}, &PasswordReset{}, &Address{}, &Role{})


	r := gin.Default()
	r.GET("/users/", GetUsers)
	r.GET("/user/:id", GetUser)
	r.POST("/user", CreateUser)
	r.PUT("/user/:id", UpdateUser)
 	r.DELETE("/people/:id", DeleteUser)

	r.Run(":9000")
}

func DeleteUser(c *gin.Context) {
	id := c.Params.ByName("id")
	var user User
	d := db.Where("id = ?", id).Delete(&user)
	fmt.Println(d)
	c.JSON(200, gin.H{"id #" + id: "deleted"})
}

func GetUser(c *gin.Context) {
	id := c.Params.ByName("id")
	var user User
	 if err := db.Where("id = ?", id).First(&user).Error; err != nil {
	    c.AbortWithStatus(404)
	    fmt.Println(err)
	 } else {
	    c.JSON(200, user)
	 }
}

func GetUsers(c *gin.Context) {
	var users []User
	 if err := db.Find(&users).Error; err != nil {
	    c.AbortWithStatus(404)
	    fmt.Println(err)
	 } else {
	    c.JSON(200, users)
	 }
}

func CreateUser(c *gin.Context) {
	var user User
	c.BindJSON(&user)
	db.Create(&user)
	c.JSON(200, user)
}

func UpdateUser(c *gin.Context) {
	var user User
	id := c.Params.ByName("id")
	if err := db.Where("id = ?", id).First(&user).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	}
	c.BindJSON(&user)
	db.Save(&user)
	c.JSON(200, user)
}