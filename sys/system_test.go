package sys

import "testing"

func TestOpen(t *testing.T) {
	err := open("E:\\MyProjects\\golang\\devTools\\sys\\fs.exe")
	t.Log(err)
}

func TestMkdir(t *testing.T) {
	mkdir("open/file")
}
