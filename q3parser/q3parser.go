package q3parser

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"bytes"
	"strconv"
	"strings"
	_ "sort"
	"path/filepath"
	"encoding/gob"
	m "github.com/gtsh77/TRIX-GO/mlib"
)

func ParseMap(path string) {
	//vars
	var (
		isBrush, isEntity, entPart, ignoreBrush, i, j, cnt, newTexelSize uint32
		brushNum, entityNum int32 = -1, -1
		fl [3]float32
		num [3]int32
		header CHEAD
		texelDup [MAXTEXDUP]CTEX
		texelFinal []CTEX
		brush []CBRUSH
		tmpBrush [1]CBRUSH
		entity []CENT
		entityCount uint32
		planes [9*MAXFACES]float64
		vertices [12*MAXFACES]float64
		rpath, base, fname, wpath string
	)

	//TEST R-OPEN --START--
	file1, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file1.Close()
	//TEST R-OPEN --END--


	//TEST W-OPEN --START--
	//get runtime path
	rpath, _ = filepath.Abs(filepath.Dir(os.Args[0]))

	//get fname
	base = filepath.Base(path)
	fname = strings.Split(base,".")[0]

	//set new path
	wpath = filepath.Join(rpath,CMAPDIR+fname+CMAPEXT)

	//try create outer-file
	file2, err := os.Create(wpath)	
	if err != nil {
		log.Fatal(err)
	}
	defer file2.Close()
	//TEST W-OPEN --END--

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

	//alloc space for brushes
	brush = make([]CBRUSH, header.BrushCount)

	//alloc space for entities
	entityCount = header.EntityCount
	if entityCount == 0 {
		entityCount = 1
	}
	entity = make([]CENT, entityCount)

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
			processBrush(string(scanner.Text()), &brush, tmpBrush, brushNum, fl, num, texelDup, header)
		} else if isEntity == 1 {
			processEntity(string(scanner.Text()), &entity, entityNum)
		}
	}

	fmt.Printf("\n*** TEXTURES HANDLER ***\n\n")

	//form unique texture list
	texelFinal = make([]CTEX, header.TexelCount)

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

	//calculate vertices
	for i := uint32(0); i < header.BrushCount; i++ {
		//prepare planes
		for j := uint8(0); j < brush[i].PlaneCount; j++ {
			for k := uint8(0); k < 9; k++ {
				planes[9*j + k] = float64(brush[i].Planes[9*j +k])
			}
		}
		//doMapCalc(planes, brush[i].plane_count, brush[i].faces, brush[i].face_count, vertices);
		getShape(planes, brush[i].PlaneCount, brush[i].Faces, brush[i].FaceCount, vertices)

		//store
		for j:= uint8(0); j < brush[i].FaceCount*12; j++ {
			brush[i].Vertices[j] = vertices[j]
		}

		//write bin
		file2.Write(SafeBytes(header))
		file2.Write(SafeBytes(texelFinal))
		file2.Write(SafeBytes(brush))
		file2.Write(SafeBytes(entity))
		// w := bufio.NewWriterSize(file2,)
		// n4, _ := w.Write(header)
	}
}

