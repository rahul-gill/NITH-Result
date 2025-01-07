CREATE TABLE student
(
    roll_number     text primary key,
    name            text not null ,
    fathers_name    text not null ,
    batch           text not null ,
    branch          text not null ,
    latest_semester integer not null ,
    cgpi            real not null
);

CREATE TABLE subject_result_data
(
    student_roll_number text not null ,
    semester            integer not null ,
    subject_code        text not null ,
    grade               text not null ,
    sub_gp              int not null ,
  	subject_name        text not null ,
  	subject_credits     int not null ,
    primary key (student_roll_number, subject_code),
    foreign key (student_roll_number) references student (roll_number)
);

CREATE TABLE semester_result_data
(
    student_roll_number text not null ,
    semester            integer not null ,
    cgpi                real not null ,
    sgpi                real not null ,
    primary key (student_roll_number, semester),
    foreign key (student_roll_number) references student (roll_number)
);
