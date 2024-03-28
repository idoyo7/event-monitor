package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "os"
)


var teamsWebhookURL = os.Getenv("WEBHOOK_URL")

func SendTeamsWebhook(data Alert) error {

    payload := map[string]interface{}{
        "@type":     "MessageCard",
        "@context":  "http://schema.org/extensions",
        "themeColor": "0076D7",
    }

    switch data.MessageType {
    case "CPU":
        payload["title"] = "Insufficient CPU"
        payload["text"] = fmt.Sprintf("**Pod Name:** %s is pending lack of CPU resource on kubernetes cluster\n\n", data.PodName)
    case "GPU":
        payload["title"] = "Insufficient GPU"
        payload["text"] = fmt.Sprintf("**Pod Name:** %s is pending lack of GPU resource on kubernetes cluster\n\n", data.PodName)
    case "NoResourceReady":
        payload["title"] = "Still Not Ready"
        payload["text"] = fmt.Sprintf("Node has still not Ready\n\n", data.PodName)
    default:
        payload["title"] = "Unknown Message Type"
        payload["text"] = "An unknown event type was received."
    }

    payloadBytes, err := json.Marshal(payload)
    if err != nil {
        return err
    }

    req, err := http.NewRequest("POST", teamsWebhookURL, bytes.NewReader(payloadBytes))
    if err != nil {
        return err
    }
    req.Header.Set("Content-Type", "application/json")
    
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("error sending message to Teams: status code %d", resp.StatusCode)
    }
    
    return nil
}    


