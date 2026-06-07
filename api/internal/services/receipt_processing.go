package services

import (
	"encoding/json"
	"fmt"
	"gopkg.in/gographics/imagick.v3/imagick"
	"gorm.io/gorm"
	"os"
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"
	"strings"
	"time"
)

type ReceiptProcessingService struct {
	BaseService
	ReceiptProcessingSettings         models.ReceiptProcessingSettings
	FallbackReceiptProcessingSettings models.ReceiptProcessingSettings
	Group                             models.Group
}

func NewSystemReceiptProcessingService(tx *gorm.DB, groupId string) (ReceiptProcessingService, error) {
	systemSettingsRepository := repositories.NewSystemSettingsRepository(tx)
	systemReceiptProcessingSettings, err := systemSettingsRepository.GetSystemReceiptProcessingSettings()
	group := models.Group{}
	if err != nil {
		return ReceiptProcessingService{}, err
	}

	if len(groupId) > 0 {
		groupRepository := repositories.NewGroupRepository(tx)
		groupToUse, err := groupRepository.GetGroupById(groupId, false, true, false)
		if err != nil {
			return ReceiptProcessingService{}, err
		}

		group = groupToUse

		if groupToUse.GroupSettings.PromptId != nil && *groupToUse.GroupSettings.PromptId > 0 {
			systemReceiptProcessingSettings.ReceiptProcessingSettings.PromptId = *groupToUse.GroupSettings.PromptId
		}

		if groupToUse.GroupSettings.FallbackPromptId != nil &&
			*groupToUse.GroupSettings.FallbackPromptId > 0 &&
			systemReceiptProcessingSettings.FallbackReceiptProcessingSettings.ID != 0 {
			systemReceiptProcessingSettings.FallbackReceiptProcessingSettings.PromptId = *groupToUse.GroupSettings.FallbackPromptId
		}
	}

	return ReceiptProcessingService{
		BaseService:                       BaseService{TX: tx},
		ReceiptProcessingSettings:         systemReceiptProcessingSettings.ReceiptProcessingSettings,
		FallbackReceiptProcessingSettings: systemReceiptProcessingSettings.FallbackReceiptProcessingSettings,
		Group:                             group,
	}, nil

}

func NewReceiptProcessingService(tx *gorm.DB, receiptProcessingSettingsId string, fallbackReceiptProcessingSettingsId string) (ReceiptProcessingService, error) {
	service := ReceiptProcessingService{BaseService: BaseService{
		DB: repositories.GetDB(),
		TX: tx,
	}}

	receiptProcessingSettingsRepository := repositories.NewReceiptProcessingSettings(nil)
	receiptProcessingSettings, err := receiptProcessingSettingsRepository.GetReceiptProcessingSettingsById(receiptProcessingSettingsId)
	if err != nil {
		return service, err
	}
	service.ReceiptProcessingSettings = receiptProcessingSettings

	if len(fallbackReceiptProcessingSettingsId) > 0 && fallbackReceiptProcessingSettingsId != "0" {
		fallbackReceiptProcessingSettings, err := receiptProcessingSettingsRepository.GetReceiptProcessingSettingsById(fallbackReceiptProcessingSettingsId)
		if err != nil {
			return service, err
		}
		service.FallbackReceiptProcessingSettings = fallbackReceiptProcessingSettings
	}

	return service, nil
}

func (service ReceiptProcessingService) ReadReceiptImage(
	imagePath string,
) (commands.UpsertReceiptCommand, commands.ReceiptProcessingMetadata, error) {
	return service.readReceipt([]string{imagePath}, "", false)
}

func (service ReceiptProcessingService) ReadReceiptImageWithEmailBody(
	imagePath string,
	emailBody string,
) (commands.UpsertReceiptCommand, commands.ReceiptProcessingMetadata, error) {
	return service.readReceipt([]string{imagePath}, emailBody, false)
}

