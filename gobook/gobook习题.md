# 习题

## 第一章

### 	// 修改echo程序，使其能够打印os.Args[0]，即被执行命令本身的名字。

```go
	// 修改echo程序，使其能够打印os.Args[0]，即被执行命令本身的名字。
	fmt.Println(strings.Join(os.Args[1:], " "))
	
	// 修改echo程序，使其打印每个参数的索引和值，每个一行。
	for i, v := range os.Args[:] {
		fmt.Println(i, v)
	}
```

### // 出现重复的行时打印文件名称。


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

### // 函数调用io.Copy(dst, src)会从src中读取内容，并将读到的结果写入到dst中，使用这个函数替代掉例子中的ioutil.ReadAll来拷贝响应结构体到os.Stdout，避免申请一个缓冲区（例子中的b）来存储。记得处理io.Copy返回结果中的错误。


```go
// 函数调用io.Copy(dst, src)会从src中读取内容，并将读到的结果写入到dst中，使用这个函数替代掉例子中的ioutil.ReadAll来拷贝响应结构体到os.Stdout，避免申请一个缓冲区（例子中的b）来存储。记得处理io.Copy返回结果中的错误。
// 修改一下题目，从网上下载一个图片，如果图片过大肯定会导致内存占用过多，这种情况使用 io.Copy 方式可以
package main

import (
	"io"
	"net/http"
	"os"
)

func main() {
	resp, err := http.Get("https://www.twle.cn/static/i/img1.jpg")
	if err != nil {
		panic(err)
	}
	file, err := os.Create("imagefile.jpg")
	if err != nil {
		panic(err)
	}
	io.Copy(file, resp.Body)
	defer resp.Body.Close()
}
```

### // 修改fetch打印出HTTP协议的状态码，可以从resp.Status变量得到该状态码。


```go
// 修改fetch打印出HTTP协议的状态码，可以从resp.Status变量得到该状态码。
package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	for _, v := range os.Args[1:] {

		resp, err := http.Get(v)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		str, err := io.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		fmt.Println(string(str))
		fmt.Println(resp.StatusCode)
	}
}
```

### // 修改fetch这个范例，如果输入的url参数没有 http:// 前缀的话，为这个url加上该前缀。你可能会用到strings.HasPrefix这个函数。

```go
// 修改fetch这个范例，如果输入的url参数没有 http:// 前缀的话，为这个url加上该前缀。你可能会用到strings.HasPrefix这个函数。
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func main() {
	for _, v := range os.Args[1:] {
		if !strings.HasPrefix(v, "http") {
			v = "http://" + v
		}
		resp, err := http.Get(v)
		if err != nil {
			panic(err)
		}
		str, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		fmt.Println(string(str))
	}
}
```

### // 找一个数据量比较大的网站，用本小节中的程序调研网站的缓存策略，对每个URL执行两遍请求，查看两次时间是否有较大的差别，并且每次获取到的响应内容是否一致，修改本节中的程序，将响应结果输出，以便于进行对比。


```go
// 找一个数据量比较大的网站，用本小节中的程序调研网站的缓存策略，对每个URL执行两遍请求，查看两次时间是否有较大的差别，并且每次获取到的响应内容是否一致，修改本节中的程序，将响应结果输出，以便于进行对比。
package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// 协程获取多个 URL 内容
func main() {
	start := time.Now()
	ch := make(chan string)
	for _, v := range os.Args[1:] {
		if !strings.HasPrefix(v, "http") {
			v = "http://" + v
		}
		go fetch(v,ch)
	}
	for _, _ = range os.Args[1:] {
		fmt.Println(<-ch)
	}
	fmt.Printf("%.2f\n", time.Since(start).Seconds())
}

func fetch(url string, ch chan<- string) {
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprint(err)
	}

	nbytes, err := io.Copy(io.Discard, resp.Body)
	defer resp.Body.Close()
	if err != nil {
		ch <- fmt.Sprint(err)
	}
	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%.2fs  %7d  %s", secs, nbytes, url)
}





/*
访问结果
0.28s    70172  https://m.gome.com.cn/
0.51s   301278  http://www.baidu.com
0.66s   113616  http://www.taobao.com
1.48s   122773  http://www.jd.com
1.48




0.09s   301246  http://www.baidu.com
0.13s   124586  http://www.jd.com
0.32s   113616  http://www.taobao.com
0.49s    70173  https://m.gome.com.cn/
0.49
*/



// 在fetchall中尝试使用长一些的参数列表，比如使用在alexa.com的上百万网站里排名靠前的。如果一个网站没有回应，程序将采取怎样的行为？
// 答：因为是使用了 channel 通道的机制，所以当访问的网站没有回应时，那么程序会一直在等待网站的响应，直到访问时间超时，随后会因为访问超时被 panic 掉
```

