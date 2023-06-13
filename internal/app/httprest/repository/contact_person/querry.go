package contactperson

const (
	syncContactPersonCompaniesQuery = `
	INSERT INTO institutions (
		code,
		name,
		address,
		website,
		postal_code,
		fax,
		telephone,
		business_permit_ojk,
		permit_bursa,
		other_business_permit_ojk,
		type,
		status,
		license,
		operational_status,
		created_by,
		created_at,
		updated_by,
		updated_at,
		deleted_by,
		deleted_at,
		is_deleted
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21
	) ON CONFLICT (code) DO UPDATE SET
		name = EXCLUDED.name,
		address = EXCLUDED.address,
		website = EXCLUDED.website,
		postal_code = EXCLUDED.postal_code,
		fax = EXCLUDED.fax,
		telephone = EXCLUDED.telephone,
		business_permit_ojk = EXCLUDED.business_permit_ojk,
		permit_bursa = EXCLUDED.permit_bursa,
		other_business_permit_ojk = EXCLUDED.other_business_permit_ojk,
		type = EXCLUDED.type,
		status = EXCLUDED.status,
		license = EXCLUDED.license,
		operational_status = EXCLUDED.operational_status,
		created_by = EXCLUDED.created_by,
		created_at = EXCLUDED.created_at,
		updated_by = EXCLUDED.updated_by,
		updated_at = EXCLUDED.updated_at,
		deleted_by = EXCLUDED.deleted_by,
		deleted_at = EXCLUDED.deleted_at,
		is_deleted = EXCLUDED.is_deleted
	`
	getAllDivisionByCompanyQuerry = `
	SELECT 
		id, name, created_at, created_by, updated_at, updated_by
	FROM institution_division 
	WHERE deleted_by IS NULL
	AND deleted_at IS NULL`
	getAllDivisionQuerry = `
	SELECT 
		id, name
	FROM institution_division 
	WHERE deleted_by IS NULL
	AND deleted_at IS NULL`
	addDivisionQuerry = `
	INSERT INTO public.institution_division (
		is_default,
		name,
		created_at,
		created_by
	)VALUES ($1, $2, $3, $4 )`
	getMemberByCompanyQuerry = `
	SELECT 
		 m.id,
		i.status AS institute_status,
		i.id AS institute_id,
		i.type AS institute_type,
		m.name,
		i.name AS company_name,
		i.code AS company_code,
		d.name AS division,
		m.position,
		m.email,
		m.phone,
		m.telephone
	FROM institution_members m
	JOIN institutions i
		ON m.institution_id = i.id
	LEFT JOIN institution_division d
		ON m.division_id = d.id
	WHERE i.is_deleted = false
	AND m.deleted_at IS NULL
	AND m.deleted_by IS NULL
	AND d.deleted_at IS NULL
	AND d.deleted_by IS NULL
	ORDER BY i.code ASC`
	crateMemberQuerry = `
	INSERT INTO public.institution_members (
		institution_id,
		division_id,
		name,
		phone,
		telephone,
		email,
		position,
		created_at,
		created_by
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		ON CONFLICT DO NOTHING RETURNING id`
	addMemberToDivisionQuerry = `
	INSERT INTO institution_member_division (
		division_id,
		institution_member_id,
		created_at,
		created_by
	)VALUES ($1,$2,$3,$4)`
	editMemberDivisionQuerry = `
	UPDATE institution_member_division
	SET 
		division_id = $3,
		updated_at = $4,
		updated_by = $5
	WHERE institution_member_id = $1
	AND division_id = $2
	AND deleted_by IS NULL
	AND deleted_at IS NULL`
	checkDivisionDeleteAvaliabilityQuerry = `
	SELECT 
		COUNT(*)
	FROM institution_members m 
	JOIN institution_division d
		ON (m.division_id = d.id)
	WHERE d.id = $1 
	AND d.deleted_at IS NULL
	AND d.deleted_by IS NULL
	AND m.deleted_at IS NULL
	AND m.deleted_by IS NULL`
	checkMemberViewAvaliabilityQuerry = `
	SELECT 
		COUNT(*)
	FROM institution_members m 
	WHERE m.id =$1
	AND deleted_by IS NULL
	AND deleted_at IS NULL`
	getAllCompanyQuerry = `
	SELECT 
		id,
		type,
		code,
		name, 
		status 
	FROM institutions
	WHERE is_deleted = false`
	getAllCompanyQuerryWithKeyword = `
	SELECT 
		id,
		type,
		code,
		name, 
		status 
	FROM institutions
	WHERE is_deleted = false
	AND (
		LOWER(code) LIKE LOWER('%' || $1 || '%' )
		OR 
		LOWER("name") LIKE LOWER('%' || $1 || '%' )
	)`
	getMembersDetailQuerry = `
	SELECT 
		m.id,
		m.name,
		m.institution_id,
		m.position,
		m.phone,
		m.telephone,
		m.email,
		m.created_at,
		m.created_by,
		m.updated_at,
		m.updated_by,
		d.name AS division,
		d.id AS division_id
	FROM institution_members m
	JOIN institution_division d
		ON m.division_id = d.id
	WHERE m.id = $1
	AND m.deleted_by IS NULL
	AND m.deleted_at IS NULL`
	getMembersDivisionQuerry = `
	SELECT 
		md.id, 
		d.id AS division_id, 
		d.name AS division_name 
	FROM institution_division d 
	JOIN institution_member_division md
		ON md.division_id = d.id
	JOIN institution_members m
		ON m.id = md.institution_member_id
	WHERE m.id = $1`
	editMemberQuerry = `
	UPDATE institution_members
	SET 
		name = $2,
		position = $3,
		email = $4,
		phone = $5,
		telephone = $6,
		updated_at = $7,
		updated_by = $8,
		division_id = $9
	WHERE id = $1
	AND deleted_by IS NULL
	AND deleted_at IS NULL`
	editDivisionQuerry = `
	UPDATE institution_division
	SET 
		name = $2,
		updated_at= $3,
		updated_by= $4
	WHERE id = $1
	AND is_default = false
	AND deleted_by IS NULL
	AND deleted_at IS NULL`
	deleteCompanyMemberQuerry = `
	UPDATE institution_members
	SET 
		deleted_by = $2,
		deleted_at = $3
	WHERE id = $1
	AND deleted_by IS NULL
	AND deleted_at IS NULL`
	deleteCompanyDivisionQuerry = `
	UPDATE institution_division
	SET 
		deleted_by = $2,
		deleted_at = $3
	WHERE id = $1
	AND is_default = false
	AND created_by IS NOT NULL
	AND deleted_by IS NULL
	AND deleted_at IS NULL`
	checkDivisionEditAvailabilityQuerry = `
	SELECT COUNT(*)
	FROM institution_division
	WHERE id = $1
	AND is_default = false
	AND created_by IS NOT NULL`
	checkDivisionViewAvailabilityQuerry = `
	SELECT 
		COUNT(*)
	FROM institution_division
	WHERE id = $1
	AND deleted_by IS NULL
	AND deleted_at IS NULL`
	getCompanyDetailQuery = `
	SELECT 
		id, 
		name, 
		type, 
		code, 
		status, 
		created_at, 
		created_by, 
		updated_at, 
		updated_by
	FROM institutions 
	WHERE id = $1
	AND deleted_by IS NULL
	AND deleted_at IS NULL`
	getCompanyDivisionByCompanyCode = `
	SELECT 
		d.id, 
		d.name,
		p.id AS company_id,
		p.name AS company_name,
		p.code AS company_code
	FROM institution_division d 
	JOIN institutions p 
		ON (d.institution_id = p.id OR d.institution_id IS NULL) 
	AND LOWER(p.code) = LOWER($1)
	AND d.deleted_at IS NULL
	AND d.deleted_by IS NULL`
	getMemberDetailByDivisionBaseQuery = `
	SELECT 
		DISTINCT ON (m.id) m.id,
		m.name,
		m.phone,
		m.email,
		m.telephone,
		m.position, 
		d.name AS division 
	FROM institution_members m  
	JOIN institution_division d 
		ON m.division_id = d.id 
	JOIN institutions i
		ON m.institution_id = i.id
	WHERE m.division_id  IN( %s ) 
	AND m.deleted_by IS NULL
	AND m.deleted_at IS NULL`
	getMemberDetailByDivisionAndCompanyCodeBaseQuery = `
	SELECT 
		DISTINCT ON (m.id) m.id,
		m.name,
		m.phone,
		m.email,
		m.telephone,
		m.position, 
		d.name AS division 
	FROM institution_members m  
	JOIN institution_division d 
		ON m.division_id = d.id 
	WHERE m.deleted_by IS NULL
	AND m.deleted_at IS NULL`
)
