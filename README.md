# isucon10-qualify

## ディレクトリ構成

```
.
├── bench           # ベンチマーカー
├── initial-data    # 初期データの生成
├── provisioning    # セットアップ用
└── webapp          # 各言語の参考実装
```

## 問題の起動方法

1. `initial-data` で初期データを生成する
2. `webapp` で Docker を用いて問題サーバーを立ち上げる
3. `bench` で問題サーバーへのベンチマークを実行する

実際のコマンド例については、各ディレクトリの README を参照してください。


## 使用データの取得元

- [Faker](https://faker.readthedocs.io/)
- [いらすとや](https://www.irasutoya.com/)

## ISUCON10 予選のインフラ構成について

- 問題用 (3台)
    - CPU: 1 Core (AMD EPYC 7352)
    - Memory: 2 GiB
    - IO throughput: 100 Mbps
    - IOPS limit: 200 (Read / Write)
    - Interface: 1 Gbps
- ベンチマーカ用 (1台)
    - CPU: 1 Core  (AMD EPYC 7352)
    - Memory: 16 GiB
    - IO throughput: 100 Mbps
    - IOPS limit: 200 (Read / Write)
    - Interface: 100 Mbps

## Links

- [ISUCON10 予選レギュレーション](http://isucon.net/archives/54753430.html)
- [ISUCON10 予選当日マニュアル](https://gist.github.com/progfay/25edb2a9ede4ca478cb3e2422f1f12f6)
