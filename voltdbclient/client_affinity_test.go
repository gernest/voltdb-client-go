package voltdbclient

import (
	"context"
	"database/sql/driver"
	"testing"
)

func TestClientAffinity(t *testing.T) {
	servers := "localhost:21212,localhost:21222,localhost:21232"
	conn, err := OpenConn(servers)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	t.Run("should enable client affinity by default", func(ts *testing.T) {
		if !conn.useClientAffinity {
			ts.Error("expected useClientAffinity to be true")
		}
	})
	t.Run("must get topology stats", func(ts *testing.T) {
		nc := conn.getConn()
		_, err := conn.MustGetTopoStatistics(context.Background(), nc)
		if err != nil {
			ts.Fatal(err)
		}
	})
	t.Run("must get partition info ", func(ts *testing.T) {
		nc := conn.getConn()
		_, err := conn.MustGetPTInfo(context.Background(), nc)
		if err != nil {
			ts.Fatal(err)
		}
	})
	t.Run("must get pertition details ", func(ts *testing.T) {
		nc := conn.getConn()
		_, err := conn.GetPartitionDetails(nc)
		if err != nil {
			ts.Fatal(err)
		}
	})
	t.Run("must set partiion details", func(ts *testing.T) {
		r, err := conn.Query("@SystemCatalog", []driver.Value{"PROCEDURES"})
		if err != nil {
			ts.Fatal(err)
		}
		r.Close()
		if conn.PartitionDetails == nil {
			t.Error("expected partition details to be set")
		}
		query := "Vote"
		args := []driver.Value{
			int64(9136958696),
			int32(4),
			int64(2),
		}
		c, err := conn.getConnByCA(conn.PartitionDetails, query, args)
		if err != nil {
			ts.Fatal(err)
		}
		if c == nil {
			t.Error("expected connection")
		}
		// t.Error(pretty.Sprint(conn.PartitionDetails.Procedures))
	})
}