// ReadReceiptImagesWithEmailBody processes one or more receipt images alongside
// an email body. When bodySentAsImage is true, the body is also represented as
// one of the images (e.g. a chromedp-rendered PDF page) and the body text is
// not duplicated into the prompt to avoid sending the same content twice.
func (service ReceiptProcessingService) ReadReceiptImagesWithEmailBody(
	imagePaths []string,
	emailBody string,
	bodySentAsImage bool,
) (commands.UpsertReceiptCommand, commands.ReceiptProcessingMetadata, error) {
	return service.readReceipt(imagePaths, emailBody, bodySentAsImage)
}

func (service ReceiptProcessingService) ReadReceiptText(
	emailBody string,
) (commands.UpsertReceiptCommand, commands.ReceiptProcessingMetadata, error) {
	return service.readReceipt(nil, emailBody, false)
}

func (service ReceiptProcessingService) readReceipt(
	imagePaths []string,
	emailBody string,
	bodySentAsImage bool,
) (commands.UpsertReceiptCommand, commands.ReceiptProcessingMetadata, error) {
	var receipt commands.UpsertReceiptCommand
	metadata := commands.ReceiptProcessingMetadata{}

	result, err := service.processImages(
		imagePaths,
		emailBody,
		bodySentAsImage,
		service.ReceiptProcessingSettings,
	)
	metadata.OcrSystemTaskCommand = result.OcrSystemTaskCommand
	metadata.PromptSystemTaskCommand = result.PromptSystemTaskCommand
	metadata.ChatCompletionSystemTaskCommand = result.ChatCompletionSystemTaskCommand
	metadata.ReceiptProcessingSettingsIdRan = service.ReceiptProcessingSettings.ID
	if err != nil {
		metadata.DidReceiptProcessingSettingsSucceed = false
		metadata.RawResponse = err.Error()

		if service.FallbackReceiptProcessingSettings.ID > 0 {
			fallbackResult, fallbackErr := service.processImages(
				imagePaths,
				emailBody,
				bodySentAsImage,
				service.FallbackReceiptProcessingSettings,
			)
			metadata.FallbackReceiptProcessingSettingsIdRan = service.FallbackReceiptProcessingSettings.ID
			metadata.FallbackOcrSystemTaskCommand = fallbackResult.OcrSystemTaskCommand
			metadata.FallbackPromptSystemTaskCommand = fallbackResult.PromptSystemTaskCommand
			metadata.FallbackChatCompletionSystemTaskCommand = fallbackResult.ChatCompletionSystemTaskCommand
			receipt = fallbackResult.Receipt
			err = fallbackErr

			if err != nil {
				metadata.DidFallbackReceiptProcessingSettingsSucceed = false
				metadata.FallbackRawResponse = err.Error()
			} else {
				metadata.DidFallbackReceiptProcessingSettingsSucceed = true
				metadata.FallbackRawResponse = fallbackResult.RawResponse
			}
		}
	} else {
		metadata.DidReceiptProcessingSettingsSucceed = true
		metadata.RawResponse = result.RawResponse
		receipt = result.Receipt
	}

	return receipt, metadata, err
}

// FormatEmailContent combines image-extracted text and email body into a formatted string
// for the AI prompt. This is used when processing emails to give the AI full context.
func FormatEmailContent(imageData string, hasImage bool, emailBody string) string {
	imageDataText := imageData
	if len(imageDataText) == 0 && !hasImage {
		imageDataText = "No attachments found"
	}

	emailBodyText := emailBody
	if len(emailBodyText) == 0 {
		emailBodyText = "No email body found"
	}

	return fmt.Sprintf("Image Data:\n%s\n\nEmail Body:\n%s", imageDataText, emailBodyText)
}

