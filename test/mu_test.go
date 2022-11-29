package test

import (
	"context"
	"github.com/995933447/distribmu/factory"
	"github.com/etcd-io/etcd/client"
	"testing"
	"time"
)

func TestEtcdMuLock(t *testing.T) {
	t.Log("start")
	etcdCli, err := client.New(client.Config{
		Endpoints:               []string{"http://127.0.0.1:2379"},
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: 3 * time.Second,
	})
	if err != nil {
		t.Error(err)
		return
	}
	newMuConf := factory.NewMuConf(
		factory.MuTypeEtcd,
		"/abcd",
		time.Second * 10,
		factory.NewEtcdMuDriverConf(etcdCli, "123"),
		)
	mu := factory.MustNewMu(newMuConf)
	success, err := mu.Lock(context.Background())
	if err != nil {
		t.Error(err)
	}
	t.Logf("bool:%v", success)
	err = mu.RefreshTTL(context.Background())
	if err != nil {
		t.Error(err)
	}
	success, err = mu.LockWait(context.Background(), time.Second * 10)
	if err != nil {
		t.Error(err)
	}
	t.Logf("bool:%v", success)
	getResp, err := client.NewKeysAPI(etcdCli).Get(context.TODO(), "abcd", nil)
	if err != nil {
		t.Error(err)
	}
	t.Logf("node val:%v", getResp.Node.Value)
	t.Logf("node val ttl:%v", getResp.Node.TTL)
}