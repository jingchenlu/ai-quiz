package ai

import (
	"aiquiz/config"
	"aiquiz/models/dto"
	"aiquiz/utils/enums"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Message 定义请求和响应结构体
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Input struct {
	Messages []Message `json:"messages"`
}

type Parameters struct {
	ResultFormat string `json:"result_format"`
}

type RequestBody struct {
	Model      string     `json:"model"`
	Input      Input      `json:"input"`
	Parameters Parameters `json:"parameters"`
}

// GenerateResponse 生成题目响应结构体
type GenerateResponse struct {
	Questions []dto.Question `json:"questions"`
}

// GenerateQuestions 生成编程题目
// language: 编程语言（从配置的支持语言中选择）
// questionType: 题目类型（"single" 或 "multiple"）
// keywords: 关键词（如"Gin 框架"、"数据库操作"等）
// count: 题目数量
func GenerateQuestions(aiModel, language, questionType, keywords string, count int) (*GenerateResponse, error) {

	// 构建提示词
	prompt, err := buildPrompt(language, questionType, keywords, count)
	if err != nil {
		return nil, err
	}

	// 构建请求体
	requestBody := buildRequestBody(aiModel, language, questionType, prompt)

	// 发送请求
	respBody, err := sendRequest(requestBody)
	if err != nil {
		return nil, err
	}

	// 根据不同的模型解析API响应
	var cleanedJson string
	switch aiModel {
	case string(enums.AiModelQwenPlus):
		cleanedJson, err = parseQwenApiResponse(respBody)
	case string(enums.AiModelDeepSeek):
		cleanedJson, err = parseV3ApiResponse(respBody)
	}
	if err != nil {
		return nil, err
	}

	// 解析题目
	questions, err := parseQuestions(cleanedJson)
	if err != nil {
		return nil, err
	}

	// 验证题目
	if err := validateQuestions(questions, questionType); err != nil {
		return nil, err
	}

	return &GenerateResponse{Questions: questions}, nil
}

// 构建提示词
func buildPrompt(language, questionType, keywords string, count int) (string, error) {
	questionTypeName := map[string]string{
		"single":   "单项",
		"multiple": "多项",
	}[questionType]

	if questionType == "single" {
		return fmt.Sprintf(`请严格按照以下要求生成%d道关于%s编程语言的%s选择题，主题围绕"%s"：

1. 输出格式：
   - 仅返回一个JSON数组，不包含任何额外文本、解释或说明
   - 数组中的每个元素必须符合以下结构：
   {
     "title": "题目标题（必须是完整的问题）",
     "options": [
       { "content": "选项内容", "value": 1 },
       { "content": "选项内容", "value": 2 },
       { "content": "选项内容", "value": 4 },
       { "content": "选项内容", "value": 8 }
     ],
     "answer": 2,  // 正确选项的value值（仅一个正确选项）
     "explanation": "详细解释正确答案的原因及错误选项的问题，不要包含value等信息"
   }

2. 内容要求：
   - 题目必须与%s编程语言和"%s"主题直接相关
   - 所有题目必须为%s选择题（只有一个正确答案）
   - 每个题目必须有4个选项
   - 选项应具有迷惑性，避免明显错误
   - 题目相互独立，不得重复

3. 格式约束：
   - 确保JSON格式完全正确
   - 选项value严格遵循2的次幂规则
   - answer字段必须是唯一正确选项的value值
   - 特别注意: 不允许包含任何Markdown格式标记，如标识json的代码块`,
			count, language, questionTypeName, keywords,
			language, keywords, questionTypeName), nil
	}

	// 多选题提示词
	return fmt.Sprintf(`请严格按照以下要求生成%d道关于%s编程语言的%s选择题，主题围绕"%s"：

1. 输出格式：
   - 仅返回一个JSON数组，不包含任何额外文本、解释或说明
   - 数组中的每个元素必须符合以下结构：
   {
     "title": "题目标题（必须是完整的问题）",
     "options": [
       { "content": "选项内容", "value": 1 },
       { "content": "选项内容", "value": 2 },
       { "content": "选项内容", "value": 4 },
       { "content": "选项内容", "value": 8 }
     ],
     "answer": 7,  // 正确选项的value之和（至少2个正确选项）
     "explanation": "详细解释正确答案的原因及错误选项的问题，不要包含value等信息"
   }

2. 内容要求：
   - 题目必须与%s编程语言和"%s"主题直接相关
   - 所有题目必须为%s选择题（至少2个正确答案）
   - 每个题目必须有4个选项
   - 选项应具有迷惑性，避免明显错误
   - 题目相互独立，不得重复

3. 格式约束：
   - 确保JSON格式完全正确
   - 选项value严格遵循2的次幂规则
   - answer字段必须是所有正确选项的value总和
   - 特别注意: 不允许包含任何Markdown格式标记,如标识json的代码块`,
		count, language, questionTypeName, keywords,
		language, keywords, questionTypeName), nil
}

// 构建请求体
func buildRequestBody(aiModel, language, questionType, prompt string) RequestBody {
	questionTypeName := map[string]string{
		"single":   "单项",
		"multiple": "多项",
	}[questionType]

	return RequestBody{
		Model: aiModel,
		Input: Input{
			Messages: []Message{
				{
					Role:    "system",
					Content: fmt.Sprintf("你是专业的编程题目生成助手，专注生成%s编程语言的%s选择题。", language, questionTypeName),
				},
				{
					Role:    "user",
					Content: prompt,
				},
			},
		},
		Parameters: Parameters{
			ResultFormat: "message",
		},
	}
}

// 发送HTTP请求
func sendRequest(requestBody RequestBody) ([]byte, error) {
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("序列化请求体失败: %v", err)
	}

	req, err := http.NewRequest("POST",
		"https://dashscope.aliyuncs.com/api/v1/services/aigc/text-generation/generation",
		bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	appConfig := config.GetConfig(false)
	apiKey := appConfig.DashScopeApiKey
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("请求失败，状态码: %d，响应内容: %s", resp.StatusCode, string(bodyText))
	}

	return bodyText, nil
}

