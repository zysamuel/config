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
	"regexp"
	"strings"
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

func writeAttributeJson(attrName string, attrType string, dstFile *os.File, fld *ast.Field) {
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
	description, isRquired := getSpecialTagsForAttribute(fld)
	dstFile.WriteString(fourTabs + "\"type\":" + "\"" + attrTypeVal + "\"" + "," + "\n")
	dstFile.WriteString(fourTabs + "\"description\":" + "\"" + description + "\"" + "," + "\n")
	if isRquired {
		dstFile.WriteString(fourTabs + "\"required\":" + "true," + "\n")
	} else {
		dstFile.WriteString(fourTabs + "\"required\":" + "false," + "\n")
	}

	dstFile.WriteString(fourTabs + "}," + "\n")
}

func writePathCompletion(dstFile *os.File) {
	dstFile.WriteString(twoTabs + " }, " + "\n")
}

func getSpecialTagsForAttribute(fld *ast.Field) (description string, isRequired bool) {
	reg, err := regexp.Compile("[`\"]")
	if err != nil {
		fmt.Println("Error in regex ", err)
	}
	if fld.Tag != nil {
		tags := reg.ReplaceAllString(fld.Tag.Value, "")
		splits := strings.Split(tags, ",")
		for _, part := range splits {
			keys := strings.Split(part, ":")
			for idx, key := range keys {
				alphas, err := regexp.Compile("[^A-Za-z]")
				if err != nil {
					fmt.Println("Error in regex ", err)
				}
				key = alphas.ReplaceAllString(key, "")
				switch key {
				case "DESCRIPTION":
					description = keys[idx+1]
					description = strings.Replace(description, "\n", " ", -1)
				case "DEFAULT":
					isRequired = false
				}
			}
		}
	}
	return description, isRequired
}

func writeResourceOperation(structName string, operation string, docJsFile *os.File, str *ast.StructType) {
	writeResourceHdr(structName, operation, docJsFile)
	for _, fld := range str.Fields.List {
		if fld.Names != nil {
			switch fld.Type.(type) {

			case *ast.ArrayType:
				//fmt.Println("### Array Type attribute ", fld.Names[0].Name)
				//arrayInfo := fld.Type.(*ast.ArrayType)
				//info := ObjectMembersInfo{}
				//info.IsArray = true
				//objMembers[varName] = info
				//idntType := arrayInfo.Elt.(*ast.Ident)
				//varType := idntType.String()
				//info.VarType = varType
				//objMembers[varName] = info
				//if fld.Tag != nil {
				//	getSpecialTagsForAttribute(fld.Tag.Value, &info)
				//}
			case *ast.Ident:
				//fmt.Printf("-- %s \n", fld.Names[0])
				idnt := fld.Type.(*ast.Ident)
				writeAttributeJson(fld.Names[0].Name, idnt.String(), docJsFile, fld)
			}
		}
	}
	docJsFile.WriteString(twoTabs + " ], " + "\n")
	writeEpilogueForStruct(structName, docJsFile)
}

func WriteConfigObject(structName string, docJsFile *os.File, str *ast.StructType) {

	docJsFile.WriteString(twoTabs + "\"/" + structName + "\": { \n")
	writeResourceOperation(structName, "post", docJsFile, str)
	writePathCompletion(docJsFile)

	docJsFile.WriteString(twoTabs + "\"/" + structName + "/{object-id}\": { \n")
	writeResourceOperation(structName, "get", docJsFile, str)
	writeResourceOperation(structName, "delete", docJsFile, str)
	writeResourceOperation(structName, "patch", docJsFile, str)
	writePathCompletion(docJsFile)
}

func WriteStateObject(structName string, docJsFile *os.File, str *ast.StructType) {
	docJsFile.WriteString(twoTabs + "\"/" + structName + "\": { \n")
	writeResourceOperation(structName, "get", docJsFile, str)
	writePathCompletion(docJsFile)
}

func WriteGlobalStateObject(structName string, docJsFile *os.File, str *ast.StructType) {
	docJsFile.WriteString(twoTabs + "\"/" + structName + "\": { \n")
	writeResourceOperation(structName, "get", docJsFile, str)
	writePathCompletion(docJsFile)
}

func WriteRestResourceDoc(docJsFile *os.File, structName string, inputFile string, objInfo ObjectInfoJson) {
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
						if ok {
							if objInfo.Access == "w" || objInfo.Access == "rw" {
								WriteConfigObject(typ.Name.Name, docJsFile, str)
							} else if objInfo.Access == "r" || objInfo.Access == "rw" {
								if strings.Contains(typ.Name.Name, "Global") {
									WriteGlobalStateObject(typ.Name.Name, docJsFile, str)
								} else {
									WriteStateObject(typ.Name.Name+"s", docJsFile, str)
								}
							}
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
				base+"/snaproute/src/models/"+objInfo.SrcFile, objInfo)
		}
	}

	docJsFile.WriteString(twoTabs + " } " + "\n")
	docJsFile.WriteString(twoTabs + " }; " + "\n")
	writeStaticPart("part2.txt", docJsFile)
}
