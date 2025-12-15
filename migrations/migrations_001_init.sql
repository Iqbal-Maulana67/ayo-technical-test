CREATE TABLE
  teams (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    logo_url TEXT,
    founded_year INT,
    home_address TEXT,
    home_city VARCHAR(100),
    created_at TIMESTAMP DEFAULT now (),
    updated_at TIMESTAMP DEFAULT now (),
    deleted_at TIMESTAMP NULL
  );

CREATE UNIQUE INDEX IF NOT EXISTS idx_unique_team_name ON teams (name)
WHERE
  deleted_at IS NULL;

CREATE TABLE
  players (
    id SERIAL PRIMARY KEY,
    team_id INT REFERENCES teams (id),
    name VARCHAR(255) NOT NULL,
    height REAL,
    weight REAL,
    position VARCHAR(255) NOT NULL,
    jersey_number INT NOT NULL,
    created_at TIMESTAMP DEFAULT now (),
    updated_at TIMESTAMP DEFAULT now (),
    deleted_at TIMESTAMP NULL,
    CONSTRAINT uq_team_jersey UNIQUE (team_id, jersey_number)
  );

CREATE TABLE
  matches (
    id SERIAL PRIMARY KEY,
    match_date TIMESTAMP NOT NULL,
    home_team_id INT REFERENCES teams (id),
    away_team_id INT REFERENCES teams (id),
    created_at TIMESTAMP DEFAULT now (),
    updated_at TIMESTAMP DEFAULT now (),
    deleted_at TIMESTAMP NULL
  );

CREATE TABLE
  match_results (
    id SERIAL PRIMARY KEY,
    match_id INT UNIQUE REFERENCES matches (id),
    final_score_home INT,
    final_score_away INT,
    created_at TIMESTAMP DEFAULT now (),
    updated_at TIMESTAMP DEFAULT now (),
    deleted_at TIMESTAMP NULL
  );

CREATE TABLE
  match_goals (
    id SERIAL PRIMARY KEY,
    match_id INT REFERENCES matches (id),
    player_id INT REFERENCES players (id),
    minute INT,
    created_at TIMESTAMP DEFAULT now (),
    updated_at TIMESTAMP DEFAULT now (),
    deleted_at TIME
  );

CREATE TABLE
  user_admins (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT now (),
    updated_at TIMESTAMP DEFAULT now (),
    deleted_at TIME
  );