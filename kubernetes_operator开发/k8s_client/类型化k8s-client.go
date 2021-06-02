package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"path/filepath"

	// 存储这个 api 资源组下特有的资源内容
	appsv1 "k8s.io/api/apps/v1"
	// core/v1 核心组的资源，例如 node
	v1 "k8s.io/api/core/v1"

	// 存储所有资源共有的元数据
	metadatav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	// 是 client-go 的核心库，创建 go 与 kubernetes 的连接
	"k8s.io/client-go/kubernetes"
	// 解析 config 文件
	"k8s.io/client-go/tools/clientcmd"
	// 检测自己的 home 目录下是否有 config 文件
	"k8s.io/client-go/util/homedir"
	// 辅助工具，例如在更新时需要重试的操作，就可以使用 RetryOnConflict() 方法实现
	"k8s.io/client-go/util/retry"
)

func getnode(clientset *kubernetes.Clientset) {
	retry.RetryOnConflict()
	nodeclient, err := clientset.CoreV1().Nodes().List(context.TODO(), metadatav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range nodeclient.Items {
		fmt.Printf("%s\t%s\t%s\t%s\n", "Status.Addresses", "v.Name", "v.Kind", "v.Status.NodeInfo.OSImage")
		fmt.Printf("%s\t%s\t%s\t%s\n", v.Status.Addresses[0].Address, v.Name, v.Kind, v.Status.NodeInfo.OSImage)
	}
}

func getdeployment(clientset *kubernetes.Clientset, namespace string) {
	deploymentclient, err := clientset.AppsV1().Deployments(namespace).List(context.TODO(), metadatav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}
	if len(deploymentclient.Items) == 0 {
		log.Fatal("this is ", namespace, "no resources")
	}
	for _, v := range deploymentclient.Items {
		fmt.Println("v.namespace", "v.kind", "v.name")
		fmt.Println(v.Namespace, v.Kind, v.Name)
	}
}

func createdeployment(clientset *kubernetes.Clientset, namespace string) {
	deployment := &appsv1.Deployment{
		// Kind: "deployment",
		// APIVersion: "apps/v1",
		TypeMeta: metadatav1.TypeMeta{
			Kind:       "deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metadatav1.ObjectMeta{
			Name:      "test-deployment",
			Namespace: namespace,
		},
		Spec: appsv1.DeploymentSpec{
			// 使用格式转换函数进行转换格式
			Replicas: int32Ptr(2),
			Selector: &metadatav1.LabelSelector{
				MatchLabels: map[string]string{"app": "nginx"},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metadatav1.ObjectMeta{
					Name:      "nginx",
					Namespace: namespace,
					Labels:    map[string]string{"app": "nginx"},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						v1.Container{
							Name:  "nginx",
							Image: "nginx",
						},
					},
				},
			},
		},
	}
	deploymentcreate, err := clientset.AppsV1().Deployments(namespace).Create(context.TODO(), deployment, metadatav1.CreateOptions{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s is the %s create is successfull\n", deploymentcreate.GetName(), deploymentcreate.GetNamespace())
}

func updatdeploy(clientset *kubernetes.Clientset, namespace string) {
	// retry.RetryOnConflict 用于在您必须担心其他代码同时对资源进行不相关的更新而导致冲突时，对资源进行更新。
	// 当冲突时，不会立刻退出，而是会等待一段时间后重试
	// 超过重试次数后退出并返回错误，默认重试在 retry.DefaultRetry 中
	// 用于修改资源时的重试
	deploymentclient := clientset.AppsV1().Deployments(namespace)
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		result, getErr := deploymentclient.Get(context.TODO(), "test-deployment", metadatav1.GetOptions{})
		if getErr != nil {
			log.Fatal("failed get latest version of deployment err is", getErr)
		}
		result.Spec.Replicas = int32Ptr(1)
		result.Spec.Template.Spec.Containers[0].Image = "nginx:1.13"
		_, updateErr := deploymentclient.Update(context.TODO(), result, metadatav1.UpdateOptions{})
		return updateErr
	})
	if retryErr != nil {
		log.Fatal(retryErr)
	}
}

func listdeploy(clientset *kubernetes.Clientset, namespace string) {
	deployment := clientset.AppsV1().Deployments(namespace)
	deploylist, err := deployment.List(context.TODO(), metadatav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range deploylist.Items {
		fmt.Println(v.Name, v.Namespace, v.Status.Replicas, v.Spec.Selector, v.Spec.Template.Spec.Containers[0].Image)
	}
}

func deletedeploy(clisentset *kubernetes.Clientset, namespace string) {
	deployment := clisentset.AppsV1().Deployments(namespace)
	err := deployment.Delete(context.TODO(), "test-deployment", metadatav1.DeleteOptions{})
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	var kubeconfig *string
	var option string
	var namespace string
	flag.StringVar(&option, "option", "node", "get node list")
	flag.StringVar(&namespace, "namespace", "default", "get namespace resource")
	// homedir.HomeDir 提出当前目录的家目录路径
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "this is kube config")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "this is kubeconfig")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		log.Fatal(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}
	// 使用命令查看 kubectl api-resources 查看资源所在的资源组，然后调用即可
	switch option {
	case "node":
		getnode(clientset)
	case "deployment":
		getdeployment(clientset, namespace)
	case "createdeploy":
		createdeployment(clientset, namespace)
	case "updatedeploy":
		updatdeploy(clientset, namespace)
	case "listdeploy":
		listdeploy(clientset, namespace)
	case "deletedeploy":
		deletedeploy(clientset, namespace)
	}
}

func int32Ptr(i int32) *int32 {
	return &i
}
