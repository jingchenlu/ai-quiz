package enums

// AiModel 定义AI模型类型
type AiModel string

// AiModelQwenPlus 预定义支持的AI模型常量
const (
	AiModelQwenPlus AiModel = "qwen-plus"
)

// SupportedAiModels 所有支持的AI模型集合
var SupportedAiModels = map[AiModel]struct{}{
	AiModelQwenPlus: {},
}

// IsSupportedAiModel 检查模型是否支持
func IsSupportedAiModel(model AiModel) bool {
	_, exists := SupportedAiModels[model]
	return exists
}
