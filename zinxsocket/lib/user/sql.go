package user

const getUserByEmailSql = `
SELECT
	u.seller_id,
	u.user_id,
	u.user_code,
	u.user_name,
	u.user_qq,
	u.user_email,
	u.user_phone,
	u.user_pwd,
	u.user_status,
	u.dept_id,
	u.user_role_id,
	u.create_date,
	u.last_login_date,
	u.login_times,
	u.login_error_times,
	u.reffer
FROM	sys_user u
WHERE u.user_email = ?
`

//bs_seller_info 以下sql
const getSellerIdByTokenSql = `
SELECT si.seller_id
FROM	bs_seller_info si
WHERE	si.token = ?
`
