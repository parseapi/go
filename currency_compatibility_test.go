package parseapi
import("encoding/json";"os";"testing")
func TestCurrencyRestorationPreservesUnknownZeroAndServerPrecision(t *testing.T){
 b,e:=os.ReadFile("testdata/currency-compatibility.json");if e!=nil{t.Fatal(e)};var rows map[string]json.RawMessage;if e=json.Unmarshal(b,&rows);e!=nil{t.Fatal(e)}
 var p Currency;if e=json.Unmarshal(rows["profile"],&p);e!=nil{t.Fatal(e)};if p.Digits==nil||*p.Digits!=0||p.Symbol!=nil||p.Deep==nil||len(p.Countries)!=1||p.Countries[0]!="JP"{t.Fatalf("profile lost values: %#v",p)}
 for _,n:=range []string{"locked","nulls"}{var q Currency;if e=json.Unmarshal(rows[n],&q);e!=nil{t.Fatal(e)};if q.Deep==nil||q.ISONumeric!=nil||q.NamePlural!=nil||q.Countries!=nil{t.Fatalf("unknown changed: %#v",q)}}
 var r,z CurrencyRate;if e=json.Unmarshal(rows["rate"],&r);e!=nil{t.Fatal(e)};if e=json.Unmarshal(rows["zero"],&z);e!=nil{t.Fatal(e)};if r.Rate!=158.33550065||r.Converted==nil||*r.Converted!=1||r.Amount==nil||*r.Amount!=1.005{t.Fatal("client changed server precision")};if z.Rate!=0||z.Amount==nil||*z.Amount!=0||z.Converted==nil||*z.Converted!=0{t.Fatal("lost known zero")}
}
