# このサイトについて

[青空文庫](https://www.aozora.gr.jp/)で公開されている作品の検索および電子書籍リーダー用ファイル形式への変換をするためのサイトです。

本の探し方はほぼ自明と思われますが、特徴として各作品のページから、関連書物（同作家、同翻訳者、同NDC分類）にリンクされています。単純な検索も右上の検索ボックスからできます。

単純なブラウザ内での読書機能もあります。縦書き、ふりがな付きですが、基本的には本屋でパラパラっと中身を見る程度の利用を想定しています。全文を１ページで表示するため、文庫本で数ページ程度の文章なら読めるかもしれませんが、長文には向きません。

電子書籍リーダー用ファイル形式はEPUBとキンドル用のAZW3です。各作品のページにダウンロード用ボタンがあります。

MacOS、iOSの場合はEPUBをダウンロードしブック（Books）で開けば読めます。Androidの場合はEPUBでダウンロードした上Playブックで開けば読めます。

残念ながらキンドルのブラウザではこのサイトは機能しないので、パソコンからAZW3形式のファイルを転送する必要があります。リーダー機器への転送は[Calibre](https://calibre-ebook.com/ja)の使用をおすすめします。Calibreを使ってEPUB、AZW3を読むことも可能です。


書籍情報と、各作品は[青空文庫](https://www.aozora.gr.jp/)によって公開されているものであり、このサイトは青空文庫の皆様のご尽力の上に成り立っています。

なお、著作権が存続している作品は本サイトに含まれていません。


## 技術的な仕様

このサイトはいわゆる[シングルページアプリケーション](https://ja.wikipedia.org/wiki/%E3%82%B7%E3%83%B3%E3%82%B0%E3%83%AB%E3%83%9A%E3%83%BC%E3%82%B8%E3%82%A2%E3%83%97%E3%83%AA%E3%82%B1%E3%83%BC%E3%82%B7%E3%83%A7%E3%83%B3)です。各ページの構築及びEPUB/AZW3への変換はすべてブラウザ内で行われます。

コードはGo、それをWebAssemblyにコンパイルしたものです。ソースコードはGitHubで公開しています：[https://github.com/adamay909/AozoraBookcase](https://github.com/adamay909/AozoraBookcase)。不具合の報告等はそちらへお願いします。

AozoraBookcase. Copyright (C) 2024 Masahiro Yamada. Licensed under AGPL-3.0. Source code available at [https://github.com/adamay909/AozoraBookcase](https://github.com/adamay909/AozoraBookcase).

