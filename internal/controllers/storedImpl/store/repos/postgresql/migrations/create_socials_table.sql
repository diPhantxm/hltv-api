CREATE TABLE socials(
    id serial NOT NULL,
    name text NOT NULL,
    link text NOT NULL,
    playerid int,
    teamid int,
    FOREIGN KEY (playerid) REFERENCES players (id),
    FOREIGN KEY (teamid) REFERENCES teams (id)
);