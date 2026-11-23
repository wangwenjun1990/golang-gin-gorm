package controllers

import (
	"exchangeapp/global"
	"exchangeapp/models"
	"exchangeapp/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

/**
	1、go语言中if的写法：

	if 条件 {
    	条件成立时执行的代码
	}

	2、ctx.ShouldBindJSON(&user)
	表示 读取 HTTP 请求体中的 JSON，并将其解析到指定的 Go 变量中，同时执行适用的绑定校验。

*/

// Register 处理注册请求。ctx 是 Gin 传入的指针，用来读取请求和发送响应。
func Register(ctx *gin.Context) {
	// 声明一个 User 结构体，用来接收请求里的用户名、密码等字段。
	var user models.User

	// 把请求体中的 JSON 填进 user。&user 取出地址，让 Gin 能修改这个变量。
	// if 里的 := 声明局部变量 err，分号后检查 err。err != nil 表示发生错误。
	if err := ctx.ShouldBindJSON(&user); err != nil {
		// 返回 HTTP 400。gin.H 用来构造 JSON 对象，err.Error() 取得错误文本。
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 用 bcrypt 计算密码哈希。函数返回两个值，分别存到 hashedPwd 和 err。
	hashedPwd, err := utils.HashPassword(user.Password)

	// 检查密码哈希是否计算失败。
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 将 user 中的密码换成哈希值，稍后写入数据库。
	user.Password = hashedPwd

	// 用用户名生成签名后的 JWT。:= 声明新变量 token，并给已有的 err 重新赋值。
	token, err := utils.GenerateJWT(user.Username)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	/**
		最重要的一句话：函数名前面的括号表示“这个方法属于谁”，函数名后面的括号表示“调用这个方法需要传入什么参数”。

		func (接收者) 方法名(参数) 返回值 {
	    	// 方法体
		}
	*/

	// global.Db 是项目的 GORM 数据库对象。AutoMigrate 根据 User 模型创建或更新表结构。
	if err := global.Db.AutoMigrate(&user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	/**
	这里的 .Error 之所以可以直接使用，是因为 GORM 的 Create() 方法返回的是一个 *gorm.DB 对象，而这个对象里面有一个叫 Error 的字段。
	注意：.Error 是访问结构体字段，不是调用方法。
	*/

	// Create 插入一条用户记录，&user 让 GORM 可以回填生成的 ID。
	if err := global.Db.Create(&user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 用户记录插入成功后，返回 HTTP 200，并把 JWT 放进 JSON 的 token 字段。
	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

func Login(ctx *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User

	if err := global.Db.Where("username = ?", input.Username).First(&user).Error; err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "wrong credentials"})
		return
	}

	if !utils.CheckPassword(input.Password, user.Password) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "wrong credentials"})
		return
	}

	token, err := utils.GenerateJWT(user.Username)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}
