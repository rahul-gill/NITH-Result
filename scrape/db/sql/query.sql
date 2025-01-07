-- name: CreateStudent :one
INSERT OR REPLACE INTO student(roll_number, name, fathers_name, batch, branch, latest_semester, cgpi)
VALUES (?, ?, ?, ?, ?, ?, ?) RETURNING *;

-- name: CreateSubjectResultData :one
INSERT OR REPLACE INTO subject_result_data(student_roll_number, subject_code, grade, sub_gp, semester, subject_name, subject_credits)
VALUES (?, ?, ?, ?, ?, ?, ?) RETURNING *;


-- name: CreateSemesterResultData :one
INSERT OR REPLACE INTO semester_result_data(student_roll_number, semester, cgpi, sgpi)
VALUES (?, ?, ?, ?) RETURNING *;



-- name: GetStudent :one
SELECT * FROM student where roll_number = ? LIMIT 1;

-- name: UpdateStudentBranch :exec
UPDATE student set branch = ? where roll_number = ?;


-- name: GetAllStudent :many
SELECT * FROM student;


-- name: GetStudentSemestersResult :many
SELECT * FROM semester_result_data where student_roll_number = ?;



-- name: GetStudentSubjectsResult :many
SELECT * FROM subject_result_data where student_roll_number = ? and semester = ?;

-- name: GetStudentSubjectsResultAll :many
SELECT res.semester as semester, res.subject_name as subject_name, res.subject_code as subject_code, res.grade as grade, subject_credits as credits, res.sub_gp as sub_gp
FROM subject_result_data as res
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
WHERE student.batch = ? AND srd.semester = ? AND srd.cgpi > ?;

-- name: GetRanksData :many
SELECT
    student.roll_number,
    semester_result_data.cgpi AS curr_cg
FROM
    student,
    semester_result_data,
    (SELECT
         inner_res.student_roll_number,
        MAX(inner_res.semester) AS latest_semester
    FROM
        semester_result_data inner_res
    GROUP BY
         inner_res.student_roll_number
    ) AS LatestSemester
WHERE
    student.roll_number = semester_result_data.student_roll_number
    AND semester_result_data.student_roll_number = LatestSemester.student_roll_number
    AND semester_result_data.semester = LatestSemester.latest_semester
    AND student.batch like ? AND student.branch like ?
ORDER BY semester_result_data.cgpi DESC, student.roll_number;
