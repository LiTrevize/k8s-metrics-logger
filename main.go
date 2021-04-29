package main

import "fmt"

func main() {
	kc := NewKubeletClient()

	summary := kc.GetSummary()

	for _, container := range summary.Node.SystemContainers {
		fmt.Println(container.Name)
	}

	for _, pod := range summary.Pods {
		fmt.Println(pod.Metadata.Name)
	}
}
