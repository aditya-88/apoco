# APOCO #

**APO**E reads **Co**unter

## Introduction ##

The program counts the number of reads pertaining to specific `APOE` alleles from `BAM` files in a given folder.

## Installation ##

Go to the releases section and downlod the binary file for your system.

## Usage ##

*The default threads are counted based on the system. The number shown as default here pertains to the system this program was built on.*

```bash
Usage of apoco:
  -chr
        Use this flag if the reference chromosomes are named with "chr" in the names (e.g. chr1, chr2, chrX, etc.)
  -f string
        Folder containing BAM files
  -hg int
        Human genome version (19 or 38) (default 38)
  -max int
        Maximum read length (default 150)
  -min int
        Minimum read length (default 100)
  -qual int
        Minimum mapping quality (default 30)
  -t int
        Number of threads to use (default 32)
```

## Results ##

The program outputs a tab delimited `apoeCounts.txt` file with the sample name and raw read counts for each `APOE` allele.
Tis file is created in the same directory as the one provided to the program, so, please make sure that you have R/W permission.

Example:

```tsv
Sample	APOE1	APOE2	APOE3	APOE4
TestA	0	2	9498	78
TestB	0	0	6284	4875
```
