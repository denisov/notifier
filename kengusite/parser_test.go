package kengusite

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/jarcoal/httpmock.v1"
)

func TestGetData(t *testing.T) {

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", formURL,
		httpmock.NewBytesResponder(200, getTestResponse(t, "kengu.html")))

	parser := NewParser("zz", "zz")
	data, err := parser.GetData()
	assert.NoError(t, err)

	assert.Equal(t, "1 015.37 руб.", data)
}

// getTestResponse возвращает содержимое файла которое является ответом
func getTestResponse(t *testing.T, fileName string) []byte {
	content, err := ioutil.ReadFile("testdata/" + fileName)
	if err != nil {
		t.Errorf("Не могу прочитать файл %s", fileName)
	}

	return content
}
