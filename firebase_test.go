package firebase

import (
	"os"
	"testing"
)

func Test_Init(t *testing.T){
	t.Log(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	if err:= Init(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"));err!=nil{
		t.Fatal(err)
	}

}

func Test_GetCredentialToken(t *testing.T)  {
	if err:= Init(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"));err!=nil{
		t.Fatal(err)
	}
	s,err:=GetCredentialToken()
	if err!=nil{
		t.Fatal(err)
	}
	t.Log(s)
}