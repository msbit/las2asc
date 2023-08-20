# las2asc

A simple utility to convert the output of [`las2txt`][1] to ESRI ASCII files.

## extents

Runs over any file with the `x`, `y` and `z` columns space separated and returns
the extent of each, as two triplets.

Expects only the input filename as an argument.

## asc

Converts any file with the `x`, `y` and `z` columns space separated into an ESRI
ASCII file, with a cell size of 1.0m. As part of the sampling, the minimum value
in each cell is used.

Expects the following arguments:

  * `-ll` the lower left coordinate of the input file, as a comma-separated pair
  * `-tr` the top right coordinate of the input file, as a comma-separated pair
  * `-in` the input filename
  * `-out` the output filename

[1]: https://lastools.github.io
