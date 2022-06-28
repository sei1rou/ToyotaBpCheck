package main

import (
	"encoding/csv"
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func failOnError(err error) {
	if err != nil {
		log.Fatal("Error:", err)
	}
}

func main() {
	flag.Parse()

	//ログファイル準備
	/*
		logfile, err := os.OpenFile("./log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
		failOnError(err)
		defer logfile.Close()

		log.SetOutput(logfile)
	*/

	log.Print("Start\r\n")

	// ファイルを読み込んで二次元配列に入れる
	records := readfile(flag.Arg(0))

	// ファイルへ書き出す
	savefile(flag.Arg(0), records)

	log.Print("Finesh !\r\n")

}

func readfile(filename string) [][]string {
	//入力ファイル準備
	infile, err := os.Open(filename)
	failOnError(err)
	defer infile.Close()

	reader := csv.NewReader(transform.NewReader(infile, japanese.ShiftJIS.NewDecoder()))
	reader.Comma = '\t'

	//CSVファイルを２次元配列に展開
	readrecords := make([][]string, 0)
	for {
		record, err := reader.Read() // 1行読み出す
		if err == io.EOF {
			break
		} else {
			failOnError(err)
		}

		readrecords = append(readrecords, record)
	}

	return readrecords
}

func savefile(filename string, saverecords [][]string) {
	//出力ファイル準備
	outDir, outfileName := filepath.Split(filename)
	pos := strings.LastIndex(outfileName, ".")
	// outfile, err := os.Create(outDir + outfileName[:pos] + "d.txt")
	outfile, err := os.Create(outDir + outfileName[:pos] + ".txt")
	failOnError(err)
	defer outfile.Close()

	writer := csv.NewWriter(transform.NewWriter(outfile, japanese.ShiftJIS.NewEncoder()))
	writer.Comma = '\t'
	writer.UseCRLF = true

	for i, out_record := range saverecords {

		if i == 0 {
			out_record = append(out_record, "T判定")
			out_record = append(out_record, "判定の変更")
			writer.Write(out_record)
		} else {

			JNo := out_record[1]
			out_record[1] = JNo[8:]

			h1 := hantei(out_record[5], out_record[9])
			h2 := hantei(out_record[7], out_record[11])
			hh := hantei(h1, h2)
			out_record = append(out_record, hh)

			nh := string(norm.NFKC.Bytes([]byte(out_record[12])))
			if hh == nh || nh == "G" {
				out_record = append(out_record, "")
			} else {
				out_record = append(out_record, "※")
			}

			writer.Write(out_record)
		}

	}

	writer.Flush()

}

func hantei(v1, v2 string) string {

	var h string
	if v1 == "" {
		h = v2
	} else if v2 == "" {
		h = v1
	} else {
		if v1 == "A" {
			if v2 == "B" || v2 == "C" || v2 == "D" || v2 == "E" || v2 == "F" {
				h = v2
			} else {
				h = v1
			}
		} else if v1 == "B" {
			if v2 == "C" || v2 == "D" || v2 == "E" || v2 == "F" {
				h = v2
			} else {
				h = v1
			}
		} else if v1 == "C" {
			if v2 == "D" || v2 == "E" || v2 == "F" {
				h = v2
			} else {
				h = v1
			}
		} else if v1 == "D" {
			if v2 == "E" || v2 == "F" {
				h = v2
			} else {
				h = v1
			}
		} else if v1 == "E" {
			if v2 == "F" {
				h = v2
			} else {
				h = v1
			}
		} else if v1 == "F" {
			h = v1
		}
	}

	return h

}
