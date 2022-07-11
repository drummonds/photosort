package main

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

var SampleFiles = [...]string{
	"2020/2020-02/2020-02-18/sample2.JPG",
	"2021/2021-08/2021-08-06/sample3.jpg",
	"2021/2021-08/2021-08-06/sample4.jpg",
	"2021/2021-08/2021-08-12/20210812_092905.dng",
	"2021/2021-08/2021-08-14/Sample_1.jpg",
}

func TestNosParam(t *testing.T) {
	var err error
	cmd := exec.Command("go", "run", ".")
	if err = cmd.Run(); err != nil {
		t.Log(err)
		return
	}
	t.Fail()
}

func IsDestFile(name string) bool {
	_, err := os.Stat(name)
	if err != nil {
		return false
	}
	return true
}

func TestParam(t *testing.T) {
	var err error
	err = os.RemoveAll("./test/sorted")
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	// TODO test dir actually works
	cmd := exec.Command("go", "run", ".", "-s", "./test/test1", "-d", "./test/sorted")
	if err = cmd.Run(); err != nil {
		t.Log(err)
		t.Fail()
	}
	for i, fn := range SampleFiles {
		if !IsDestFile("./test/sorted/" + fn) {
			t.Log(fmt.Sprintf("File %v absent = %s, %v", i, fn, err))
			t.Fail()
		}
	}
}
