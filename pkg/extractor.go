package extractor

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gocolly/colly"
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

func Extract(url, auth string) {
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
			fmt.Println("Error parsing JSON:", err)
			return
		}

		for i, question := range responseData.Data.Form.Questions {
			cleanedQuestion := strings.ReplaceAll(question.FormsProRTQuestionTitle, "&nbsp;", " ")
			cleanedQuestion = strings.TrimSpace(cleanedQuestion)
			fmt.Printf("Question %d: %s\n", i+1, cleanedQuestion)

			var questionInfo QuestionInfo
			err := json.Unmarshal([]byte(question.QuestionInfo), &questionInfo)
			if err != nil {
				fmt.Println("Error parsing questionInfo JSON:", err)
				continue
			}

			fmt.Println("  Choices:")
			for i, choice := range questionInfo.Choices {
				fmt.Printf("%v %s\n", i+1, choice.Description)
			}
		}
	})

	err := scraper.Visit(url)
	if err != nil {
		fmt.Println("Error visiting URL:", err)
	}
}
