package util

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"k8s.io/client-go/rest"
)

type KubeletClient struct {
	Client *http.Client
	Config *rest.Config
	Url    string
}

func NewKubeletClient() *KubeletClient {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	kc := KubeletClient{
		Client: &http.Client{Transport: tr},
		Config: config,
		Url:    "https://" + os.Getenv("NODE_IP") + ":10250",
	}
	_, err = kc.Request("GET", "/stats/summary")
	if err != nil {
		fmt.Println("Warning: Cannot reach kubelet, switch to edgecore")
		kc.Url = "http://" + os.Getenv("NODE_IP") + ":10350"
		_, err = kc.Request("GET", "/stats/summary")
		if err != nil {
			panic(fmt.Sprintf("Connection error: Cannot reach either kubelet or edgecore, %+v", err))
		}
	}
	return &kc
}

func (kc *KubeletClient) GetSecretsToken() string {
	return kc.Config.BearerToken
}

func (kc *KubeletClient) Request(method, urlPrefix string) (*http.Response, error) {
	req, err := http.NewRequest(method, kc.Url+urlPrefix, nil)
	if err != nil {
		fmt.Println("request creation error: ", err)
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+kc.GetSecretsToken())
	return kc.Client.Do(req)
}

func (kc *KubeletClient) GetSummary() *SummaryType {
	rsp, err := kc.Request("GET", "/stats/summary")
	if err != nil {
		fmt.Println("request error: ", err)
		return nil
	}
	defer rsp.Body.Close()
	summary := SummaryType{}
	data, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		fmt.Println("Read body error: ", err)
		return nil
	}
	err = json.Unmarshal(data, &summary)
	if err != nil {
		fmt.Println("Json parse error: ", err)
	}
	return &summary

}

func (kc *KubeletClient) LogMetrics() {
	summary := kc.GetSummary()
	if summary == nil {
		return
	}

	nodeTag := map[string]string{"node": summary.Node.NodeName}

	(&MetricsLog{"node_cpu_usage_nano_cores", nodeTag, summary.Node.CPU.UsageNanoCores, summary.Node.CPU.Time, nil}).Log()

	(&MetricsLog{"node_memory_usage_bytes", nodeTag, summary.Node.Memory.UsageBytes, summary.Node.Memory.Time, nil}).Log()

	for _, pod := range summary.Pods {
		podTag := map[string]string{
			"node":      summary.Node.NodeName,
			"pod":       pod.Metadata.Name,
			"namespace": pod.Metadata.Namespace}
		(&MetricsLog{"pod_cpu_usage_nano_cores", podTag, pod.CPU.UsageNanoCores, pod.CPU.Time, nil}).Log()
		(&MetricsLog{"pod_memory_usage_bytes", podTag, pod.Memory.UsageBytes, pod.Memory.Time, nil}).Log()
	}

}

type SummaryType struct {
	Node struct {
		NodeName         string
		StartTime        time.Time
		CPU              CPUType
		Memory           MemType
		SystemContainers []ContainerType
	}
	Pods []PodType
}

type CPUType struct {
	Time           time.Time
	UsageNanoCores int
}

type MemType struct {
	Time           time.Time
	AvailableBytes int
	UsageBytes     int
}

type ContainerType struct {
	Name      string
	StartTime time.Time
	CPU       CPUType
	Memory    MemType
}

type PodType struct {
	Metadata struct {
		Name      string
		Namespace string
		UID       string
	} `json:"PodRef"`
	StartTime  time.Time
	CPU        CPUType
	Memory     MemType
	Containers []ContainerType
}
