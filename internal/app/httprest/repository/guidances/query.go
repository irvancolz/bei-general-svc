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
		"order",
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
	"order",
	version,
	created_by,
	created_at,
	updated_by,
	updated_at
	FROM public.guidance_file_and_regulation
	WHERE deleted_at IS NULL
	AND deleted_by IS NULL
	ORDER BY updated_at, created_at DESC`

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
	"order" = $14
	WHERE id = $1
	AND category = $2`
	querryDelete = `UPDATE public.guidance_file_and_regulation 
	SET deleted_at  = $1,
	deleted_by = $2
	WHERE id = $3`
	updateOrderQuery = `
	UPDATE public.guidance_file_and_regulation 
		SET "order" = "order" + 1
	WHERE "order" >= $1
	AND category = $2
	`
	checkIsOrderFilledQuery = `
	SELECT 
		COUNT(*)
	FROM public.guidance_file_and_regulation
	WHERE "order" = $1
	AND category = $2
	AND deleted_by IS NULL
	AND deleted_at IS NULL
	`
	getCurrentOrderQuery = `
	SELECT 
		"order"
	FROM public.guidance_file_and_regulation
	WHERE id = $1
		`
	getFileSavedPathQuery = `
	SELECT 
		file_path
	FROM public.guidance_file_and_regulation
	WHERE id = $1
	`
	orderQuery = ` ORDER BY
		CASE
			WHEN updated_at IS NOT NULL 
				THEN updated_at
			ELSE created_at
		END DESC`
)
