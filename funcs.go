package main

import (
	"io"
	"log"
	"os"
	"strings"

	"github.com/biogo/hts/bam"
	"github.com/biogo/hts/bgzf"
)

func getBamFiles(folder string) []string {
	var files []string
	f, err := os.Open(folder)
	if err != nil {
		log.Fatalln("Error opening folder:", err)
	}
	defer f.Close()
	fileInfo, err := f.Readdir(-1)
	if err != nil {
		log.Fatalln("Error reading folder:", err)
	}
	for _, file := range fileInfo {
		if strings.HasSuffix(file.Name(), ".bam") {
			files = append(files, folder+"/"+file.Name())
		}
	}
	return files
}
func processBam(path string, threads int, rs7412 int, rs429358 int, qual int, chr bool, minReadLength int, maxReadLength int) APOE {
	chrName := "19"
	chrExcess := "20"
	if chr {
		chrName = "chr19"
		chrExcess = "chr20"
	}
	// Get the sample name from the path
	sampleName := strings.Split(strings.Split(path, "/")[len(strings.Split(path, "/"))-1], ".")[0]
	// Define apoe variable of type APOE
	apoe := APOE{
		SampleName: sampleName,
		APOE1:      0,
		APOE2:      0,
		APOE3:      0,
		APOE4:      0,
	}
	var r io.Reader
	f, err := os.Open(path)
	if err != nil {
		log.Fatalln("Error opening BAM file:", err)
	}
	defer f.Close()
	ok, err := bgzf.HasEOF(f)
	if err != nil {
		log.Fatalln("Error checking EOF:", err)
	}
	if !ok {
		log.Fatalf("EOF not found, attempting to read anyway")
	}
	r = f
	b, err := bam.NewReader(r, threads)
	if err != nil {
		log.Fatalln("Error reading BAM file:", err)
	}
	defer b.Close()
	for {
		rec, err := b.Read()
		rs429358Status := ""
		rs7412Status := ""
		trigger := false
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln("Error reading BAM file:", err)
		}
		if rec.Ref.Name() == chrName && rec.Len() <= maxReadLength && rec.Len() >= minReadLength {
			// Check if QUAL and SEQUENCE are of same length
			if rec.Start() <= rs429358 && rec.End() >= rs429358 {
				relPos := rs429358 - rec.Start()
				trigger = true
				if len(rec.Qual) <= relPos {
					continue
				}
				currQual := int(rec.Qual[relPos])
				if currQual < qual {
					continue
				}
				alleleA := ((strings.Split(strings.Split(rec.String(), " ")[9], ""))[relPos])
				if alleleA == "C" || alleleA == "G" {
					rs429358Status = "wildtype"
				} else if alleleA == "T" || alleleA == "A" {
					rs429358Status = "mutant"
				}
				rs7412Status = "wildtype"
			}
			if rec.Start() <= rs7412 && rec.End() >= rs7412 {
				relPos := rs7412 - rec.Start()
				trigger = true
				if len(rec.Qual) <= relPos {
					continue
				}
				currQual := int(rec.Qual[relPos])
				if currQual < qual {
					continue
				}
				alleleB := ((strings.Split(strings.Split(rec.String(), " ")[9], ""))[relPos])
				if alleleB == "C" || alleleB == "G" {
					rs7412Status = "wildtype"
				} else if alleleB == "T" || alleleB == "A" {
					rs7412Status = "mutant"
				}
				rs429358Status = "wildtype"
			}
			if trigger {
				if rs429358Status == "wildtype" && rs7412Status == "wildtype" {
					apoe.APOE3++
				} else if rs429358Status == "mutant" && rs7412Status == "wildtype" {
					apoe.APOE4++
				} else if rs429358Status == "wildtype" && rs7412Status == "mutant" {
					apoe.APOE2++
				} else if rs429358Status == "mutant" && rs7412Status == "mutant" {
					apoe.APOE1++
				}
			}
		} else if rec.Ref.Name() == chrExcess {
			break
		}
	}
	return apoe
}

func writeResult(result string, file string) {
	f, err := os.Create(file)
	if err != nil {
		log.Fatalln("Error creating result file:", err)
	}
	defer f.Close()
	_, err = f.WriteString(result)
	if err != nil {
		log.Fatalln("Error writing to result file:", err)
	}
}
