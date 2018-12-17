package DB

import (
	"log"

	"upper.io/db.v3/lib/sqlbuilder"
	"upper.io/db.v3/mysql"
)

type dbServer struct {
	sess sqlbuilder.Database
}

func (d *dbServer) CloseSession() {
	d.sess.Close()
}

func (d *dbServer) GetSession() sqlbuilder.Database {
	return d.sess
}

//Dbinit : creates data base connection and returns database object
func Dbinit(settings mysql.ConnectionURL) (*dbServer, error) {
	var d dbServer
	var err error

	d.sess, err = mysql.Open(settings)
	if err != nil {
		log.Fatal("connection error", err.Error())
		return nil, err
	}
	return &d, nil
}
