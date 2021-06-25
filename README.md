[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](http://golang.org)

# ThanosDecryptor

ThanosDecryptor is an project to decrypt files encrypted by Thanos ransomware.

## Command Arguments
```
Usage of ./thanos_decrypt:
  -b string
        Custom search with byte value. (i.e. \xde\xad\xbe\xef -> deadbeef)
        Please use ?? to match any byte (i.e. de??beef)
  -c    Use current tickcount. (only support in Windows)
  -e string
        Search file extension.
  -i string
        Input encrypted file.
  -k string
        Decrypt with this key.
  -o string
        Output decrypted file.
  -p int
        Use n thread. (default 1)
  -r    Reversed tickcount.
  -s string
        Custom search with regular expression.
  -t int
        Start tickcount. (default 0)
```

## Usage
### Guess password
Guess the password of a png image from tickcount 0.
```bash
./thanos_decrypt -i ./sample/CyCraft.png.PROM\[prometheushelp@mail.ch\] -o ./output/CyCraft.png -e png -p 16
```

In this command, there are 4 arguments:
- i: input encrypted file
- o: output file
- e: search file format
- p: thread count

### Reversed Tickcount
Guess the password of a png image from tickcount 100000 in reversed order.
```bash
./thanos_decrypt -i ./sample/CyCraft.png.PROM\[prometheushelp@mail.ch\] -o ./output/CyCraft.png -e png -p 16 -t 100000 -r
```

There are 2 additional arguments:
- t: start from 100000
- r: reversed order (100000...0)

### Guess from current tickcount (only for Windows)
Guess the password of a png image from the current tickcount in reversed order. This feature is usually used with reversed order.
```bash
./thanos_decrypt -i ./sample/CyCraft.png.PROM\[prometheushelp@mail.ch\] -o ./output/CyCraft.png -e png -p 16 -c -r
```

There is an additional argument:
- c: start from the current tickcount

### Decrypt (Encrypt) with a key
Decrypt (Encrypt) a file with a provided key.
```bash
./thanos_decrypt -i ./sample/CyCraft.png.PROM\[prometheushelp@mail.ch\] -o ./output/CyCraft.png -k "+@[%T-mZSh+E[^^i{W:dpwnhdL4<b8D4}]]"
```

There is an additional argument:
- k: provided key

### Guess password with custom format (regular expression)
Guess the password of a text file with a known string "we had another great".
```bash
./thanos_decrypt -i ./sample/test.txt.enc -o ./output/test.txt -p 16 -s "we had another great"
```

There is an additional argument:
- s: regular expression to match the decrypted file

### Guess password with custom format (bytes pattern)
Guess the password of a png file with its header in hex.
```bash
./thanos_decrypt -i ./sample/test.txt.enc -o ./output/test.txt -p 16 -b '89??4e??0d??1a0a??00'
```

There is an additional argument:
- b: PNG header in hex format.
  - The full bytes are "8950 4e47 0d0a 1a0a 0000".
  - We can use ?? to match any byte.

Custom search with bytes pattern is much more convenient than regular expression, since there are lots of file format that it can't be performed by visible characters.



### Output
The output should like this. Since we match the file with magic number, it might be matched even a wrong key is provided. Therefore, we keep the decryption process continued to guess. You can terminate it anytime if you find the correct decrypted file.
```bash
 % ./thanos_decrypt -i ./sample/test.txt.enc -o ./output/test.txt -p 16 -s "we had another great"
 Decrypt file with seed 615750, key: +@[%T-mZSh+E[^^i{W:dpwnhdL4<b8D4, path: ./output/615750_test.txt
 2795306...
```

### GUI
We provide a GUI version for windows users. All features is supported in the GUI version. If you know nothing about programming, please follow the steps below to decrypt your files:

1. Choose a file to decrypt.
2. Choose the output file name.
3. Select "Use thread" and fill in 16. (Threads usually make the decryption routine faster, but it actually depends on amount of your cpu cores)
4. Select "Search extension" and fill in your file type. (For instance, PNG)
5. Click decrypt.
6. There is a counter, which shows the current guessing tickcount.
7. The decrypting result will show in the text block below. (There may be multiple possible key, so the decryption routine will continue to decrypt even find a possible key. You can terminate it at any time.)
8. Since the tickcounts (seeds) used to encrypt are near, you can try to record the seed above and select "Start tickcount" with value `seed-10000` next time. It may be faster.

![GUI](https://raw.githubusercontent.com/cycraft-corp/ThanosDecryptor/master/GUI.png)

## Build
```bash
make win32    # windows 32 bits
make win64    # windows 64 bits
make linux    # linux
make win32GUI # windows 32 bits GUI (built on windows)
make win64GUI # windows 64 bits GUI (build on windows)
```

## Supported File Format
We match the magic number with https://github.com/h2non/filetype. 
Here is the file type we currently support:

### Image

- **jpg** - `image/jpeg`
- **png** - `image/png`
- **gif** - `image/gif`
- **webp** - `image/webp`
- **cr2** - `image/x-canon-cr2`
- **tif** - `image/tiff`
- **bmp** - `image/bmp`
- **heif** - `image/heif`
- **jxr** - `image/vnd.ms-photo`
- **psd** - `image/vnd.adobe.photoshop`
- **ico** - `image/vnd.microsoft.icon`
- **dwg** - `image/vnd.dwg`

### Video

- **mp4** - `video/mp4`
- **m4v** - `video/x-m4v`
- **mkv** - `video/x-matroska`
- **webm** - `video/webm`
- **mov** - `video/quicktime`
- **avi** - `video/x-msvideo`
- **wmv** - `video/x-ms-wmv`
- **mpg** - `video/mpeg`
- **flv** - `video/x-flv`
- **3gp** - `video/3gpp`

### Audio

- **mid** - `audio/midi`
- **mp3** - `audio/mpeg`
- **m4a** - `audio/m4a`
- **ogg** - `audio/ogg`
- **flac** - `audio/x-flac`
- **wav** - `audio/x-wav`
- **amr** - `audio/amr`
- **aac** - `audio/aac`

### Archive

- **epub** - `application/epub+zip`
- **zip** - `application/zip`
- **tar** - `application/x-tar`
- **rar** - `application/vnd.rar`
- **gz** - `application/gzip`
- **bz2** - `application/x-bzip2`
- **7z** - `application/x-7z-compressed`
- **xz** - `application/x-xz`
- **zstd** - `application/zstd`
- **pdf** - `application/pdf`
- **exe** - `application/vnd.microsoft.portable-executable`
- **swf** - `application/x-shockwave-flash`
- **rtf** - `application/rtf`
- **iso** - `application/x-iso9660-image`
- **eot** - `application/octet-stream`
- **ps** - `application/postscript`
- **sqlite** - `application/vnd.sqlite3`
- **nes** - `application/x-nintendo-nes-rom`
- **crx** - `application/x-google-chrome-extension`
- **cab** - `application/vnd.ms-cab-compressed`
- **deb** - `application/vnd.debian.binary-package`
- **ar** - `application/x-unix-archive`
- **Z** - `application/x-compress`
- **lz** - `application/x-lzip`
- **rpm** - `application/x-rpm`
- **elf** - `application/x-executable`
- **dcm** - `application/dicom`

### Documents

- **doc** - `application/msword`
- **docx** - `application/vnd.openxmlformats-officedocument.wordprocessingml.document`
- **xls** - `application/vnd.ms-excel`
- **xlsx** - `application/vnd.openxmlformats-officedocument.spreadsheetml.sheet`
- **ppt** - `application/vnd.ms-powerpoint`
- **pptx** - `application/vnd.openxmlformats-officedocument.presentationml.presentation`

### Font

- **woff** - `application/font-woff`
- **woff2** - `application/font-woff`
- **ttf** - `application/font-sfnt`
- **otf** - `application/font-sfnt`

### Application

- **wasm** - `application/wasm`
- **dex** - `application/vnd.android.dex`
- **dey** - `application/vnd.android.dey`

## How it work ?
Thanos ransomware use salsa20 with a tickcount-based random password to encrypt. The size of the random password is 32 bytes, and every character is visible character. Since the password use tickcount as the key, we can guess it brutally.
