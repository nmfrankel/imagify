# Imagify

Imagify is a command-line tool designed to convert PDF pages into images. It allows users to extract specific pages or all pages from a PDF file and save them as images in various formats such as PNG, JPG, PDF, or WEBP. The tool also provides options to customize the output by specifying scaling factors, dimensions, and output directories.

## Features

- Extract specific pages or all pages from a PDF file.
- Save output images in multiple formats: **PNG**, **JPG**, **PDF**, or **WEBP**.
- Resize output images by specifying a scaling factor or custom dimensions.
- Specify the output directory for saving the generated images.

## Flags

- `--pdf_path` (Required): Path to the input PDF file.
- `--output_path`: Directory where output files will be saved. Defaults to the current directory.
- `--scale`: Scaling factor for the output image as a percentage. Defaults to 100%.
- `--width`: Width of the output image in pixels. Ignored if scale is provided.
- `--height`: Height of the output image in pixels. Ignored if scale is provided.
- `--file_type`: Output image format. Supported formats: `png`, `jpg`, `pdf`, `webp`. Defaults to `png`.
- `--pages`: List of page numbers to process (e.g., `--pages=1,2,3` or `--pages=[1,2,3]`).

## Usage

Imagify is ideal for converting PDF documents into image formats for use in presentations, web applications, or other scenarios where image formats are required.

### Example

```bash
imagify --pdf_path=Original-Microsoft-Source-Code.pdf --output_path=./images --scale=150 --file_type=jpg --pages=1,2,3
```

This command converts pages 1, 2, and 3 of `input.pdf` into JPG images, scales them to 150%, and saves them in the `./images` directory.
