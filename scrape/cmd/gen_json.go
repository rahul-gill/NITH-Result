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
	for index, student := range studentDetails {
		classRank, _ := queries.GetStudentClassRank(ctx, db.GetStudentClassRankParams{
			Batch:    student.Batch,
			Branch:   student.Branch,
			Semester: student.LatestSemester,
			Cgpi:     student.Cgpi,
		})
		yearRank, _ := queries.GetStudentYearRank(ctx, db.GetStudentYearRankParams{
			Batch:    student.Batch,
			Semester: student.LatestSemester,
			Cgpi:     student.Cgpi,
		})
		branchRank, _ := queries.GetStudentBranchRank(ctx, db.GetStudentBranchRankParams{
			Branch:   student.Branch,
			Semester: student.LatestSemester,
			Cgpi:     student.Cgpi,
		})

		studentsOut[index] = resultNITH.StudentResultWithRanks{
			RollNumber:  student.RollNumber,
			Name:        student.Name,
			FathersName: student.FathersName,
			CGPI:        student.Cgpi,
			Branch:      student.Branch,
			Batch:       student.Batch,
			BranchRank:  branchRank,
			YearRank:    yearRank,
			ClassRank:   classRank,
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
	branches = append(branches, "All")
	for branch := range branchSet {
		branches = append(branches, branch)
	}

	batchSet := make(map[string]bool)
	for _, result := range data2 {
		batchSet[result.Batch] = true
	}
	batches := make([]string, 0, len(batchSet))
	batches = append(batches, "All")
	for batch := range batchSet {
		batches = append(batches, batch)
	}
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
