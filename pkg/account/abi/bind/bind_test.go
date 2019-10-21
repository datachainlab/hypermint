package bind

import "testing"

func TestBind(t *testing.T) {
	src, err := Bind("main", "test", `
[
  {
    "type": "function",
    "name": "hoge",
    "inputs": [
      {
        "name": "piyo",
        "type": "address"
      },
      {
        "name": "unko",
        "type": "u64"
      }
    ],
    "outputs": [
      {
        "name": "",
        "type": "u64"
      }
    ]
  }
]
`, true)
	if err != nil {
		t.Error(err)
	}
	t.Log(src)
}
