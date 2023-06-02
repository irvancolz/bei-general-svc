package unggahberkas

const (
	uploadNewFilesQuery = `
	INSERT INTO uploaded_files (
        report_type,
        company_code,
        company_name,
		company_id,
        is_uploaded,
        file_name,
        file_path,
        file_size,
        created_by,
        created_at
    )
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	getUploadedFilesQuery = `
	SELECT 
		id,
		report_type AS type,
        company_code,
        company_name,
		company_id,
        is_uploaded,
        file_name,
        file_path,
        file_size,
		created_by,
		created_at,
		updated_by,
		updated_at
	FROM uploaded_files
	WHERE deleted_at IS NULL
		AND deleted_by IS NUll
		AND is_uploaded = true`
	deleteUploadedFilesQuery = `
	UPDATE uploaded_files
	SET deleted_at  = $3,
	deleted_by = $2
	WHERE id = $1
	AND deleted_by IS NULL
	AND deleted_at IS NULL`
	checkDataAvaliabilityQuery = `
	SELECT COUNT(*)
	FROM uploaded_files
	WHERE id = $1
	AND deleted_by IS NULL
	AND deleted_at IS NULL`
	getUploadedFilesPathQuery = `
	SELECT 
	file_path
	FROM uploaded_files
	WHERE id = $1
	AND deleted_by IS NULL
	AND deleted_at IS NULL
	`
)
