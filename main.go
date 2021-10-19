package main

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)
type postJson struct {
	Type string `json:"type" binding:"required"`
	Username string `json:"user" binding:"required"`
	Password string `json:"password" binding:"required"`//前端应先判断信息是否完整

}

type User struct {
	gorm.Model
	Username string
	Pwd string
}

func main() {
	var user User
	db, errConnectDatabase := gorm.Open(mysql.Open("root:020103@tcp(127.0.0.1:3306)/usersinformation?parseTime=true"), &gorm.Config{})
	if errConnectDatabase != nil {
		panic("failed to connect database")
	}
	// 迁移 schema
	db.AutoMigrate(&User{})


	r := gin.Default()
	r.POST("/loginData/", func(context *gin.Context) {
		var postData postJson
		var msg string

		// Create
		//db.Create(&User{Username: "admin", Pwd: "123456"})

		errJson := context.ShouldBindJSON(&postData)
		if postData.Type=="signIn"{

			if errJson != nil {
				fmt.Print("errJson")
			}
			db.First(&user, "username = ?", postData.Username)
			serverPwd := user.Pwd

			if serverPwd != postData.Password {
				msg = "账户名或密码错误"
			} else {
				msg = "欢迎"
			}
		}else if postData.Type=="signUp" {
			errNameIsNotExisting :=(db.First(&user,"username = ?", postData.Username)).Error
			fmt.Print(errNameIsNotExisting)
			if errors.Is(errNameIsNotExisting, gorm.ErrRecordNotFound){
				db.Create(&User{Username: postData.Username, Pwd: postData.Password})
				db.Commit()
				msg = "注册成功"
			}else {msg = "用户名已存在"}
		}

		context.JSON(200, gin.H{"msg": msg, "time": time.Now()})
		//context.String(200, "\n")
		fmt.Print(db)
	})
	errWebRun := r.Run()
	if errWebRun != nil {
		fmt.Print("failed run")
	}
}