### //修改Lissajour服务，从URL读取变量，比如你可以访问 http://localhost:8000/?cycles=20 这个URL，这样访问可以将程序里的cycles默认的5修改为20。字符串转换为数字可以调用strconv.Atoi函数。你可以在godoc里查看strconv.Atoi的详细说明。


```go
//修改Lissajour服务，从URL读取变量，比如你可以访问 http://localhost:8000/?cycles=20 这个URL，这样访问可以将程序里的cycles默认的5修改为20。字符串转换为数字可以调用strconv.Atoi函数。你可以在godoc里查看strconv.Atoi的详细说明。

package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

var mu = sync.Mutex{}

var countint int

type handle struct {
	cou int
}

func (h *handle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error
	mu.Lock()
	countint++
	mu.Unlock()
	if value := r.FormValue("cycles") ; value != "" {
		h.cou, err = strconv.Atoi(value)
		if err != nil {
			log.Print(err)
		}
	}
	w.Write([]byte(fmt.Sprintf("你最棒了，加油哦！, cou value is ：%s", strconv.Itoa(h.cou))))
}

func count(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	fmt.Fprintf(w, "countint is ：%s", countint)
	mu.Unlock()
}

func main() {
	handle := &handle{
		cou: 20,
	}
	http.Handle("/", handle)
	http.HandleFunc("/count", count)

	log.Fatal(http.ListenAndServe("127.0.0.1:9000", nil))
}
```

### // 编写一个程序，默认情况下打印标准输入的SHA256编码，并支持通过命令行flag定制，输出SHA384或SHA512哈希算法。编写一个程序，默认情况下打印标准输入的SHA256编码，并支持通过命令行flag定制，输出SHA384或SHA512哈希算法。


```go
// 编写一个程序，默认情况下打印标准输入的SHA256编码，并支持通过命令行flag定制，输出SHA384或SHA512哈希算法。编写一个程序，默认情况下打印标准输入的SHA256编码，并支持通过命令行flag定制，输出SHA384或SHA512哈希算法。

package main

import (
	"crypto/sha256"
	"crypto/sha512"
	"flag"
	"fmt"
	"log"
)
```

### //编写一个程序，默认情况下打印标准输入的SHA256编码，并支持通过命令行flag定制，输出SHA384或SHA512哈希算法。

```go
//编写一个程序，默认情况下打印标准输入的SHA256编码，并支持通过命令行flag定制，输出SHA384或SHA512哈希算法。

func getargs(args []string) ([]byte, bool) {
	if len(args) == 0 {
		return nil, false
	}
	return []byte(args[0]), true
}

func main() {
	var option string
	var help bool
	flag.StringVar(&option, "X", "sha256", "default is sha256, option is sha384 or sha512")
	flag.BoolVar(&help, "h", false, "help")
	flag.Usage = func() {
			fmt.Println("usage flagargs [-X sha384|sha512]")
			flag.PrintDefaults()
	}
	flag.Parse()
	value,err  := getargs(flag.Args())
	if !err {
		log.Print("value is nil")
	}
	if option == "sha256" {
		byte := sha256.Sum256(value)
		fmt.Printf("%v\n",string(byte[:]))
	} else if option == "sha384" {
		byte := sha512.Sum384(value)
		fmt.Printf("%v\n",string(byte[:]))
	} else if option == "sha512" {
		byte := sha512.Sum512(value)
		fmt.Printf("%#v\n",string(byte[:]))
	} else {
		flag.Usage()
	}
}
```


