package config

// AIConfig 存储AI大模型的配置信息
type AIConfig struct {
	APIKey  string
	ModelID string
	BaseURL string
}

// GetAIConfig 获取AI配置
func GetAIConfig() AIConfig {
	return AIConfig{
		APIKey:  GetEnv("AI_API_KEY", "e41deff8-b5e8-4000-a588-ea5171dba541"),
		ModelID: GetEnv("AI_MODEL_ID", "deepseek-r1-distill-qwen-7b-250120"),
		BaseURL: GetEnv("AI_BASE_URL", "https://api.volcengine.com/v1/llm"),
	}
}
