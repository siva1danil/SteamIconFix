# SteamIconFix

**SteamIconFix** is a small program designed to fix missing Steam game icons. It aims to provide a quick user-friendly solution for ensuring that Steam shortcut icons are displayed correctly. It's especially useful after OS or Steam reinstallation.

## Features

* Quickly fixes Steam icons for all games (works with the ones without existing shortcuts, unlike some existing tools).
* Easy to use.
* Works locally without using any third-party services.
* Lightweight and (almost) portable*.

\* - The program cleans up created files on normal exit on Windows. In case of force exit, the files might remain in %AppData%.

## Usage

1. [Download the latest release](https://github.com/siva1danil/SteamIconFix/releases) (or compile the program yourself) and run the executable.
2. Type in your Steam directory (if not detected automatically).
3. Click "Load games from Steam".
4. Click "Fix icons" and wait.
5. Press "Refresh / F5" on your Desktop to reload the icons.
6. Done!

## System Requirements

* Windows (tested on Windows 11; might work on other OSes).  
* Steam must be installed.
* Any Steam game must be launched at least once.

## Demo

<img src="img/demo.gif" width="500">

## TODO

* Check cross-platform compatibility.
* Check game compatibility.

## License

```
The MIT License (MIT)

Copyright (c) 2025 siva1danil

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```