package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"bitbucket.org/llg/vcard"
	"github.com/jessevdk/go-flags"
)

var opts struct {
	VCardFolder  string `long:"vcard-folder" description:"The import folder of the vcards" required:"true"`
	ExportFolder string `long:"export-folder" description:"The export folder for the images" required:"true"`
}

func main() {
	p := flags.NewNamedParser("vcard-export-images", flags.HelpFlag)
	p.ShortDescription = "Export all images of a vcard file"
	p.AddGroup("Arguments", "", &opts)

	if _, err := p.ParseArgs(os.Args); err != nil {
		if e, ok := err.(*flags.Error); !ok || e.Type != flags.ErrHelp {
			panic(err)
		} else {
			p.WriteHelp(os.Stdout)

			return
		}
	}

	var addressBook vcard.AddressBook

	files, err := ioutil.ReadDir(opts.VCardFolder)
	if err != nil {
		log.Printf("Cannot read vcard folder %s\n", opts.VCardFolder)
		return
	}

	for _, fi := range files {
		fp := opts.VCardFolder + "/" + fi.Name()

		f, err := os.Open(fp)
		if err != nil {
			log.Printf("Cannot read vcard file %s\n", fp)
			return
		}
		defer f.Close()

		reader := vcard.NewDirectoryInfoReader(f)
		addressBook.ReadFrom(reader)
	}

	for i, c := range addressBook.Contacts {
		if c.Photo.Data != "" {
			var ext string

			c.Photo.Type = strings.ToLower(c.Photo.Type)

			switch c.Photo.Type {
			case "":
				// ignore the need for an extension
			case "image/jpeg", "image/png":
				ext = strings.Split(c.Photo.Type, "/")[1]
			default:
				ext = c.Photo.Type
			}

			if ext != "" {
				ext = "." + ext
			}

			data, err := base64.StdEncoding.DecodeString(c.Photo.Data)
			if err != nil {
				log.Printf("Could not decode data")
				continue
			}

			if err := ioutil.WriteFile(fmt.Sprintf("%s/%d%s", opts.ExportFolder, i, ext), data, 0644); err != nil {
				log.Printf("Could not write file %v", err)
				continue
			}
		}
	}

}
