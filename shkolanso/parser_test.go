package shkolanso

import (
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gopkg.in/jarcoal/httpmock.v1"
)

// в этом файле пробую использвать testify/require, который в отличие от assert падает при первой ошибке в тесте

func TestComposeTextMessage(t *testing.T) {
	loc, err := time.LoadLocation("Asia/Novosibirsk")
	require.NoError(t, err)

	sm := summaryMarks{
		DisciplineMarks: []disciplineMark{
			disciplineMark{
				Discipline: "Математика",
				Marks: []mark{
					mark{
						Date: shkolaTime{Time: time.Date(2018, 11, 15, 0, 0, 0, 0, loc)},
						Mark: "3",
					},
					mark{
						Date: shkolaTime{Time: time.Date(2018, 11, 14, 0, 0, 0, 0, loc)},
						Mark: "5",
					},
					mark{
						Date: shkolaTime{Time: time.Date(2018, 11, 15, 0, 0, 0, 0, loc)},
						Mark: "2",
					},
				},
			},
			disciplineMark{
				Discipline: "Литературное чтение",
				Marks: []mark{
					mark{
						Date: shkolaTime{Time: time.Date(2018, 11, 14, 0, 0, 0, 0, loc)},
						Mark: "5",
					},
				},
			},
			disciplineMark{
				Discipline: "Физическая культура",
				Marks: []mark{
					mark{
						Date: shkolaTime{Time: time.Date(2018, 11, 14, 0, 0, 0, 0, loc)},
						Mark: "5",
					},
				},
			},
		},
	}
	message, err := composeTextMessage(sm, time.Date(2018, 11, 14, 0, 0, 0, 0, loc))
	require.NoError(t, err)
	require.Equal(t, "Математика: 5\nЛитературное чтение: 5\nФизическая культура: 5", message)

	// нет оценок
	message, err = composeTextMessage(sm, time.Date(1999, 1, 1, 0, 0, 0, 0, loc))
	require.NoError(t, err)
	require.Equal(t, "За сегодня нет оценок", message)
}

func TestParseJSONResponse(t *testing.T) {
	sm, err := parseJSONResponse(getTestResponse(t, "summaryMarks.json"))
	require.NoError(t, err)
	require.Len(t, sm.DisciplineMarks, 7)
	require.NotEmpty(t, sm.DisciplineMarks[0].Discipline)
	require.Len(t, sm.DisciplineMarks[0].Marks, 4)

	loc, err := time.LoadLocation("Asia/Novosibirsk")
	require.NoError(t, err)
	require.Equal(t,
		time.Date(2018, 11, 14, 0, 0, 0, 0, loc),
		sm.DisciplineMarks[0].Marks[0].Date.Time,
	)
	require.Equal(t, "5", sm.DisciplineMarks[0].Marks[0].Mark)
}

func TestGetSummaryMarksRaw(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", formURL,
		func(req *http.Request) (*http.Response, error) {
			err := req.ParseForm()
			require.NoError(t, err)
			require.Equal(t, "ll", req.Form.Get("login_login"))
			require.Equal(t, "pp", req.Form.Get("login_password"))

			return httpmock.NewStringResponse(http.StatusOK, `{"success": true,"redirect": "/"}`), nil
		},
	)

	httpmock.RegisterResponder("GET", mainURL, httpmock.NewStringResponder(http.StatusOK, ""))

	httpmock.RegisterResponder("GET", summaryMarksURL,
		func(req *http.Request) (*http.Response, error) {
			require.Contains(t, req.URL.Query(), "date")
			require.Contains(t, req.URL.Query()["date"], time.Now().Format("2006-01-02"))
			return httpmock.NewBytesResponse(http.StatusOK, getTestResponse(t, "summaryMarks.json")), nil
		},
	)

	summaryMarksRaw, err := getSummaryMarksRaw("ll", "pp")

	require.NoError(t, err)
	defer summaryMarksRaw.Close()

	summaryMarksRawString, err := ioutil.ReadAll(summaryMarksRaw)
	require.NoError(t, err)

	require.Equal(t, getTestResponse(t, "summaryMarks.json"), summaryMarksRawString)
}

func TestGetSummaryMarksRawWrongLogin(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", formURL,
		httpmock.NewStringResponder(
			http.StatusUnauthorized,
			`{redirect: "", message: "Неверный логин или пароль.", success: false}`))

	_, err := getSummaryMarksRaw("ll", "pp")
	require.EqualError(t, err, `Неуспешный ответ страницы логина. Status:401 Body:{redirect: "", message: "Неверный логин или пароль.", success: false}`)
}

// getTestResponse возвращает содержимое файла которое является ответом
func getTestResponse(t *testing.T, fileName string) []byte {
	content, err := ioutil.ReadFile("testdata/" + fileName)
	if err != nil {
		t.Errorf("Не могу прочитать файл %s", fileName)
	}

	return content
}
