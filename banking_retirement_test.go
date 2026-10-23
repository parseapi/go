package parseapi
import ("reflect"; "testing")
func TestBankingRetirement(t *testing.T) {
 typ := reflect.TypeOf((*Client)(nil))
 for _, name := range []string{"Bank", "IBAN", "ACH", "Routing", "BankRequirements", "IBANRequirements", "ACHRequirements"} {
  if _, found := typ.MethodByName(name); found { t.Fatalf("retired method remains: %s", name) }
 }
 for _, name := range []string{"Card", "Country"} { if _, found := typ.MethodByName(name); !found { t.Fatalf("surviving method missing: %s", name) } }
}
