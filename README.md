# AozoraBookcase

青空文庫のテキストを手っ取り早くEPUB3、Kindle、その他のフォーマットで入手できるようにするものです。

青空文庫のテキストを複数のサイトやアプリを開かずに、EPUB3やKindle用の形式で入手できるように開発しました。

##  実装サイト

本アプリを実装したサイトを公開しています：

https://aozora.orihasam.com/

これはAWSのS3の静的ウエブサイト機能を利用したものです。EPUB等への変換はすべてローカルブラウザ内で行われます。

## 実装に必要なもの

- 青空文庫のサイトのコピー：[https://github.com/aozorabunko/aozorabunko](https://github.com/aozorabunko/aozorabunko)のクローンでOK。実際に使うのは  
cards/\*/files/\*.zip   
のパターンに一致するファイルと  
index_pages/list_person_all_extended_utf8.zip   
のみなので、残りは削除してもかまいません。

- httpsに対応し、静的ウェブサイトを配信できるサーバー。

以上の要件を満たしていれば、 Makefile.exampleをMakefileとファイル名を変え、
```
go mod tidy
go get -u
make build
```

の後、`dist`ディレクトリ内のファイルをすべて青空文庫サイトコピーのトップフォルダにコピーすれば、青空文庫のサイトコピーをウェブサーバーで配信開始するだけです。

あとは　https://\[*青空文庫のコピーのURL*\]/aozobookcase.htmlにブラウザで行けば青空文庫の検索とEPUB、AZW3（Kindle）その他の形式のファイルのダウンロードができます。アプリ内のナビゲーションの仕方は、見ればわかると思います。

GoからWebAssemblyにコンパイルしており、バイナリが巨大（9MB）なためもあり、起動に少し時間がかかります（2017年版のThinkpadX1Carbonで５秒弱、iPadAir M1で２秒弱、Google Pixel 7aも同程度）。そこから先は個人的感想では重いという印象は受けません。例えばかなり長い作品の部類に入る谷崎の『細雪』の上巻は、AZW3の変換に2017年版のThinkpadX1Carbonで二秒弱かかります。大多数のテキストではページ移動に通常かかる時間と違いは感じられません。


### EPUB/AZW3/その他フォーマットへの 変換ライブラリ

変換に使うコードはもともとコマンドラインツールとして作ったもので、別途[https://github.com/adamay909/AozoraConvert](https://github.com/adamay909/AozoraConvert)で公開しています。本アプリで使用しているのはレポジトリ内のv2の方です。


