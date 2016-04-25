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

const (
	NO_ATTTRS = iota
	SELECTED_ATTR
	KEY_ATTRS
	ALL_ATTRS
)

var opDescMap map[string]string = make(map[string]string, 1)

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
	dstFile.WriteString(twoTabs + "\"summary\": \"" + opDescMap[operation] + strName + "\"," + "\n")
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

func writeAttributeJson(attrInfo AttributeListItem, dstFile *os.File) {
	var attrTypeVal string
	dstFile.WriteString(fourTabs + "{" + "\n")
	dstFile.WriteString(fourTabs + "\"in\": \"formData\"," + "\n")
	//fmt.Println("### AttrInfo ", attrInfo)
	dstFile.WriteString(fourTabs + "\"name\":" + "\"" + attrInfo.AttrName + "\"" + "," + "\n")
	switch attrInfo.VarType {
	case "string":
		attrTypeVal = "string"
	case "int32", "uint32":
		attrTypeVal = "integer"

	case "bool":
		attrTypeVal = "boolean"
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
		dstFile.WriteString(fourTabs + "\"default\":" + "\"" + attrInfo.DefaultVal + "\",\n")
	}
	if attrInfo.Selections != "" {
		enumVals := strings.Split(attrInfo.Selections, "/")
		enumStr := ""
		for idx, val := range enumVals {
			if idx != 0 {
				enumStr = enumStr + "," + "\"" + val + "\""
			} else {
				enumStr = enumStr + "\"" + val + "\""
			}
		}
		dstFile.WriteString(fourTabs + "\"enum\":" + "[" + enumStr + "],\n")
	}

	dstFile.WriteString(fourTabs + "}," + "\n")
}

func writePathCompletion(dstFile *os.File) {
	dstFile.WriteString(twoTabs + " }, " + "\n")
}

func writeResourceOperation(structName string, operation string, docJsFile *os.File, membersInfo []AttributeListItem, mode int) {
	writeResourceHdr(structName, operation, docJsFile)
	switch mode {
	case ALL_ATTRS, KEY_ATTRS, SELECTED_ATTR:
		for _, attrInfo := range membersInfo {
			if mode == ALL_ATTRS || mode == SELECTED_ATTR {
				writeAttributeJson(attrInfo, docJsFile)
			}
			if mode == KEY_ATTRS && attrInfo.IsKey == true {
				writeAttributeJson(attrInfo, docJsFile)
			}
		}
	}
	docJsFile.WriteString(twoTabs + " ], " + "\n")
	writeEpilogueForStruct(structName, docJsFile)
}

func WriteConfigObject(structName string, docJsFile *os.File, membersInfo []AttributeListItem) {

	docJsFile.WriteString(twoTabs + "\"/Config/" + structName + "\": { \n")
	writeResourceOperation(structName, "post", docJsFile, membersInfo, ALL_ATTRS)
	writePathCompletion(docJsFile)

	docJsFile.WriteString(twoTabs + "\"/Config/" + structName + "/{object-id}\": { \n")
	writeResourceOperation(structName, "get", docJsFile, membersInfo, NO_ATTTRS)
	writeResourceOperation(structName, "delete", docJsFile, membersInfo, NO_ATTTRS)
	writeResourceOperation(structName, "patch", docJsFile, membersInfo, ALL_ATTRS)
	writePathCompletion(docJsFile)

	docJsFile.WriteString(twoTabs + "\"/Config/" + structName + "/\": { \n")
	writeResourceOperation(structName, "get", docJsFile, membersInfo, KEY_ATTRS)
	writeResourceOperation(structName, "delete", docJsFile, membersInfo, KEY_ATTRS)
	writeResourceOperation(structName, "patch", docJsFile, membersInfo, ALL_ATTRS)
	writePathCompletion(docJsFile)
}

func WriteStateObject(structName string, docJsFile *os.File, membersInfo []AttributeListItem) {
	for idx, attrInfo := range membersInfo {
		if strings.Contains(strings.ToLower(attrInfo.QueryParam), "optional") {
			docJsFile.WriteString(twoTabs + "\"/state/" + structName + "/{" + attrInfo.AttrName + "}\" : {\n")
			mbrInfo := make([]AttributeListItem, 1)
			mbrInfo[0] = membersInfo[idx]
			writeResourceOperation(structName, "get", docJsFile, mbrInfo, SELECTED_ATTR)
			writePathCompletion(docJsFile)
		}
	}
	docJsFile.WriteString(twoTabs + "\"/state/" + structName + "s\": { \n")
	writeResourceOperation(structName, "get", docJsFile, membersInfo, NO_ATTTRS)
	writePathCompletion(docJsFile)
}

func WriteGlobalStateObject(structName string, docJsFile *os.File, membersInfo []AttributeListItem) {
	docJsFile.WriteString(twoTabs + "\"/state/" + structName + "\": { \n")
	writePathCompletion(docJsFile)
}

func WriteRestResourceDoc(docJsFile *os.File, structName string, membersInfo []AttributeListItem, objInfo ObjectInfoJson) {
	if objInfo.Access == "w" || objInfo.Access == "rw" {
		WriteConfigObject(structName, docJsFile, membersInfo)
	} else if objInfo.Access == "r" || objInfo.Access == "rw" {
		if strings.Contains(structName, "Global") {
			WriteGlobalStateObject(structName, docJsFile, membersInfo)
		} else {
			WriteStateObject(structName, docJsFile, membersInfo)
		}
	}
}

func WriteObjectList(docJsFile *os.File, objMap map[string]ObjectInfoJson, objList []string, membersInfoBase string) {
	idx := 0
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

		attrList := make([]AttributeListItem, len(memberMap))
		for key, info := range memberMap {
			var item AttributeListItem
			item.AttrName = key
			item.VarType = info.VarType
			item.IsKey = info.IsKey
			item.IsArray = info.IsArray
			item.Description = info.Description
			item.DefaultVal = info.DefaultVal
			item.Position = info.Position
			item.Selections = info.Selections
			item.QueryParam = info.QueryParam
			if item.Position > 0 {
				attrList[item.Position-1] = item
			}
		}
		WriteRestResourceDoc(docJsFile, objName, attrList, objInfo)
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
	Selections  string `json:"selections"`
	QueryParam  string `json:"queryparam"`
	Accelerated bool   `json:"accelerated"`
	Min         int    `json:"min"`
	Max         int    `json:"max"`
	Len         int    `json:"len"`
}

type AttributeListItem struct {
	ObjectMembersInfo
	AttrName string
}

func main() {
	outFileName := "flexApis.json"
	opDescMap["post"] = "Create New "
	opDescMap["delete"] = "Delete "
	opDescMap["get"] = "Query "
	opDescMap["patch"] = "Update existing "
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
		cfgObjList := make([]string, 0)
		stateObjList := make([]string, 0)
		for key, objInfo := range objMap {
			idx++
			if objInfo.Access == "w" || objInfo.Access == "rw" {
				cfgObjList = append(cfgObjList, key)
			} else if objInfo.Access == "r" {
				stateObjList = append(stateObjList, key)
			}
		}
		sort.Strings(cfgObjList)
		sort.Strings(stateObjList)
		WriteObjectList(docJsFile, objMap, cfgObjList, membersInfoBase)
		WriteObjectList(docJsFile, objMap, stateObjList, membersInfoBase)
	}
	fmt.Println("Total Objects ", idx)
	docJsFile.WriteString(twoTabs + " } " + "\n")
	docJsFile.WriteString(twoTabs + " }; " + "\n")
	writeStaticPart("part2.txt", docJsFile)
}
