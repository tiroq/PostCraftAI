package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tiroq/postcraftai/backend/models"
)

// GeneratePost calls the OpenAI API, enforcing rate limiting and logging.
func GeneratePost(c *gin.Context) {
	username := c.GetString("username")
	user, exists := models.Users[username]
	if !exists || !user.Allowed {
		log.Printf("GeneratePost: unauthorized access attempt by %s", username)
		c.JSON(http.StatusForbidden, gin.H{"error": "User not enabled"})
		return
	}
	if time.Now().After(user.AccessExpiresAt) {
		log.Printf("GeneratePost: access expired for user %s", username)
		c.JSON(http.StatusForbidden, gin.H{"error": "User access expired"})
		return
	}
	if !models.RateLimiterInstance.Allow(username, models.OpenAIRateLimit) {
		log.Printf("GeneratePost: rate limit exceeded for user %s", username)
		c.JSON(http.StatusTooManyRequests, gin.H{"error": fmt.Sprintf("Rate limit exceeded (max %d req/min)", models.OpenAIRateLimit)})
		return
	}

	var reqData struct {
		Article string `json:"article"`
	}
	if err := c.BindJSON(&reqData); err != nil {
		log.Printf("GeneratePost: bind error for user %s: %v", username, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	prompt := fmt.Sprintf(`Transform the following article into a short post.
Start with an engaging hook and summarize the key points succinctly.

Article:
%s`, reqData.Article)

	openAIReq := models.OpenAIRequest{
		Model:       "gpt-4",
		Prompt:      prompt,
		MaxTokens:   250,
		Temperature: 0.7,
	}

	apiKey := models.GetenvOrFail("OPENAI_API_KEY")
	responseText, err := callOpenAI(apiKey, openAIReq)
	if err != nil {
		log.Printf("GeneratePost: OpenAI call failed for user %s: %v", username, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate post"})
		return
	}
	models.LogRequest(username)
	log.Printf("GeneratePost: post generated for user %s", username)
	c.JSON(http.StatusOK, gin.H{"post": responseText})
}

func callOpenAI(apiKey string, request models.OpenAIRequest) (string, error) {
	url := "https://api.openai.com/v1/completions"
	body, err := json.Marshal(request)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("OpenAI API error: %s", resp.Status)
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var openAIResp models.OpenAIResponse
	if err := json.Unmarshal(respBody, &openAIResp); err != nil {
		return "", err
	}
	if len(openAIResp.Choices) > 0 {
		return openAIResp.Choices[0].Text, nil
	}
	return "", fmt.Errorf("No response from OpenAI")
}
