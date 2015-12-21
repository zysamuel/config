package main

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	//"io/ioutil"
	"log"
	"os"
)

var oneTabs string = "\t"
var twoTabs string = "\t\t"
var threeTabs string = "\t\t\t"
var fourTabs string = "\t\t\t\t"

func writeHeaders(dstFile *os.File) {
	hdrFile := "header.txt"
	hdr, err := os.Open(hdrFile)
	if err != nil {
		log.Fatal(err)
	}
	defer hdr.Close()

	scanner := bufio.NewScanner(hdr)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
		dstFile.WriteString(line + "\n")
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func writeResourceHdr(strName string, dstFile *os.File) {
	dstFile.WriteString(twoTabs + "\"/" + strName + "\": { \n")
	dstFile.WriteString(twoTabs + "\"post\": { " + "\n")
	dstFile.WriteString(threeTabs+ "\"tags\": [ " + "\n")
	dstFile.WriteString(fourTabs+ "\"" + strName + "\"" + "\n")
	dstFile.WriteString(threeTabs+ "]," + "\n")
	dstFile.WriteString(twoTabs + "\"summary\": \"Create New" + strName +"\","+ "\n")
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
	dstFile.WriteString(twoTabs + " } " + "\n")
}

func writeAttributeJson (attrName string, dstFile *os.File) {
	dstFile.WriteString(fourTabs + "{" + "\n")
	dstFile.WriteString(fourTabs + "\"in\": \"body\","  + "\n")
	dstFile.WriteString(fourTabs + "\"name\":"  + "\""+  attrName+ "\""+ ","+  "\n")
	dstFile.WriteString(fourTabs + "\"description\":"  + "\""+  attrName+ "\""+ ","+  "\n")
	dstFile.WriteString(fourTabs + "\"required\":" + "true," + "\n")
	dstFile.WriteString(fourTabs + "}," + "\n")
}

func writePathCompletion(dstFile *os.File) {
	dstFile.WriteString(twoTabs + " }, " + "\n")
}
func main() {
	fset := token.NewFileSet() // positions are relative to fset
	outFileName := "flexApis.json"

	inputFile := "../../models/objects.go"

	docJsFile, err := os.Create(outFileName)
	if err != nil {
		fmt.Println("Failed to open the file")
		return
	}
	defer docJsFile.Close()

	// Write Header by copying each line from header file
	writeHeaders(docJsFile)
	docJsFile.Sync()

	// Parse the object file.
	f, err := parser.ParseFile(fset,
		inputFile,
		nil,
		parser.ParseComments)

	if err != nil {
		fmt.Println("Failed to parse input file ", inputFile, err)
		return
	}
	/*
	   "/pet": {
	        "post": {
	          "tags": [
	            "pet"
	          ],
	          "summary": "Add a new pet to the store",
	          "description": "",
	          "operationId": "addPet",
	          "consumes": [
	            "application/json",
	            "application/xml"
	          ],
	          "produces": [
	            "application/json",
	            "application/xml"
	          ],
	          "parameters": [
	            {
	              "in": "body",
	              "name": "pet",
	              "description": "Pet object that needs to be added to the store",
	              "required": false,
	              "schema": {
	                "$ref": "#/definitions/Pet"
	              }
	            }
	          ],
	          "responses": {
	            "405": {
	              "description": "Invalid input"
	            }
	          },
	          "security": [
	            {
	              "petstore_auth": [
	                "write:pets",
	                "read:pets"
	              ]
	            }
	          ]
	        },
	*/
	for _, dec := range f.Decls {
		tk, ok := dec.(*ast.GenDecl)
		if ok {
			for _, spec := range tk.Specs {
				switch spec.(type) {
				case *ast.TypeSpec:
					typ := spec.(*ast.TypeSpec)
					fmt.Printf("%s \n", typ.Name)

					str, ok := typ.Type.(*ast.StructType)
					if (typ.Name.Name == "BGPNeighborConfig") || (typ.Name.Name == "IPv4Intf") {
						if ok {
							writeResourceHdr(typ.Name.Name, docJsFile)
							for _, fld := range str.Fields.List {
								if fld.Names != nil {
								   writeAttributeJson(fld.Names[0].Name,  docJsFile)
									fmt.Printf("-- %s : %s \n", fld.Names[0], fld.Type)
								}
							}
							docJsFile.WriteString(twoTabs + " ], " + "\n")
							writeEpilogueForStruct(typ.Name.Name, docJsFile)
							writePathCompletion(docJsFile)
						}

					}
				}

			}
		}
	}
	docJsFile.WriteString(twoTabs + " } " + "\n")
	docJsFile.WriteString(twoTabs + " } " + "\n")

}
