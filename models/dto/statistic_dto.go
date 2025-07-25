package dto

// UserStatisticsRes 用户统计信息响应
type UserStatisticsRes struct {
	UserID                   int                    `json:"user_id"`           // 用户ID
	Username                 string                 `json:"username"`          // 用户名
	QuestionCount            int                    `json:"question_count"`    // 出题次数
	PaperCount               int                    `json:"paper_count"`       // 试卷数量
	QuestionTypeDistribution []TypeDistribution     `json:"type_distribution"` // 题目类型分布
	LanguageDistribution     []LanguageDistribution `json:"language_distribution"`
	ActiveTimeAnalysis       ActiveTimeAnalysis     `json:"active_time"` // 活跃时间分析
}

// TypeDistribution 题目类型分布
type TypeDistribution struct {
	Type  string `json:"type"`  // 题目类型(如:单选、多选)
	Count int    `json:"count"` // 数量
}

// LanguageDistribution 语言类型
type LanguageDistribution struct {
	Language string `json:"language"` // 语言类型（Go、Java）
	Count    int    `json:"count"`    // 数量
}

// ActiveTimeAnalysis 活跃时间分析
type ActiveTimeAnalysis struct {
	Daily   map[string]int `json:"daily"`   // 按天统计(最近30天)
	Monthly map[string]int `json:"monthly"` // 按月统计(最近12个月)
	Hourly  [24]int        `json:"hourly"`  // 按小时统计(0-23点)
}

// SystemStatisticsRes 系统统计结果响应
type SystemStatisticsRes struct {
	TotalUserCount            int                         `json:"total_user_count"`            // 总用户数
	TotalQuestionCount        int                         `json:"total_question_count"`        // 总题目数
	TotalPaperCount           int                         `json:"total_paper_count"`           // 总试卷数
	LanguageDistribution      []LanguageDistribution      `json:"language_distribution"`       // 编程语言分布
	AIModelUsage              []AIModelDistribution       `json:"ai_model_usage"`              // AI模型使用情况
	PaperQuestionDistribution []PaperQuestionDistribution `json:"paper_question_distribution"` // 试卷题目数量分布
	ActivityAnalysis          SystemActivityAnalysis      `json:"activity_analysis"`           // 活跃度分析
}

// AIModelDistribution AI模型使用分布
type AIModelDistribution struct {
	ModelName string `json:"model_name"` // 模型名称
	Count     int    `json:"count"`      // 使用次数
}

// PaperQuestionDistribution 试卷题目数量分布
type PaperQuestionDistribution struct {
	Range string `json:"range"` // 题目数量范围（如"1-10"）
	Count int    `json:"count"` // 该范围的试卷数量
}

// SystemActivityAnalysis 系统活跃度分析
type SystemActivityAnalysis struct {
	Daily  map[string]int `json:"daily"`  // 每日活跃度（key:日期 val:次数）
	Weekly map[string]int `json:"weekly"` // 每周活跃度（key:周数 val:次数）
}
