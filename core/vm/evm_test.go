package vm

import (
	"fmt"
	"math/big"
	"testing"
)

func TestUnstakingHeight(t *testing.T) {
	type args struct {
		stakeAmt  *big.Int
		appendAmt *big.Int
		sno       uint64
		cno       uint64
		lockedNo  uint64
	}
	tests := []struct {
		name    string
		args    args
		want    uint64
		wantErr bool
	}{
		{
			name: "TestUnstakingHeight successfully1",
			args: args{stakeAmt: big.NewInt(100000), appendAmt: big.NewInt(300000), sno: 1, cno: 4, lockedNo: 10},
			want: 9,
		},
		{
			name: "TestUnstakingHeight successfully2",
			args: args{stakeAmt: big.NewInt(100000), appendAmt: big.NewInt(300000), sno: 1, cno: 10, lockedNo: 10},
			want: 7,
		},
		{
			name: "TestUnstakingHeight successfully3",
			args: args{stakeAmt: big.NewInt(100000), appendAmt: big.NewInt(300000), sno: 1, cno: 15, lockedNo: 10},
			want: 7,
		},
		{
			name:    "TestUnstakingHeight illegal amount",
			args:    args{stakeAmt: big.NewInt(0), appendAmt: big.NewInt(0), sno: 1, cno: 2, lockedNo: 10},
			wantErr: true,
		},
		{
			name:    "TestUnstakingHeight illegal height",
			args:    args{stakeAmt: big.NewInt(100000), appendAmt: big.NewInt(300000), sno: 2, cno: 1, lockedNo: 10},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnstakingHeight(
				tt.args.stakeAmt,
				tt.args.appendAmt,
				tt.args.sno,
				tt.args.cno,
				tt.args.lockedNo)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnstakingHeight() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UnstakingHeight() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBigFloat(t *testing.T) {
	total := big.NewFloat(0).Add(new(big.Float).SetInt(big.NewInt(100000)), new(big.Float).SetInt(big.NewInt(300000)))
	h1 := big.NewFloat(0).Mul(big.NewFloat(0).Quo(new(big.Float).SetInt(big.NewInt(100000)), total), new(big.Float).SetInt(big.NewInt(int64(0))))
	h2 := big.NewFloat(0).Mul(big.NewFloat(0).Quo(new(big.Float).SetInt(big.NewInt(300000)), total), new(big.Float).SetInt(big.NewInt(int64(10))))
	delayHeight, _ := big.NewFloat(0).Add(h1, h2).Uint64()
	fmt.Println(total)
	fmt.Println(h1)
	fmt.Println(h2)
	fmt.Println(delayHeight)
	a1 := uint64(1)
	a2 := uint64(11)
	a3 := uint64(15)
	fmt.Println(a1 + a2 - a3)
}
