-- User table

CREATE TYPE user_role AS ENUM ('jobseeker', 'admin', 'recruiter');  -- Define your enum options here

CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  full_name VARCHAR(255) NOT NULL,
  email VARCHAR(255) UNIQUE NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  user_type user_role DEFAULT 'jobseeker' NOT NULL,
  avatar_url VARCHAR(255),
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Profile table

CREATE TABLE profiles (
  id          SERIAL PRIMARY KEY,
  user_id     INT NOT NULL UNIQUE,
  bio         TEXT,
  company     VARCHAR(255),
  job_role       VARCHAR(255) NOT NULL,
  description TEXT,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);


CREATE TABLE skills (
  id    SERIAL PRIMARY KEY,
  name  VARCHAR(255) UNIQUE NOT NULL
);

-- Intermediate table for User-Skill relationship

CREATE TABLE profile_skills (
  id SERIAL PRIMARY KEY,
  user_id     INT NOT NULL,
  skill_id    INT NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY (skill_id) REFERENCES skills(id) ON DELETE CASCADE
);

-- Career table 

CREATE TABLE careers (
  id SERIAL PRIMARY KEY,
  user_id     INT NOT NULL,
  title       VARCHAR(255) NOT NULL,
  company     VARCHAR(255),
  description TEXT,
  skill_id    INT NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY (skill_id) REFERENCES skills(id) ON DELETE CASCADE


);
