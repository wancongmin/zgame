package srvconfig

const (
	insertConfigSql = `
INSERT INTO sys_configure(
	seller_id,
	config_key,
	config_value,
	look_type,
	data_type,
	remark,
	invalid_date
)
VALUES(
	:seller_id,
	:config_key,
	:config_value,
	:look_type,
	:data_type,
	:remark,
	:invalid_date
)
`
	updateConfig = `
UPDATE sys_configure SET
	config_value = :config_value
WHERE 	seller_id = :seller_id
	AND config_key = :config_key
`

	//优先取自己的，取不到取默认的
	getConfigSql = `
SELECT	cfg.config_value
FROM	sys_configure cfg
WHERE	cfg.seller_id IN (0,?)
		AND cfg.config_key = ?
ORDER BY cfg.seller_id DESC
LIMIT 1
`
)
