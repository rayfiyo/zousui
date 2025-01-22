package consts

const (
	SpecifyingResponseFormat string = `このコミュニティの文化を新しい方向に進化させるアイデアを提案してください。あなたの出力は **必ず次のJSON形式** で返してください。
{
    "newCulture": "string",
    "populationChange": 0
}"`
	GeminiModel string = "gemini-2.0-flash-exp"
	DALLEModel  string = "dall-e-3"
	ImageSize   string = "1024x1024"
)
