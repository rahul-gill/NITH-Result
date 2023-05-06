import { SortingType } from "./types";

// export const branchOptionsUI = [];
// export const batchOptionsUI = [] 
// export const sortOptionsUI = ['Sort by CGPI', 'Sort by Roll Number', 'Sort by Name']


// export const branchOptionsDB = ['All branches', 'Computer Science', 'Electrical Engineering', 'Mechanical Engineering', 'Civil Engineering'];
// export const batchOptionsDB = ['All batches', '2020 Batch', '2021  Batch', '2022 Batch'] 
// export const sortOptionsDB = ['Sort by CGPI', 'Sort by Roll Number', 'Sort by Name']

// export function sortOptionToDb(sortOptionUI: string): SortingType{
//     let sorting = SortingType.cgpi;
//     switch(sortOptionUI!.toString()){
//         case "Sort by CGPI": sorting = SortingType.cgpi; break;
//         case "Sort by Roll Number": sorting = SortingType.rollNumber; break;
//         case "Sort by Name": sorting = SortingType.name; break;
//     }
//     return sorting;
// }

// export function branchOptionsToDB(branchOptionUI: string){
//     return branchOptionsDB[branchOptionsUI.indexOf(branchOptionUI)]
// }

// export function batchOptionsToDB(batchOptionUI: string){
//     return batchOptionUI[batchOptionsDB.indexOf(batchOptionUI)]
// }