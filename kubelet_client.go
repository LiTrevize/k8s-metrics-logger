package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func getSecretsToken() string {
	return "eyJhbGciOiJSUzI1NiIsImtpZCI6ImN2Q29EOGgyRHZoME8zemo0Zlg0RlZyLTVLS0k1cWJJc1ltYUsxUnVoMFEifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJrdWJlLXN5c3RlbSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VjcmV0Lm5hbWUiOiJtZXRyaWNzLXNlcnZlci10b2tlbi1rYzlqeCIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJtZXRyaWNzLXNlcnZlciIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6IjAyNmY5OWY1LWMxNTYtNDRkZC04MmU5LWM0Yjk1M2E2MzBlMCIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDprdWJlLXN5c3RlbTptZXRyaWNzLXNlcnZlciJ9.fpnrEPTniLLeZtyRvUy10AC3YTFBbNoL_BwEbY-z3o4gdYs_tRKVf3eBat3KScnh80m_JrE3WLs5RepynAoAcuksuvtMyL1RApSigGppi8zglRt2pHjooMP3cVAH8h6894gukLrseL6Cwc3C6yOQw0pDb9amvxsCyO4LLw_7aT2T3pwope09-mf_iJ63w0TAfvbE5DzPVajJF2Wlkav904LEK3dB6wt813pK6TA8TsqIOS74AUNoTyktBrjZWB-Mun2lg0coG4Ywki8x5k7Hp4FMIVo_4_zY2GdxViNyb2EUw-27bBZXjSGi8QNxfo2FOkuA7YFRHvhlS8w5rg-W6A"
}

type KubeletClient struct {
	Client *http.Client
	Url    string
}

func NewKubeletClient() *KubeletClient {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &KubeletClient{
		Client: &http.Client{Transport: tr},
		Url:    "https://127.0.0.1:10250",
	}
}

func (kc *KubeletClient) GetSummary() *SummaryType {
	req, err := http.NewRequest("GET", kc.Url+"/stats/summary", nil)
	if err != nil {
		fmt.Println("request creation error: ", err)
	}
	req.Header.Add("Authorization", "Bearer "+getSecretsToken())
	rsp, err := kc.Client.Do(req)
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
