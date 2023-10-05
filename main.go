package main

import (
	"flag"
	"fmt"
	"runtime"
)

var (
	software string = "APOCO - APOE allele read counter"
	version  string = "0.2.1-beta"
	dev      string = "Aditya Singh"
	gitHub   string = "https://www.github.com/aditya-88"
	folder   string
	threads  int
	rs7412   int = 44908822
	rs429358 int = 44908684
	minQual  int
	hg       int
	chr      bool
	min      int
	max      int
	result   string = "Sample\tAPOE1\tAPOE2\tAPOE3\tAPOE4\n"
)

type APOE struct {
	SampleName string
	APOE1      int
	APOE2      int
	APOE3      int
	APOE4      int
}

func flagsProcess() {
	flag.StringVar(&folder, "f", "", "Folder containing BAM files")
	flag.IntVar(&threads, "t", runtime.NumCPU(), "Number of threads to use")
	flag.IntVar(&hg, "hg", 38, "Human genome version (19 or 38)")
	flag.IntVar(&minQual, "qual", 30, "Minimum mapping quality")
	flag.BoolVar(&chr, "chr", false, "Use this flag if the reference chromosomes are named with \"chr\" in the names (e.g. chr1, chr2, chrX, etc.)")
	flag.IntVar(&min, "min", 100, "Minimum read length")
	flag.IntVar(&max, "max", 150, "Maximum read length")
	flag.Parse()
}
func main() {
	flagsProcess()
	fmt.Printf("Welcome to %s v%s\nMake sure that the BAM file(s) are coordinate sorted!\n", software, version)
	fmt.Println("Developed by:", dev)
	fmt.Println("GitHub:", gitHub)
	// Check if required flags are set
	if folder == "" {
		fmt.Println("\n>>Missing required flag(s)<<")
		flag.Usage()
		return
	}
	// Get a list of BAM files in the folder and process them
	files := getBamFiles(folder)
	// If no BAM files are found, exit
	if len(files) == 0 {
		fmt.Println("No BAM files found in folder:", folder)
		return
	}
	if hg == 19 || hg == 37 {
		rs7412 = 45412079
		rs429358 = 45411941
	}
	// Print the setup of the program
	fmt.Println("##################################################")
	fmt.Printf("Folder: %s\nFound BAM files: %d\nThreads: %d\nAssembly: hg%d\nrs7412: %d\nrs429358: %d\nMinimum mapping quality: %d\nMinimum read length: %d\nMaxmimum read length: %d\n", folder, len(files), threads, hg, rs7412, rs429358, minQual, min, max)
	fmt.Println("##################################################")
	// Process each file
	for _, file := range files {
		fmt.Println("Processing file:", file)
		apoe := processBam(file, threads, rs7412, rs429358, minQual, chr, min, max)
		result += fmt.Sprintf("%s\t%d\t%d\t%d\t%d\n", apoe.SampleName, apoe.APOE1, apoe.APOE2, apoe.APOE3, apoe.APOE4)
	}
	fmt.Println("##################################################")
	fmt.Println("Writing results to file...")
	writeResult(result, folder)
	fmt.Println("Done!")
}
