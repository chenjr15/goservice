package prog

import (
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	yaml "github.com/goccy/go-yaml"
	"github.com/jander/golog/logger"
	"github.com/kardianos/service"
)

// * main() -> service.Run()-> prog.Start() -> prog.run()
type Program struct {
	ServiceName string
	Basedir     string   `json:"basedir,omitempty"`
	Bin         string   `json:"bin,omitempty"`
	Args        []string `json:"args,omitempty"`
	stop        chan struct{}
	proc        *exec.Cmd
}

func (p *Program) Start(s service.Service) error {
	// 系统启动服务回调
	logger.Println("[Start] start running")
	// 异步启动任务然后直接返回
	go p.run(s)
	return nil
}

func (p *Program) run(s service.Service) {
	// 内部逻辑
	logger.Println("[Run] in running")

	if p.Bin == "" {
		time.Sleep(time.Second * 2)
		err := s.Stop()
		logger.Fatalf("no binpath, stop", err)
		os.Exit(-1)
		return
	}
	p.proc = exec.Command(p.Bin, p.Args...)
	if p.proc == nil {
		return
	}
	p.proc.Dir = p.Basedir
	logger.Printf("[Run] exec: %v \n", p.proc)
	if err := p.proc.Start(); err != nil {

		logger.Fatalf("Fail to run %v:%v\n", p.Bin, err)
	}
	p.proc.Wait()
	logger.Println("[Run] proc exicted.")
	// notify system
	s.Stop()

}

func (p *Program) Stop(s service.Service) error {

	logger.Println("[Stop] Stop running")
	if p.proc != nil {
		p.proc.Process.Kill()
	}
	return nil
}
func (p *Program) LoadConfig(configPath string) (err error) {
	var bytes []byte
	bytes, err = ioutil.ReadFile(configPath)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(bytes, p)
	return
}
