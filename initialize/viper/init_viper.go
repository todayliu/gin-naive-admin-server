package viper

import (
	"flag"
	"fmt"
	"gin-admin-server/global"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func Viper(path ...string) *viper.Viper {
	var config string
	if len(path) == 0 {
		flag.StringVar(&config, "r", "", "choose config file.")
		flag.Parse()
		if config == "" { // 判断命令行参数是否为空
			if configEnv := os.Getenv(ConfigEnv); configEnv == "" {
				switch gin.Mode() {
				case gin.DebugMode:
					config = ConfigDefaultFile
					fmt.Printf("您正在使用gin模式的%s环境名称,config的路径为%s\n", gin.EnvGinMode, ConfigDefaultFile)
				case gin.ReleaseMode:
					config = ConfigReleaseFile
					fmt.Printf("您正在使用gin模式的%s环境名称,config的路径为%s\n", gin.EnvGinMode, ConfigReleaseFile)
				case gin.TestMode:
					config = ConfigTestFile
					fmt.Printf("您正在使用gin模式的%s环境名称,config的路径为%s\n", gin.EnvGinMode, ConfigTestFile)
				}
			} else {
				config = configEnv
				fmt.Printf("您正在使用%s环境变量,config的路径为%s\n", ConfigEnv, config)
			}
		} else {
			fmt.Printf("您正在使用命令行的-r参数传递的值,config的路径为%s\n", config)
		}
	} else {
		config = path[0]
		fmt.Printf("您正在使用func Viper()传递的值,config的路径为%s\n", config)
	}

	v := viper.New()
	v.SetConfigFile(config)
	v.SetConfigType("yaml")
	err := v.ReadInConfig()

	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		if err = v.Unmarshal(&global.GNA_CONFIG); err != nil {
			fmt.Println(err)
		}
	})
	if err = v.Unmarshal(&global.GNA_CONFIG); err != nil {
		panic(err)
	}

	// root 适配性 根据root位置去找到对应迁移位置,保证root路径有效
	//global.GNA_CONFIG.AutoCode.Root, _ = filepath.Abs("..")
	return v
}
