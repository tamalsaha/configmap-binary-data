package main

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/klog/v2/klogr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
)

func NewClient() (client.Client, error) {
	scheme := runtime.NewScheme()
	_ = clientgoscheme.AddToScheme(scheme)

	ctrl.SetLogger(klogr.New())
	cfg := ctrl.GetConfigOrDie()
	cfg.QPS = 100
	cfg.Burst = 100

	mapper, err := apiutil.NewDynamicRESTMapper(cfg)
	if err != nil {
		return nil, err
	}

	return client.New(cfg, client.Options{
		Scheme: scheme,
		Mapper: mapper,
		//Opts: client.WarningHandlerOptions{
		//	SuppressWarnings:   false,
		//	AllowDuplicateLogs: false,
		//},
	})
}

func main() {
	c, err := NewClient()
	if err != nil {
		panic(err)
	}

	var cms core.ConfigMapList
	err = c.List(context.TODO(), &cms, client.InNamespace(""))
	if err != nil {
		panic(err)
	}
	for _, obj := range cms.Items {
		fmt.Println(obj.Name)
	}

	im := false
	cm := core.ConfigMap{
		// TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "bin-cm2",
			Namespace: "default",
		},
		Immutable: &im,
		Data: map[string]string{
			"data.txt": `{"name": "tamal"}`,
		},
		BinaryData: map[string][]byte{
			// "data.txt": []byte(`{"name": "tamal"}`),
			"d2.txt": []byte(`{"name": "tamal"}`),
		},
	}
	err = c.Update(context.TODO(), &cm)
	if err != nil {
		panic(err)
	}
}
