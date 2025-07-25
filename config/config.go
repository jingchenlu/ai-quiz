package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type AppConfig struct {
	ServerPort         string
	Mode               string
	DBPath             string
	DashScopeApiKey    string
	SupportedLanguages map[string]interface{}
}

func init() {
	loadEnvFiles()
}

func loadEnvFiles() {
	// 获取当前工作目录
	workDir, err := os.Getwd()
	if err != nil {
		log.Println("警告: 无法获取工作目录:", err)
	}

	// 基础配置文件路径
	baseEnvPath := filepath.Join(workDir, ".env")

	// 首先加载基础配置
	err = godotenv.Load(baseEnvPath)
	if err != nil {
		log.Println("警告: 未找到.env文件或无法加载:", err)
	} else {
		log.Println("成功加载基础配置文件:", baseEnvPath)
	}

	// 加载本地覆盖配置
	localEnvPath := filepath.Join(workDir, ".env.local")
	if _, err := os.Stat(localEnvPath); err == nil {
		// .env.local 文件存在，加载它以覆盖基础配置
		err = godotenv.Overload(localEnvPath)
		if err != nil {
			log.Println("警告: 无法加载覆盖配置文件:", err)
		} else {
			log.Println("成功加载覆盖配置文件:", localEnvPath)
		}
	}

}

// GetConfig 获取应用配置
func GetConfig(LanguageFlag bool) *AppConfig {
	supportedLanguages := make(map[string]interface{})
	// 根据flag参数来决定是否加载语言支持
	if LanguageFlag {
		languagesStr := getEnv("SUPPORTED_LANGUAGES", "Go,Python,Java,JavaScript,C++,C#,PHP,Ruby")
		languages := strings.Split(languagesStr, ",")
		for _, lang := range languages {
			supportedLanguages[lang] = nil
		}
	}
	return &AppConfig{
		// 从环境变量中读取，如果不存在则使用默认值
		ServerPort:         getEnv("SERVER_PORT", ":8080"),
		Mode:               getEnv("GIN_MODE", "debug"),
		DBPath:             getEnv("DB_PATH", "./aiquiz.db"),
		DashScopeApiKey:    getEnv("DASHSCOPE_API_KEY", ""),
		SupportedLanguages: supportedLanguages,
	}
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
