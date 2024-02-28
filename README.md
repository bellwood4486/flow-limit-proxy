# Local Flow Proxy

Local Flow Proxyは、ローカルホスト内のポート間のHTTP通信を制御するためのGo製コマンドラインツールです。  

## 特徴

- HTTPプロトコルに基づく通信のプロキシ
- ローカルホスト内でのポート間の通信管理
- 同時通信数の上限設定
- 通信エラー時のリトライ

## インストール

```bash
go install github.com/bellwood4486/local-flow-proxy@latest
```

## 使い方

./local-flow-proxy [-limit=<number>] <fromPort>:<toPort>

## ライセンス

MIT
