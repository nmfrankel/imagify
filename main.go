package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/sunshineplan/imgconv"
)

type intSlice []int

var PAGES intSlice
var PDF_PATH string
var OUTPUT_PATH string
var scale float64
var WIDTH int
var HEIGHT int
var FILE_TYPE string

func (i *intSlice) String() string {
	return fmt.Sprintf("%v", *i)
}

func (i *intSlice) Set(value string) error {
	// Trim brackets if they exist
	value = strings.Trim(value, "[]")
	if value == "" {
		return nil
	}

	// Split string by comma
	strValues := strings.Split(value, ",")
	for _, str := range strValues {
		num, err := strconv.Atoi(strings.TrimSpace(str))
		if err != nil {
			return err
		}
		*i = append(*i, num)
	}
	return nil
}

func main() {
	flag.StringVar(&PDF_PATH, "pdf_path", "", "Specify the path to the input PDF file. (Required)")
	flag.StringVar(&OUTPUT_PATH, "output_path", "./", "Specify the directory where output files will be saved. Defaults to the current directory.")
	flag.Float64Var(&scale, "scale", 100, "Set the scaling factor for the output image as a percentage. Defaults to 100 for 100%.")
	flag.IntVar(&WIDTH, "width", 0, "Set the width of the output image in pixels. Ignored if scale is provided.")
	flag.IntVar(&HEIGHT, "height", 0, "Set the height of the output image in pixels. Ignored if scale is provided.")
	flag.StringVar(&FILE_TYPE, "file_type", "png", "Specify the output image format. Supported formats: png, jpg, pdf, webp. Defaults to png.")
	flag.Var(&PAGES, "pages", "Specify the list of page numbers to process (e.g., --pages=1,2,3 or --pages=[1,2,3]).")
	flag.Parse()

	if PDF_PATH == "" {
		fmt.Println("Must provide a PDF file path using --pdf_path flag")
		return
	}

	_, err := os.Stat(PDF_PATH)
	if os.IsNotExist(err) {
		fmt.Println("Provided pdf_path does not exist:", PDF_PATH)
		return
	}

	if OUTPUT_PATH == "./" {
		// pwd, _ := os.Getwd()
		OUTPUT_PATH = strings.Split(PDF_PATH, ".")[0]
	}

	if err := os.MkdirAll(OUTPUT_PATH, os.ModePerm); err != nil {
		fmt.Println("Error creating specified output directory:", err)
		return
	}

	if len(PAGES) == 0 {
		fmt.Println("[WARN] No pages provided. Defaulting to all pages.")

		pageCount, err := api.PageCountFile(PDF_PATH)
		if err != nil {
			fmt.Println("Error getting page count:", err)
			return
		}
		for i := 0; i < pageCount; i++ {
			PAGES = append(PAGES, i+1)
		}
	}

	ctx, err := api.ReadContextFile(PDF_PATH)
	if err != nil {
		fmt.Println("Error reading PDF context:", err)
		return
	}

	for _, i := range PAGES {
		r, err := api.ExtractPage(ctx, i)
		if err != nil {
			fmt.Println("Error extracting page:", err)
			return
		}

		img, err := imgconv.Decode(r)
		if err != nil {
			fmt.Println("Error decoding image:", err)
			return
		}

		if scale != 100 {
			img = imgconv.Resize(img, &imgconv.ResizeOption{Percent: scale})
		} else if WIDTH != 0 && HEIGHT != 0 {
			img = imgconv.Resize(img, &imgconv.ResizeOption{Width: WIDTH, Height: HEIGHT})
		}

		formatMap := map[string]imgconv.Format{
			"png":  imgconv.PNG,
			"jpg":  imgconv.JPEG,
			"pdf":  imgconv.PDF,
			"webp": imgconv.WEBP,
		}

		format, ok := formatMap[strings.ToLower(FILE_TYPE)]
		if !ok {
			fmt.Println("Unsupported file type:", FILE_TYPE)
			return
		}

		filename := fmt.Sprintf("%s/%d.%s", OUTPUT_PATH, i, strings.ToLower(FILE_TYPE))
		if err := imgconv.Save(filename, img, &imgconv.FormatOption{Format: format}); err != nil {
			fmt.Println("Error saving image:", err)
			return
		}
	}
}
