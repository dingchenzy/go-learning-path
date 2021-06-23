# 习题

```go
	// 修改echo程序，使其能够打印os.Args[0]，即被执行命令本身的名字。
	fmt.Println(strings.Join(os.Args[1:], " "))
	
	// 修改echo程序，使其打印每个参数的索引和值，每个一行。
	for i, v := range os.Args[:] {
		fmt.Println(i, v)
	}
```

```go
// 出现重复的行时打印文件名称。
package main

import (
	"bufio"
	"fmt"
	"os"
)

type file struct {
	filename string
	count int
}

func main() {
	filevalue := make(map[string]*file)
	for _, v := range os.Args[1:] {
		f, err := os.OpenFile(v, os.O_RDWR, os.ModePerm)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		count(f, v, filevalue)
		defer f.Close()
	}
	for i, v := range filevalue {
		if v.count >1 {
			fmt.Println(i, "filename is :", v.filename)
		}
	}
}

func count(file2 *os.File, filename1 string, filevalue map[string]*file) {
	scanner := bufio.NewScanner(file2)
	for scanner.Scan() {
		_, ok := filevalue[scanner.Text()]
		if ok {
			filevalue[scanner.Text()].count++
		} else {
			filevalue[scanner.Text()]= new(file)
			filevalue[scanner.Text()].count=1
			filevalue[scanner.Text()].filename = filename1
		}
	}
}
```
