package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/tiroq/postcraftai/backend/models"

	"github.com/gin-gonic/gin"
)

// GenerateLinkedInPost calls the OpenAI API, enforcing rate limiting and logging.
func GenerateLinkedInPost(c *gin.Context) {
	username := c.GetString("username")
	user, exists := models.Users[username]
	if !exists || !user.Allowed {
		log.Printf("GenerateLinkedInPost: unauthorized access attempt by %s", username)
		c.JSON(http.StatusForbidden, gin.H{"error": "User not enabled"})
		return
	}
	if time.Now().After(user.AccessExpiresAt) {
		log.Printf("GenerateLinkedInPost: access expired for user %s", username)
		c.JSON(http.StatusForbidden, gin.H{"error": "User access expired"})
		return
	}
	if !models.RateLimiterInstance.Allow(username, models.OpenAIRateLimit) {
		log.Printf("GenerateLinkedInPost: rate limit exceeded for user %s", username)
		c.JSON(http.StatusTooManyRequests, gin.H{"error": fmt.Sprintf("Rate limit exceeded (max %d req/min)", models.OpenAIRateLimit)})
		return
	}

	var reqData struct {
		Article string `json:"article"`
	}
	if err := c.BindJSON(&reqData); err != nil {
		log.Printf("GenerateLinkedInPost: bind error for user %s: %v", username, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	prompt := fmt.Sprintf(`Transform the following article into a LinkedIn post.
Start with a hook, summarize key insights in 2-3 short paragraphs, and include a call-to-action.

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
		log.Printf("GenerateLinkedInPost: OpenAI call failed for user %s: %v", username, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate post"})
		return
	}
	models.LogRequest(username)
	log.Printf("GenerateLinkedInPost: post generated for user %s", username)
	c.JSON(http.StatusOK, gin.H{"linkedin_post": responseText})
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
