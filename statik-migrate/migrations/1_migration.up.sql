CREATE TABLE jobs (
    job_id varchar(30) NOT NULL PRIMARY KEY,
    update_date integer NOT NULL,
    status varchar(10) NOT NULL
);

CREATE TABLE projects (
    project_id varchar(30) NOT NULL PRIMARY KEY,
    path varchar(30)  NOT NULL
);
