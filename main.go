// roast.go
//
// ビルド例: go build -o roast roast.go
// 実行例:  OPENAI_API_KEY=sk-xxx echo "遅れてすみません" | ./roast

package main

import (
	"bytes"
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

/* ---------- 1. ローカル毒舌 ---------- */

func localRoast(text string) string {
	text = strings.TrimSpace(text)
	l := strings.ToLower(text)
	switch {
	case strings.Contains(l, "遅れ"):
		return "また遅刻かよ、反省する脳みそついてんの？"
	case strings.Contains(l, "出られません"):
		return "また逃げか？いい加減にしろよな。"
	case strings.Contains(l, "すみません"):
		return "謝れば済むと思ってんのか？"
	case strings.Contains(l, "やる気"):
		return "やる気出ない？じゃあ帰れ。"
	default:
		return fmt.Sprintf("%s？は？何甘えてんだ。", text)
	}
}

/* ---------- 2. ChatGPT で毒舌生成 ---------- */

func chatRoast(ctx context.Context, key, prompt string) (string, error) {
	client := openai.NewClient(option.WithAPIKey(key))
	resp, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: openai.ChatModelGPT4oMini,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(`あなたは辛辣な日本語芸人です。渡された文章を一言の暴言に変換してください（意味は保持）`),
			openai.UserMessage(prompt),
		},
		MaxTokens: openai.Int(120),
	})
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(resp.Choices[0].Message.Content), nil
}

/* ---------- 3. Text-to-Speech (tts-1) ---------- */

func speakTTS(ctx context.Context, key, text string) error {
	// ① TTS に POST
	body, _ := json.Marshal(map[string]any{
		"model":          "tts-1",
		"input":          text,
		"voice":          "alloy",        // 他: nova / onyx / shimmer / fable / echo
		"response_format": "wav",         // WAV は ffplay が自動判定しやすい
	})
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost,
		"https://api.openai.com/v1/audio/speech", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+key)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("TTS API error: %s", b)
	}

	// ② ffplay を起動してストリーム再生
	ff, err := exec.LookPath("ffplay")
	if err != nil {
		return fmt.Errorf("ffplay not found: %w", err)
	}
	cmd := exec.Command(ff, "-nodisp", "-autoexit", "-loglevel", "quiet", "-")
	in, _  := cmd.StdinPipe()
	if err := cmd.Start(); err != nil {
		return err
	}
	// Body → ffplay へコピー
	if _, err := io.Copy(in, resp.Body); err != nil {
		return err
	}
	in.Close()
	return cmd.Wait()
}

/* ---------- 4. main ---------- */

func main() {
	// 標準入力をまとめて取得
	sc := bufio.NewScanner(os.Stdin)
	var lines []string
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	original := strings.Join(lines, " ")

	apiKey := os.Getenv("OPENAI_API_KEY")
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	/* A. API キーある → Chat + TTS */
	if apiKey != "" {
		roast, err := chatRoast(ctx, apiKey, original)
		if err == nil && roast != "" {
			if err := speakTTS(ctx, apiKey, roast); err == nil {
				return // 音声再生成功
			}
			log.Println("TTS 再生失敗:", err)
			// フォールバックしてテキスト出力
			fmt.Println(roast)
			return
		}
		log.Println("ChatCompletions 失敗:", err)
	}

	/* B. ローカルフォールバック */
	fmt.Println(localRoast(original))
}

