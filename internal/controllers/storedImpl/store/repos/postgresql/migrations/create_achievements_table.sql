CREATE TABLE achievements(
    id serial NOT NULL,
    name text NOT NULL,
    placement text NOT NULL,
    playerid int REFERENCES players (id),
    teamid int REFERENCES teams (id),
    CHECK (playerid IS NULL OR teamid IS NULL),
    CHECK (playerid IS NOT NULL OR teamid IS NOT NULL)
);