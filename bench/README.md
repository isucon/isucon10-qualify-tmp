# bench

問題サーバーに対して負荷を掛けるベンチマーカー

### ビルド

```sh
make
```


### ベンチマーク

```sh
# localhost:1323 (デフォルト) に向けてベンチマーク実行
./bench

# ターゲットサーバーを指定する
TARGET_SERVER_URL = http://localhost:1323
./bench --target-url $TARGET_SERVER_URL

# 初期データのディレクトリを指定する
./bench --data-dir ../initial-data

# fixtureのディレクトリを指定する
./bench --fixture-dir ../webapp/fixture
```
