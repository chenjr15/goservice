package main

import (
	"flag"
	"fmt"
	"goservice/prog"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"

	"github.com/goccy/go-yaml"
	"github.com/jander/golog/logger"
	"github.com/kardianos/service"
)

var u, err = user.Current()
var homeDir string = u.HomeDir
var baseDir string = filepath.Join(homeDir, ".goservice")

var action string
var configPath string
var loggerPath string
var program = &prog.Program{
	ServiceName: "GoServiceWrapper",
	Basedir:     "C:\\",
	Args:        []string{},
	Bin:         "example.exe",
}

const LogFileName = "goservice.log"
const ConfigFileName = "goservice.yml"

/**
* MAIN函数，程序入口,
* main() -> service.Run()-> prog.Start() -> prog.run()
 */

func main() {

	os.Chdir(baseDir)

	svcConfig := &service.Config{
		Name:        program.ServiceName,                 //服务显示名称
		DisplayName: program.ServiceName + "[GoService]", //服务名称
		Description: fmt.Sprintf("%v  [basedir]:%v \t[exe]:%v %v,",
			configPath, program.Basedir, program.Bin, program.Args), //服务描述
		Arguments: []string{
			"-config", configPath,
		},
	}

	s, err := service.New(program, svcConfig)
	if err != nil {
		logger.Fatal(err)
	}

	if action == "install" {
		err = s.Install()
		if err != nil {
			logger.Fatalf("服务%v安装失败！%v \n", svcConfig.DisplayName, err)
		}
		logger.Printf("服务%v安装成功", svcConfig.DisplayName)
		return
	} else if action == "uninstall" {
		err = s.Uninstall()
		if err != nil {
			logger.Fatalf("服务%v卸载失败！%v\n", program.ServiceName, err)
		}
		logger.Printf("服务%v卸载成功！", program.ServiceName)
		return
	}

	// 一直等到操作系统给停止服务的信号才会退出
	err = s.Run()
	if err != nil {
		logger.Error(err)
	}
	logger.Println("[main] System required stop")
}

func init() {

	configPath = filepath.Join(baseDir, ConfigFileName)
	flag.StringVar(&action, "action", "run", "动作：run,install,uninstall")
	flag.StringVar(&configPath, "config", configPath, "服务配置文件, 建议和需要的exe放一起")
	flag.Parse()
	absPath, err := filepath.Abs(configPath)
	if err != nil {
		logger.Fatalln("无法配置绝对路径：", configPath)
	}
	configPath = absPath
	if err := program.LoadConfig(configPath); err != nil {
		logger.Warn("无法加载配置文件，将写入默认配置。%v", err)
		data, _ := yaml.Marshal(program)
		ioutil.WriteFile(configPath, data, 0777)
	}
	baseDir = program.Basedir
	os.MkdirAll(baseDir, 0777)
	// 用path.join 会出问题，永远都是 / 而不会在windows下变成 \
	loggerPath = filepath.Join(baseDir, LogFileName)

	logger.SetHandlers(logger.NewFileHandler(loggerPath), logger.NewConsoleHandler())
	logger.Println("loggerPath", loggerPath)
	logger.Println("configPath", configPath)

}
