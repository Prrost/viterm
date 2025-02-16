package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

var asciiChars = "@%#*+=-:. "

func main() {
	videoPath := os.Args[1]

	tempDir, err := VideoSetup(videoPath)
	if err != nil {
		fmt.Println(err)
	}
	defer os.RemoveAll(tempDir)

	frames, err := LoadFrames(tempDir)
	if err != nil {
		fmt.Println(err)
		return
	}

	StartVideo(frames)

}

func VideoSetup(videoPath string) (string, error) {
	projectDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("could not get project dir: %w", err)
	}

	tempDir, err := os.MkdirTemp(projectDir, "temp")
	if err != nil {
		return "", fmt.Errorf("could not create temp dir: %w", err)
	}

	err = VideoToPNG(tempDir, videoPath)
	if err != nil {
		return "", err
	}

	return tempDir, nil
}

func VideoToPNG(tempDir, videoPath string) error {
	outputPattern := filepath.Join(tempDir, "frame_%04d.png")
	cmd := exec.Command("ffmpeg", "-i", videoPath, "-vf", "fps=10,scale=70:-1", outputPattern)

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("could not convert %s to PNG: %w", outputPattern, err)
	}
	return nil
}

func LoadFrames(tempDir string) ([]image.Image, error) {
	files, err := filepath.Glob(filepath.Join(tempDir, "frame_*.png"))
	if err != nil {
		return nil, err
	}

	var frames []image.Image
	for _, f := range files {
		img, err := LoadImage(f)
		if err != nil {
			fmt.Println("Ошибка загрузки кадра:", f, err)
			continue
		}
		frames = append(frames, img)
	}
	return frames, nil
}

func LoadImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("could not decode image: %w", err)
	}

	return img, nil
}

func ImageToAscii(img image.Image) string {
	bounds := img.Bounds()
	var ascii string

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			gray := color.GrayModel.Convert(img.At(x, y)).(color.Gray)
			index := int(gray.Y) * (len(asciiChars) - 1) / 255
			ascii += string(asciiChars[index])
			ascii += " "
		}
		ascii += "\n"
	}
	return ascii
}

func StartVideo(frames []image.Image) {
	for _, frame := range frames {
		ClearTerminal()

		ascii := ImageToAscii(frame)
		fmt.Println(ascii)

		time.Sleep(100 * time.Millisecond)
	}
}

func ClearTerminal() {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	default:
		cmd = exec.Command("clear")
	}

	cmd.Stdout = os.Stdout
	cmd.Run()
}
