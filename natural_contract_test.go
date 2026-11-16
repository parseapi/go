package parseapi
import ("context";"net/http";"net/http/httptest";"sync/atomic";"encoding/json";"os";"testing";"reflect";"strings")
func TestApproved21RootPayloads(t *testing.T) {
 data,err:=os.ReadFile("testdata/approved-2.1.json");if err!=nil{t.Fatal(err)}
 var rows []struct{ ID string `json:"id"`;Product string `json:"product"`;Fields []string `json:"fields"`;Body json.RawMessage `json:"body"` };if err=json.Unmarshal(data,&rows);err!=nil{t.Fatal(err)}
 read:=func(value any,path string) any { for _,key:=range strings.Split(path,"."){obj,ok:=value.(map[string]any);if !ok{return nil};value=obj[key]};return value }
 for _,row:=range rows{t.Run(row.ID,func(t *testing.T){var model any;switch row.Product {
case "ip": model=&IP{}
case "country": model=&Country{}
case "state": model=&State{}
case "postal": model=&Postal{}
case "email": model=&Email{}
case "carrier": model=&Carrier{}
case "hlr": model=&HLR{}
case "useragent": model=&UserAgent{}
case "weather": model=&Weather{}
}
 if err:=json.Unmarshal(row.Body,model);err!=nil{t.Fatal(err)};encoded,err:=json.Marshal(model);if err!=nil{t.Fatal(err)}
 var before,after map[string]any;json.Unmarshal(row.Body,&before);json.Unmarshal(encoded,&after)
 if _,ok:=after["deep"];ok{t.Fatal("2.1 model emitted a wrapper")}
 // Historical State ISO inputs re-encode using the current canonical key.
 for _,field:=range row.Fields {encodedField:=field;if row.Product=="state" && field=="iso_subdivision"{encodedField="iso_code"};want:=read(before,field);got:=read(after,encodedField);if !reflect.DeepEqual(got,want){t.Fatalf("%s: got %#v want %#v",field,got,want)}}
 })}
}

func TestEmailVerifyExplicitPostAndNoDefaultRetry(t *testing.T) {
 var calls atomic.Int64
 server:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  calls.Add(1)
  if r.Method!="POST" || r.URL.Path!="/email/verify" || r.URL.RawQuery!="" || r.Header.Get("X-API-Version")!="2.2.0" { t.Errorf("wrong verification request: %s %s",r.Method,r.URL) }
  var body map[string]string;if err:=json.NewDecoder(r.Body).Decode(&body);err!=nil{t.Error(err)}
  if !reflect.DeepEqual(body,map[string]string{"email":"a+tag@example.com"}) {t.Errorf("wrong JSON body: %#v",body)}
  w.Header().Set("Content-Type","application/json");w.WriteHeader(503);w.Write([]byte(`{"code":"service_unavailable","message":"pending"}`))
 }));defer server.Close()
 client,err:=New("fixture",WithBaseURL(server.URL));if err!=nil{t.Fatal(err)}
 if _,err=client.EmailVerify(context.Background(),"a+tag@example.com");err==nil{t.Fatal("missing service error")}
 if calls.Load()!=1{t.Fatalf("metered verification retried %d times",calls.Load())}
}
