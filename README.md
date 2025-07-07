# basay

日本語テキストを毒舌に変換し、音声で再生するコマンドラインツール

## 概要

basayは、入力されたテキスト（主に日本語）を毒舌・辛辣な返答に変換するツールです。OpenAI APIを使用して創造的な毒舌を生成し、Text-to-Speechで音声出力することができます。API キーがない場合は、ローカルのパターンマッチングによる毒舌生成にフォールバックします。

## 機能

- **ローカル毒舌生成**: よくある日本語フレーズに対するパターンマッチング
- **ChatGPT統合**: OpenAI GPT-4o-miniを使用した創造的な毒舌生成
- **音声出力**: OpenAI TTS APIを使用した音声再生
- **グレースフルデグラデーション**: API利用不可時のローカル処理への自動切り替え

## インストール

### 必要な依存関係

- Go 1.24.4以上
- FFmpeg（音声再生用の`ffplay`コマンド）

### ビルド

```bash
go build -o basay main.go
```

## 使用方法

### OpenAI APIを使用する場合

```bash
export OPENAI_API_KEY=your-api-key-here
echo "遅れてすみません" | ./basay
```

### ローカルモードのみ（API キーなし）

```bash
echo "遅れてすみません" | ./basay
```

### 使用例

```bash
# 遅刻の言い訳
echo "電車が遅れて..." | ./basay

# 会議を欠席
echo "会議に出られません" | ./basay

# やる気がない
echo "やる気が出ません" | ./basay

# 一般的な謝罪
echo "申し訳ございません" | ./basay
```

## 技術仕様

- **言語**: Go
- **OpenAI API**: GPT-4o-mini（テキスト生成）、TTS-1（音声生成）
- **音声形式**: WAV
- **音声再生**: ffplay（FFmpeg）

## 環境変数

- `OPENAI_API_KEY`: OpenAI APIキー（オプション）

## ライセンス

MIT License
