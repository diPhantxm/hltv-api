CREATE TABLE maps(
    id serial PRIMARY KEY,
    matchid int NOT NULL,
    name text NOT NULL,
    teamascore int NOT NULL,
    teambscore int NOT NULL,
    FOREIGN KEY (matchid) REFERENCES matches (id)
);