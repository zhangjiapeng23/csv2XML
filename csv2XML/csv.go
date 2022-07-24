package csv2xml

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

var (
	nodeField = []string{"Summary", "Action", "Expected Result"}
)

type Csv struct {
	Filename     string
	NodeList     []*Node
	NodePosition map[string]int
	CsvBody      []string
	FileDir      string
}

func NewCsv(filepath string) *Csv {
	csv := &Csv{
		Filename:     getFileName(filepath),
		NodeList:     make([]*Node, 0),
		NodePosition: make(map[string]int),
		CsvBody:      make([]string, 0),
		FileDir:      getFileDir(filepath),
	}
	csv.parseNodePosition(filepath)
	return csv
}

func (csv *Csv) parseNodePosition(filePath string) {
	data, err := ioutil.ReadFile(filePath)
	newLine := ""
	lineList := make([]string, 0)
	flag := true
	csvTitle := ""
	if err != nil {
		fmt.Fprintf(os.Stderr, "file error: %v", err)
		os.Exit(1)
	}
	for _, line := range strings.Split(string(data), "\n") {
		for _, s := range line {
			if s == '"' {
				flag = !flag
			}
		}
		newLine += line + "\n"
		if flag {
			lineList = append(lineList, newLine)
			newLine = ""
		}
	}
	for index, line := range lineList {
		if index == 0 {
			csvTitle = line
		} else {
			csv.CsvBody = append(csv.CsvBody, line)
		}
	}
	csvTitleList := strings.Split(csvTitle, ",")
	for _, field := range nodeField {
		for index, title := range csvTitleList {
			if field == title {
				csv.NodePosition[field] = index
				break
			}
		}
	}
}

func (csv *Csv) GetNodeTree() {
	moduleRecord := make(map[string]string)
	root := NewNode(csv.Filename)
	csv.NodeList = append(csv.NodeList, root)
	parentId := root.Id
	for _, line := range csv.CsvBody {
		lineList := splitLine(line)
		actoinRecord := make([]string, 0)
		for _, field := range nodeField {
			index := csv.NodePosition[field]
			if field == "Summary" {
				moduleNamePattern, _ := regexp.Compile("^【(.*)】(.*)")
				match := moduleNamePattern.FindStringSubmatch(lineList[index])
				moduleName, summaryName := "", ""
				if len(match) > 1 {
					moduleName = match[1]
					summaryName = match[2]
				} else {
					moduleName = "Common"
					summaryName = lineList[index]
				}
				if id, ok := moduleRecord[moduleName]; !ok {
					parentId = root.Id
					node := NewNode(moduleName)
					node.ParentId = parentId
					parentId = node.Id
					csv.NodeList = append(csv.NodeList, node)
					moduleRecord[node.Name] = node.Id
				} else {
					parentId = id
				}
				summaryName = fmt.Sprintf("(ID %v) %v", lineList[0], summaryName)
				node := NewNode(summaryName)
				node.ParentId = parentId
				parentId = node.Id
				csv.NodeList = append(csv.NodeList, node)
			} else if field == "Action" {
				actions := lineList[index]
				for _, action := range parseStr(actions, "步骤") {
					node := NewNode(action)
					node.ParentId = parentId
					actoinRecord = append(actoinRecord, node.Id)
					csv.NodeList = append(csv.NodeList, node)
				}
			} else if field == "Expected Result" {
				expects := lineList[index]
				for index, expect := range parseStr(expects, "预期结果") {
					node := NewNode(expect)
					if index < len(actoinRecord) {
						node.ParentId = actoinRecord[index]
					} else {
						node.ParentId = actoinRecord[len(actoinRecord)-1]
					}
					csv.NodeList = append(csv.NodeList, node)
				}
			}
		}
	}
}

func getFileName(filePath string) (fileName string) {
	fileList := strings.Split(filePath, "/")
	file := fileList[len(fileList)-1]
	fileNameList := strings.Split(file, ".")
	fileName = strings.Join(fileNameList[:len(fileNameList)-1], ".")
	fileName = strings.Split(fileName, "_")[0]
	return
}

func splitLine(line string) (lineList []string) {
	flag := true
	str := ""
	for _, w := range strings.Split(line, ",") {
		str += w
		for _, s := range w {
			if s == '"' {
				flag = !flag
			}
		}
		if flag {
			lineList = append(lineList, str)
			str = ""
		}
	}
	return
}

func parseStr(s string, split string) (splitList []string) {
	linePattern, _ := regexp.Compile(".+\n?")
	splitPattern, _ := regexp.Compile(fmt.Sprintf(".*(%s[0-9]).*", split))
	lineList := linePattern.FindAllString(s, -1)
	single := ""
	for _, line := range lineList {
		line = strings.Trim(line, " ")
		line = strings.ReplaceAll(line, "\"", "")
		if splitPattern.MatchString(line) {
			if len(single) > 0 {
				splitList = append(splitList, single)
				single = line
			} else {
				single += line
			}
		} else {
			single += line
		}
	}
	splitList = append(splitList, single)
	return

}

func getFileDir(filePath string) (fileDir string) {
	fileList := strings.Split(filePath, "/")
	fileDir = strings.Join(fileList[:len(fileList)-1], "/")
	return
}
