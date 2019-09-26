package main

import (
	"fmt"
	"os"
	"testing"
)

func TestRange(t *testing.T)  {
	strs := []string{"1","2","3"}
	param := make([]string, 5)
	for _,v := range strs{
		param = append(param, v)
	}

	fmt.Println(param)
}

func TestEnv(t  *testing.T){
	err := os.Setenv("TEST_KEY", "iamstring")
	if err !=nil{
		t.Error(err)
	}

	env :=os.Getenv("TEST_KEY")
	if env != "iamstring" {
		t.Fail()
	}
}