func processBrush(line string, brush *[]CBRUSH, tmpBrush [1]CBRUSH, brushNum int32, fl [3]float32, num [3]int32, TexelDup [MAXTEXDUP]CTEX, header CHEAD) {
	//parse main line into tmp struct
	cnt, err := fmt.Sscanf(line, "( %d %d %d ) ( %d %d %d ) ( %d %d %d ) %s %d %d %f %f %f", &tmpBrush[0].Planes[0], &tmpBrush[0].Planes[1], &tmpBrush[0].Planes[2], &tmpBrush[0].Planes[3], &tmpBrush[0].Planes[4], &tmpBrush[0].Planes[5], &tmpBrush[0].Planes[6], &tmpBrush[0].Planes[7], &tmpBrush[0].Planes[8], &tmpBrush[0].Texel[0], &num[0], &num[1], &fl[2], &fl[0], &fl[1])
	if err != nil {
		log.Fatal(err)
	}
	if cnt == 15 {
		//store planes
		for i := 0; i < 9; i++ {
			(*brush)[brushNum].Planes[((*brush)[brushNum].PlaneCount*9)+uint8(i)] = tmpBrush[0].Planes[i]
		}
		//check if it has valid face
		if strings.Contains(tmpBrush[0].Texel[0], "common/caulk") == false && strings.Contains(tmpBrush[0].Texel[0], "common/nodraw") == false {
			//store face id
			(*brush)[brushNum].Faces[(*brush)[brushNum].FaceCount] = int((*brush)[brushNum].PlaneCount)
			//store texel name
			(*brush)[brushNum].Texel[(*brush)[brushNum].FaceCount] = tmpBrush[0].Texel[0]
			//store texel shift
			(*brush)[brushNum].ShiftX[(*brush)[brushNum].FaceCount] = num[0]
			(*brush)[brushNum].ShiftY[(*brush)[brushNum].FaceCount] = num[1]
			//store texel scale
			(*brush)[brushNum].ScaleX[(*brush)[brushNum].FaceCount] = fl[0]
			(*brush)[brushNum].ScaleY[(*brush)[brushNum].FaceCount] = fl[1]
			//update global texels array
			TexelDup[header.TexelCount].Path = tmpBrush[0].Texel[0]
			//upd tx cnt
			header.TexelCount++
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

// ===== MAP =====
//get normals and distances
func GetND(planes [9*MAXFACES]float64, planeNum uint8, normals []float64, distances *float64) {
	var a, b, c, ab, cb, normal, normal_n []float64 = []float64{0, 0, 0}, []float64{0, 0, 0}, []float64{0, 0, 0}, []float64{0, 0, 0}, []float64{0, 0, 0}, []float64{0, 0, 0}, []float64{0, 0, 0}
	var distance float64

	a[0] = planes[planeNum*9+0]
	a[1] = planes[planeNum*9+1]
	a[2] = planes[planeNum*9+2]

	b[0] = planes[planeNum*9+3]
	b[1] = planes[planeNum*9+4]
	b[2] = planes[planeNum*9+5]

	c[0] = planes[planeNum*9+6]
	c[1] = planes[planeNum*9+7]
	c[2] = planes[planeNum*9+8]

	m.SubV(a, b, 3, ab)
	m.SubV(c, b, 3, cb)
	m.CrossV3(ab, cb, normal)
	m.Normalize(normal, 3, normal_n)

	distance = m.Scalar(normal_n, a, 3)
	*distances = distance
	normals[0] = normal_n[0]
	normals[1] = normal_n[0]
	normals[2] = normal_n[0]
}

//get intersection
func GetIntersection(planes [9*MAXFACES]float64, i uint8, j uint8, k uint8, intersection [3]float64) {
	var sum float64
	var distances [3]float64
	var n, t map[uint8][]float64

	for i := range [9]uint8{} {
		if i == 0 {
			n, t = make(map[uint8][]float64), make(map[uint8][]float64)
		}
		if i < 3 {
			n[uint8(i)] = []float64{0, 0, 0}
		}
		t[uint8(i)] = []float64{0, 0, 0}
	}

	//get normals and distances
	GetND(planes, i, n[0], &distances[0])
	GetND(planes, j, n[1], &distances[1])
	GetND(planes, k, n[2], &distances[2])

	//calc intersection
	m.CrossV3(n[1], n[2], t[0])
	m.MulV(t[0], -distances[0], 3, t[1])

	//-d2*Cross(n3, n1)
	m.CrossV3(n[2], n[0], t[2])
	m.MulV(t[2], -distances[1], 3, t[3])

	//-d3*Cross(n1, n2)
	m.CrossV3(n[0], n[1], t[4])
	m.MulV(t[4], -distances[2], 3, t[5])

	//summ p1 p2 p3
	m.AddV(t[1], t[3], 3, t[6])
	m.AddV(t[6], t[5], 3, t[7])

	// //denominator
	sum = m.Scalar(n[0], t[0], 3)

	//set
	if sum == 0 {
		intersection[0] = -1
	} else {
		intersection[0] = t[7][0]
		intersection[1] = t[7][1]
		intersection[2] = t[7][2]
	}
}

func getShape(planes [9*MAXFACES]float64, pc uint8, faces [MAXFACES]int, fc uint8, vertices [12*MAXFACES]float64) {
	var stored uint8
	var distances [6]float64
	var intersection [3]float64
	var n map[uint8][]float64

	for i := range [6]uint8{} {
		if i == 0 {
			n = make(map[uint8][]float64)
		}
		n[uint8(i)] = []float64{0, 0, 0}
	}

	//get all normals and distances
	for i := uint8(0); i < pc; i++ {
		GetND(planes, i, n[i], &distances[i])
	}

	//get vertices
	for i := uint8(0); i <  pc; i++ {
		//proccess only vis faces
		if intIn(int(i), faces) {
			for j := uint8(0); j < pc; j++ {
				for k := uint8(0); k < pc; k++ {
					if i != j && i != k && j != k {
						//get intersetion
						GetIntersection(planes, i, j, k, intersection)
						//tmp chk if legal by denominator and x < 0
						if (intersection[0] != -1 && intersection[0] >= 0) {
							vertices[stored] = intersection[0]
							vertices[stored+1] = intersection[1]
							vertices[stored+2] = intersection[2]
							stored += 3
						}						
					}
				}
			}
		}

	}
}

func intIn(i int, fa [MAXFACES]int) bool {
	for _, v := range fa {
		if v == i {
			return true
		} else {
			return false
		}
	}
	return false
}

func SafeBytes(s interface{}) []byte {
	var b *bytes.Buffer = bytes.NewBuffer(nil)
	var g *gob.Encoder = gob.NewEncoder(b)

	g.Encode(s)

	return b.Bytes()
}
