package config

import (
	"log"

	"github.com/spf13/viper"
)

// 1、
// 这里是定义一个叫 config的结构体，类似java的实体类
// 这里结构体内嵌套结构体，类似Java 中的对象套对象

// 2、
// &：我要地址。
// *：我要指针指向的实际内容（用于表达式时）。
// *Config：表示指针类型。

// 3、go中创建对象的方法
//Java	                     Go	                含义
//User user = new User()	user := User{}	    创建结构体
//User user = new User()	user := &User{}	    创建结构体并取得指针
//User user = null	        var user *User	    只有引用/指针变量，没有创建结构体
//user.getName()	        user.Name	        访问字段

type Config struct {
	App struct {
		Name string
		Port string
	}
	Database struct {
		Dsn          string
		MaxIdleConns int
		MaxOpenConns int
	}
	Redis struct {
		Addr     string
		DB       int
		Password string
	}
}

// 声明一个名为 AppConfig 的变量，它的类型是 *Config，也就是指向 Config 结构体的指针。

var AppConfig *Config

func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	// &Config{} 会创建一个 Config 结构体，并取它的地址，得到一个类型为 *Config 的指针。
	// & - 是取地址符
	AppConfig = &Config{}

	if err := viper.Unmarshal(AppConfig); err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
	}

	initDB()
	initRedis()
}
