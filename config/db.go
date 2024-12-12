package config

import (
	"os"

	"github.com/gocroot/helper/atdb"
)

var MongoString string = os.Getenv("MONGOSTRING")
var mongoinfo = atdb.DBInfo{
	DBString: MongoString,
	DBName:   "jajankuy",
}
var Mongoconn, ErrorMongoconn = atdb.MongoConnect(mongoinfo)

// Geospacial Database
var MongoStringGeo string = "mongodb+srv://map:admin123@map.9ieis.mongodb.net/"

var mongoinfoGeo = atdb.DBInfo{
	DBString: MongoStringGeo,
	DBName:   "maps",
}

var MongoconnGeo, ErrorMongoconnGeo = atdb.MongoConnect(mongoinfoGeo)