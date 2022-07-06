package ethhelper

import (
	"fmt"
	"testing"
)

func TestBalanceOfWeth(t *testing.T) {
	c, err := BalanceOfWeth("0x077d34394Ed01b3f31fBd9816cF35d4558146066")
	fmt.Println(c, err)
}

func TestAllowanceOfWeth(t *testing.T) {
	c, err := AllowanceOfWeth("0x10CEc672c6BB2f6782BEED65987E020902B7bD15")
	fmt.Println(c, err)
}
//0xdd62ed3e00000000000000000000000010cec672c6bb2f6782beed65987e020902b7bd15000000000000000000000000077d34394ed01b3f31fbd9816cf35d4558146066
//0xdd62ed3e00000000000000000000000010cec672c6bb2f6782beed65987e020902b7bd15000000000000000000000000077d34394ed01b3f31fbd9816cf35d4558146066