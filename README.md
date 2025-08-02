# AozoraBookcase

青空文庫のテキストを手っ取り早くEPUB3及びKindle用の形式で入手できるようにするものです。

青空文庫のテキストを複数のサイトやアプリを開かずに、手っ取り早くEPUB3及びKindle用の形式で入手できるように開発しました。


必要なもの：
- 青空文庫のサイトのコピー：[https://github.com/aozorabunko/aozorabunko](https://github.com/aozorabunko/aozorabunko)のクローンでOK。実際に使うのは  
cards/\*/files/\*.{html,png}   
のパターンに一致するファイルと  
index_pages/list_person_all_extended_utf8.zip   
のみなので、残りは削除してもかまいません。

- httpsに対応し、静的ウェブサイトを配信できるサーバー。

以上の要件を満たしていれば、go get -u した後、make buildを実行し、serverfilesディレクトリ内のファイルをすべて青空文庫サイトコピーのトップフォルダにコピーすれば、青空文庫のサイトコピーをウェブサーバーで配信開始するだけです。

あとは　https://\[*青空文庫のコピーのURL*\]/aozobookcase.htmlにブラウザで行けば青空文庫の検索とEPUBおよびAZW3ファイルのダウンロードができます。アプリ内のナビゲーションの仕方は、見れば解ると思います。

GoからWebAssemblyにコンパイルしており、バイナリが巨大（22MB）なためもあり、起動に少し時間がかかります（２０１７年版のThinkpadX1Carbonで５秒弱、iPadAir M1で２秒弱、Google Pixel 7aも同程度）。そこから先は個人的感想では重いという印象は受けません。例えばかなり長い作品の部類に入る谷崎の『細雪』の上巻は、AZW3の変換に２０１７年版のThinkpadX1Carbonで二秒弱かかります。


### EPUB/AZW3 変換ライブラリ

EPUB、AZW3への変換に使うコードはもともとコマンドラインツールとして作ったもので、別途[https://github.com/adamay909/AozoraConvert](https://github.com/adamay909/AozoraConvert)で公開しています。


