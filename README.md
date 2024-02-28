# Flow Limit Proxy

Flow Limit Proxyは、HTTPの同時通信制限をかけられるGo製のHTTPプロキシです。

## 特徴

- HTTP通信のプロキシ
- 同時通信数の上限設定
- 通信エラー時のリトライ

※ 対応しているのは、localhostのポート間のみです。

## インストール

```bash
go install github.com/bellwood4486/flow-limit-proxy@latest
```

## 使い方

```
Usages:
  flow-limit-proxy [-limit=<number>] <fromPort>:<toPort>
Options:
  -limit int
        concurrent transfer limit (default 10)
```

## ライセンス

MIT
