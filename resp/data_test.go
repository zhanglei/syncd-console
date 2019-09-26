package resp_test

import (
	"encoding/json"
	"fmt"
	"syncd-console/resp"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	body := `{"code":0,"data":{"token":"f71LpjExKRhxiV4vNKcj2TAu3AoDURscmgM9cvPy_R3NwS39VC43kckNfaj5vbWxlIlN5L-_Hbd2ALFu1erNoA=="},"message":"success"}`
	respData := resp.DataResponse{}
	err := json.Unmarshal([]byte(body), &respData)
	if err != nil {
		t.Fatalf("Unmarshal failed")
	}

	fmt.Println(respData.Data["token"])
}
