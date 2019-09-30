package main

import (
	"flag"
	"fmt"
	z "github.com/nutzam/zgo"
	"runtime/debug"
	"strconv"
	"time"
)

func help() {
	fmt.Print(`usage:
	syncd-console submit test-admin-server -m ''
	syncd-console projects
	syncd-console tasks`)
}

func Recover() {
	if r := recover(); r != nil {
		fmt.Println("执行错误：", r)
		fmt.Println(string(debug.Stack()))
	}
}

func ParseFlag(flags []string) (map[string]string, error) {
	ret := make(map[string]string)
	var k string
	flagArr := []string{"-p", "-m"}
	for _, f := range flags {
		if z.IndexOfStrings(flagArr, f) != -1 {
			k = f
		} else {
			ret[k] = f
		}
	}

	return ret, nil
}

/*
app submit -p test-admin-server -m ''
app projects
app tasks
*/
func main() {
	defer Recover()
	InitConfig()

	flag.Parse()
	//fmt.Printf("%v",flag.Args())
	//println(flag.NArg(),flag.Arg(0))

	var cmd string
	if flag.NArg() == 0 {
		cmd = "help"
	} else {
		cmd = flag.Arg(0)
	}
	switch cmd {
	case "login":
		request := NewRequest(syncdCfg.access)
		request.Login()
		fmt.Println("登录成功")
	case "submit":
		//检查参数 -p -m
		//fmt.Printf("%v", flag.Args()[1:])
		params, err := ParseFlag(flag.Args()[1:])
		if err != nil {
			panic(err)
		}

		if params["-p"] == "" || params["-m"] == "" {
			panic("参数错误,请输入 -p project_name -m description")
		}

		//fmt.Printf("%v", params)
		request := NewRequest(syncdCfg.access)
		err = request.Submit(params["-p"], params["-m"], params["-m"])
		if err != nil {
			panic("任务提交失败")
		}

		time.Sleep(time.Second * 1)

		//读取任务列表，找到任务id
		respData := request.ApplyList()
		list := respData["list"]
		var taskId int
		for _, v := range list.([]interface{}) {
			//fmt.Printf("%v", v)
			username := v.(map[string]interface{})["username"].(string)
			projectname := v.(map[string]interface{})["project_name"].(string)
			id := int(v.(map[string]interface{})["id"].(float64)) //任务id
			status := int(v.(map[string]interface{})["status"].(float64))

			if username == syncdCfg.access.Username && projectname == params["-p"] && status == TASK_STATUS_WAIT {
				taskId = id
				break;
			}
			//fmt.Println(z.AlignLeft(strconv.Itoa(id), 10, ' '), z.AlignLeft(projectname, 80, ' '), z.AlignLeft(username, 80, ' '), status)
		}

		if taskId == 0 {
			panic("未找到任务")
		}

		var build, deploy chan int
		build = make(chan int)
		deploy = make(chan int)

		defer func() {
			close(build)
			close(deploy)
		}()

		//build
		go func(taskId int) {
			fmt.Println("开始构建")
			err := request.BuildStart(taskId)
			if err != nil {
				panic("构建启动失败:" + err.Error())
			}

			for {
				select {
				case <-time.After(time.Second * 30):
					panic("构建超时，请重试")
				default:
					status := request.BuildStatus(taskId)
					switch status {
					case BUILD_STATUS_ERROR:
						panic("构建出错")
					case BUILD_STATUS_DONE:
						build <- 1
					case BUILD_STATUS_RUNNING:
						fmt.Print('.')
					}

					time.Sleep(time.Second * 2)
				}
			}
		}(taskId)
		<-build
		//构建结束，开始部署

		go func(taskId int) {
			fmt.Println("开始部署")
			err := request.DeployStart(taskId)
			if err != nil {
				panic("部署启动失败")
			}

			for {
				select {
				case <-time.After(time.Second * 30):
					panic("部署超时，请重试")
				default:
					status := request.DeployStatus(taskId)
					switch status {
					case DEPLOY_STATUS_DONE:
						deploy <- 1
					case DEPLOY_STATUS_FAIL:
						panic("部署失败")
					case DEPLOY_STATUS_RUNNING:
						fmt.Print('.')
					}

					time.Sleep(time.Second * 2)
				}
			}
		}(taskId)
		<-deploy
		println("部署成功！")

	case "projects":
		request := NewRequest(syncdCfg.access)
		projectJson := request.Projects()
		projects := NewProjects(projectJson)
		fmt.Printf("%s - %s\n", z.AlignLeft("项目名称", 40, ' '), "空间名称")
		for _, v := range projects.data {
			fmt.Printf("%s - %s\n", z.AlignLeft(v.ProjectName, 40, ' '), v.SpaceName)
		}
	case "tasks":
		request := NewRequest(syncdCfg.access)
		respData := request.ApplyList()
		list := respData["list"]
		fmt.Println(z.AlignLeft("任务id", 10, ' '), z.AlignLeft("项目名称", 40, ' '), z.AlignLeft("用户", 30, ' '), "状态")
		for _, v := range list.([]interface{}) {
			username := v.(map[string]interface{})["username"].(string)
			projectname := v.(map[string]interface{})["project_name"].(string)
			id := int(v.(map[string]interface{})["id"].(float64))
			status := int(v.(map[string]interface{})["status"].(float64))

			fmt.Println(z.AlignLeft(strconv.Itoa(id), 10, ' '), z.AlignLeft(projectname, 40, ' '), z.AlignLeft(username, 30, ' '), GetTaskStatusText(status))
		}
	case "test":
		request := NewRequest(syncdCfg.access)
		println(request.BuildStatus(10887))
	case "help":
	default:
		help()
	}
}
