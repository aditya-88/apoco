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
		log.Fatalln("\nError opening BAM file:", err)
	}
	defer f.Close()
	ok, err := bgzf.HasEOF(f)
	if err != nil {
		// Report error and continue
		log.Println("\nError checking for EOF:", err)
		log.Println("\nSkipping file")
		return apoe
	}
	if !ok {
		log.Println("\nEOF not found in file:", path)
		log.Println("Skipping file")
		return apoe
	}
	r = f
	b, err := bam.NewReader(r, threads)
	if err != nil {
		log.Println("\nError reading BAM file:", err)
		log.Println("Skipping file")
		return apoe
	}
	defer b.Close()
	for {
		rec, err := b.Read()
		rs429358Status := "wildtype"
		rs7412Status := "wildtype"
		trigger := false
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln("\nError reading BAM file:", err)
		}
		if rec.Ref.Name() == chrName && rec.Len() <= maxReadLength && rec.Len() >= minReadLength {
			// Check if QUAL and SEQUENCE are of same length
			if rec.Start() <= rs429358 && rec.End() >= rs429358 {
				relPos := rs429358 - rec.Start()
				if relPos > 0 {
					relPos = relPos - 1
				}
				if len(rec.Qual) <= relPos {
					continue
				}
				trigger = true
				currQual := int(rec.Qual[relPos])
				if currQual < qual {
					continue
				}
				alleleA := ((strings.Split(strings.Split(rec.String(), " ")[9], ""))[relPos])
				if alleleA == "T" {
					rs429358Status = "wildtype"
				} else if alleleA == "C" {
					rs429358Status = "mutant"
				}
			}
			if rec.Start() <= rs7412 && rec.End() >= rs7412 {
				relPos := rs7412 - rec.Start()
				if relPos > 0 {
					relPos = relPos - 1
				}
				if len(rec.Qual) <= relPos {
					continue
				}
				trigger = true
				currQual := int(rec.Qual[relPos])
				if currQual < qual {
					continue
				}
				alleleB := ((strings.Split(strings.Split(rec.String(), " ")[9], ""))[relPos])
				if alleleB == "C" {
					rs7412Status = "wildtype"
				} else if alleleB == "T" {
					rs7412Status = "mutant"
				}
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
	//Check if output file exists, if not, create it and add a header
	if _, err := os.Stat(file); os.IsNotExist(err) {
		f, err := os.Create(file)
		if err != nil {
			log.Fatalln("\nError creating output file:", err)
		}
		defer f.Close()
		_, err = f.WriteString("Sample\tAPOE1\tAPOE2\tAPOE3\tAPOE4\n")
		if err != nil {
			log.Fatalln("\nError writing to output file:", err)
		}
	}
	f, err := os.OpenFile(file, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatalln("\nError opening output file:", err)
	}
	defer f.Close()
	_, err = f.WriteString(result)
}
