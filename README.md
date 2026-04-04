# go-stations

## 概要
Go言語を用いて実装したWeb APIサーバです。  
DBを利用したTODO管理機能に加え、Middleware機構によるログ出力やPanic Recoveryなど、実運用を意識した設計を行っています。

「単に動作するアプリ」ではなく、**保守性・拡張性を考慮したバックエンド設計の理解**を目的として開発しました。

---

## 主な機能

- TODOのCRUD API
- データベースによるデータ永続化
- Middlewareによるリクエストログ出力
- Panic Recovery（サーバのクラッシュ防止）
- 環境情報（OSなど）の取得API

---

## 工夫した点

### ① Middlewareによる関心の分離
ログ出力やエラーハンドリングといった横断的関心を、handlerから分離しました。  
これにより、コードの可読性と再利用性を向上させています。

---

### ② レイヤ構造の設計
以下のように責務ごとにディレクトリを分割しています。

/handler // HTTPリクエスト処理
/router // ルーティング定義
/db // DB操作
/middleware // 共通処理

これにより、機能追加時の影響範囲を限定し、保守しやすい構成としています。

---

### ③ 標準ライブラリ中心の実装
`net/http` を用いてサーバを構築し、フレームワークに依存しない設計とすることで、HTTP処理の基礎理解を重視しました。

---

### ④ Panic Recoveryの実装
`defer` と `recover()` を用いて、サーバがクラッシュせずにエラーレスポンスを返せるよう設計しています。

---

## 技術スタック

- Go
- net/http
- （使用していれば）SQLite / MySQL など

---

## セットアップ方法

```bash
git clone https://github.com/take08042004/go-stations.git
cd go-stations
go run main.go
