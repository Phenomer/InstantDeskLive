# InstantDeskLive
Windowsのデスクトップを手軽にライブ配信するツールです。
HLSで、遅延数秒程度のライブ配信ができます。

LAN内でデスクトップ映像をシュッと共有したい時などに便利です。


## 必要なもの

* [ffmpegのバイナリ(デスクトップキャプチャ・エンコード)](https://www.ffmpeg.org/download.html)
* [Caddyのバイナリ(Webサーバー)](https://caddyserver.com/download)

まとめて`winget`でインストールできます。

```
C:\>winget install Gyan.FFmpeg CaddyServer.Caddy
```

### 動作確認済みバージョン

ffmpeg-7.1.1と、caddy-2.10.0で動作しました。

```
C:\>ffmpeg -version
ffmpeg version 7.1.1-essentials_build-www.gyan.dev Copyright (c) 2000-2025 the FFmpeg developers
...
```

```
C:\>caddy -v
v2.10.0 h1:fonubSaQKF1YANl8TXqGcn4IbIRUDdfAkpcsfI/vX5U=
```

`ffmpeg`と`caddy`のインストールができたら、後は`desktop_live.exe`を実行するだけです。


## つかいかた

`desklive.exe`をダブルクリックするだけです。


## エンコード・配信対象画面の設定

`process.json`のffmpegの引数を変更し再起動して調整します。

例えば、FullHD(1920x1080)のディスプレイが2台ある時、右側の2台目の映像だけを配信したい場合は、
`-offset_y`を`1920`にします。

```
            "-offset_x",
            "0",
            "-offset_y",
            "1920",
```
