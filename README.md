# AozoraBookcase

<<<<<<< HEAD
青空文庫のテキストを手っ取り早くEPUB3及びKindle用の形式で入手できるようにするサーバーです。AozoraBookcaseのパッケージがサーバー、AozoraFSのパッケージがバックエンドです。
=======


青空文庫のテキストを複数のサイトやアプリを開かずに、手っ取り早くEPUB3及びKindle用の形式で入手できるように開発しました。

シングルページアプリケーション版とサーバー版があります。

## シングルページアプリケーション (SPA)

メインのソースコードはabSPAフォルダにあります。

必要なもの：
- 青空文庫のサイトのコピー：[https://github.com/aozorabunko/aozorabunko](https://github.com/aozorabunko/aozorabunko)のクローンでOK。実際に使うのは  
cards/\*/files/\*.{html,png}   
のパターンに一致するファイルと  
index_pages/list_person_all_extended_utf8.zip   
のみなので、残りは削除してもかまいません。

- httpsに対応し、静的ウェブサイトを配信できるサーバー。

以上の要件を満たしていれば、abSPAディレクトリでmakeを実行後、serverfilesディレクトリ内のファイルをすべて青空文庫サイトコピーのトップフォルダにコピーし、青空文庫のサイトコピーをウェブサーバーで配信開始するだけです。

あとは　https://\[*青空文庫のコピーのURL*\]/aozobookcase.htmlにブラウザで行けば青空文庫の検索とEPUBおよびAZW3ファイルのダウンロードができます。アプリ内のナビゲーションの仕方は、見れば解ると思います。

GoからWebAssemblyにコンパイルしており、バイナリが巨大（22MB）なためもあり、起動に少し時間がかかります（２０１７年版のThinkpadX1Carbonで５秒弱、iPadAir M1で２秒弱、Google Pixel 7aも同程度）。そこから先は個人的感想では重いという印象は受けません。例えばかなり長い作品の部類に入る谷崎の『細雪』の上巻は、AZW3の変換に２０１７年版のThinkpadX1Carbonで二秒弱かかります。


### サーバー版

メインのソースコードはaozoraBookcaseフォルダにあります。

自前のサーバーを運用できる環境が必要となります。
>>>>>>> f8f7ea4 (edit README)

青空文庫からの必要フォルダ、ファイルの作成及びダウンロードはすべて自動で行われるので、サーバーを立ち上げるための準備は特に必要ありません。AozoraBookcaseをbuildとinstall後に

     $ aozoraBookcase

で問題なくスタートできるはずです。最初は作品データベースのみのダウンロードです。残りのファイルは必要に応じて随時ダウンロードされます。

localhost:3333よりサーバーにアクセスできます。メインページは作者リスト。各作者名をクリックするとその作者の本のリストにうつり、各本のリンクをクリックすればEPUB、AZW3(Kindle)それぞれのダウンロードリンクが表示されます。そのままブラウザで縦書き表示することもできます。

簡単な検索機能もあります。

サーバー立ち上げ時のコマンドオプション:

	-c bool
        すべてのデータを青空文庫からダウンロードし直す。

	-children bool
	  	児童書のみ表示する。児童書かどうかは青空文庫のデータベースのNDC分類にKが追加されているかどうかで判別する。ただし、児童書でも新字新仮名出ないものは除外する。

	-d string
	  	サーバーファイルのローカルディレクトリ。既定値は $HOME/aozorabunko。作成及びダウンロードされるファイルは全てこれのサブフォルダーにある。

	-i string
	  	サーバーのインターフェース。既定値はすべてのインターフェースを利用する。

	-p string
	  	サーバーのポート番号。既定値は3333。

	-refresh string
	  	青空文庫本家の更新有無をチェックする間隔。

	-strict bool
	  	著作権が切れているもののみを表示するかどうか。

	-v bool
        ログをスクリーンにも出力する。


自宅サーバーで家族の使用を念頭に作ったので、サーバーを一般公開した場合の負荷に耐えるかなどは試していません。





