package i

import (
	"reflect"
	"testing"
)

func TestP(t *testing.T) {
	testCases := []struct {
		s    s
		p, r v
	}{
		{"", nil, nil},
		{"1", z1, z1},
		{"1;2\n3", l{";", z1, z2, z3}, z3},
		{"1 2 3", zv{1, 2, 3}, zv{1, 2, 3}},
		{"1.2 3a90", zv{1.2, c(0, 3)}, zv{1.2, c(0, 3)}},
		{" -1  3a180 0i2 ", zv{-1, -3, c(0, 2)}, zv{-1, -3, c(0, 2)}},
		{"`a", l{"`", "a"}, "a"},
		{"`", l{"`", ""}, ""},
		{"`a`b", sv{"a", "b"}, sv{"a", "b"}},
		{"`a_1 `bZ3\"+-,:\"", sv{"a_1", "bZ3", "+-,:"}, sv{"a_1", "bZ3", "+-,:"}},
		{`"x"`, l{"`", "x"}, "x"},
		{`""`, l{"`", ""}, ""},
		{`"xy"`, l{"`", "xy"}, "xy"},
		{`"x\ny"`, l{"`", "x\ny"}, "x\ny"},
		{`[a:1;beta:2 3]`, l{"!", sv{"a", "beta"}, l{nil, z1, zv{2, 3}}}, [2]l{l{"a", "beta"}, l{z1, zv{2, 3}}}},
		{"x_3:1 2i3", l{":", "x_3", zv{1, c(2, 3)}}, zv{1, c(2, 3)}},
		{"+", "+", "fn"},
		{"1+2", l{"+", z1, z2}, z3},
		{"1+", l{"+", z1, nil}, "fn"},
		{"*1", l{"*", z1}, z1},
		{"a:1;a", l{";", l{":", "a", z1}, "a"}, z1},
		{"x_3:2", l{":", "x_3", z2}, z2},
		{"x:1;x+1", l{";", l{":", "x", z1}, l{"+", "x", z1}}, z2},
		{"1+(1;2;3)", l{"+", z1, l{nil, z1, z2, z3}}, l{z2, z3, z4}},
		{"(1+2)-3", l{"-", l{"+", z1, z2}, z3}, z0},
		{"1 2 3[0 2]", l{"@", zv{1, 2, 3}, zv{0, 2}}, zv{1, 3}},
		{"`a`b`c[2]", l{"@", sv{"a", "b", "c"}, z2}, "c"},
		{"(1;2;3)[0 1]", l{"@", l{nil, z1, z2, z3}, zv{0, 1}}, zv{1, 2}},
		{"`a`b!3 4", l{"!", sv{"a", "b"}, zv{3, 4}}, [2]l{l{"a", "b"}, l{z3, z4}}},
		{"[a:(1+2);b:4]", l{"!", sv{"a", "b"}, l{nil, l{"+", z1, z2}, z4}}, [2]l{l{"a", "b"}, l{z3, z4}}},
		{"+/1 2 3", l{l{"/", "+"}, zv{1, 2, 3}}, zi(6)},
		{"+/3", l{l{"/", "+"}, z3}, z3},
		{"+/,3", l{l{"/", "+"}, l{",", z3}}, z3},
		{"2-/3 7 8", l{l{"/", "-"}, z2, zv{3, 7, 8}}, zi(-16)},
		{`-\3 7 8`, l{l{`\`, "-"}, zv{3, 7, 8}}, zv{3, -4, -12}},
		{`2-\3 7 8`, l{l{`\`, "-"}, z2, zv{3, 7, 8}}, zv{-1, -8, -16}},
		{`1 2 3-\3 7 8`, l{l{`\`, "-"}, zv{1, 2, 3}, zv{3, 7, 8}}, l{zv{-2, -1, -0}, zv{-9, -8, -7}, zv{-17, -16, -15}}},
		// {`{%x}'4 5`, l{l{"'", "%"}, zv{4, 5}}, zv{0.25, 0.2}}, // TODO
		{`15%'3 5`, l{l{"'", "%"}, c(15, 0), zv{3, 5}}, zv{5, 3}},
		{`-':3 5 2`, l{l{"':", "-"}, zv{3, 5, 2}}, zv{3, 2, -3}},
		{`=':3 3 4 4 5`, l{l{"':", "="}, zv{3, 3, 4, 4, 5}}, zv{3, 1, 0, 1, 0}},
		{`99,':2 3 4`, l{l{"':", ","}, c(99, 0), zv{2, 3, 4}}, l{zv{2, 99}, zv{3, 2}, zv{4, 3}}},
		{`2 3,/:4 5 6`, l{l{"/:", ","}, zv{2, 3}, zv{4, 5, 6}}, l{zv{2, 3, 4}, zv{2, 3, 5}, zv{2, 3, 6}}},
		{`2 3 4 ,\: 5 6 7`, l{l{`\:`, ","}, zv{2, 3, 4}, zv{5, 6, 7}}, l{zv{2, 5, 6, 7}, zv{3, 5, 6, 7}, zv{4, 5, 6, 7}}},
		{`2,:/3`, l{l{"/", ",:"}, z2, z3}, l{zv{3}}},
		{`f:(100>);f (2*)/1`,
			l{";", l{":", "f", l{">", c(100, 0), nil}},
				l{l{"/", l{"*", z2, nil}}, "f", z1}},
			c(128, 0)},
	}
	for _, tc := range testCases {
		p := P(tc.s)
		tt(t, tc.p, p, "P: %q %+v\n", tc.s, tc.p)
		r := E(p, nil)
		if tc.r == "fn" && rval(r).Kind() == reflect.Func {
		} else {
			tt(t, tc.r, r, "E: %q %+v\n", tc.s, tc.r)
		}
	}
}

func TestScan(t *testing.T) {
	type iv [10]int //0     1     2     3     4     5     6     7     8     9
	var f = [10]sf{sNum, sNam, sStr, sVrb, sAsn, sIov, sAdv, sViw, sDct, sWsp}
	var testCases = []struct {
		s s
		r iv
	}{
		{` `, iv{0, 0, 0, 0, 0, 0, 0, 0, 0, 1}},
		{`0`, iv{1, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
		{`23`, iv{2, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
		{`+1`, iv{2, 0, 0, 1, 0, 0, 0, 0, 0, 0}}, // +1 is a number
		{`-1`, iv{2, 0, 0, 1, 0, 0, 0, 0, 0, 0}},
		{`1e`, iv{1, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
		{`1.`, iv{1, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
		{`.5`, iv{0, 0, 0, 1, 0, 0, 0, 0, 0, 0}},  // no number: .5 use 0.5
		{`1i`, iv{1, 0, 0, 0, 0, 0, 0, 0, 0, 0}},  // no number: 1i
		{`0i1`, iv{3, 0, 0, 0, 0, 0, 0, 0, 0, 0}}, // complex i: 0i1
		{`-1i0`, iv{4, 0, 0, 1, 0, 0, 0, 0, 0, 0}},
		{`i`, iv{0, 1, 0, 0, 0, 0, 0, 0, 0, 0}},
		{`-i`, iv{0, 0, 0, 1, 0, 0, 0, 0, 0, 0}},
		{`1.23e+06i-1.23e-13`, iv{18, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
		{`13.275a275.2`, iv{12, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
		{`π`, iv{0, 1, 0, 0, 0, 0, 0, 0, 0, 0}}, // name!
		{`a`, iv{0, 1, 0, 0, 0, 0, 0, 0, 0, 0}},
		{`a2`, iv{0, 2, 0, 0, 0, 0, 0, 0, 0, 0}},
		{`a2/`, iv{0, 2, 0, 0, 0, 0, 0, 0, 0, 0}},
		{"`a", iv{0, 0, 2, 0, 0, 0, 0, 0, 0, 0}},
		{"`a3", iv{0, 0, 3, 0, 0, 0, 0, 0, 0, 0}},
		{"`a3.", iv{0, 0, 3, 0, 0, 0, 0, 0, 0, 0}},
		{`"a"`, iv{0, 0, 3, 0, 0, 0, 0, 0, 0, 0}},
		{"`a_3", iv{0, 0, 4, 0, 0, 0, 0, 0, 0, 0}},
		{`"a"b`, iv{0, 0, 3, 0, 0, 0, 0, 0, 0, 0}},
		{`"x\ny"`, iv{0, 0, 6, 0, 0, 0, 0, 0, 0, 0}},
		{`"a\"b\n"b`, iv{0, 0, 8, 0, 0, 0, 0, 0, 0, 0}},
		{`+`, iv{0, 0, 0, 1, 0, 0, 0, 0, 0, 0}},
		{`⍟3`, iv{0, 0, 0, 1, 0, 0, 0, 0, 0, 0}},
		{`⍟:`, iv{0, 0, 0, 2, 0, 0, 0, 0, 0, 0}},
		{`+:`, iv{0, 0, 0, 2, 0, 0, 0, 0, 0, 0}},
		{`1:`, iv{1, 0, 0, 0, 0, 2, 0, 0, 0, 0}},
		{`/`, iv{0, 0, 0, 0, 0, 0, 1, 0, 0, 0}},
		{`':`, iv{0, 0, 0, 0, 0, 0, 2, 0, 0, 0}},
		{`⍨`, iv{0, 0, 0, 0, 0, 0, 1, 0, 0, 0}},
		{`::x`, iv{0, 0, 0, 0, 0, 0, 0, 2, 0, 0}},
		{`[`, iv{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
		{`[:`, iv{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
		{`[3:`, iv{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
		{`[a:`, iv{0, 0, 0, 0, 0, 0, 0, 0, 3, 0}},
		{`[a3:`, iv{0, 0, 0, 0, 0, 0, 0, 0, 4, 0}},
		{`[ab3:`, iv{0, 0, 0, 0, 0, 0, 0, 0, 5, 0}},
		{`  \n `, iv{0, 0, 0, 0, 0, 0, 0, 0, 0, 2}},
		{"  \t\r ", iv{0, 0, 0, 0, 0, 0, 0, 0, 0, 5}},
	}
	for _, tc := range testCases {
		for k, f := range f {
			if n := f([]rune(tc.s)); n != tc.r[k] {
				t.Fatalf("%q: f[%d] got %d, exp %d", tc.s, k, n, tc.r[k])
			}
		}
	}
}
func TestBeg(t *testing.T) {
	testCases := [][2]string{
		{"xyz", "xyz"},
		{"/xyz", ""},
		{"/x\nyz", "\nyz"},
		{"ab/x\nyz", "ab/x\nyz"},
		{"ab /x\nyz", "ab \nyz"},
		{"ab /x;yz", "ab "},
		{`ab" /x;"yz`, `ab" /x;"yz`},
		{"1+2", "1+ 2"},
		{"a-2", "a- 2"},
		{"a-b", "a-b"},
		{"-13", "-13"},
		{"2.5e-03", "2.5e-03"},
		{"a+b", "a+b"},
		{"a*2.3i-5", "a*2.3i-5"},
	}
	for _, tc := range testCases {
		r := string(beg(rv(tc[0])))
		tt(t, tc[1], r, "beg %q %q\n", tc[0], tc[1])
	}
}
