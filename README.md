# Imagify

Imagify is a cli tool to split & convert PDF pages into images. The tool also provides options to customize the output by specifying scaling factors, dimensions, and output directories.

## Features

- Extract specific pages or all pages from a PDF file.
- Save output images in multiple formats: **PNG**, **JPG**, **PDF**, or **WEBP**.
- Resize output images by specifying a scaling factor or custom dimensions.
- Specify the output directory for saving the generated images.

## Example

```bash
.\imagify --pdf "Original-Microsoft-Source-Code.pdf" --out "images" --scale=50 --pages=1,2,3
```

This command converts pages 1, 2, and 3 of `Original-Microsoft-Source-Code.pdf` into JPG images, scales them to 150%, and saves them in the `./images` directory.

## Flags

| Option        | Description                                                                                     |
|-------------- |-------------------------------------------------------------------------------------------------|
| `--pdf`       | Path to the input PDF file.                                                                     |
| `--out`       | Directory where output files will be saved. Defaults to the current directory.                  |
| `--scale`     | Scaling factor for the output image as a percentage. Defaults to `100%`.                        |
| `--width`     | Width of the output image in pixels. Ignored if `--scale` is provided.                          |
| `--height`    | Height of the output image in pixels. Ignored if `--scale` is provided.                         |
| `--encoding`  | Output image format. Supported formats: `png`, `jpg`, `pdf`, `webp`. Defaults to `png`.         |
| `--pages`     | List of page numbers to process (e.g., `--pages=1,2,3` or `--pages=[1,2,3]`).                   |
| `--debug`     | Print debug logs. Defaults to `false`.                                                          |

> [Guide to setup CGO_ENABLED=1](https://github.com/go101/go101/wiki/CGO-Environment-Setup)
