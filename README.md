# duplicate_image_detector
重複している画像をコマンドで検出し、結果をCSV, html出力する

dhashで画像をハッシュ化して類似度を算出

Goの勉強用で作った

削除処理はこれから作る



## Setup
```
go get
```


## 実行
```
go run main.go -dir <探索対象フォルダのパス> -threshold 1
```

## ビルド
```
go build main.go
```


## テスト実行
```
go test -v ./...
```


## カバレッジ付きでテスト実行
```
go test -v -cover ./...

# カバレッジをwebで見る
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## ベンチマーク計測
Benchmark~とつく関数のパフォーマンスを計測できる。
```
cd <計測したいパッケージ>
go test -bench=. -benchmem
```

## プロファイリング
ベンチマーク計測して各処理のメモリ・CPU消費などを計測できる。
```
cd <計測したいパッケージ>

# 性能計測
go test -memprofile=mem.out -bench=.
go test -blockprofile=block.out -bench=.
go test -cpuprofile=cpu.out -bench=.

# pprofツールで可視化
go tool pprof -text -nodecount=10 <計測したいパッケージ> block.out
```


## コマンドライン上でドキュメンテーションを見る
パッケージや公開関数などをコマンドライン上で確認できる
```
go doc <パッケージ名>

go doc <パッケージ名>.<関数/構造体など>

# 例
go doc ./struct_sample
go doc ./struct_sample.TestStruct
```


## ドキュメンテーション一覧
Web上で標準パッケージ、自作パッケージのドキュメントが見れる。
```
go get -v golang.org/x/tools/cmd/godoc
godoc -http ":3000"
```
