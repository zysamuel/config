package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
)

var oneTabs string = "\t"
var twoTabs string = "\t\t"
var threeTabs string = "\t\t\t"
var fourTabs string = "\t\t\t\t"

func writeStaticPart(inputFile string, dstFile *os.File) {
	hdr, err := os.Open(inputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer hdr.Close()

	scanner := bufio.NewScanner(hdr)
	for scanner.Scan() {
		line := scanner.Text()
		dstFile.WriteString(line + "\n")
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func writeResourceHdr(strName string, operation string, dstFile *os.File) {
	//dstFile.WriteString(twoTabs + "\"/" + strName + "\": { \n")
	dstFile.WriteString(twoTabs + "\"" + operation + "\": { " + "\n")
	dstFile.WriteString(threeTabs + "\"tags\": [ " + "\n")
	dstFile.WriteString(fourTabs + "\"" + strName + "\"" + "\n")
	dstFile.WriteString(threeTabs + "]," + "\n")
	dstFile.WriteString(twoTabs + "\"summary\": \"Create New" + strName + "\"," + "\n")
	dstFile.WriteString(twoTabs + "\"description\":" + "\"" + strName + "\"," + "\n")
	//dstFile.WriteString(twoTabs + "\"operationId\": \"add" + strName + ",\n")
	dstFile.WriteString(twoTabs + "\"consumes\": [" + "\n")
	dstFile.WriteString(threeTabs + "\"application/json\"," + "\n")
	dstFile.WriteString(twoTabs + "]," + "\n")
	dstFile.WriteString(twoTabs + "\"produces\": [ " + "\n")
	dstFile.WriteString(threeTabs + "\"application/json\"," + "\n")
	dstFile.WriteString(twoTabs + "]," + "\n")
	dstFile.WriteString(twoTabs + "\"parameters\": [ " + "\n")
}

func writeEpilogueForStruct(strName string, dstFile *os.File) {
	dstFile.WriteString(twoTabs + "\"responses\": { " + "\n")
	dstFile.WriteString(threeTabs + "\"405\": {" + "\n")
	dstFile.WriteString(fourTabs + "\"description\": \"Invalid input\"" + "\n")
	dstFile.WriteString(threeTabs + " }" + "\n")
	dstFile.WriteString(twoTabs + " }" + "\n")
	dstFile.WriteString(twoTabs + " }," + "\n")
}

func writeAttributeJson(attrName string, attrType string, dstFile *os.File) {
	var attrTypeVal string
	dstFile.WriteString(fourTabs + "{" + "\n")
	dstFile.WriteString(fourTabs + "\"in\": \"formData\"," + "\n")
	dstFile.WriteString(fourTabs + "\"name\":" + "\"" + attrName + "\"" + "," + "\n")
	switch attrType {
	case "string":
		attrTypeVal = "string"
	case "int32", "uint32":
		attrTypeVal = "integer"
	default:
		attrTypeVal = "string"
	}
	dstFile.WriteString(fourTabs + "\"type\":" + "\"" + attrTypeVal + "\"" + "," + "\n")
	dstFile.WriteString(fourTabs + "\"description\":" + "\"" + attrName + "\"" + "," + "\n")
	dstFile.WriteString(fourTabs + "\"required\":" + "true," + "\n")
	dstFile.WriteString(fourTabs + "}," + "\n")
}

func writePathCompletion(dstFile *os.File) {
	dstFile.WriteString(twoTabs + " }, " + "\n")
}

func writeResourceOperation(structName string, operation string, docJsFile *os.File, str *ast.StructType) {
	writeResourceHdr(structName, operation, docJsFile)
	for _, fld := range str.Fields.List {
		if fld.Names != nil {
			switch fld.Type.(type) {
			case *ast.Ident:
				fmt.Printf("-- %s \n", fld.Names[0])
				idnt := fld.Type.(*ast.Ident)
				writeAttributeJson(fld.Names[0].Name, idnt.String(), docJsFile)
			}
		}
	}
	docJsFile.WriteString(twoTabs + " ], " + "\n")
	writeEpilogueForStruct(structName, docJsFile)
}

func WriteRestResourceDoc(docJsFile *os.File, structName string, inputFile string) {
	fset := token.NewFileSet() // positions are relative to fset

	// Parse the object file.
	f, err := parser.ParseFile(fset,
		inputFile,
		nil,
		parser.ParseComments)

	if err != nil {
		fmt.Println("Failed to parse input file ", inputFile, err)
		return
	}
	for _, dec := range f.Decls {
		tk, ok := dec.(*ast.GenDecl)
		if ok {
			for _, spec := range tk.Specs {
				switch spec.(type) {
				case *ast.TypeSpec:
					typ := spec.(*ast.TypeSpec)
					str, ok := typ.Type.(*ast.StructType)
					if typ.Name.Name == structName {
						fmt.Printf("%s \n", typ.Name.Name)
						if ok {
							docJsFile.WriteString(twoTabs + "\"/" + typ.Name.Name + "\": { \n")
							writeResourceOperation(typ.Name.Name, "post", docJsFile, str)
							writeResourceOperation(typ.Name.Name, "get", docJsFile, str)
							writeResourceOperation(typ.Name.Name, "delete", docJsFile, str)
							writePathCompletion(docJsFile)
						}
					}
				}

			}
		}
	}
}

type ObjectInfoJson struct {
	Access       string `json:"access"`
	Owner        string `json:"owner"`
	SrcFile      string `json:"srcfile"`
	Multiplicity string `json:"multiplicity"`
}

func main() {
	outFileName := "flexApis.json"

	docJsFile, err := os.Create(outFileName)
	if err != nil {
		fmt.Println("Failed to open the file")
		return
	}
	defer docJsFile.Close()

	// Write Header by copying each line from header file
	writeStaticPart("part1.txt", docJsFile)
	docJsFile.Sync()

	var jsonFilesList []string
	base := os.Getenv("SR_CODE_BASE")
	if len(base) <= 0 {
		fmt.Println(" Environment Variable SR_CODE_BASE has not been set")
		return
	}
	jsonFilesList = append(jsonFilesList, base+"/snaproute/src/models/genObjectConfig.json")
	jsonFilesList = append(jsonFilesList, base+"/snaproute/src/models/handCodedObjInfo.json")

	for _, infoFile := range jsonFilesList {
		var objMap map[string]ObjectInfoJson
		objMap = make(map[string]ObjectInfoJson, 1)
		bytes, err := ioutil.ReadFile(infoFile)
		if err != nil {
			fmt.Println("Error in reading Object configuration file", infoFile)
			return
		}
		err = json.Unmarshal(bytes, &objMap)
		for objName, objInfo := range objMap {
			WriteRestResourceDoc(docJsFile, objName,
				base+"/snaproute/src/models/"+objInfo.SrcFile)
		}
	}

	docJsFile.WriteString(twoTabs + " } " + "\n")
	docJsFile.WriteString(twoTabs + " }; " + "\n")
	writeStaticPart("part2.txt", docJsFile)
}
