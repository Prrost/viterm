# Video to ASCII Converter

This project converts a black-and-white video into an ASCII representation and displays it in the terminal frame by frame.

## Requirements
- Go 1.20+
- FFmpeg (for extracting frames from video)

## FFmpeg
if you dont have it:
- Macos:
```sh
brew install ffmpeg
```
- Linux:
```sh
sudo apt install ffmpeg
```
- Windows:
```sh
echo real way to download: winget install Gyan.FFmpeg > nul & echo Windows sucks
```

## Installation

Clone the repository and navigate to the project folder:
```sh
git clone https://github.com/Prrost/viterm.git
cd viterm
go run main.go <path to your video>
```
