package main

func main() {
	InitConfig()

	request := NewRequest(syncdCfg.access)
	request.Login()

	_, _ = request.Projects()

	//request := NewRequest(*syncdCfg)
	//if err := request.Login(); err != nil {
	//	panic(err)
	//}
}
