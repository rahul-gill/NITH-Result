package Result_NITH

import "fmt"

func GenRollNumbers() (rollNumbers []string) {
	for i := 0; i < 150; i++ {
		//2023
		rollNumbers = append(rollNumbers, fmt.Sprintf("191%.3d", i), fmt.Sprintf("192%.3d", i), fmt.Sprintf("193%.3d", i), fmt.Sprintf("194%.3d", i), fmt.Sprintf("195%.3d", i), fmt.Sprintf("196%.3d", i), fmt.Sprintf("197%.3d", i), fmt.Sprintf("198%.3d", i))
		if i < 100 {
			rollNumbers = append(rollNumbers, fmt.Sprintf("1955%.2d", i), fmt.Sprintf("1945%.2d", i))
		}
		//2022
		rollNumbers = append(rollNumbers, fmt.Sprintf("181%.3d", i), fmt.Sprintf("182%.3d", i), fmt.Sprintf("183%.3d", i), fmt.Sprintf("184%.3d", i), fmt.Sprintf("185%.3d", i), fmt.Sprintf("186%.3d", i), fmt.Sprintf("187%.3d", i), fmt.Sprintf("188%.3d", i))
		if i < 100 {
			rollNumbers = append(rollNumbers, fmt.Sprintf("1855%.2d", i), fmt.Sprintf("1845%.2d", i))
		}
		//TODO: i don't understand 2024+ scheme of roll numbers
	}
	return
}

func GetUrlForRollNumber(rollNumber string) string {
	scheme := rollNumber[:2]
	return fmt.Sprintf("http://14.139.56.19/scheme%s/studentresult/result.asp", scheme)
}

func GetBatchAndBranch(rollNumber string) (batch string, branch string) {
	yearStr := rollNumber[:2]
	branchCode := rollNumber[2:3]

	switch yearStr {
	case "18":
		batch = "2022"
	case "19":
		batch = "2023"
	case "20":
		batch = "2024"
	case "21":
		batch = "2025"
	default:
		batch = "Unset"
	}

	switch branchCode {
	case "1":
		branch = "Civil"
	case "2":
		branch = "Electrical"
	case "3":
		branch = "Mechanical"
	case "4":
		branch = "Electronics"
	case "5":
		branch = "ComputerScience"
	case "6":
		branch = "Arch"
	case "7":
		branch = "Chemical"
	case "8":
		branch = "Chemical"
	default:
		branch = "Unset"
	}
	if rollNumber[3] == '5' {
		if rollNumber[2] == '4' {
			branch = "ElectronicsDual"
		}
		if rollNumber[2] == '6' {
			branch = "ComputerScienceDual"
		}
	}

	return batch, branch
}
