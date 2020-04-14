package seq

import (
	"bangseller.com/lib/exception"
	"bangseller.com/lib/fun"
	"bangseller.com/lib/mdb"
	"fmt"
	"strconv"
)

type Seq struct {
	SellerId int    `db:"seller_id"`
	SeqName  string `db:"seq_name"`
	SeqYmd   string `db:"seq_ymd"`
	Len      int    //序号长度,默认为4，如果不输入的话
}

//生成唯一ID, YYmmddx..x
func GetIDByDate(s *Seq) int64 {
	if s.Len == 0 {
		s.Len = 4
	}
	s.SeqYmd = fun.Now().Format("060102")

	tx := mdb.Db.MustBegin()
	defer tx.Rollback()
	res, err := tx.NamedExec(updateSeqSql, s)
	exception.CheckSqlError(err)
	rows, err := res.RowsAffected()
	var id int64 = 1
	if rows > 0 {
		stmt, err := tx.PrepareNamed(selectSeqSql)
		exception.CheckSqlError(err)
		defer stmt.Close()
		err = stmt.Get(&id, s)
		exception.CheckSqlError(err)
	} else {
		//这儿不严谨，如果同时多个请求来，将会重复插入，Insert语句无事务特性，所以此处应该用其他锁方式
		//不过可以对数据库加唯一主键，这样可以避免重复插入，但是将会对客户端不友好
		//测试发现，好像不管是否更新到数据，都会被锁住，执行 mdb.Db.NamedExec(insertSeqSql,s) 会被锁住，这样问题就解决了
		_, err = tx.NamedExec(insertSeqSql, s)
		exception.CheckSqlError(err)
	}
	tx.Commit()

	sid := s.SeqYmd + fmt.Sprintf("%04d", id)
	id, _ = strconv.ParseInt(sid, 10, 64)
	return id
}

//生成唯一ID,By Seller,不区分日期，序号从1开始
func GetIDBySeller(s *Seq) int {
	//	s.SeqYmd = "600101" //不单独用表保存，直接保存到 sys_sequence 表中，但日期设置靠后

	tx := mdb.Db.MustBegin()
	defer tx.Rollback()
	res, err := tx.NamedExec(updateSeqBySellerSql, s)
	exception.CheckSqlError(err)
	rows, err := res.RowsAffected()
	var id int64 = 1
	if rows > 0 {
		stmt, err := tx.PrepareNamed(selectSeqBySellerSql)
		exception.CheckSqlError(err)
		defer stmt.Close()
		err = stmt.Get(&id, s)
		exception.CheckSqlError(err)
	} else {
		//这儿不严谨，如果同时多个请求来，将会重复插入，Insert语句无事务特性，所以此处应该用其他锁方式
		//不过可以对数据库加唯一主键，这样可以避免重复插入，但是将会对客户端不友好
		//测试发现，好像不管是否更新到数据，都会被锁住，执行 mdb.Db.NamedExec(insertSeqSql,s) 会被锁住，这样问题就解决了
		_, err = tx.NamedExec(insertSeqBySellerSql, s)
		exception.CheckSqlError(err)
	}
	tx.Commit()
	return int(id)
}
