import {type ResultDataConfig, SortingType, type StudentResult, type StudentResultCompact} from "./types";
import _result_config from '../../../static/result_config.json';
import _json_data from '../../../static/ranks_result.json';
import _json_data_detailed from '../../../static/detailed_result.json';


const jsonData: StudentResultCompact[] = _json_data;
const jsonDataDetailed: StudentResult[] = _json_data_detailed;
export const resultConfig: ResultDataConfig = _result_config;

export function getFilteredAndSortedResults(
    nameOrRollNumberQuery: string,
    isSortOrderAscending: boolean,
    sortingType: SortingType,
    branch: string,
    batch: string,
    pageSize = 50,
    pageIndex = 0,
    minCGPI = 0,
    maxCGPI = 11
): StudentResultCompact[] {
  
    // Filter the data based on the provided arguments
    const filteredData = jsonData.filter((item) => {
      const rollNumberMatch = item.roll_number.toLowerCase().includes(nameOrRollNumberQuery.toLowerCase());
      const nameMatch = item.name.toLowerCase().includes(nameOrRollNumberQuery.toLowerCase());
  
      return rollNumberMatch || nameMatch;
    }).filter((item) => 
        (branch === '' || item.branch === branch) && (batch === '' || item.batch === batch)
    )
      .filter((item) => item.cgpi >= minCGPI && item.cgpi <= maxCGPI);
  
    // Sort the filtered data based on the provided sorting type
    const sortedData = filteredData.sort((a, b) => {

        if(sortingType == SortingType.name){
            const aValue = a.name;
            const bValue = b.name;
            if (aValue < bValue) {
                return isSortOrderAscending ? -1 : 1;
            } else if (aValue > bValue) {
                return isSortOrderAscending ? 1 : -1;
            } else {
                return 0;
            }
        } else if(sortingType == SortingType.rollNumber) {
            const aValue = a.roll_number;
            const bValue = b.roll_number;
            if (aValue < bValue) {
                return isSortOrderAscending ? -1 : 1;
            } else if (aValue > bValue) {
                return isSortOrderAscending ? 1 : -1;
            } else {
                return 0;
            }
        } else {
            if(a.cgpi < b.cgpi){
                return isSortOrderAscending ? -1 : 1;
            } else if(a.cgpi > b.cgpi){
                return  isSortOrderAscending ? 1 : -1;
            }
            const aValue = a.roll_number;
            const bValue = b.roll_number;
            if (aValue < bValue) {
                return isSortOrderAscending ? 1 : -1;
            } else if (aValue > bValue) {
                return isSortOrderAscending ? -1 : 1;
            } else {
                return 0;
            }
        }

    });
  
    // Paginate the sorted data based on the provided page size and index
    const startIndex = pageIndex * pageSize;
    const endIndex = startIndex + pageSize;
  
    return sortedData.slice(startIndex, endIndex);
  }


export function getStudent(roll_number: string): StudentResult | null {
  const student = jsonDataDetailed.find(item => item.roll_number === roll_number);
  if(typeof student === 'undefined')
      return null;
  return student;
}



// export function getStudentsList(
//     searchString: string,
//     isSortOrderAscending: boolean,
//     sortingType: SortingType,
//     branches: string[],
//     batches: string[],
//     pageSize = 50,
//     pageIndex = 0,
//     minCG = 0,
//     maxCG = 11
// ): StudentResultCompact[] {
//     let query = "SELECT * FROM student WHERE 1 = 1 "
//     if (searchString.length != 0) {
//         query += "AND (name like '%" + searchString + "%' or roll_number like '%" + searchString + "%') "
//     }
//     if (branches.length != 0) {
//         for (const branch in branches) {
//             query += "AND branch = " + branch + " "
//         }
//     }
//     if (batches.length != 0) {
//         for (const batch in batches) {
//             query += "AND batch = " + batch + " "
//         }
//     }
//     query += "AND cgpi >= " + minCG.toFixed(2) + " "
//     query += "AND cgpi <= " + maxCG.toFixed(2) + " "
//     query += "ORDER BY " + sortingType + " " + ((isSortOrderAscending) ? "asc" : "desc") + " "
//     query += " LIMIT " + pageSize.toFixed(0) + " OFFSET " + pageIndex.toFixed(0)
    
//     const stmnt = db.prepare(query);
//     const rows = stmnt.all({ pageSize });
//     return rows as StudentResultCompact[];
// }

// export function getStudent(roll_number: string): StudentResult {
//     const studentDetails = db.prepare(`SELECT * FROM student where roll_number = ${roll_number}`).all() as StudentResultCompact[];
//     const semesterDetails = db.prepare(`SELECT * FROM semester_result_data where student_roll_number = '${roll_number}'`).all() as SemesterResultDBType[];
//     const subjectDetails = db.prepare(`
//         SELECT res.semester, sbj.name, sbj.code, res.grade, sbj.credits, res.sub_gp FROM 
//         subject_result_data as res JOIN subject sbj ON res.subject_code = sbj.code
//         where student_roll_number = ${roll_number}
//     `).all() as SubjectResultDBType[];

//     return {
//         roll_number: roll_number,
//         name: studentDetails[0].name,
//         fathers_name: studentDetails[0].fathers_name,
//         semester_results: semesterDetails.map(sem => {
//             return {
//                 semesterNumber: sem.semester,
//                 subjectResults: subjectDetails.filter(sbj => sbj.semester == sem.semester).map(item => {
//                     return {
//                         SubjectName: item.name,
//                         SubjectCode: item.code,
//                         SubPoint: item.credits,
//                         Grade: item.grade,
//                         SubGP: item.sub_gp
//                     }
//                 }),
//                 SGPI: sem.sgpi,
//                 CGPI: sem.cgpi   
//             }
//         }).reverse(),
//         CGPI: '',
//         branch: studentDetails[0].branch,
//         batch: studentDetails[0].batch
//     }
// }

// // function binarySearchRollNumber(rollNumber: string, start = 0, end = jsonDataDetailed.length - 1): number {
// //     const mid = Math.floor((start + end) / 2);
  
// //     if (rollNumber === jsonDataDetailed[mid].rollNumber) {
// //       return mid;
// //     }
  
// //     if (start >= end) {
// //       return -1;
// //     }
  
// //     return rollNumber < jsonDataDetailed[mid].rollNumber
// //       ? binarySearchRollNumber(rollNumber, start, mid - 1)
// //       : binarySearchRollNumber(rollNumber, mid + 1, end);
// //   }
  
