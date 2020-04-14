package seq

const insertSeqSql = `
INSERT INTO sys_sequence(
seller_id,
seq_name,
seq_ymd,
seq_id
)
VALUES(
:seller_id,
:seq_name,
:seq_ymd,
1
)
`

const updateSeqSql = `
UPDATE sys_sequence seq
SET	seq.seq_id = seq.seq_id + 1
WHERE	seq.seller_id = :seller_id
		AND seq.seq_ymd = :seq_ymd
		AND seq.seq_name = :seq_name
`

const selectSeqSql = `
SELECT
	seq_id
FROM	sys_sequence seq
WHERE	seq.seller_id = :seller_id
		AND seq.seq_ymd = :seq_ymd
		AND seq.seq_name = :seq_name
`

const insertSeqBySellerSql = `
INSERT INTO sys_sequence_seller(
seller_id,
seq_name,
seq_id
)
VALUES(
:seller_id,
:seq_name,
1
)
`

const updateSeqBySellerSql = `
UPDATE sys_sequence_seller seq
SET		seq.seq_id = seq.seq_id + 1
WHERE	seq.seller_id = :seller_id
		AND seq.seq_name = :seq_name
`

const selectSeqBySellerSql = `
SELECT
	seq_id
FROM	sys_sequence_seller seq
WHERE	seq.seller_id = :seller_id
		AND seq.seq_name = :seq_name
`
