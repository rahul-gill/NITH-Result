-- name: CreateStudent :one
INSERT OR REPLACE INTO student(roll_number, name, fathers_name, batch, branch, latest_semester, cgpi)
VALUES (?, ?, ?, ?, ?, ?, ?) RETURNING *;

-- name: CreateSubjectResultData :one
INSERT OR REPLACE INTO subject_result_data(student_roll_number, subject_code, grade, sub_gp, semester)
VALUES (?, ?, ?, ?, ?) RETURNING *;


-- name: CreateSemesterResultData :one
INSERT OR REPLACE INTO semester_result_data(student_roll_number, semester, cgpi, sgpi)
VALUES (?, ?, ?, ?) RETURNING *;


-- name: CreateSubject :one
INSERT OR REPLACE INTO subject(code, name, credits)
VALUES (?, ?, ?) RETURNING *;

-- name: GetStudent :one
SELECT * FROM student where roll_number = ? LIMIT 1;


-- name: GetAllStudent :many
SELECT * FROM student;


-- name: GetStudentSemestersResult :many
SELECT * FROM semester_result_data where student_roll_number = ?;



-- name: GetStudentSubjectsResult :many
SELECT * FROM subject_result_data where student_roll_number = ? and semester = ?;

-- name: GetStudentSubjectsResultAll :many
SELECT res.semester as semester, sbj.name as subject_name, sbj.code as subject_code, res.grade as grade, sbj.credits as credits, res.sub_gp as sub_gp
FROM subject_result_data as res JOIN subject sbj ON res.subject_code = sbj.code
where student_roll_number = ?;

-- name: GetStudentCGPI :one
SELECT cgpi FROM semester_result_data WHERE student_roll_number = ? AND semester = ?;

-- name: GetStudentClassRank :one
SELECT COUNT(*) + 1
FROM semester_result_data AS srd INNER JOIN student
ON srd.student_roll_number = student.roll_number
WHERE student.batch = ? AND student.branch = ? AND srd.semester = ? AND srd.cgpi > ?;

-- name: GetStudentBranchRank :one
SELECT COUNT(*) + 1
FROM semester_result_data AS srd INNER JOIN student
ON srd.student_roll_number = student.roll_number
WHERE student.branch = ? AND srd.semester = ? AND srd.cgpi > ?;

-- name: GetStudentYearRank :one
SELECT COUNT(*) + 1
FROM semester_result_data AS srd INNER JOIN student
ON srd.student_roll_number = student.roll_number
WHERE student.batch = ? AND srd.semester = ? AND srd.cgpi > ?