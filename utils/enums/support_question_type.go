package enums

type QuestionType string

const (
	SingleType   QuestionType = "single"
	MultipleType QuestionType = "multiple"
)

// SupportedQuestionType 所有支持的题型
var SupportedQuestionType = map[QuestionType]struct{}{
	SingleType:   {},
	MultipleType: {},
}

// IsSupportedQuestionType 检查模型是否支持
func IsSupportedQuestionType(qType QuestionType) bool {
	_, exists := SupportedQuestionType[qType]
	return exists
}
