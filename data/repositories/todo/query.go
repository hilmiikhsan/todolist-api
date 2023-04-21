package todo

const (
	queryCreateTodo = `
	INSERT INTO todos (title, activity_group_id, is_active, priority, updated_at) VALUES (?, ?, ?, ?, ?)
	`

	queryGetAllTodo = `
	SELECT
		todo_id as id,
		title,
		activity_group_id,
		is_active,
		priority,
		updated_at,
		created_at
	FROM todos
	`

	queryGetOneTodo = `
	SELECT
		todo_id as id,
		title,
		activity_group_id,
		is_active,
		priority,
		updated_at,
		created_at
	FROM todos
	WHERE todo_id = ?
	`

	queryUpdateTodo = `
	UPDATE todos
	SET
		title = ?,
		is_active = ?,
		priority = ?,
		updated_at = ?
	WHERE todo_id = ?
	`

	queryDeleteTodo = `
	DELETE FROM todos WHERE todo_id = ?
	`
)
