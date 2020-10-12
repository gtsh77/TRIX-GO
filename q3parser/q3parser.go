package q3parser

import (
	"bufio"
	"fmt"
	"os"
	"log"
	"strings"
)

func ParseMap(path string) string {

	//alloc space for header
	header := CHEAD{BrushCount: 0, TexelCount: 0, EntityCount: 0}
	//brushNum := -1
	//entityNum := -1

	//get count of all brushes

	//TEST OPEN --START--
	file1, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file1.Close()
	//TEST OPEN --END--

	scanner := bufio.NewScanner(file1)
	for scanner.Scan() {
		if (strings.Contains(scanner.Text(),"// brush")){
			header.BrushCount++

		} else if (strings.Contains(scanner.Text(),"  patchDef2")){
			header.BrushCount--

		} else if (strings.Contains(scanner.Text(),"// entity")){
			header.EntityCount++
		}
	}
	
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

    fmt.Printf("\ntotal_brushes: %d\n",header.BrushCount);
    fmt.Printf("total_entities: %d\n\n",header.EntityCount);
    //fmt.Printf("*** VERTICES HANDLER ***\n\n");

	return "q3parser called.\n"
}