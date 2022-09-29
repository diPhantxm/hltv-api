CREATE TABLE players(
    id int NOT NULL PRIMARY KEY,
    age int NOT NULL,
    nickname text NOT NULL,
    team text NOT NULL,
    firstname text NOT NULL,
    lastname text NOT NULL,
    country text NOT NULL
);

CREATE TABLE stats(
    id int PRIMARY KEY,
    mapsplayed int NOT NULL,
    rating real NOT NULL,
    killsperround real NOT NULL,
    headshots real NOT NULL,
    deathsperround real NOT NULL,
    roundscontributed real NOT NULL,
    FOREIGN KEY (id) REFERENCES players (id)
);

CREATE TABLE events(
    id int NOT NULL PRIMARY KEY,
    name text NOT NULL,
    startdate bigint NOT NULL,
    enddate bigint NOT NULL,
    prizepool text NOT NULL,
    teams text NOT NULL,
    location text NOT NULL
);

CREATE TABLE teams(
    id int NOT NULL PRIMARY KEY,
    ranking int NOT NULL,
    weeksintop30 int NOT NULL,
    name text NOT NULL,
    country text NOT NULL,
    averageage real NOT NULL
);

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

CREATE TABLE socials(
    id serial NOT NULL,
    name text NOT NULL,
    link text NOT NULL,
    playerid int REFERENCES players (id),
    teamid int REFERENCES teams (id),
    CHECK (playerid IS NULL OR teamid IS NULL),
    CHECK (playerid IS NOT NULL OR teamid IS NOT NULL)
);

CREATE TABLE achievements(
    id serial NOT NULL,
    name text NOT NULL,
    placement text NOT NULL,
    playerid int REFERENCES players (id),
    teamid int REFERENCES teams (id),
    CHECK (playerid IS NULL OR teamid IS NULL),
    CHECK (playerid IS NOT NULL OR teamid IS NOT NULL)
);

CREATE TABLE matches(
    id int NOT NULL PRIMARY KEY,
    teama text NOT NULL,
    teamb text NOT NULL,
    starttime bigint NOT NULL,
    viewers int NOT NULL,
    playerofthematch text NOT NULL,
    isover boolean NOT NULL
);

CREATE TABLE maps(
    id serial PRIMARY KEY,
    matchid int NOT NULL,
    name text NOT NULL,
    teamascore int NOT NULL,
    teambscore int NOT NULL,
    FOREIGN KEY (matchid) REFERENCES matches (id)
);
