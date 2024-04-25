package main

import (
	resultNITH "Result-NITH"
	"Result-NITH/db"
	"context"
	"encoding/json"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"sort"
	"time"
)

func getRanksData(queries *db.Queries) []resultNITH.StudentResultWithRanks {
	ctx := context.Background()
	studentDetails, err := queries.GetAllStudent(ctx)
	if err != nil {
		log.Fatal(err)
	}

	studentsOut := make([]resultNITH.StudentResultWithRanks, len(studentDetails))

	var branchBasedRanksMap = map[string]int{}
	var batchBasedRanksMap = map[string]int{}
	var batchBranchBasedRanksMap = map[string]int{}

	for _, eachBranch := range resultNITH.BranchCodesToNames {
		branchRanks, _ := queries.GetRanksData(ctx, db.GetRanksDataParams{
			Batch:  "%",
			Branch: eachBranch,
		})
		for i, item := range branchRanks {
			branchBasedRanksMap[item.RollNumber] = i + 1
		}
	}
	for _, eachBatch := range []string{"2022", "2023", "2024", "2025", "2026", "2027"} {
		batchRanks, _ := queries.GetRanksData(ctx, db.GetRanksDataParams{
			Batch:  eachBatch,
			Branch: "%",
		})
		for i, item := range batchRanks {
			batchBasedRanksMap[item.RollNumber] = i + 1
		}
	}
	for _, eachBranch := range resultNITH.BranchCodesToNames {
		for _, eachBatch := range []string{"2022", "2023", "2024", "2025", "2026", "2027"} {
			classRanks, _ := queries.GetRanksData(ctx, db.GetRanksDataParams{
				Batch:  eachBatch,
				Branch: eachBranch,
			})
			for i, item := range classRanks {
				batchBranchBasedRanksMap[item.RollNumber] = i + 1
			}
		}
	}

	for index, student := range studentDetails {

		studentsOut[index] = resultNITH.StudentResultWithRanks{
			RollNumber:  student.RollNumber,
			Name:        student.Name,
			FathersName: student.FathersName,
			CGPI:        student.Cgpi,
			Branch:      student.Branch,
			Batch:       student.Batch,
			BranchRank:  int64(branchBasedRanksMap[student.RollNumber]),
			YearRank:    int64(batchBasedRanksMap[student.RollNumber]),
			ClassRank:   int64(batchBranchBasedRanksMap[student.RollNumber]),
		}
	}
	sort.Slice(studentsOut, func(i, j int) bool {
		return studentsOut[i].CGPI > studentsOut[j].CGPI
	})
	return studentsOut
}

func getDetailedResults(queries *db.Queries) []resultNITH.StudentResult {
	ctx := context.Background()
	studentDetails, err := queries.GetAllStudent(ctx)
	if err != nil {
		log.Fatal(err)
	}

	studentsOut := make([]resultNITH.StudentResult, len(studentDetails))
	for index, student := range studentDetails {
		semesterDetails, err := queries.GetStudentSemestersResult(ctx, student.RollNumber)
		if err != nil {
			log.Fatal(err)
		}
		subjectDetails, err := queries.GetStudentSubjectsResultAll(ctx, student.RollNumber)
		if err != nil {
			log.Fatal(err)
		}
		studentsOut[index] = resultNITH.StudentResult{
			RollNumber:  student.RollNumber,
			Name:        student.Name,
			FathersName: student.FathersName,
			SemesterResults: func() []resultNITH.SemesterResult {
				var results []resultNITH.SemesterResult
				for i := len(semesterDetails) - 1; i >= 0; i-- {
					sem := semesterDetails[i]
					results = append(results, resultNITH.SemesterResult{
						SemesterNumber: sem.Semester,
						SubjectResults: func() []resultNITH.SubjectResult {
							var subjects []resultNITH.SubjectResult
							for _, item := range subjectDetails {
								if item.Semester == sem.Semester {
									subjects = append(subjects, resultNITH.SubjectResult{
										SubjectName: item.SubjectName,
										SubjectCode: item.SubjectCode,
										SubPoint:    item.Credits,
										Grade:       item.Grade,
										SubGP:       item.SubGp,
									})
								}
							}
							return subjects
						}(),
						SGPI: sem.Sgpi,
						CGPI: sem.Cgpi,
					})
				}
				return results
			}(),
			CGPI:   student.Cgpi,
			Branch: student.Branch,
			Batch:  student.Batch,
		}
	}
	return studentsOut
}

func main() {
	// Detailed result data
	_, queries := resultNITH.GetDbQueries()
	data := getDetailedResults(queries)
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile("detailed_result.json", jsonBytes, 0644)
	if err != nil {
		log.Fatal(err)
	}

	// Ranks result data
	data2 := getRanksData(queries)
	jsonBytes2, err := json.Marshal(data2)
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile("ranks_result.json", jsonBytes2, 0644)
	if err != nil {
		log.Fatal(err)
	}

	// Other useful data
	branchSet := make(map[string]bool)
	for _, result := range data2 {
		branchSet[result.Branch] = true
	}

	branches := make([]string, 0, len(branchSet))
	branches = append(branches, "All branches")
	for branch := range branchSet {
		branches = append(branches, branch)
	}
	sort.Strings(branches[1:])

	batchSet := make(map[string]bool)
	for _, result := range data2 {
		batchSet[result.Batch] = true
	}
	batches := make([]string, 0, len(batchSet))
	batches = append(batches, "All batches")
	for batch := range batchSet {
		batches = append(batches, batch)
	}
	sort.Strings(batches[1:])
	jsonBytes3, err := json.Marshal(map[string]interface{}{
		"last_update_date":   time.Now().UTC().Format("02 Jan, 2006"),
		"available_branches": branches,
		"available_batches":  batches,
	})
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile("result_config.json", jsonBytes3, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
