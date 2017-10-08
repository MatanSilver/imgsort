package main

import (
	"crypto/sha256"
	"encoding/hex"
	//"fmt"
	//"github.com/urfave/cli"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	//"syscall"
	"gopkg.in/h2non/filetype.v1"
	//"github.com/djherbis/atime"
	"github.com/rwcarlsen/goexif/exif"
)

// Copy the src file to dst. Any existing file will be overwritten and will not
// copy file attributes.
func Copy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

type FileInfoWrapper struct {
	Info    os.FileInfo
	Path    string
	Hash    string
	Created time.Time
}

// Takes in a directory path. Recursively crawls the directory and outputs a
// list of paths of files in that directory and subdirectories
func ls_imgs(dir string) []FileInfoWrapper {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		//log.Fatalf("dir: %s, err: %s", dir, err)
		log.Println(err)
	}
	var fileinfos []FileInfoWrapper
	for _, file := range files {
		if file.IsDir() {
			//if the file is a directory, recursively add the directory's contents
			fileinfos = append(fileinfos, ls_imgs(file.Name())...)
		} else {
			fullpath := strings.Join([]string{dir, file.Name()}, "/")
			buf, err := ioutil.ReadFile(fullpath)
			if err != nil { //something went terribly wrong with reading the file
				log.Fatal(err)
			}
			if filetype.IsImage(buf) { //only list image files, reject others
				f, err := os.Open(fullpath)
				if err != nil { //something went terribly wrong with reading the file
					log.Fatal(err)
				}
				loc, err := time.LoadLocation("") //use utc time by default
				if err != nil {
					log.Fatal(err)
				}
				//preload tm with a dummy date, before when we'd have many digital photos
				tm := time.Date(2000, time.January, 1, 1, 1, 1, 1, loc)
				x, err := exif.Decode(f)
				if err != nil { //if exif loads improperly (i.e. header missing)
												//we will keep the default date
					log.Printf("Error in file %s: %s\n", fullpath, err)
				} else { //if exif loads properly, get the date
					tm, err = x.DateTime()
					if err != nil {
						log.Printf("Error in file %s: %s\n", fullpath, err)
					}
				}

				//now we generate a hash, which might be useful for checking for
				//duplicates
				buff, err := ioutil.ReadFile(fullpath)
				if err != nil {
					log.Fatal(err)
				}
				hasher := sha256.New()
				hasher.Write(buff)
				if err != nil {
					log.Fatal(err)
				}
				fileinfo := FileInfoWrapper{file, fullpath, hex.EncodeToString(hasher.Sum(nil)), tm}
				fileinfos = append(fileinfos, fileinfo)
			}
		}
	}
	return fileinfos
}

func main() {
	//fmt.Println("test")
	/*_ = sha256.New()
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "lang",
			Value: "english",
			Usage: "language for the greeting",
		},
	}

	app.Action = func(c *cli.Context) error {
		return nil
	}

	app.Run(os.Args)
	*/
	fileinfos := ls_imgs(".")
	//fmt.Printf("%v\n", fileinfos)
	//return
	for _, fileinfo := range fileinfos {
		//we want to read image data
		//fmt.Println(fileinfo.Info.ModTime())
		//year, month, day := fileinfo.Info.ModTime().Date()
		year, month, day := fileinfo.Created.Date()
		newpath := strings.Join([]string{strconv.Itoa(year), month.String(), strconv.Itoa(day), fileinfo.Info.Name()}, "/")
		mode := os.FileMode(0777)
		os.MkdirAll(strings.Join([]string{strconv.Itoa(year), month.String(), strconv.Itoa(day)}, "/"), mode)

		log.Printf("Copying %s to %s", fileinfo.Path, newpath)
		//we don't care about waiting for the copy to finish, so we dispatch
		//to a goroutine
		go Copy(fileinfo.Path, newpath)
	}
}
