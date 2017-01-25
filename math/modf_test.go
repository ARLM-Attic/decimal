package math

import (
	"fmt"
	"testing"

	"github.com/ericlagergren/decimal"
)

func TestBig_Modf(t *testing.T) {
	tests := [...]struct {
		dec  string
		intg string
		frac string
	}{
		0:  {"296474.3772789836", "296474", "0.3772789836"},
		1:  {"-317556.65040295396", "-317556", "-0.65040295396"},
		2:  {"832564.9806826359", "832564", "0.9806826359"},
		3:  {"934740.1797369102", "934740", "0.1797369102"},
		4:  {"520774.2085155975", "520774", "0.2085155975"},
		5:  {"487789.5461755025", "487789", "0.5461755025"},
		6:  {"-562938.1338471398", "-562938", "-0.1338471398"},
		7:  {"77234.6153124352", "77234", "0.6153124352"},
		8:  {"-190793.61336112698", "-190793", "-0.61336112698"},
		9:  {"-117612.82472773292", "-117612", "-0.82472773292"},
		10: {"562490.1936370123", "562490", "0.1936370123"},
		11: {"-454463.82465413434", "-454463", "-0.82465413434"},
		12: {"-468197.66342017427", "-468197", "-0.66342017427"},
		13: {"-45385.99390947411", "-45385", "-0.99390947411"},
		14: {"108976.52613041783", "108976", "0.52613041783"},
		15: {"-883838.5535850908", "-883838", "-0.5535850908"},
		16: {"710903.7632352882", "710903", "0.7632352882"},
		17: {"-647550.5910254063", "-647550", "-0.5910254063"},
		18: {"929294.3969446819", "929294", "0.3969446819"},
		19: {"367607.3691864349", "367607", "0.3691864349"},
		20: {"377847.0999681826", "377847", "0.0999681826"},
		21: {"-921604.174125825", "-921604", "-0.174125825"},
		22: {"76410.7862004172", "76410", "0.7862004172"},
		23: {"-104096.36392393638", "-104096", "-0.36392393638"},
		24: {"940700.46632265", "940700", "0.46632265"},
		25: {"-536862.4033232268", "-536862", "-0.4033232268"},
		26: {"675503.2444769265", "675503", "0.2444769265"},
		27: {"737754.1066881085", "737754", "0.1066881085"},
		28: {"-812094.9541646338", "-812094", "-0.9541646338"},
		29: {"577545.4240398374", "577545", "0.4240398374"},
		30: {"-573554.9376775054", "-573554", "-0.9376775054"},
		31: {"-546642.96324421", "-546642", "-0.96324421"},
		32: {"162519.7570301781", "162519", "0.7570301781"},
		33: {"612961.6010606149", "612961", "0.6010606149"},
		34: {"196102.13226522505", "196102", "0.13226522505"},
		35: {"832033.3624345269", "832033", "0.3624345269"},
		36: {"-344966.9944758781", "-344966", "-0.9944758781"},
		37: {"-933325.047928792", "-933325", "-0.047928792"},
		38: {"-202012.0305155943", "-202012", "-0.0305155943"},
		39: {"424408.7393911381", "424408", "0.7393911381"},
		40: {"667823.2651242441", "667823", "0.2651242441"},
		41: {"270847.74042637344", "270847", "0.74042637344"},
		42: {"-829152.1560524672", "-829152", "-0.1560524672"},
		43: {"-666573.9418051436", "-666573", "-0.9418051436"},
		44: {"905488.3131974128", "905488", "0.3131974128"},
		45: {"439093.1743144549", "439093", "0.1743144549"},
		46: {"-357757.9398403404", "-357757", "-0.9398403404"},
		47: {"-705790.0738368309", "-705790", "-0.0738368309"},
		48: {"565130.6182316178", "565130", "0.6182316178"},
		49: {"-358703.4697589793", "-358703", "-0.4697589793"},
		50: {"783249.3845349478", "783249", "0.3845349478"},
		51: {"0", "0", "0"},
		52: {"0.0", "0", "0"},
		53: {"1.0", "1", "0"},
		54: {"0.1", "0", "0.1"},
		55: {"1", "1", "0"},
		56: {"100000000000000000000", "100000000000000000000", "0"},
		57: {"100000000000000000000.1", "100000000000000000000", "0.1"},
		58: {"100000000000000000000.0", "100000000000000000000", "0"},
		59: {"0.000000000000000000001", "0", "0.000000000000000000001"},
	}
	for i, v := range tests {
		dec := newbig(v.dec)
		integ, frac := Modf(new(decimal.Big), dec)
		m := new(decimal.Big).Add(integ, frac)
		vig := newbig(v.intg)
		vfr := newbig(v.frac)
		if m.Cmp(dec) != 0 || integ.Cmp(vig) != 0 || vfr.Cmp(frac) != 0 {
			t.Fatalf("#%d: Modf(%s) wanted (%s, %s), got (%s, %s)",
				i, v.dec, v.intg, v.frac, integ, frac)
		}
		testFormZero(t, integ, fmt.Sprintf("#%d: integral part", i))
		testFormZero(t, frac, fmt.Sprintf("#%d: fractional part", i))
	}
}

// testFormZero verifies that if z == 0, z.form == zero.
func testFormZero(t *testing.T, z *decimal.Big, name string) {
	iszero := z.Cmp(zero) == 0
	if iszero && z.Sign() == 0 {
		t.Fatalf("%s: z == 0, but form not marked zero: %+v", name, z)
	}
	if !iszero && z.Sign() == 0 {
		t.Fatalf("%s: z != 0, but form marked zero", name)
	}
}
