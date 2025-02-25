package extractorimport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type ResponseData struct {
	Data struct {
		Form struct {
			Questions []struct {
				FormsProRTQuestionTitle string `json:"formsProRTQuestionTitle"`
				QuestionInfo            string `json:"questionInfo"`
			} `json:"questions"`
		} `json:"form"`
	} `json:"data"`
}

type QuestionInfo struct {
	Choices []struct {
		Description string `json:"Description"`
	} `json:"Choices"`
}

type APIRequestPayload struct {
	Question string   `json:"question"`
	Choices  []string `json:"choices"`
}

func welcome() {
	fmt.Println("Welcome to the Extraction and API Integration Tool!")
}

func Extract(url, auth string) {
	welcome() // Call the welcome function
	fmt.Println("Extracting...")
	scraper := colly.NewCollector()

	scraper.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Accept", "*/*")
		r.Headers.Set("Accept-Language", "en-US,en;q=0.9")
		r.Headers.Set("sec-ch-ua-platform", `"Windows"`)
		r.Headers.Set("Authorization", auth)
		r.Headers.Set("x-fsw-page", "/pages/assignmentsresponsepage.aspx")
		r.Headers.Set("sec-ch-ua", `"Not?A_Brand";v="99", "Chromium";v="130"`)
		r.Headers.Set("sec-ch-ua-mobile", "?0")
		r.Headers.Set("x-fsw-ring", "Business")
		r.Headers.Set("x-fsw-enable", "1")
		r.Headers.Set("x-fsw-server", "16.0.18619.42500")
		r.Headers.Set("x-fsw-startup", "1")
		r.Headers.Set("Content-Type", "application/json")
		r.Headers.Set("Referer", "https://forms.office.com/Pages/AssignmentsResponsePage.aspx?id=lBwpK7Bet0SJ6kuvFuzEqYzhZnsU-WZHlw3Gwxuot4dUM0o2TUVFU0pPOVM3SkM4WVpER0tXQVA5WiQlQCN0PWcu&tid=2b291c94-5eb0-44b7-89ea-4baf16ecc4a9")
		r.Headers.Set("x-ms-form-request-ring", "business")
		r.Headers.Set("x-fsw-baseclient", "formweekly_cd_20250218.2")
		r.Headers.Set("x-ms-form-muid", "22E27181B2D663352E42641CB6D668B6")
		r.Headers.Set("x-ms-form-request-source", "ms-assignments")
		r.Headers.Set("x-fsw-cdn", "https://forms.office.com/cdn")
		r.Headers.Set("x-fsw-client", "forms_hotfix_cd_20250219.1")
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36")
		r.Headers.Set("DNT", "1")
		r.Headers.Set("X-CorrelationId", "6985eeff-4fe6-4fd9-873b-24c471c723b3")
		fmt.Println("Requesting:", r.URL)
	})

	scraper.OnResponse(func(r *colly.Response) {
		var responseData ResponseData
		err := json.Unmarshal(r.Body, &responseData)
		if err != nil {
			log.Println("Error parsing JSON:", err)
			time.Sleep(time.Second * 10) // Sleep for 10 seconds
			return
		}
		for i, question := range responseData.Data.Form.Questions {
			cleanedQuestion := strings.ReplaceAll(question.FormsProRTQuestionTitle, "&nbsp;", " ")
			cleanedQuestion = strings.TrimSpace(cleanedQuestion)
			fmt.Printf("Question %d: %s\n", i+1, cleanedQuestion)
			var questionInfo QuestionInfo
			err := json.Unmarshal([]byte(question.QuestionInfo), &questionInfo)
			if err != nil {
				log.Println("Error parsing questionInfo JSON:", err)
				time.Sleep(time.Second * 10) // Sleep for 10 seconds
				continue
			}
			fmt.Println("  Choices:")
			for i, choice := range questionInfo.Choices {
				fmt.Printf("%v %s\n", i+1, choice.Description)
			}
			var choices []string
			for _, choice := range questionInfo.Choices {
				choices = append(choices, choice.Description)
			}
			sendToApi(cleanedQuestion, choices)
		}
	})

	scraper.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
		time.Sleep(time.Second * 10) // Sleep for 10 seconds
	})

	err := scraper.Visit(url)
	if err != nil {
		log.Println("Error visiting URL:", err)
		time.Sleep(time.Second * 10) // Sleep for 10 seconds
	}
}

func sendToApi(question string, choices []string) {
	var api string = "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent?key=AIzaSyCg3Nxtt4HRorR_AEcpFf9taiX9Ad4YUEo"
	requestBody := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]string{
					{"text": fmt.Sprintf("Solve The Answers: %s\nChoices: %v", question, choices)},
				},
			},
		},
	}
	body, err := json.Marshal(requestBody)
	if err != nil {
		log.Println("Error marshalling request body:", err)
		time.Sleep(time.Second * 10) // Sleep for 10 seconds
		return
	}
	req, err := http.NewRequest("POST", api, bytes.NewBuffer(body))
	if err != nil {
		log.Println("Error creating request:", err)
		time.Sleep(time.Second * 10) // Sleep for 10 seconds
		return
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending request:", err)
		time.Sleep(time.Second * 10) // Sleep for 10 seconds
		return
	}
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		time.Sleep(time.Second * 10) // Sleep for 10 seconds
		return
	}
	var apiResponse map[string]interface{}
	err = json.Unmarshal(responseBody, &apiResponse)
	if err != nil {
		log.Println("Error unmarshalling API response:", err)
		time.Sleep(time.Second * 10) // Sleep for 10 seconds
		return
	}
	// Check if "candidates" key is present and not empty
	candidates, ok := apiResponse["candidates"].([]interface{})
	if !ok || len(candidates) == 0 {
		log.Printf("No candidates found in response: %v", apiResponse)
		return // Exit early if candidates are missing or empty
	}
	// Extract the 'text' from the response (Answer)
	candidate := candidates[0].(map[string]interface{})
	content := candidate["content"].(map[string]interface{})
	parts := content["parts"].([]interface{})
	if len(parts) > 0 {
		text := parts[0].(map[string]interface{})["text"].(string)
		// Format the output (only question and answer)
		formattedResponse := fmt.Sprintf("Question: %s\nAnswer: %s\n\n", question, text)
		// Append the response to response.txt
		f, err := os.OpenFile("response.txt", os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Println("Error opening file to append:", err)
			time.Sleep(time.Second * 10) // Sleep for 10 seconds
			return
		}
		defer f.Close()
		// Write the formatted response to the file
		_, err = f.WriteString(formattedResponse)
		if err != nil {
			log.Println("Error appending to file:", err)
			time.Sleep(time.Second * 10) // Sleep for 10 seconds
			return
		}
		fmt.Println("Response saved to response.txt")
	}
}
