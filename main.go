package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
	"sync"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	logger "github.com/sirupsen/logrus"
	"github.com/sunshineplan/imgconv"
)

type intSlice []int

var formatMap = map[string]imgconv.Format{
	"png":  imgconv.PNG,
	"jpg":  imgconv.JPEG,
	"jpeg": imgconv.JPEG,
	"pdf":  imgconv.PDF,
	"webp": imgconv.WEBP,
}

var PAGES intSlice
var PDF_PATH string
var OUTPUT_PATH string
var SCALE float64
var WIDTH int
var HEIGHT int
var FILE_TYPE string
var DEBUG bool

func init() {
	flag.StringVar(&PDF_PATH, "pdf_path", "", "Path to the input PDF file. (Required)")
	flag.StringVar(&OUTPUT_PATH, "output_path", "", "Directory where output files will be saved. Defaults to the current directory.")
	flag.Float64Var(&SCALE, "scale", 100, "Scaling factor for the output image as a percentage. Defaults to 100%.")
	flag.IntVar(&WIDTH, "width", 0, "Width of the output image in pixels. Ignored if scale is provided.")
	flag.IntVar(&HEIGHT, "height", 0, "Height of the output image in pixels. Ignored if scale is provided.")
	flag.StringVar(&FILE_TYPE, "file_type", "png", "Output image format. Supported formats: png, jpg, pdf, webp. Defaults to png.")
	flag.Var(&PAGES, "pages", "List of page numbers to process (e.g., --pages=1,2,3 or --pages=[1,2,3]).")
	flag.BoolVar(&DEBUG, "debug", false, "Print debug logs.")
	flag.Parse()

	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logger.TextFormatter{
		ForceColors:     true,
		TimestampFormat: "02-01-2006 15:04:05",
		FullTimestamp:   true,
		PadLevelText:    true,
	})
	if DEBUG {
		logger.SetLevel(logger.DebugLevel)
	}
}

func main() {
	logger.Info("Starting Imagify...")
	if PDF_PATH == "" {
		logger.Error("Please provide the path to the PDF file using the --pdf_path flag.")
		return
	}

	_, err := os.Stat(PDF_PATH)
	if os.IsNotExist(err) {
		logger.Errorf("The specified PDF file does not exist (%s).", PDF_PATH)
		logger.Debug(err)
		return
	}

	if OUTPUT_PATH == "" {
		wd, err := os.Getwd()
		if err != nil {
			logger.Error("Unable to retrieve the current working directory.")
			logger.Debug(err)
			return
		}

		fn := path.Base(strings.Replace(PDF_PATH, "\\", "/", -1))
		baseFn := strings.TrimSuffix(fn, path.Ext(PDF_PATH))
		OUTPUT_PATH = path.Join(wd, baseFn)
		logger.Warnf("No output path specified. Defaulting to the current directory (%s).", OUTPUT_PATH)
	}

	err = os.MkdirAll(OUTPUT_PATH, os.ModePerm)
	if err != nil {
		logger.Errorf("Unable to create the output directory (%s).", OUTPUT_PATH)
		logger.Debug(err)
		return
	}

	if len(PAGES) == 0 {
		pageCount, err := api.PageCountFile(PDF_PATH)
		if err != nil {
			logger.Errorf("Could not retrieve the page count from the PDF file (%s).", PDF_PATH)
			logger.Debug(err)
			return
		}
		logger.Warnf("No pages specified. Defaulting to all %d pages.", pageCount)
		for i := 1; i <= pageCount; i++ {
			PAGES = append(PAGES, i)
		}
	}

	if SCALE != 100 && (WIDTH != 0 || HEIGHT != 0) {
		logger.Warnf("Both scale and width/height are specified. Only scale will be applied.")
	}

	FILE_TYPE = strings.ToLower(FILE_TYPE)
	format, ok := formatMap[FILE_TYPE]
	if !ok {
		logger.Errorf("Unsupported file type specified (%s).", FILE_TYPE)
		logger.Debug(err)
		return
	}

	ctx, err := api.ReadContextFile(PDF_PATH)
	if err != nil {
		logger.Error("Failed to read the PDF context.")
		logger.Debug(err)
		return
	}

	wg := sync.WaitGroup{}
	cpuCores := runtime.NumCPU()
	ch := make(chan struct{}, min(cpuCores, len(PAGES)))

	for _, v := range PAGES {
		wg.Add(1)
		ch <- struct{}{}
		go func(v int) {
			extractPage(ctx, v, &format)
			wg.Done()
			<-ch
		}(v)
	}
	wg.Wait()

	logger.Info("PDF to image conversion completed successfully.")
}

func extractPage(ctx *model.Context, pageNum int, format *imgconv.Format) {
	logger.Debugf("-- Processing page %d --", pageNum)

	r, err := api.ExtractPage(ctx, pageNum)
	if err != nil {
		logger.Errorf("Unable to extract page %d from the provided PDF.", pageNum)
		logger.Debug(err)
		return
	}

	img, err := imgconv.Decode(r)
	logger.Debug(err)
	if err != nil {
		logger.Errorf("Failed to decode page %d from the provided PDF.", pageNum)
		logger.Debug(err)
		return
	}

	if SCALE != 100 {
		img = imgconv.Resize(img, &imgconv.ResizeOption{Percent: SCALE})
	} else if WIDTH != 0 && HEIGHT != 0 {
		img = imgconv.Resize(img, &imgconv.ResizeOption{Width: WIDTH, Height: HEIGHT})
	} else if WIDTH != 0 {
		img = imgconv.Resize(img, &imgconv.ResizeOption{Width: WIDTH})
	} else if HEIGHT != 0 {
		img = imgconv.Resize(img, &imgconv.ResizeOption{Height: HEIGHT})
	}

	filename := fmt.Sprintf("%s/%d.%s", OUTPUT_PATH, pageNum, FILE_TYPE)
	err = imgconv.Save(filename, img, &imgconv.FormatOption{Format: *format})
	if err != nil {
		logger.Errorf("Could not save page %d to file (%s).", pageNum, filename)
		logger.Debug(err)
		return
	}
}
