---
title: koikoi
import_path: go.ngs.io/koikoi
repo_url: https://github.com/ngs/go-koikoi
description: Go 製の CUI 花札こいこいゲーム（任天堂ルール準拠）
version: v0.0.6
documentation_url: https://pkg.go.dev/go.ngs.io/koikoi
license: MIT
author: Atsushi Nagase
created_at: 2026-03-08T16:59:29Z
updated_at: 2026-04-11T02:01:32Z
---

# koikoi

[![CI](https://github.com/ngs/go-koikoi/actions/workflows/ci.yml/badge.svg)](https://github.com/ngs/go-koikoi/actions/workflows/ci.yml)

Go 製の CUI 花札こいこいゲーム（任天堂ルール準拠）。

![Screenshot](https://github.com/ngs/go-koikoi/raw/master/docs/screenshot.jpg)

## 機能

- 🎴 任天堂ルール準拠のこいこい
- 🖱️ マウスクリック対応
- 🎰 親決め演出（ルーレット風アニメーション）
- 🤖 3段階の CPU 難易度
- 💾 ゲーム進捗の自動セーブ
- 🏆 リーチ（テンパイ）表示

## インストール

### Go

```bash
go install go.ngs.io/koikoi@latest
```

### Homebrew

```bash
brew tap ngs/tap
brew install ngs/tap/koikoi
```

### ソースからビルド

```bash
git clone https://github.com/ngs/go-koikoi.git
cd go-koikoi
go build -o koikoi .
./koikoi
```

## 使い方

```bash
koikoi
```

設定ディレクトリを指定する場合:

```bash
koikoi -config /path/to/config
```

## 操作方法

### キーボード

| キー | 操作 |
|------|------|
| ↑↓ / j k | 手札カーソル移動 |
| ←→ / h l | 場札・選択肢カーソル移動 |
| Enter / Space | 選択・決定 |
| ? | ヘルプ表示 |
| l | 行動履歴表示 |
| o | オプション設定 |
| q | 終了（確認あり） |
| Esc | ポップアップを閉じる |

### マウス

- 手札・場札・ボタンをクリックして選択・決定

## ルール

花札こいこいのルール詳細は [docs/rules.md](docs/rules.md) を参照してください。

### 役一覧

| 役名 | 点数 | 条件 |
|------|------|------|
| 五光 | 10文 | 光札5枚すべて |
| 四光 | 8文 | 柳を除く光札4枚 |
| 雨四光 | 7文 | 柳を含む光札4枚 |
| 三光 | 5文 | 柳を除く光札3枚 |
| 猪鹿蝶 | 5文+ | 萩に猪＋紅葉に鹿＋牡丹に蝶 |
| 赤短・青短 | 10文+ | 赤短と青短の両方成立 |
| 赤短 | 5文+ | 松・梅・桜の赤短冊3枚 |
| 青短 | 5文+ | 牡丹・菊・紅葉の青短冊3枚 |
| 花見で一杯 | 5文 | 桜に幕＋菊に盃 |
| 月見で一杯 | 5文 | 芒に月＋菊に盃 |
| タネ | 1文+ | 種札5枚以上 |
| タン | 1文+ | 短冊札5枚以上 |
| カス | 1文+ | カス札10枚以上 |

## CPU 難易度

| 難易度 | 説明 |
|--------|------|
| かんたん | ランダム要素が多い。こいこいしない |
| ふつう | 札の価値を評価して最善手を選ぶ |
| つよい | 役への近さも考慮。積極的にこいこいする |

## 設定ファイル

| パス | 内容 |
|------|------|
| `~/.koikoi/settings.json` | ラウンド数、CPU 難易度 |
| `~/.koikoi/game.json` | ゲーム進捗のセーブデータ |

## 開発

### 必要環境

- Go 1.21 以上

### テスト

```bash
go test -v ./...
```

カバレッジ付きテスト:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Lint

```bash
golangci-lint run
```

### プロジェクト構成

| ファイル | 概要 |
|----------|------|
| `main.go` | エントリポイント。CLI フラグ解析、設定・セーブデータ読込、UI 起動 |
| `card.go` | Card 構造体、全48枚の定義、月・札種の定数 |
| `yaku.go` | 役判定ロジック（任天堂ルール準拠、排他制御含む） |
| `game.go` | Game 構造体、ラウンド管理、マッチング処理、親決め |
| `cpu.go` | CPU AI（3段階の難易度: かんたん/ふつう/つよい） |
| `ui.go` | gocui ベースの TUI。画面描画、キーバインド、マウス操作、ゲームフェーズ管理 |
| `settings.go` | 設定の永続化（ラウンド数、難易度） |
| `save.go` | ゲーム進捗のセーブ・ロード |

## ライセンス

[MIT](LICENSE)
