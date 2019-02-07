package DB

import (
	db "upper.io/db.v3"
)

func (d *dbServer) WaitingforOther(userid string) (uint64, error) {
	Collection := d.sess.Collection(ContractCollection)
	res := Collection.Find(db.Cond{"Creator": userid, "delStatus": 0, "status": "in progress"})
	total, err := res.Count()

	if err != nil {
		return 0, err
	}
	return total, nil
}

func (d *dbServer) WaitingforMe(userid string) (uint64, error) {
	Collection := d.sess.Collection(SignerCollection)
	res := Collection.Find(db.Cond{"userID": userid, "Access": 1, "SignStatus": "needs to sign", "DeleteApprove": 0})
	total, err := res.Count()

	if err != nil {
		return 0, err
	}
	return total, nil
}