func (service ReceiptProcessingService) processImages(
	imagePaths []string,
	emailBody string,
	bodySentAsImage bool,
	receiptProcessingSettings models.ReceiptProcessingSettings,
) (commands.ReceiptProcessingResult, error) {
	aiMessages := []structs.AiClientMessage{}
	receipt := commands.UpsertReceiptCommand{}
	result := commands.ReceiptProcessingResult{}
	ocrText := ""
	encodedImages := []string{}
	hasImage := len(imagePaths) > 0

	if hasImage {
		if receiptProcessingSettings.IsVisionModel {
			for _, imagePath := range imagePaths {
				encoded, err := service.encodeImageForAi(imagePath, receiptProcessingSettings)
				if err != nil {
					return result, err
				}
				encodedImages = append(encodedImages, encoded)
			}
		} else {
			ocrService := NewOcrService(service.TX, receiptProcessingSettings)
			ocrResults := make([]ocrImageResult, 0, len(imagePaths))
			for _, imagePath := range imagePaths {
				resultText, ocrSystemTaskCommand, ocrErr := ocrService.ReadImage(imagePath)
				ocrResults = append(ocrResults, ocrImageResult{
					Text:    resultText,
					Command: ocrSystemTaskCommand,
					Err:     ocrErr,
				})
				if ocrErr != nil {
					break
				}
			}
			combinedText, combinedCmd, combinedErr := combineOcrResults(ocrResults)
			result.OcrSystemTaskCommand = combinedCmd
			if combinedErr != nil {
				return result, combinedErr
			}
			ocrText = combinedText
		}
	}

	promptBody := emailBody
	if bodySentAsImage {
		promptBody = ""
	}
	promptText := FormatEmailContent(ocrText, hasImage, promptBody)

	prompt, promptSystemTask, err := service.buildPrompt(receiptProcessingSettings, promptText)
	result.PromptSystemTaskCommand = promptSystemTask
	if err != nil {
		return result, err
	}

	message := structs.AiClientMessage{
		Role:    "user",
		Content: prompt,
	}
	if len(encodedImages) > 0 {
		message.Images = encodedImages
	}

	aiMessages = append(aiMessages, message)

	aiClient := AiService{
		ReceiptProcessingSettings: receiptProcessingSettings,
	}

	response, chatCompletionSystemTaskCommand, err := aiClient.CreateChatCompletion(structs.AiChatCompletionOptions{
		Messages:   aiMessages,
		DecryptKey: true,
	})
	result.ChatCompletionSystemTaskCommand = chatCompletionSystemTaskCommand
	result.RawResponse = response
	if err != nil {
		return result, err
	}

	cleanedResponse := service.cleanResponse(response)

	err = json.Unmarshal([]byte(cleanedResponse), &receipt)
	if err != nil {
		return result, err
	}

	result.Receipt = receipt
	return result, nil
}

func (service ReceiptProcessingService) cleanResponse(response string) string {
	response = strings.ReplaceAll(response, "```json", "")
	response = strings.ReplaceAll(response, "```", "")
	response = stripTrailingCommas(response)
	return response
}

// stripTrailingCommas removes trailing commas before } or ] so that output
// from LLMs that emit them (notably OpenAI gpt-4o) parses under Go's strict
// encoding/json. It is string-aware: it tracks whether the scan position is
// inside a JSON string (honoring backslash escapes) and only removes a comma
// when the comma is outside any string and the next non-whitespace character
// is a closing brace or bracket. This avoids corrupting values that legitimately
// contain sequences like ",}" or ",]" inside a quoted string.
func stripTrailingCommas(s string) string {
	var b strings.Builder
	b.Grow(len(s))

	inString := false
	escaped := false

	for i := 0; i < len(s); i++ {
		c := s[i]

		if inString {
			b.WriteByte(c)
			switch {
			case escaped:
				escaped = false
			case c == '\\':
				escaped = true
			case c == '"':
				inString = false
			}
			continue
		}

		if c == '"' {
			inString = true
			b.WriteByte(c)
			continue
		}

		if c == ',' {
			// Look ahead past whitespace for a closing brace/bracket.
			j := i + 1
			for j < len(s) && (s[j] == ' ' || s[j] == '\t' || s[j] == '\r' || s[j] == '\n') {
				j++
			}
			if j < len(s) && (s[j] == '}' || s[j] == ']') {
				// Trailing comma: drop it.
				continue
			}
		}

		b.WriteByte(c)
	}

	return b.String()
}