## slice

### // 反转数组中的元素

```go
// 反转数组中的元素
package main

import "fmt"

func reverse(sliceint []int) []int {
	for i, j := 0 ,len(sliceint) -1; i<j; i, j = i +1, j-1{
		sliceint[i], sliceint[j] = sliceint[j], sliceint[i]
	}
	return sliceint
}

func main() {
	fmt.Println(reverse([]int{1,2,3,4,5,6}))
}
```

### // 测试 slice 是否相同


```go
// 测试 slice 是否相同
package main

import (
	"fmt"
)

func equal(x, y []string) bool {
	if len(x) != len(y) {
		return false
	}

	for i:=0;i<len(x);i++ {
		if x[i] != y[i] {
			return false
		}
	}
	return true
}

func main() {
	x := []string{"a", "b"}
	y := []string{"a", "b"}
	if !equal(x,y) {
		fmt.Println("array is not equal")
	}
}
```

### // append 函数重写


```go
// append 函数重写
package main

import "fmt"

func appendint(x []int, y int) []int {
	var z []int
	// x 的长度并 + 1
	zlen := len(x) +1
	// 判断 cap 容量是否可以容纳新的容量
	if cap(x) >= len(x)+1 {
		// 如果可以容纳，那么将会将数组切片并赋值给 z
		z = x[:zlen]
	} else {
		// 如果不能容纳，则会创建新的可以容纳 +1 元素的新容量
		zcap := zlen
		// 扩容
		if zcap < 2*len(x) {
			zcap = 2 * len(x)
		}
		// 对 z 进行扩容
		z= make([]int, zlen, zcap)
		// 将 x slice 中的数据 copy 到 z 中，默认追加到最前
		copy(z, x)
	}
	// 判断 x 的长度并赋值
	z[len(x)] = y
	return z
}

func main() {
	var x, y []int
	for i := 0; i< 10; i++{
		y= appendint(x, i)
		fmt.Printf("%d  cap=%d\t%v\n", i, cap(y), y)
		x= y
	}
}
```


### // 返回不为空的值，并使用底层的同一个 slice，减少内存使用


```go
// 返回不为空的值，并使用底层的同一个 slice，减少内存使用
package main

import "fmt"

// 传入 slice，使用 for range 判断值是否有空，不为空则保存，如果为空则忽略
func nonempty(x []int) []int {
	y := 0
	for i, v := range x {
		if v != 0 {
			// 这里不可以写 i，因为如果使用 i，则判断为 0 时，那么索引也不会被调用到
			//x[i] = v
			x[y] = v
			y++
		}
	}
	return x[:y]
}

func main() {
	x := []int{1,2,0,4,5}
	fmt.Println(nonempty(x))
}
```


### // 删除中间的切片


```go
// 删除中间的切片
package main

import "fmt"

// 需要顺序的删除，使用 copy 进行 copy 后然后返回取出最后一个元素的切片
func sortdelete(x []int, y int) []int {
	copy(x[y:], x[y+1:])
	return x[:len(x)-1]
}

// 不需要顺序的删除，直接覆盖最后一个
func nosortdelete(x []int, i int) []int {
	x[i] = x[len(x)-1]
	return x[:len(x)-1]
}

func main() {
	var x = []int{1,2,3,4,5}

	fmt.Println(nosortdelete(x,3))
	var y = []int{1,2,3,4,5}

	fmt.Println(sortdelete(y, 2))
}
```

### // 统计文本中的单词数量


```go
// 统计文本中的单词数量
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	countvalue := make(map[string]int)

	file, err := os.Open("e.txt")
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		countvalue[scanner.Text()]++
	}
	for i, v := range countvalue {
		fmt.Printf("value is %s\tcount is %d\n", i,v)
	}
}
```