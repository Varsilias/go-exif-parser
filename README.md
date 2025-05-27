## CS-MONITOR

An Exchangeable Image File Format(EXIF) Parser Written in Golang

## Usage

- Clone the repository `git clone git@github.com:Varsilias/go-exif-parser.git`
- Navigate into the project directory `cd go-exif-parser`
- Built the project `go build -o exif-parser  .`
- Run the built binary with the command `./exif-parser --image=/path/to/image --output=output.json`
- The `--output` flag is optional. When not specifiied, the result of parsing the image will be stored in a file `output.json`
- You can also use `-i` in place of `--image` and `-o` in place of `--output`
- Run the built binary like so `./exif-parser --image=./images/Canon_iR1018.jpg --output=output.json`
- The command above uses one of the images included in the project directory, you can use your custom JPEG or JPG image
- **Windows Powersehll:** `.\exif-parser.exe --image="$(Get-Location)\images\Canon_iR1018.jpg" --output=test.json`
