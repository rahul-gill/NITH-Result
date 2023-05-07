package main

import (
	resultNITH "Result-NITH"
	"Result-NITH/db"
	"context"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	students := []resultNITH.StudentHtmlParsed{
		{
			RollNumber:  "192087",
			Name:        "Rahul",
			FathersName: "Rahul",
			SemesterResults: []resultNITH.SemesterResult{
				{
					SemesterNumber: 0,
					SubjectResults: []resultNITH.SubjectResult{
						{
							SubjectName: "",
							SubjectCode: "",
							SubPoint:    0,
							Grade:       "",
							SubGP:       0,
						},
					},
					SGPI: 0,
					CGPI: 0,
				},
			},
			CGPI: 0,
		},
	}
	println("\n\nFinished fetching students\n")

	db, queries := resultNITH.GetDbQueriesForNewDb("sample.db")

	err := StoreStudentInDbX(db, queries, students)
	if err != nil {
		println("Error in storeTheStudentInDb")
		println(err)
	}
}

func StoreStudentInDbX(sqlDB *sql.DB, queries *db.Queries, students []resultNITH.StudentHtmlParsed) error {
	ctx := context.Background()

	tx, err := sqlDB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := queries.WithTx(tx)
	for studentIndex, student := range students {
		//student table
		batch, branch := resultNITH.GetBatchAndBranch(student.RollNumber)
		_, err := qtx.CreateStudent(ctx, db.CreateStudentParams{
			RollNumber:     student.RollNumber,
			Name:           student.Name,
			FathersName:    student.FathersName,
			Batch:          batch,
			Branch:         branch,
			LatestSemester: int64(len(student.SemesterResults)),
			Cgpi:           student.CGPI,
		})
		if err != nil {
			println("Error when creating student entry for rollNumber ", student.RollNumber, ": ", err)
		}
		//subject table
		for _, sem := range student.SemesterResults {
			for _, subject := range sem.SubjectResults {
				_, err := qtx.CreateSubject(ctx, db.CreateSubjectParams{
					Code:    subject.SubjectCode,
					Name:    subject.SubjectName,
					Credits: subject.SubPoint,
				})
				if err != nil {
					println("Error when creating subject entry for rollNumber ", student.RollNumber, "and subject ", subject.SubjectName, ": ", err)
				}
			}
		}
		//semester result
		for semNumber, sem := range student.SemesterResults {
			_, err := qtx.CreateSemesterResultData(ctx, db.CreateSemesterResultDataParams{
				StudentRollNumber: student.RollNumber,
				Semester:          int64(semNumber + 1),
				Cgpi:              sem.CGPI,
				Sgpi:              sem.SGPI,
			})
			if err != nil {
				println("Error when creating semester entry for rollNumber ", student.RollNumber, "and semNumber ", semNumber, ": ", err)
			}
		}
		//subject result data
		for semNumber, sem := range student.SemesterResults {
			for _, subject := range sem.SubjectResults {
				_, err := qtx.CreateSubjectResultData(ctx, db.CreateSubjectResultDataParams{
					StudentRollNumber: student.RollNumber,
					SubjectCode:       subject.SubjectCode,
					Grade:             subject.Grade,
					SubGp:             subject.SubGP,
					Semester:          int64(semNumber + 1),
				})
				if err != nil {
					println("Error when creating subject result entry for rollNumber ", student.RollNumber, "and subject ", subject.SubjectName, ": ", err.Error())
				}
			}
		}
		println("Finished student Number ", studentIndex, "out of ", len(students))
	}
	return tx.Commit()
}
