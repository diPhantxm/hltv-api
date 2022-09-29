CREATE TABLE matches(
    id int NOT NULL PRIMARY KEY,
    teama text NOT NULL,
    teamb text NOT NULL,
    starttime bigint NOT NULL,
    viewers int NOT NULL,
    playerofthematch text NOT NULL,
    isover boolean NOT NULL
);