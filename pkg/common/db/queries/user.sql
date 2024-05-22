-- name: GetAllUsers :many
SELECT * FROM users ;


-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: CreateUser :one
INSERT INTO users (full_name, email, password_hash, user_type, avatar_url)
VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET full_name = $2,
  password_hash = $3,
  avatar_url = $4
WHERE id = $1
RETURNING *;

-- name: DeleteUserByID :exec
DELETE FROM users WHERE id = $1;

-- name: GetUsersByUserType :many
SELECT * FROM users WHERE user_type = $1;

-- name: GetUsersByEmail :many
SELECT * FROM users WHERE email = $1;

-- name: GetUsersByEmailAndUserType :many
SELECT * FROM users WHERE email = $1 AND user_type = $2;


-------- profile --------------
-- name: GetAllProfile :many
SELECT * FROM profiles ;


-- name: GetProfileByUserID :one
SELECT * FROM profiles WHERE user_id = $1;

-- name: GetProfileByID :one
SELECT * FROM profiles WHERE id = $1;

-- name: CreateProfile :one
INSERT INTO profiles (user_id, bio, company, job_role, description)
VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: UpdateProfile :one
UPDATE profiles
SET bio = $2,
  company = $3,
  job_role = $4,
  description = $5
WHERE id = $1
RETURNING *;

-- name: DeleteProfileByUserID :exec
DELETE FROM profiles WHERE user_id = $1;

-- name: DeleteProfileByID :exec
DELETE FROM profiles WHERE id = $1;


-- --------- skill --------------------

-- name: GetAllskills :many
SELECT * FROM skills ;

-- name: GetSkillByID :one
SELECT * FROM skills WHERE id = $1;


-- name: GetSkillByName :one
SELECT * FROM skills WHERE name ILIKE $1;

-- name: CreateSkill :one
INSERT INTO skills (name) VALUES (LOWER($1)) RETURNING *;

-- name: UpdateSkill :one
UPDATE skills
SET name = LOWER($2)
WHERE id = $1
RETURNING *;

-- name: DeleteSkillByID :exec
DELETE FROM skills WHERE id = $1;

----------- user skills ----------

-- name: AddSkillToUser :exec
INSERT INTO profile_skills (user_id, skill_id) VALUES ($1, $2);

-- name: GetAllUserSkillsByUserID :many
SELECT * FROM profile_skills WHERE user_id = $1;


-- name: GetUserSkillsByUserIDAndSkillId :one
SELECT * FROM profile_skills WHERE user_id = $1 AND skill_id= $2;


-- name: GetAllUserBySkillsName :many
SELECT u.*
FROM users u
INNER JOIN profile_skills ps ON u.id = ps.user_id
INNER JOIN skills s ON ps.skill_id = s.id
WHERE s.name ILIKE $1;

------- career ------------------------
-- name: GetAllCareers :many
SELECT * FROM careers;

-- name: GetCareersByUserID :many
SELECT * FROM careers WHERE user_id = $1;

-- name: GetCareerByID :one
SELECT * FROM careers WHERE id = $1;

-- name: CreateCareer :one
INSERT INTO careers (user_id, title, company, description, skill_id)
VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: UpdateCareer :one
UPDATE careers
SET title = $2,
  company = $3,
  description = $4,
  skill_id = $5
WHERE id = $1
RETURNING *;

-- name: DeleteCareerByID :exec
DELETE FROM careers WHERE id = $1;

