# ascii

![GitHub](https://img.shields.io/github/license/jibstack64/ascii) ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/jibstack64/ascii) ![GitHub release (latest by date)](https://img.shields.io/github/v/release/jibstack64/ascii) ![GitHub all releases](https://img.shields.io/github/downloads/jibstack64/ascii/total)

A standalone executable that converts a standard image into viable ASCII art.

Syntax: `./ascii <--in> photo.png [--out] out.txt [--scale] 0.5 [--stretch] 0.5 [--print]`
> Often, you may require a very low scale factor (>0.1), particularly for larger images.

### Arguments

#### Required:
- `--in <image.png>` - Specifies the input image. Can be `png` or  `jpg/jpeg`.
  
#### Optional:
- `--out <out.txt>` - Specifies the output file. Defaults to `out.txt`.
- `--scale <number>` - Scales the result, making it more viable to be printed in a console. Must be above `0` (e.g. `0.5` would halve the size of the result). Defaults to `0.5`.
- `--stretch <integer>` - Stretches the result horizontally. This is useful for larger images. Must be above `0`. Defaults to `1`.
- `--print` - Prints the result to the console once finished.
- `--pretty` - When `--print` is passed, output is printed layer-by-layer.
- `--colour` - Colours the output. Only works on terminals with RGB ANSI colour support.

> ### Example
> ![screenshot-2023-01-08-18:31:48](https://user-images.githubusercontent.com/107510599/211212842-5e66841c-bbc3-4c8c-bb9f-08961293682f.png)

