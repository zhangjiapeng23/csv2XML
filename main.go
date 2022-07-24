package main

import (
	csv2xml "csv2xml/csv2XML"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Fprintln(os.Stderr, "Miss csv file path parameter!")
		os.Exit(1)
	}
	filepath := os.Args[1]
	csv := csv2xml.NewCsv(filepath)
	// csv := csv2xml.NewCsv("/Users/jameszhang/python/project/csv2xml/files/webull_AMS2.0_testcase_upload_csv.csv")
	csv.GetNodeTree()
	csv2xml.WriteXML(csv)
}
