package activity

const (
	queryCreateActivity = `
	INSERT INTO activities (title, email, updated_at) VALUES (?, ?, ?)
	`

	queryGetAllActivity = `
	SELECT
		activity_id as id,
		title,
		email,
		updated_at,
		created_at
	FROM activities
	`

	queryGetOneActivity = `
	SELECT
		activity_id as id,
		title,
		email,
		updated_at,
		created_at
	FROM activities
	WHERE activity_id = ?
	`

	queryUpdateActivity = `
	UPDATE activities
	SET
		title = ?,
		updated_at = ?
	WHERE activity_id = ?
	`

	queryDeleteActivity = `
	DELETE FROM activities WHERE activity_id = ?
	`
)
