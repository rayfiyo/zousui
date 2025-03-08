package consts

const (
	SpecifyingResponseFormat string = `このコミュニティの文化を新しい方向に進化させるアイデアを提案してください。**出力は必ず日本語で返してください。** あなたの出力は **必ず次のJSON形式** で返してください。
{
    "newCulture": "string",
    "populationChange": 0
}"`
	GeminiModel   string = "gemini-2.0-flash-exp"
	DALLEModel    string = "dall-e-3"
	ImageSize     string = "1024x1024"
	DALLEEndpoint string = "https://api.openai.com/v1/images/generations"
)
