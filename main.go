package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	keep := flag.Bool("k", false, "keep image size")
	quality := flag.Int("q", 0, "quality parameter")
	size := flag.Int("s", 3000, "max size")
	format := flag.String("f", "", "forced output format")
	flag.Parse()
	if len(flag.Args()) != 1 {
		fmt.Fprintln(os.Stderr, "error: path to media is required")
		fmt.Fprintln(os.Stderr, "usage: shrink path/to/media")
		os.Exit(1)
	}
	input := flag.Arg(0)

	info, err := os.Stat(input)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	var paths []string
	if info.IsDir() {
		entries, err := os.ReadDir(input)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
		for _, e := range entries {
			if e.IsDir() {
				continue
			}
			if strings.Contains(e.Name(), "_small") {
				continue
			}
			paths = append(paths, filepath.Join(input, e.Name()))
		}
	} else {
		paths = append(paths, input)
	}

	for _, path := range paths {
		err := convert(path, *format, *keep, *quality, *size)
		if err != nil {
			fmt.Fprintf(
				os.Stderr,
				"error: failed converting %v: %v\n",
				path,
				err,
			)
			os.Exit(1)
		}
		fmt.Println("shrunk:", path)
	}
}

func convert(path, format string, keep bool, quality, size int) error {
	ext := strings.ToLower(filepath.Ext(path))
	if format != "" {
		ext = "." + format
	}

	switch ext {
	case ".jpg", ".jpeg":
		return jpg(path, keep, size)
	case ".tif":
		return jpg(path, keep, size)
	case ".png":
		return png(path, keep, size)
	case ".mp4":
		return webm(path, keep, quality)
	case ".mkv":
		return webm(path, keep, quality)
	case ".mov":
		return webm(path, keep, quality)
	default:
		return fmt.Errorf("cannot handle %v", ext)
	}
}

func magick(input, output string, keep bool, size int) error {
	args := []string{
		input,
		"-delete",
		"1--1",
	}
	if !keep {
		s := strconv.Itoa(size) + "x" + strconv.Itoa(size) + ">"
		args = append(args, "-scale", s)
	}
	args = append(args, output)
	magick := exec.Command("magick", args...)
	return magick.Run()
}

func ffmpeg(input, output string, keep bool, quality int) error {
	if quality == 0 {
		quality = 35
	}
	args := []string{
		"-i",
		input,
		"-y",
		"-c:v",
		"libvpx-vp9",
		"-crf",
		strconv.Itoa(quality),
	}
	args = append(args, []string{
		"-c:a",
		"libopus",
	}...)
	if !keep {
		args = append(args, []string{
			"-vf",
			"scale=-1:1080",
		}...)
	}
	args = append(args, output)
	fmt.Println(args)
	ffmpeg := exec.Command("ffmpeg", args...)
	ffmpeg.Stderr = os.Stderr
	return ffmpeg.Run()
}

func jpg(path string, keep bool, size int) error {
	output := strings.TrimSuffix(path, filepath.Ext(path)) + "_small.jpg"
	if err := magick(path, output, keep, size); err != nil {
		return err
	}
	jpegoptim := exec.Command(
		"jpegoptim",
		"-m",
		"92",
		output,
	)
	if err := jpegoptim.Run(); err != nil {
		return err
	}
	return nil
}

func png(path string, keep bool, size int) error {
	output := strings.TrimSuffix(path, filepath.Ext(path)) + "_small.png"
	if err := magick(path, output, keep, size); err != nil {
		return err
	}
	optipng := exec.Command(
		"optipng",
		output,
	)
	if err := optipng.Run(); err != nil {
		return err
	}
	return nil
}

func webm(path string, keep bool, quality int) error {
	output := strings.TrimSuffix(path, filepath.Ext(path)) + "_small.webm"
	return ffmpeg(path, output, keep, quality)
}
