# k8s-client

是 k8s 官方提供的操作 k8s 资源的 golang 库，可以调用其中的方法实现对 k8s 资源的 crud。

分为类型化的客户端和动态客户端。

GVR 是用来生成 REST API 的，GVK 用来声明 yaml 文件的定义。

GVK 与 GVR 中是有映射关系的，这种映射也被称为 REST 映射。

kubernetes api 服务器也提供了一个发现 api 可以查找所有可用的 REST 映射 `kubectl api-resources` 可以看到 REST 之间的映射关系。

## 类型客户端

使用 Go 结构来表示，可以使用类型安全的方式编辑资源，可以自动找到 REST 映射以发送 API 请求。

最常用的Go客户端库位于`k8s.io/client-go`软件包中。该软件包依赖于`k8s.io/api`和`k8s.io/apimachinery`，`k8s.io/api`是各种结构的集合，而`k8s.io/apimachinery`实现GVK，GVR和其他实用程序。

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

`"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"` 该包将允许非结构化的数据进行解析

使用 `unstructred.Unstructured` 非结构化之后可以通过调用这个模块实现非结构化的使用。可以适用到 CRD 资源之上，因为这种非结构化的方式可以兼容所有的资源。

与类型化的客户端不同的是需要向动态客户端提供 GVR，结构化类型的客户端只需要提供 GVK 即可。

```go
package main

import (
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	// k8s.io/apimachinery/pkg/runtime 库作用就是将泛型语句转换为结构化的数据
	// 就是通过解析将 yaml 的无类型语句解析为结构化的数据
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
)

const dsManifest = `
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: example
  namespace: default
spec:
  selector:
    matchLabels:
      name: nginx-ds
  template:
    metadata:
      labels:
        name: nginx-ds
    spec:
      containers:
      - name: nginx
        image: nginx:latest
`

func main() {
	obj := &unstructured.Unstructured{}

	fmt.Println(obj)
	// 将 yaml 解码支持添加到支持 json 序列化中
	dec := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	// _, gvk, _ := dec.Decode([]byte(dsManifest), nil, obj)
	// decode 解码，会将数据解析为 map 格式
	dec.Decode([]byte(dsManifest), nil, obj)

	// fmt.Println(obj.GetName(), gvk.String())
	fmt.Println(obj.GetNamespace(), obj.GetName())

	// enc := json.NewEncoder(os.Stdout)
	// enc.SetIndent("", "    ")
	// enc.Encode(obj)
}
```

