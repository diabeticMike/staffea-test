package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	BasePath = "./configs/tooltips"
)

type Tooltip struct {
	ID          int    `json:"id,omitempty"`
	Label       string `json:"label,omitempty"`
	Description string `json:"description,omitempty"`
}

type TooltipPage struct {
	Page     string     `json:"page,omitempty"`
	Locale   string     `json:"locale,omitempty"`
	Tooltips []*Tooltip `json:"tooltips,omitempty"`
}

type Tooltips []TooltipPage

//extract and validate information from tooltip's filename
func loadTooltipInfoFromFileName(path string) (id int, locale string, err error) {
	// get id & tooltip's label from filename
	fileParams := strings.Split(filepath.Base(path), "-")

	//	fmt.Printf(" id %s locale %s", fileParams[0], path)

	// File is expected to be {int}-{locale}, locale equal fr_FR
	id, err = strconv.Atoi(fileParams[0])
	if err != nil {
		err = errors.New(fmt.Sprintf("%s is not an Id in %s", fileParams[0], path))
		return
	}
	// locale
	if len(fileParams[1]) == 0 {
		err = errors.New(fmt.Sprintf("%s is malformated, locale is empty, expected  {id}-{label}.html", path))
		return
	}
	fileParams = strings.Split(fileParams[1], ".")
	locale = fileParams[0]

	return
}

//extract and validate information from tooltip's file
func loadTooltipFromFile(path string, id int) (elmt Tooltip, err error) {
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		err = errors.New(fmt.Sprintf("error opening file: %s", path))
		return
	}
	// read file, extract line1 as label and following as description
	r := bufio.NewScanner(f)
	counter := 0
	elmt.ID = id

	for r.Scan() {
		line := r.Text()
		if counter == 0 {
			elmt.Label = line
		} else {
			elmt.Description = elmt.Description + line
		}
		counter++
	}
	if counter < 2 {
		err = errors.New(fmt.Sprintf("%s is malformated, first line must be a label and second at least a description", path))
	}

	return

}

func loadToolTips(initialPath string) Tooltips {
	var lTooltips Tooltips

	// Read path recursively
	err := filepath.Walk(initialPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}

		//to do remove this to replace by initialPath substitution
		tp := strings.Split(filepath.Dir(path), "/")
		page := strings.Join(tp[1:], "/")

		if !info.IsDir() && filepath.Ext(path) == ".html" {

			//Extract datas from filename
			id, locale, err := loadTooltipInfoFromFileName(path)
			if err != nil {
				log.Fatal(err)
			}

			//extract information from file
			elmt, err := loadTooltipFromFile(path, id)
			if err != nil {
				log.Fatal(err)
			}

			//populating main tooltips structure
			isExist := false
			//testing if tooltips page exist
			for idx, b := range lTooltips {
				if b.Locale == locale && b.Page == page {
					isExist = true
					// appending a new tooltips page
					lTooltips[idx].Tooltips = append(lTooltips[idx].Tooltips, &elmt)
					continue
				}
			}
			if isExist == false {
				lTooltips = append(lTooltips, TooltipPage{page, locale, []*Tooltip{&elmt}})

			}

		}

		return nil
	})
	if err != nil {
		fmt.Printf("walk error [%v]\n", err)
	}

	return lTooltips

}

func getTooltips(c *gin.Context, myToolTips Tooltips) {
	// get params
	page := "tooltips/tooltips/" + strings.ToLower(c.Query("page"))
	locale := "fr_fr"
	if len(c.Request.Header["X-Locale"]) == 1 {
		locale = strings.ToLower(c.Request.Header["X-Locale"][0])
	}

	// looking for tooltips from page/locale
	for _, vx := range myToolTips {
		if strings.ToLower(vx.Page) == page && strings.ToLower(vx.Locale) == locale {
			c.JSON(200, vx.Tooltips)
			return
		}
	}
	c.JSON(404, gin.H{"message": "no tooltips found"})

}

func main() {
	// global variable to keep tooltips
	myToolTips := loadToolTips(BasePath)

	if isDebugMode() {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()

	router.GET("/front/tooltips", func(c *gin.Context) {
		getTooltips(c, myToolTips)
	})

	err := router.Run()
	if err != nil {
		return
	}
}

func isDebugMode() bool {
	return os.Getenv("GIN_MODE") == gin.DebugMode || os.Getenv("GIN_MODE") == ""
}
