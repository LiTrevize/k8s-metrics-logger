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
	return "eyJhbGciOiJSUzI1NiIsImtpZCI6ImN2Q29EOGgyRHZoME8zemo0Zlg0RlZyLTVLS0k1cWJJc1ltYUsxUnVoMFEifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJrdWJlLXN5c3RlbSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VjcmV0Lm5hbWUiOiJtZXRyaWNzLXNlcnZlci10b2tlbi1qcnhnayIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJtZXRyaWNzLXNlcnZlciIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6ImE3NTc4ZmM4LTZkMGQtNDEyMC05Y2Y2LTU1ZmVhMTY5MTZkMSIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDprdWJlLXN5c3RlbTptZXRyaWNzLXNlcnZlciJ9.AYjaVlF4n5_5y__O-boiVNe86SNYdkMPqgg7H-lG0puhcGKI4r3pSbNUuxuoaizBOj-YiES0EnY7jNKW0myPfZH6tWO1DuKX2Nu25zSYXrsD9dRvs4-9cM7BZTPw520c0e0kMkkGGH6SGu--i7hSjHUxu2opbaECnwpbDlkcltwx1Xd6VUmIJUMiU_eFwCl_f01KzKwBDbqTgdYZQe6shSBS578GQ6OALbMaieOVb7uHmhES_K8Rjif6k3CAcwZO49XQB0jbgMzJ-Z2gkLxrgU3UR5TwJXe-NVupOfZ9c6hAIDguNXea978xhNFD8o4XCYdfMrPN-SLcB0uJTerXKw"
}

type KubeletClient struct {
	Client *http.Client
	Url    string
}

func NewKubeletClient() KubeletClient {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return KubeletClient{
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