// qwen响应解析
func parseQwenApiResponse(bodyText []byte) (string, error) {
	var apiResponse struct {
		Output struct {
			Text string `json:"text"`
		} `json:"output"`
	}

	if err := json.Unmarshal(bodyText, &apiResponse); err != nil {
		return "", fmt.Errorf("解析API响应失败: %v，响应内容: %s", err, string(bodyText))
	}

	return apiResponse.Output.Text, nil
}

// DeepSeek-V3的响应解析
func parseV3ApiResponse(bodyText []byte) (string, error) {
	var apiResponse struct {
		Output struct {
			Choices []struct {
				Message struct {
					Content string `json:"content"`
				} `json:"message"`
			} `json:"choices"`
		} `json:"output"`
	}
	if err := json.Unmarshal(bodyText, &apiResponse); err != nil {
		return "", fmt.Errorf("解析API响应失败: %v，响应内容: %s", err, string(bodyText))
	}
	// 检查是否有返回结果
	if len(apiResponse.Output.Choices) == 0 {
		return "", fmt.Errorf("API响应中没有找到有效内容，响应内容: %s", string(bodyText))
	}
	// 返回第一个选择的内容
	return apiResponse.Output.Choices[0].Message.Content, nil
}

// 解析题目数组
func parseQuestions(cleanedJson string) ([]dto.Question, error) {
	var questions []dto.Question
	if err := json.Unmarshal([]byte(cleanedJson), &questions); err != nil {
		return nil, fmt.Errorf("解析题目数组失败: %v，内容: %s", err, cleanedJson)
	}
	return questions, nil
}

// 验证题目是否符合题型要求
func validateQuestions(questions []dto.Question, questionType string) error {
	for i, q := range questions {
		if questionType == "single" {
			if err := validateSingleQuestion(q, i+1); err != nil {
				return err
			}
		} else {
			if err := validateMultipleQuestion(q, i+1); err != nil {
				return err
			}
		}
	}
	return nil
}

// 验证单选题
func validateSingleQuestion(q dto.Question, index int) error {
	// 检查选项数量
	if len(q.Options) != 4 {
		return fmt.Errorf("第%d题不符合要求，单选题必须有4个选项", index)
	}

	// 检查答案是否为单个选项的值
	found := false
	for _, opt := range q.Options {
		if opt.Value == q.Answer {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("第%d题不符合单选题要求，答案不是单个选项的value", index)
	}

	return nil
}

// 验证多选题
func validateMultipleQuestion(q dto.Question, index int) error {
	// 检查选项数量
	if len(q.Options) != 4 {
		return fmt.Errorf("第%d题不符合要求，多选题必须有4个选项", index)
	}

	// 检查正确选项数量是否至少为2个
	sum := 0
	countCorrect := 0
	for _, opt := range q.Options {
		if (q.Answer & opt.Value) == opt.Value {
			sum += opt.Value
			countCorrect++
		}
	}

	if sum != q.Answer {
		return fmt.Errorf("第%d题答案计算错误，正确选项value之和与answer不匹配", index)
	}

	if countCorrect < 2 {
		return fmt.Errorf("第%d题不符合多选题要求，正确选项数量不足2个", index)
	}

	return nil
}
