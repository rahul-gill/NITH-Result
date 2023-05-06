package Result_NITH

type StudentResult struct {
	RollNumber      string           `json:"roll_number"`
	Name            string           `json:"name"`
	FathersName     string           `json:"fathers_name"`
	SemesterResults []SemesterResult `json:"semester_results"`
	CGPI            float64          `json:"cgpi"`
	Branch          string           `json:"branch"`
	Batch           string           `json:"batch"`
}

type StudentResultWithRanks struct {
	RollNumber  string  `json:"roll_number"`
	Name        string  `json:"name"`
	FathersName string  `json:"fathers_name"`
	CGPI        float64 `json:"cgpi"`
	Branch      string  `json:"branch"`
	Batch       string  `json:"batch"`
	BranchRank  int64   `json:"branch_rank"`
	YearRank    int64   `json:"year_rank"`
	ClassRank   int64   `json:"class_rank"`
}

type SemesterResult struct {
	SemesterNumber int64           `json:"semester_number"`
	SubjectResults []SubjectResult `json:"subject_results"`
	SGPI           float64         `json:"sgpi"`
	CGPI           float64         `json:"cgpi"`
}

type SubjectResult struct {
	SubjectName string `json:"subject_name"`
	SubjectCode string `json:"subject_code"`
	SubPoint    int64  `json:"sub_point"`
	Grade       string `json:"grade"`
	SubGP       int64  `json:"sub_gp"`
}

type StudentHtmlParsed struct {
	RollNumber      string           `json:"roll_number"`
	Name            string           `json:"name"`
	FathersName     string           `json:"fathers_name"`
	SemesterResults []SemesterResult `json:"semester_results"`
	CGPI            float64          `json:"cgpi"`
}