// ocrImageResult is the per-image outcome from running OCR. Used by
// combineOcrResults to aggregate multi-image OCR runs into a single text
// blob and a single system task command.
type ocrImageResult struct {
	Text    string
	Command commands.UpsertSystemTaskCommand
	Err     error
}

const ocrImageSeparator = "\n--- Image ---\n"

// combineOcrResults aggregates per-image OCR outcomes into a single text
// (joined with ocrImageSeparator) and a single UpsertSystemTaskCommand
// that spans the first start time through the latest end time. The status
// is FAILED if any individual OCR failed. Returns the first error
// encountered, if any. The combined command's ResultDescription is
// overwritten with the joined OCR text whenever any image OCR succeeded so
// the system-task UI reflects the full output; if every image failed, the
// original failure description set by OcrService.ReadImage is preserved
// (texts is empty in that case so there's nothing useful to overwrite
// with).
func combineOcrResults(results []ocrImageResult) (string, commands.UpsertSystemTaskCommand, error) {
	var combined commands.UpsertSystemTaskCommand
	var firstErr error
	texts := make([]string, 0, len(results))
	for index, r := range results {
		if index == 0 {
			combined = r.Command
		} else {
			combined.EndedAt = r.Command.EndedAt
			if r.Command.Status == models.SYSTEM_TASK_FAILED {
				combined.Status = models.SYSTEM_TASK_FAILED
			}
		}
		if r.Err != nil {
			if firstErr == nil {
				firstErr = r.Err
			}
			continue
		}
		texts = append(texts, r.Text)
	}
	combinedText := strings.Join(texts, ocrImageSeparator)
	if len(texts) > 0 {
		combined.ResultDescription = combinedText
	}
	return combinedText, combined, firstErr
}

func (service ReceiptProcessingService) encodeImageForAi(
	imagePath string,
	receiptProcessingSettings models.ReceiptProcessingSettings,
) (string, error) {
	switch receiptProcessingSettings.AiType {
	case models.OLLAMA:
		return service.getOllamaBase64Image(imagePath)
	case models.OPEN_AI_NEW, models.OPEN_AI_CUSTOM, models.OPEN_AI_CUSTOM_NEW:
		return service.getOpenAiBase64Image(imagePath)
	case models.GEMINI_NEW:
		return service.getGeminiImage(imagePath)
	}
	return "", fmt.Errorf("unsupported AI type for vision encoding: %s", receiptProcessingSettings.AiType)
}

// TODO: move to new ai client
func (service ReceiptProcessingService) getOllamaBase64Image(imagePath string) (string, error) {
	mw := imagick.NewMagickWand()
	err := mw.ReadImage(imagePath)
	if err != nil {
		return "", err
	}

	fileBytes, err := mw.GetImageBlob()
	if err != nil {
		return "", err
	}

	return utils.Base64Encode(fileBytes), nil
}

func (service ReceiptProcessingService) getOpenAiBase64Image(imagePath string) (string, error) {
	fileRepository := repositories.NewFileRepository(service.TX)
	fileBytes, err := os.ReadFile(imagePath)
	if err != nil {
		return "", err
	}

	uri, err := fileRepository.BuildEncodedImageString(fileBytes)
	if err != nil {
		return "", err
	}

	return uri, nil
}

func (service ReceiptProcessingService) getGeminiImage(imagePath string) (string, error) {
	fileBytes, err := os.ReadFile(imagePath)
	if err != nil {
		return "", err
	}

	return utils.Base64Encode(fileBytes), nil
}

