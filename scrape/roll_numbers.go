package Result_NITH

import (
	"fmt"
	"strconv"
)

func GenRollNumbersByBranchAndYear(branchCode string, yearCode int) (rollNumbers []string) {
	if yearCode == 19 {
		for i := 0; i < 150; i++ {
			if branchCode == "45" || branchCode == "55" {
				if i > 100 {
					break
				}
				rollNumbers = append(rollNumbers, fmt.Sprintf("19%s%.2d", branchCode, i))
			} else {
				rollNumbers = append(rollNumbers, fmt.Sprintf("19%s%.3d", branchCode, i))
			}

		}

	} else {
		for i := 0; i < 150; i++ {
			if branchCode == "dcs" || branchCode == "dec" && i > 100 {
				break
			}
			rollNumbers = append(rollNumbers, fmt.Sprintf("%d%s%.3d", yearCode, branchCode, i))
		}
	}

	return
}

func GenRollNumbers() (rollNumbers []string) {

	for branchCode := range BranchCodesToNamesBefore19 {
	innermost2:
		for i := 0; i < 150; i++ {
			if branchCode == "45" || branchCode == "55" {
				if i > 100 {
					break innermost2
				}
				rollNumbers = append(rollNumbers, fmt.Sprintf("19%s%.2d", branchCode, i))
			} else {
				rollNumbers = append(rollNumbers, fmt.Sprintf("19%s%.3d", branchCode, i))
			}

		}
	}

	for year := 20; year <= 22; year++ {
		for branchCode := range BranchCodesToNames {
		innermost:
			for i := 0; i < 150; i++ {
				if branchCode == "dcs" || branchCode == "dec" && i > 100 {
					break innermost
				}
				rollNumbers = append(rollNumbers, fmt.Sprintf("%d%s%.3d", year, branchCode, i))
			}
		}
	}

	return
}

func GetUrlForRollNumber(rollNumber string) string {
	scheme := rollNumber[:2]
	return fmt.Sprintf("http://results.nith.ac.in/scheme%s/studentresult/result.asp", scheme)
}

func GetBatchAndBranch(rollNumber string) (batch string, branch string) {
	year, _ := strconv.Atoi(rollNumber[:2])
	batch = "20" + strconv.Itoa(year+4)

	if year <= 19 {
		branchName, found := BranchCodesToNamesBefore19[rollNumber[2:3]]
		if found {
			branch = branchName
		} else {
			branchName = BranchCodesToNamesBefore19[rollNumber[2:4]]
		}
	} else {
		branch = BranchCodesToNames[rollNumber[2:5]]
	}
	return
}

var BranchCodesToNames = map[string]string{
	"bce": "Civil",
	"bee": "Electrical",
	"bme": "Mechanical",
	"bec": "Electronics",
	"bcs": "Computer Science",
	"bar": "Architecture",
	"bch": "Chemical",
	"bms": "Material",
	"dcs": "Computer Science Dual",
	"dec": "Electronics Dual",
	"bph": "Physics",
}

var BranchCodesToNamesBefore19 = map[string]string{
	"1":  "Civil",
	"2":  "Electrical",
	"3":  "Mechanical",
	"4":  "Electronics",
	"5":  "Computer Science",
	"6":  "Architecture",
	"7":  "Chemical",
	"8":  "Material",
	"55": "Computer Science Dual",
	"45": "Electronics Dual",
}
