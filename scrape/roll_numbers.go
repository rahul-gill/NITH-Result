package Result_NITH

import (
	"fmt"
	"strconv"
	"strings"
)

func GenRollNumbers(forOnlyBatch *int) (rollNumbers []string) {
	for branchCode := range BranchCodesToNamesBefore19 {
	innermost2:
		for i := 0; i < 150; i++ {
			if branchCode == "45" || branchCode == "55" {
				if i > 100 {
					break innermost2
				}
				if forOnlyBatch == nil || *forOnlyBatch == 18 {
					rollNumbers = append(rollNumbers, fmt.Sprintf("18%s%.2d", branchCode, i))
				}
				if forOnlyBatch == nil || *forOnlyBatch == 19 {
					rollNumbers = append(rollNumbers, fmt.Sprintf("19%s%.2d", branchCode, i))
				}
			}
			if forOnlyBatch == nil || *forOnlyBatch == 18 {
				rollNumbers = append(rollNumbers, fmt.Sprintf("18%s%.3d", branchCode, i))
			}
			if forOnlyBatch == nil || *forOnlyBatch == 19 {
				rollNumbers = append(rollNumbers, fmt.Sprintf("19%s%.3d", branchCode, i))
			}
		}
	}

	for year := 20; year <= 24; year++ {
		if forOnlyBatch != nil && *forOnlyBatch != year {
			continue
		}
		for branchCode := range BranchCodesToNames {
		innermost:
			for i := 0; i < 150; i++ {
				if (branchCode == "dcs" || branchCode == "dec") && i > 100 {
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
	rollNumber = strings.ToLower(rollNumber)
	year, _ := strconv.Atoi(rollNumber[:2])
	batch = "20" + strconv.Itoa(year+4)

	if exceptionBranch, foundException := BranchExceptionRollNumbers[rollNumber]; foundException {
		branch = exceptionBranch
		return
	}

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
	"bph": "Engineering Physics",
	"dec": "Electronics Dual",
	"dcs": "Computer Science Dual",
	"bma": "Maths & Computing",
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

// BranchExceptionRollNumbers : Put roll numbers in lower case
var BranchExceptionRollNumbers = map[string]string{
	"20bch008": "Mechanical",
	"20bms023": "Mechanical",
	"20bee039": "Electronics Dual",
	"21bce032": "Electrical",
	"21bce002": "Electrical",
	"21bch009": "Electrical",
	"21bch039": "Electrical",
	"21bme040": "Electrical",
}
