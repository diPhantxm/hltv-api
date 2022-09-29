CREATE TABLE teams(
    id int NOT NULL PRIMARY KEY,
    ranking int NOT NULL,
    weeksintop30 int NOT NULL,
    name text NOT NULL,
    country text NOT NULL,
    averageage real NOT NULL
);