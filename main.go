package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

// ---------- ローカル変換 ----------

func defaultRoast(text string) string {
	text = strings.TrimSpace(text)
	lower := strings.ToLower(text)

	switch {
	case strings.Contains(lower, "遅れ"):
		return "また遅刻かよ、反省する脳みそついてんの？"
	case strings.Contains(lower, "出られません"):
		return "また逃げか？いい加減にしろよな。"
	case strings.Contains(lower, "すみません"):
		return "謝れば済むと思ってんのか？"
	case strings.Contains(lower, "やる気"):
		return "やる気出ない？じゃあ帰れ。"
	default:
		return fmt.Sprintf("%s？は？何甘えてんだ。", text)
	}
}

// ---------- ChatGPT呼び出し ----------

func chatgptRoast(apiKey, text string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	client := openai.NewClient(
		option.WithAPIKey(apiKey), // OPENAI_API_KEY でなくても OK
	)

	resp, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: openai.ChatModelGPT4oMini,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage("あなたは辛辣な日本語作成芸人です。ユーザーから渡された文章を暴言に変換してください。伝えたい内容・意味はそのままで"),
			openai.UserMessage(text),
		},
		MaxTokens: openai.Int(1000),
		// Temperature など細かいパラメータは必要に応じて
	})
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(resp.Choices[0].Message.Content), nil
}

// ---------- main ----------

func main() {
	// 1. 標準入力をまとめて取得
	scanner := bufio.NewScanner(os.Stdin)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	input := strings.Join(lines, " ")

	// 2. APIKEY があるか判定
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey != "" {
		if roast, err := chatgptRoast(apiKey, input); err == nil && roast != "" {
			fmt.Println(roast)
			return
		}
		// 失敗時はフォールバック
	}

	fmt.Println(defaultRoast(input))
}

