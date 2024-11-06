package utils

import (
	"gogpt/models"
	"log"
	"net/http"
)

func SendGoogleAnalyticsEvent(userAgent string, clientIP string) {
	integration, err := models.GetIntegrations()
	if err != nil || integration.GoogleAnalyticsID == "" {
		log.Println("Google Analytics ID bulunamadı veya hatalı:", err)
		return
	}

	trackingID := integration.GoogleAnalyticsID
	collectURL := "https://www.google-analytics.com/collect"

	req, err := http.NewRequest("POST", collectURL, nil)
	if err != nil {
		log.Println("Google Analytics isteği oluşturulamadı:", err)
		return
	}

	// Google Analytics verileri
	params := req.URL.Query()
	params.Add("v", "1")
	params.Add("tid", trackingID)
	params.Add("cid", clientIP)
	params.Add("t", "pageview")

	req.URL.RawQuery = params.Encode()
	req.Header.Set("User-Agent", userAgent)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Google Analytics isteği başarısız:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("Google Analytics hatalı durum kodu döndü:", resp.Status)
	}
}
