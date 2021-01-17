package tests

import (
	"leanote/app/db"
	"testing"
	//	. "leanote/app/lea"
	//	"leanote/app/service"
	//	"gopkg.in/mgo.v2"
	//	"fmt"
)

func TestDBConnect(t *testing.T) {
	db.Init("mongodb://localhost:27017/leanote", "leanote")
}
