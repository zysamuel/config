package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
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

func writeAttributeJson(attrName string, attrInfo ObjectMembersInfo, dstFile *os.File) {
	var attrTypeVal string
	dstFile.WriteString(fourTabs + "{" + "\n")
	dstFile.WriteString(fourTabs + "\"in\": \"formData\"," + "\n")
	dstFile.WriteString(fourTabs + "\"name\":" + "\"" + attrName + "\"" + "," + "\n")
	switch attrInfo.VarType {
	case "string":
		attrTypeVal = "string"
	case "int32", "uint32":
		attrTypeVal = "integer"
	default:
		attrTypeVal = "string"
	}
	//description := strings.Trim(attrInfo.Description, "\n")
	description := strings.Replace(attrInfo.Description, "\n", " ", -1)

	var isRequired bool
	if attrInfo.DefaultVal == "" {
		isRequired = true
	} else {
		isRequired = false
	}
	dstFile.WriteString(fourTabs + "\"type\":" + "\"" + attrTypeVal + "\"" + "," + "\n")
	dstFile.WriteString(fourTabs + "\"description\":" + "\"" + description + "\"" + "," + "\n")
	if isRequired {
		dstFile.WriteString(fourTabs + "\"required\":" + "true," + "\n")
	} else {
		dstFile.WriteString(fourTabs + "\"required\":" + "false," + "\n")
	}

	dstFile.WriteString(fourTabs + "}," + "\n")
}

func writePathCompletion(dstFile *os.File) {
	dstFile.WriteString(twoTabs + " }, " + "\n")
}

func writeResourceOperation(structName string, operation string, docJsFile *os.File, membersMap map[string]ObjectMembersInfo) {
	writeResourceHdr(structName, operation, docJsFile)
	for attrName, attrInfo := range membersMap {
		writeAttributeJson(attrName, attrInfo, docJsFile)
	}
	docJsFile.WriteString(twoTabs + " ], " + "\n")
	writeEpilogueForStruct(structName, docJsFile)
}

func WriteConfigObject(structName string, docJsFile *os.File, membersMap map[string]ObjectMembersInfo) {

	docJsFile.WriteString(twoTabs + "\"/" + structName + "\": { \n")
	writeResourceOperation(structName, "post", docJsFile, membersMap)
	writePathCompletion(docJsFile)

	docJsFile.WriteString(twoTabs + "\"/" + structName + "/{object-id}\": { \n")
	writeResourceOperation(structName, "get", docJsFile, membersMap)
	writeResourceOperation(structName, "delete", docJsFile, membersMap)
	writeResourceOperation(structName, "patch", docJsFile, membersMap)
	writePathCompletion(docJsFile)
}

func WriteStateObject(structName string, docJsFile *os.File, membersMap map[string]ObjectMembersInfo) {
	docJsFile.WriteString(twoTabs + "\"/" + structName + "\": { \n")
	writeResourceOperation(structName, "get", docJsFile, membersMap)
	writePathCompletion(docJsFile)
}

func WriteGlobalStateObject(structName string, docJsFile *os.File, membersMap map[string]ObjectMembersInfo) {
	docJsFile.WriteString(twoTabs + "\"/" + structName + "\": { \n")
	writePathCompletion(docJsFile)
}

func WriteRestResourceDoc(docJsFile *os.File, structName string, membersMap map[string]ObjectMembersInfo, objInfo ObjectInfoJson) {
	if objInfo.Access == "w" || objInfo.Access == "rw" {
		WriteConfigObject(structName, docJsFile, membersMap)
		//fmt.Println("Config Object: ", structName)
	} else if objInfo.Access == "r" || objInfo.Access == "rw" {
		if strings.Contains(structName, "Global") {
			WriteGlobalStateObject(structName, docJsFile, membersMap)
			//fmt.Println("Global Config Object: ", structName)
		} else {
			WriteStateObject(structName, docJsFile, membersMap)
			//fmt.Println("State Object: ", structName)
		}
	}
}

type ObjectInfoJson struct {
	Access       string `json:"access"`
	Owner        string `json:"owner"`
	SrcFile      string `json:"srcfile"`
	Multiplicity string `json:"multiplicity"`
}

type ObjectMembersInfo struct {
	VarType     string `json:"type"`
	IsKey       bool   `json:"isKey"`
	IsArray     bool   `json:"isArray"`
	Description string `json:"description"`
	DefaultVal  string `json:"default"`
	Position    int    `json:"position"`
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
	idx := 0
	jsonFilesList = append(jsonFilesList, base+"/snaproute/src/models/genObjectConfig.json")
	membersInfoBase := base + "/reltools/codegentools/._genInfo/"
	for _, infoFile := range jsonFilesList {
		fmt.Println("Info File", infoFile)
		var objMap map[string]ObjectInfoJson
		objMap = make(map[string]ObjectInfoJson, 1)
		bytes, err := ioutil.ReadFile(infoFile)
		if err != nil {
			fmt.Println("Error in reading Object configuration file", infoFile)
			continue
		}
		err = json.Unmarshal(bytes, &objMap)
		//fmt.Println("Object Map ", objMap)
		//return
		objList := make([]string, 0)
		for key, _ := range objMap {
			objList = append(objList, key)
		}
		sort.Strings(objList)

		//for objName, objInfo := range objMap {
		for _, objName := range objList {
		   idx++
			objInfo := objMap[objName]
			membersFile := membersInfoBase + objName + "Members.json"
			var memberMap map[string]ObjectMembersInfo
			memberMap = make(map[string]ObjectMembersInfo, 1)
			bytes, err := ioutil.ReadFile(membersFile)
			if err != nil {
				fmt.Println("Error in reading Object configuration file", membersFile)
				continue
			}
			err = json.Unmarshal(bytes, &memberMap)

			WriteRestResourceDoc(docJsFile, objName, memberMap, objInfo)
		}
	}
	fmt.Println("Total Objects ", idx)
	docJsFile.WriteString(twoTabs + " } " + "\n")
	docJsFile.WriteString(twoTabs + " }; " + "\n")
	writeStaticPart("part2.txt", docJsFile)
}
