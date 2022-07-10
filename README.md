# photosort
A small program written in Golang to copy and sort photos by year and month for a source to a dest folder.

This is forked on code by ?.  I have moved to maore standard places.

The format of the top directory is yyyy
The format of the month directory is yyyy-mm
The format of the day directory is yyyy-mm-dd

This is my first attempt at the go language. Feel free to improve it through pull requests!

***Important: This has only been tested on Windows 10, let me know how it behaves in other OS.***

### Functionality
- Recursively visits all files and folders in the srcFolder directoy tree
- Analyses if files have Exif date information
- If so it copies the media file to the corresponding folder in the destfolder, creating the necessary folders such as Year and Month
- If the file is a media file but has no information regarding the capture date, it copies the file into a folder called Others in the destination folder
- It ignores any duplicated file or non media files

### Media files supported
```.tiff .tif .gif .jpeg .jpg .png .raw .webm .mkv .avi .mov .wmv .mp4 .MP4 .m4v .mpg .mp2 .mpeg .dng```

The case of the extension is ignored.

### Build
```
go get github.com/rwcarlsen/goexif/exif
go build photosort.go errors.go
```

### Run
```
./photosort --source-folder [srcFolder] --destination-folder [destFolder]
```
