CREATE TABLE roasters(
    id serial PRIMARY KEY,
    nickname text NOT NULL,
    status text NOT NULL,
    timeonteam text NOT NULL,
    mapsplayed int NOT NULL,
    rating real NOT NULL,
    teamid int NOT NULL,
    FOREIGN KEY (teamid) REFERENCES teams (id)
);