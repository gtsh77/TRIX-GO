package q3parser

import (
	"bufio"
	"fmt"
	"os"
	"log"
	"strings"
)

func processEntity(line string, entity *[]CENT, entityNum int32) {
	cnt, _ := fmt.Sscanf(line,"%s %[^\t\n]s",(*entity)[entityNum].Values[(*entity)[entityNum].ValueCnt].Name,(*entity)[entityNum].Values[(*entity)[entityNum].ValueCnt].Value)
	if (cnt == 2) {
		if ((*entity)[entityNum].Values[(*entity)[entityNum].ValueCnt].Name == "\"classname\"") {
			(*entity)[entityNum].ClassName = (*entity)[entityNum].Values[(*entity)[entityNum].ValueCnt].Value
		}
		(*entity)[entityNum].ValueCnt++;
	} else {
		fmt.Printf("bad_ent_format: %s\n",line);
	}
}

func processBrush(line string) {

}

func ParseMap(path string) string {

	//vars
	var isBrush, isEntity, entPart, ignoreBrush, i, j, k uint32 = 0,0,0,0,0,0,0;
	var brushNum, entityNum int32 = -1, -1;
	var fl [3]float32;
	var num [3]int32;

	//alloc space for header
	header := CHEAD{BrushCount: 0, TexelCount: 0, EntityCount: 0};

	//get count of all brushes

	//TEST OPEN --START--
	file1, err := os.Open(path)
	if err != nil {
		log.Fatal(err);
	}
	defer file1.Close();
	//TEST OPEN --END--

	scanner := bufio.NewScanner(file1);
	for scanner.Scan() {
		if (strings.Contains(scanner.Text(),"// brush")){
			header.BrushCount++;

		} else if (strings.Contains(scanner.Text(),"  patchDef2")){
			header.BrushCount--;

		} else if (strings.Contains(scanner.Text(),"// entity")){
			header.EntityCount++;
		}
	}
	
	if err := scanner.Err(); err != nil {
		log.Fatal(err);
	}

    fmt.Printf("\ntotal_brushes: %d\n",header.BrushCount);
    fmt.Printf("total_entities: %d\n\n",header.EntityCount);
	fmt.Printf("*** VERTICES HANDLER ***\n\n");

	//alloc tmp space for textures duplicates
	texelDup := make([]CTEX, 10000);

	//alloc space for brushes
	brush := make([]CBRUSH, header.BrushCount);

	//alloc tmp space for brush
	tmpBrush := make([]CBRUSH,1);

	//alloc space for entities
	entityCount := header.EntityCount;
	if (entityCount == 0) {
		entityCount = 1;
	}
	entity := make([]CENT,entityCount);

	//set file-pointer back
	file1.Seek(0, 0);

	//parse brushes
	scanner = bufio.NewScanner(file1);
	for scanner.Scan() {
		//if new brush use new struct
		if (strings.Contains(scanner.Text(),"// brush")){
			brushNum++;
			brush[brushNum].PlaneCount = 0;
			ignoreBrush = 0;

			//init struct, clear garbage
			brush[brushNum].ID = uint32(brushNum);
			brush[brushNum].FaceCount = 0;
			fmt.Printf("brush %d/%d\n",brushNum+1,header.BrushCount);

		} else if (strings.Contains(scanner.Text(),"  patchDef2")){
			//upd cntrs    
			brushNum--;
			ignoreBrush = 1;

		} else if (strings.Contains(scanner.Text(),"// entity")){
			if (entPart == 0) {
				entPart = 1;
			}
			//upd cntrs    
			entityNum++;
			//init struct, clear garbage
			entity[entityNum].ID = uint32(entityNum);
			entity[entityNum].ValueCnt = 0;
			fmt.Printf("entity %d/%d\n",entityNum+1,header.EntityCount);
		} else if (ignoreBrush == 0 && string(scanner.Text()[1]) == "("){
			isBrush = 1;
		} else if (entPart == 1 && string(scanner.Text()[1]) == "\""){
			isEntity = 1;

		} else if (string(scanner.Text()[1]) == "}"){
			isBrush = 0;
			isEntity = 0;
		}

		if (isBrush == 1){
			processBrush(scanner.Text())
		} else if (isEntity == 1){
			processEntity(scanner.Text(),&entity,entityNum)
		}
	}


	return "q3parser called.\n"
}