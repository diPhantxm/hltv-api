CREATE TABLE stats(
    id int PRIMARY KEY,
    mapsplayed int NOT NULL,
    rating real NOT NULL,
    killsperround real NOT NULL,
    headshots real NOT NULL,
    deathsperround real NOT NULL,
    roundscontributed real NOT NULL
);