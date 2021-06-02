# k8s-client

是 k8s 官方提供的操作 k8s 资源的 golang 库，可以调用其中的方法实现对 k8s 资源的 crud。

分为类型化的客户端和动态客户端。

GVR 是用来生成 REST API 的，GVK 用来声明 yaml 文件的定义。

GVK 与 GVR 中是有映射关系的，这种映射也被称为 REST 映射。

kubernetes api 服务器也提供了一个发现 api 可以查找所有可用的 REST 映射 `kubectl api-resources` 可以看到 REST 之间的映射关系。

## 类型客户端

使用 Go 结构来表示，可以使用类型安全的方式编辑资源，可以自动找到 REST 映射以发送 API 请求。

### 将 deployment 内容转换为 json

```go
package main

import (
	"encoding/json"
	"os"

	appsv1 "k8s.io/api/apps/v1"
)

func main() {
	deployment := &appsv1.Deployment{}

	enc := json.NewEncoder(os.Stdout)
	// 添加换行以及缩进
	// 将数据流转换为 json 的数据
	enc.Encode(deployment)
}
```

### 转换原理

```go
package main

import (
	"encoding/json"
	"os"
)

type teststruct1 struct {
	Address string `json: address`
}

type teststruct struct {
	Name string      `json:"Name"`
	Id   string      `json:"ID"`
	Add  teststruct1 `json:"add"`
}

func main() {
	var user *teststruct
	user = &teststruct{
		Id:   "1",
		Name: "chen",
		Add: teststruct1{
			Address: "山东",
		},
	}

	enc.SetIndent("", "    ")
	enc.Encode(user)
}


/*
{
    "Name": "chen",
    "ID": "1",
    "add": {
        "Address": "山东"
    }
}
*/
```

## 动态客户端

不使用 `k8s.io/api` 中定义的 Go 类型，而是使用 `unstrcutured.Unstructured`。非结构化。

`"k8s.io/apimachinery/pkg/runtime/serializer/yaml"` 包用来将 yaml 的数据解码编译到 json 的数据化中