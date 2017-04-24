package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

func main() {
	for i := 1; i <= 15; i++ {
		srcPkg := "k8s.io/kubernetes"
		srcPkgB := []byte(srcPkg)
		dstPkg := "github.com/sourcegraph/monorepo-test-1/kubernetes-" + strconv.Itoa(i)
		dstPkgB := []byte(dstPkg)
		srcDir := os.Getenv("GOPATH") + "/src/" + srcPkg
		dstDir := os.Getenv("GOPATH") + "/src/" + dstPkg

		fmt.Printf("%s -> %s\n", srcPkg, dstPkg)

		os.RemoveAll(dstDir)

		filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() && info.Name()[0] == '.' {
				return filepath.SkipDir
			}
			if info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
				return nil
			}

			data, err := ioutil.ReadFile(path)
			if err != nil {
				panic(err)
			}

			data = bytes.Replace(data, srcPkgB, dstPkgB, -1)

			dstFile := dstDir + path[len(srcDir):]
			os.MkdirAll(filepath.Dir(dstFile), 0777)
			if err := ioutil.WriteFile(dstFile, data, info.Mode()); err != nil {
				panic(err)
			}

			return nil
		})
	}
}
