CREATE TABLE events(
    id int NOT NULL PRIMARY KEY,
    name text NOT NULL,
    startdate bigint NOT NULL,
    enddate bigint NOT NULL,
    prizepool text NOT NULL,
    teams text NOT NULL,
    location text NOT NULL
);