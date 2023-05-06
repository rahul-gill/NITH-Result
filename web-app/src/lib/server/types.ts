export type StudentResult = {
    roll_number: string;
    name: string;
    fathers_name: string;
    semester_results: SemesterResult[];
    cgpi: string;
    branch: string;
    batch: string;
};

export type StudentResultCompact = {
    roll_number: string
    name: string;
    fathers_name: string;
    cgpi: number;
    branch: string;
    batch: string;
    branch_rank: number;
    year_rank: number;
    class_rank: number
};

export type SemesterResult = {
    semester_number: number;
    subject_results: SubjectResult[];
    sgpi: number;
    cgpi: number;
};

export type SubjectResult = {
    subject_name: string;
    subject_code: string;
    sub_point: number;
    grade: string;
    sub_gp: number;
};


export enum SortingType {
    name = "name",
    cgpi = "cgpi",
    rollNumber = "roll_number"
}


export type SemesterResultDBType = {
    student_roll_number: string;
    semester: number;
    cgpi: number;
    sgpi: number;
};

export type SubjectResultDBType = {
    semester: number;
    name: string;
    code: string;
    grade: string;
    credits: number;
    sub_gp: number;
};


export type ResultDataConfig = {
    last_update_date: string;
    available_branches: string[];
    available_batches: string[];
}