func (service ReceiptProcessingService) buildPrompt(
	receiptProcessingSettings models.ReceiptProcessingSettings,
	ocrText string,
) (string, commands.UpsertSystemTaskCommand, error) {
	systemTaskCommand := commands.UpsertSystemTaskCommand{
		Type:                 models.PROMPT_GENERATED,
		Status:               models.SYSTEM_TASK_SUCCEEDED,
		AssociatedEntityType: models.PROMPT,
		AssociatedEntityId:   receiptProcessingSettings.PromptId,
		StartedAt:            time.Now(),
	}

	promptRepository := repositories.NewPromptRepository(service.TX)

	stringPromptId := utils.UintToString(receiptProcessingSettings.PromptId)

	prompt, err := promptRepository.GetPromptById(stringPromptId)
	if err != nil {
		systemTaskCommand.Status = models.SYSTEM_TASK_FAILED
		systemTaskCommand.ResultDescription = err.Error()
		endedAt := time.Now()
		systemTaskCommand.EndedAt = &endedAt

		return "", systemTaskCommand, err
	}

	templateVariableMap, err := service.buildTemplateVariableMap(ocrText)
	if err != nil {
		systemTaskCommand.Status = models.SYSTEM_TASK_FAILED
		systemTaskCommand.ResultDescription = err.Error()
		endedAt := time.Now()
		systemTaskCommand.EndedAt = &endedAt

		return "", systemTaskCommand, err
	}

	regex := utils.GetTriggerRegex()
	realPrompt := regex.ReplaceAllStringFunc(prompt.Prompt, func(variable string) string {
		templateVariable := structs.PromptTemplateVariable(variable)
		return templateVariableMap[templateVariable]
	})

	endedAt := time.Now()
	systemTaskCommand.EndedAt = &endedAt
	systemTaskCommand.ResultDescription = realPrompt

	return realPrompt, systemTaskCommand, nil
}

func (service ReceiptProcessingService) buildTemplateVariableMap(ocrText string) (map[structs.PromptTemplateVariable]string, error) {
	result := make(map[structs.PromptTemplateVariable]string)

	categoriesString, err := service.getCategoriesString()
	if err != nil {
		return result, err
	}

	tagsString, err := service.getTagsString()
	if err != nil {
		return result, err
	}

	customFieldsString, err := service.getCustomFieldsString()
	if err != nil {
		return result, err
	}

	currentYearString := utils.UintToString(uint(time.Now().Year()))

	result[structs.CATEGORIES] = categoriesString
	result[structs.TAGS] = tagsString
	result[structs.CUSTOM_FIELDS] = customFieldsString
	result[structs.OCR_TEXT] = ocrText
	result[structs.CURRENT_YEAR] = currentYearString

	return result, nil
}

func (service ReceiptProcessingService) getCategoriesString() (string, error) {
	categoryRepository := repositories.NewCategoryRepository(nil)
	categories, err := categoryRepository.GetAllCategories("id, name, description")
	if err != nil {
		return "", err
	}

	categoriesBytes, err := json.Marshal(categories)
	if err != nil {
		return "", err
	}

	return string(categoriesBytes), nil
}

func (service ReceiptProcessingService) getTagsString() (string, error) {
	tagsRepository := repositories.NewTagsRepository(nil)
	tags, err := tagsRepository.GetAllTags("id, name")
	if err != nil {
		return "", err
	}

	tagsBytes, err := json.Marshal(tags)
	if err != nil {
		return "", err
	}

	return string(tagsBytes), nil
}

func (service ReceiptProcessingService) getCustomFieldsString() (string, error) {
	customFieldRepository := repositories.NewCustomFieldRepository(nil)
	customFields, err := customFieldRepository.GetAllCustomFields("id, name, type, description")
	if err != nil {
		return "", err
	}

	customFieldsBytes, err := json.Marshal(customFields)
	if err != nil {
		return "", err
	}

	return string(customFieldsBytes), nil
}
