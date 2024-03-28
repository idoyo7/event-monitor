package main

import (
    "encoding/json"
    "fmt"
    "log"
    "strings"
    "time"


    "k8s.io/client-go/rest"
    "k8s.io/apimachinery/pkg/fields"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/cache"
    corev1 "k8s.io/api/core/v1"
)

type Alert struct {
    PodName       string `json:"pod_name"`
    ExecutionTime string `json:"execution_time"`
    MessageType   string `json:"MessageType"`
}

func main() {
    // In-cluster config
    config, err := rest.InClusterConfig()
    if err != nil {
        panic(err.Error())
    }

    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        panic(err.Error())
    }

    watchlist := cache.NewListWatchFromClient(clientset.CoreV1().RESTClient(), "events", corev1.NamespaceAll, fields.Everything())

    // Record the current time as the start time
    startTime := time.Now()

    _, controller := cache.NewInformer(
        watchlist,
        &corev1.Event{},
        0, // Duration is ignored in this context but can be set to a particular sync time.
        cache.ResourceEventHandlerFuncs{
            AddFunc: func(obj interface{}) {
                event := obj.(*corev1.Event)
                // Compare event creation timestamp with the recorded start time
                if event.CreationTimestamp.Time.After(startTime) {
                    if event.Reason == "FailedScheduling" && strings.Contains(event.Message, "Insufficient cpu") {
                        // fmt.Printf("%s Insufficient cpu: %s\n", event.ObjectMeta.CreationTimestamp.Time, event.InvolvedObject.Name)
                        podName := event.InvolvedObject.Name
                        executionTime := event.CreationTimestamp.Time.Format(time.RFC3339) // or event.LastTimestamp.Time.Format(time.RFC3339)
                        alert := Alert{
                            PodName:       podName,
                            ExecutionTime: executionTime,
                            MessageType:   "CPU",
                        }
                        alertJson, err := json.Marshal(alert)
                        if err != nil {
                            log.Fatal(err)
                        }

                        fmt.Println(string(alertJson))
                        err = SendTeamsWebhook(alert)
                            
                        if err != nil {
                            fmt.Println("Error sending Teams webhook:", err)
                        }
                        
                    }
                }
            },
        },
    )

    stop := make(chan struct{})
    defer close(stop)
    go controller.Run(stop)

    // Keep the main thread alive to listen to events
    select {}
}

