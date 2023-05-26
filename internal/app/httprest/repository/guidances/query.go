package guidances

const (
	createNewDataQuerry = `
	INSERT INTO public.guidance_file_and_regulation(
		category,
		name,
		description,
		link,
		file,
		file_size,
		file_path,
		file_group,
		owner,
		is_visible,
		version,
		created_by,
		created_at
	) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`
	getAllDataQuerry = `SELECT 
	id, 
	category, 
	name, 
	description, 
	link,
	file,
	owner as file_owner,
	file_size,
	file_path,
	file_group,
	is_visible,
	version,
	created_by,
	created_at,
	updated_by,
	updated_at
	FROM public.guidance_file_and_regulation
	WHERE deleted_at IS NULL
	AND deleted_by IS NUll`
	querryUpdate = `UPDATE public.guidance_file_and_regulation 
	SET category  = $2,
	name = $3,
	description = $4,
	link = $5,
	file = $6,
	version = $7,
	updated_by = $8,
	updated_at = $9,
	file_size = $10,
	file_path = $11,
	file_group = $12,
	owner = $13,
	is_visible = $14
	WHERE id = $1
	AND category = $2`
	querryDelete = `UPDATE public.guidance_file_and_regulation 
	SET deleted_at  = $1,
	deleted_by = $2
	WHERE id = $3`
)