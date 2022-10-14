package core

import "testing"

func TestGetInitFrozenAccounts(t *testing.T) {
	accountList := GetInitFrozenAccounts(FrozenAccounts)

	for _, account := range accountList.FrozenAccounts {
		t.Log(account.Account.Hex(), account.Amount, account.UnfrozenTime)
	}
}
