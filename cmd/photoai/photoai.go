package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	openai "github.com/sashabaranov/go-openai"
	"github.com/shahin-bayat/go-scraper/internal/model"
	"github.com/shahin-bayat/go-scraper/internal/store"
	"github.com/shahin-bayat/go-scraper/internal/util"
)

func main() {
	store, err := store.NewPostgresStore()
	if err != nil {
		log.Fatalf(err.Error())
	}
	if err = store.Init(); err != nil {
		log.Fatalf(err.Error())
	}

	OPENAI_API_KEY := util.GetEnvVariable("OPENAI_API_KEY")
	if OPENAI_API_KEY == "" {
		log.Fatalf("OPENAI_API_KEY is not set")
	}

	client := openai.NewClient(OPENAI_API_KEY)
	if client == nil {
		log.Fatalf("OpenAi client is nil")
	}

	images, err := store.GetImages()
	if err != nil {
		log.Fatalf(err.Error())
	}

	// INFO: update images that contain another image
	// Uncomment and first this block for updating images that contain another image
	// for _, image := range images {
	// 	filename := strings.Split(image.Filename, ".")[0] + ".txt"
	// 	imgData, err := os.ReadFile(fmt.Sprintf("assets/base64/%s", filename))
	// 	base64Image := strings.Split(string(imgData), ",")[1]

	// 	if err != nil {
	// 		log.Fatalf(err.Error())
	// 	}
	// 	hasImage, err := util.HasImage(string(base64Image))
	// 	if err != nil {
	// 		log.Fatalf(err.Error())
	// 	}
	// 	if hasImage {
	// 		updatedImage := model.UpdateImage(image, &model.UpdateImageRequest{
	// 			HasImage: hasImage,
	// 		})
	// 		// fmt.Printf("Updated image: %v has image : %t\n", image.ID, updatedImage.HasImage)
	// 		if err = store.UpdateImage(image.ID, updatedImage); err != nil {
	// 			log.Fatalf(err.Error())
	// 		}
	// 	}
	// }
	// image, err := store.GetImageByFilename("17617.png")
	// if err != nil {
	// 	log.Fatalf(err.Error())
	// }

	for _, image := range images {
		if image.ExtractedText != nil {
			continue
		}
		delay := time.Duration(util.GenerateRandomDelay(1500, 3000)) * time.Millisecond
		time.Sleep(delay)

		filename := strings.Split(image.Filename, ".")[0] + ".txt"
		base64Image, err := os.ReadFile(fmt.Sprintf("assets/base64/%s", filename))
		if err != nil {
			log.Fatalf(err.Error())
		}

		prompt := `Please look at the image provided below and extract the exact text in German. The text is always at the top of the image. The response should be always in the following format: Text: [Extracted text]`

		resp, err := client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model:     openai.GPT4VisionPreview,
				MaxTokens: 300,
				Messages: []openai.ChatCompletionMessage{
					{
						Role: openai.ChatMessageRoleUser,
						MultiContent: []openai.ChatMessagePart{
							{
								Type: openai.ChatMessagePartTypeText,
								Text: prompt,
							},
							{
								Type: openai.ChatMessagePartTypeImageURL,
								ImageURL: &openai.ChatMessageImageURL{
									URL:    string(base64Image),
									Detail: openai.ImageURLDetailLow,
								},
							},
						},
					},
				},
			},
		)

		if err != nil {
			fmt.Printf("ChatCompletion error: %v\n", err)
			return
		}
		aiRes := resp.Choices[0].Message.Content
		extractedText := strings.Split(aiRes, "Text:")[1]
		fmt.Printf("Extracted text: %s\n", extractedText)
		os.WriteFile(fmt.Sprintf("assets/extracted/%s", filename), []byte(extractedText), 0644)
		updatedImage := model.UpdateImage(image, &model.UpdateImageRequest{
			ExtractedText: &extractedText,
		})
		if err = store.UpdateImage(image.ID, updatedImage); err != nil {
			log.Fatalf(err.Error())
		}
	}
}
