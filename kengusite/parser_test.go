package kengusite

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
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

func TestGetDataError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	parser := NewParser("zz", "zz")
	var err error

	// Ответ: 302 редирект без Location
	httpmock.RegisterResponder("POST", formURL,
		httpmock.NewStringResponder(http.StatusFound, ""))
	_, err = parser.GetData()
	assert.EqualError(t, err, "не могу запостить форму: Post \"https://billing.kengudetyam.ru/cabinet/Account/Login\": 302 response missing Location header")

	httpmock.RegisterResponder("POST", formURL,
		httpmock.NewStringResponder(http.StatusOK, "blabla"))
	_, err = parser.GetData()
	assert.EqualError(t, err, "не могу найти '.balance' в html коде")
}

// getTestResponse возвращает содержимое файла которое является ответом
func getTestResponse(t *testing.T, fileName string) []byte {
	content, err := ioutil.ReadFile("testdata/" + fileName)
	if err != nil {
		t.Errorf("Не могу прочитать файл %s", fileName)
	}

	return content
}
