package main

import (
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
	return err == nil
}

func copyForTest1(t *testing.T, source string, destination string, move bool) {
	var (
		err error
		cmd *exec.Cmd
	)
	err = os.RemoveAll(destination)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	if move {
		cmd = exec.Command("go", "run", ".", "-s", source, "-d", destination, "-m")
	} else {
		cmd = exec.Command("go", "run", ".", "-s", source, "-d", destination)
	}
	if err = cmd.Run(); err != nil {
		t.Log(err)
		t.Fail()
	}
	for i, fn := range SampleFiles {
		ffn := destination + "/" + fn
		if !IsDestFile(ffn) {
			t.Logf("File %v absent = %s, %v", i, ffn, err)
			t.Fail()
		}
	}
}

func TestParam(t *testing.T) {
	copyForTest1(t, "./test/test1", "./test/sorted", false)
}

func TestMove(t *testing.T) {
	copyForTest1(t, "./test/test1", "./test/test2", false)
	copyForTest1(t, "./test/test2", "./test/sorted", true)
	var err error
	// make sure move has worked
	for i, fn := range SampleFiles {
		ffn := "./test/test2" + "/" + fn
		if IsDestFile(ffn) {
			t.Logf("File %v has not been moved %s, %v", i, ffn, err)
			t.Fail()
		}
	}
}
