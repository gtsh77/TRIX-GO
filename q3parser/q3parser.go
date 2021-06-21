package q3parser

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func ParseMap(path string) {

	//vars
	var	isBrush, isEntity, entPart, ignoreBrush, i, j, cnt, newTexelSize uint32 = 0, 0, 0, 0, 0, 0, 0, 0
	var	brushNum, entityNum int32 = -1, -1
	var	fl [3]float32
	var	num [3]int32

	//alloc space for header
	header := CHEAD{BrushCount: 0, TexelCount: 0, EntityCount: 0}

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
		if strings.Contains(scanner.Text(), "// brush") {
			header.BrushCount++

		} else if strings.Contains(scanner.Text(), "  patchDef2") {
			header.BrushCount--

		} else if strings.Contains(scanner.Text(), "// entity") {
			header.EntityCount++
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\ntotal_brushes: %d\n", header.BrushCount)
	fmt.Printf("total_entities: %d\n\n", header.EntityCount)
	fmt.Printf("*** VERTICES HANDLER ***\n\n")

	//alloc tmp space for textures duplicates
	texelDup := make([]CTEX, 10000)

	//alloc space for brushes
	brush := make([]CBRUSH, header.BrushCount)

	//alloc tmp space for brush
	tmpBrush := make([]CBRUSH, 1)

	//alloc space for entities
	entityCount := header.EntityCount
	if entityCount == 0 {
		entityCount = 1
	}
	entity := make([]CENT, entityCount)

	//set file-pointer back
	file1.Seek(0, 0)

	//parse brushes
	scanner = bufio.NewScanner(file1)
	for scanner.Scan() {
		//if new brush use new struct
		if strings.Contains(scanner.Text(), "// brush") {
			brushNum++
			brush[brushNum].PlaneCount = 0
			ignoreBrush = 0

			//init struct, clear garbage
			brush[brushNum].ID = uint32(brushNum)
			brush[brushNum].FaceCount = 0
			fmt.Printf("brush %d/%d\n", brushNum+1, header.BrushCount)

		} else if strings.Contains(scanner.Text(), "  patchDef2") {
			//upd cntrs
			brushNum--
			ignoreBrush = 1

		} else if strings.Contains(scanner.Text(), "// entity") {
			if entPart == 0 {
				entPart = 1
			}
			//upd cntrs
			entityNum++
			//init struct, clear garbage
			entity[entityNum].ID = uint32(entityNum)
			entity[entityNum].ValueCnt = 0
			fmt.Printf("entity %d/%d\n", entityNum+1, header.EntityCount)
		} else if ignoreBrush == 0 && string(scanner.Text()[0:1]) == "(" {
			isBrush = 1
		} else if entPart == 1 && string(scanner.Text()[0:1]) == "\"" {
			isEntity = 1

		} else if string(scanner.Text()[0:1]) == "}" {
			isBrush = 0
			isEntity = 0
		}

		if isBrush == 1 {
			processBrush(string(scanner.Text()), &brush, &tmpBrush, brushNum, &fl, &num, &texelDup, &header)
		} else if isEntity == 1 {
			processEntity(string(scanner.Text()), &entity, entityNum)
		}
	}

	fmt.Printf("\n*** TEXTURES HANDLER ***\n\n")

	//form unique texture list
	texelFinal := make([]CTEX, header.TexelCount)

	for i, cnt = 0, 0; i < header.TexelCount; i, cnt = i+1, 0 {
		for j = 0; j < newTexelSize; j, cnt = j+1, cnt+1 {
			if texelDup[i].Path == texelFinal[j].Path {
				break
			}
		}

		if cnt == newTexelSize {
			texelFinal[newTexelSize].Path = texelDup[i].Path
			newTexelSize++
		}
	}

	fmt.Printf("total_vis_faces: %d\n", header.TexelCount)
	fmt.Printf("total_unique_textures: %d\n", newTexelSize)

	//now set proper unique texel count
	header.TexelCount = newTexelSize
}

func processBrush(line string, brush *[]CBRUSH, tmpBrush *[]CBRUSH, brushNum int32, fl *[3]float32, num *[3]int32, TexelDup *[]CTEX, header *CHEAD) {
	//parse main line into tmp struct
	cnt, err := fmt.Sscanf(line, "( %d %d %d ) ( %d %d %d ) ( %d %d %d ) %s %d %d %f %f %f", &(*tmpBrush)[0].Planes[0], &(*tmpBrush)[0].Planes[1], &(*tmpBrush)[0].Planes[2], &(*tmpBrush)[0].Planes[3], &(*tmpBrush)[0].Planes[4], &(*tmpBrush)[0].Planes[5], &(*tmpBrush)[0].Planes[6], &(*tmpBrush)[0].Planes[7], &(*tmpBrush)[0].Planes[8], &(*tmpBrush)[0].Texel[0], &(*num)[0], &(*num)[1], &(*fl)[2], &(*fl)[0], &(*fl)[1])
	if err != nil {
		log.Fatal(err)
	}
	if cnt == 15 {
		//store planes
		for i := 0; i < 9; i++ {
			(*brush)[brushNum].Planes[((*brush)[brushNum].PlaneCount*9)+uint8(i)] = (*tmpBrush)[0].Planes[i]
		}
		//check if it has valid face
		if strings.Contains((*tmpBrush)[0].Texel[0], "common/caulk") && strings.Contains((*tmpBrush)[0].Texel[0], "common/nodraw") {
			//store face id
			(*brush)[brushNum].Faces[(*brush)[brushNum].FaceCount] = (*brush)[brushNum].PlaneCount
			//store texel name
			(*brush)[brushNum].Texel[(*brush)[brushNum].FaceCount] = (*tmpBrush)[0].Texel[0]
			//store texel shift
			(*brush)[brushNum].ShiftX[(*brush)[brushNum].FaceCount] = num[0]
			(*brush)[brushNum].ShiftY[(*brush)[brushNum].FaceCount] = num[1]
			//store texel scale
			(*brush)[brushNum].ScaleX[(*brush)[brushNum].FaceCount] = fl[0]
			(*brush)[brushNum].ScaleY[(*brush)[brushNum].FaceCount] = fl[1]
			//update global texels array
			(*TexelDup)[(*header).TexelCount].Path = (*tmpBrush)[0].Texel[0]
			//upd tx cnt
			(*header).TexelCount++
			//upd struct fc cnt
			(*brush)[brushNum].FaceCount++
		}
		//upd plane num
		(*brush)[brushNum].PlaneCount++

	} else {
		fmt.Printf("bad_brush_format: %s\n", line)
	}
}

func processEntity(line string, entity *[]CENT, entityNum int32) {
	cnt, err := fmt.Sscanf(line, "%q %q", &(*entity)[entityNum].Values[(*entity)[entityNum].ValueCnt].Name, &(*entity)[entityNum].Values[(*entity)[entityNum].ValueCnt].Value)
	if err != nil {
		log.Fatal(err)
	}
	if cnt == 2 {
		// "%s %[^\t\n]s" workaround --start--
		strconv.Quote((*entity)[entityNum].Values[(*entity)[entityNum].ValueCnt].Name)
		strconv.Quote((*entity)[entityNum].Values[(*entity)[entityNum].ValueCnt].Value)
		// "%s %[^\t\n]s" workaround --end--
		if (*entity)[entityNum].Values[(*entity)[entityNum].ValueCnt].Name == "\"classname\"" {
			(*entity)[entityNum].ClassName = (*entity)[entityNum].Values[(*entity)[entityNum].ValueCnt].Value
		}
		(*entity)[entityNum].ValueCnt++
	} else {
		fmt.Printf("bad_ent_format: %s\n", line)
	}
}