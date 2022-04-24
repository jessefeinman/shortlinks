package shortlinks

import (
	"net/http"
	"strings"
)

type edit struct {
	Shortlink
	Submit string

	History []History
}

func (e edit) Title() string {
	if e.From == "" {
		return "Create"
	}

	return "Edit " + e.From
}

func editHandler(db DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		from := strings.TrimPrefix(r.URL.Path, "/_edit/")

		if r.Method == "POST" {
			if err := r.ParseForm(); err != nil {
				_500(w, err)
				return
			}

			if from == "" {
				from = r.Form.Get("from")
			}

			if err := db.InsertHistory(History{From: from, To: r.Form.Get("to")}); err != nil {
				_500(w, err)
				return
			}
			if err := db.CreateShortlink(Shortlink{
				To:   r.Form.Get("to"),
				From: from,
			}); err != nil {
				_500(w, err)
				return
			}
			w.Header().Add("Location", "/")
			w.WriteHeader(302)
			return
		}

		sl, err := db.Shortlink(from)
		if err != nil {
			_500(w, err)
			return
		}

		h, err := db.History(from)
		if err != nil {
			_500(w, err)
			return
		}

		v := edit{
			Shortlink: sl,
			History:   h,

			Submit: "Update",
		}

		if err := tpl.ExecuteTemplate(w, "edit.html", v); err != nil {
			_500(w, err)
			return
		}
	})
}
