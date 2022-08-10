package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type task struct {
	Tasks    string
	Assignee string
	Deadline string
	Status   int
}

type response struct {
	Status bool
	Pesan  string
	Data   []task
}

func koneksi() (*sql.DB, error) {
	db, salahe := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/tugasgolang")
	if salahe != nil {
		return nil, salahe
	}

	return db, nil
}

func tampil(pesane string) response {
	db, salahe := koneksi()

	if salahe != nil {
		return response{
			Status: false,
			Pesan:  "Gagal Koneksi: " + salahe.Error(),
			Data:   []task{},
		}
	}
	defer db.Close()

	dataTsk, salahe := db.Query("select * from task")
	if salahe != nil {
		return response{
			Status: false,
			Pesan:  "Gagal Query: " + salahe.Error(),
			Data:   []task{},
		}
	}
	defer dataTsk.Close()

	var hasil []task

	for dataTsk.Next() {
		var tsk = task{}
		var salahe = dataTsk.Scan(&tsk.Tasks, &tsk.Assignee, &tsk.Deadline, &tsk.Status)

		if salahe != nil {
			return response{
				Status: false,
				Pesan:  "Gagal Baca: " + salahe.Error(),
				Data:   []task{},
			}
		}

		hasil = append(hasil, tsk)
	}

	salahe = dataTsk.Err()

	if salahe != nil {
		return response{
			Status: false,
			Pesan:  "Kesalahan: " + salahe.Error(),
			Data:   []task{},
		}
	}

	return response{
		Status: true,
		Pesan:  pesane,
		Data:   hasil,
	}

}

func getTsk(tasks string) response {
	db, salahe := koneksi()
	if salahe != nil {
		return response{
			Status: false,
			Pesan:  "Gagal Koneksi: " + salahe.Error(),
			Data:   []task{},
		}
	}
	defer db.Close()

	dataTsk, salahe := db.Query("select * from task where tasks=?", tasks)
	if salahe != nil {
		return response{
			Status: false,
			Pesan:  "Gagal Query: " + salahe.Error(),
			Data:   []task{},
		}
	}
	defer dataTsk.Close()

	var hasil []task

	for dataTsk.Next() {
		var tsk = task{}
		var salahe = dataTsk.Scan(&tsk.Tasks, &tsk.Assignee, &tsk.Deadline, &tsk.Status)

		if salahe != nil {
			return response{
				Status: false,
				Pesan:  "Gagal Baca: " + salahe.Error(),
				Data:   []task{},
			}
		}

		hasil = append(hasil, tsk)
	}

	salahe = dataTsk.Err()

	if salahe != nil {
		return response{
			Status: false,
			Pesan:  "Kesalahan: " + salahe.Error(),
			Data:   []task{},
		}
	}

	return response{
		Status: true,
		Pesan:  "Berhasil Tampil",
		Data:   hasil,
	}

}

func tambah(tasks string, assignee string, deadline string, status string) response {
	db, salahe := koneksi()
	if salahe != nil {
		return response{
			Status: false,
			Pesan:  "Gagal Koneksi: " + salahe.Error(),
			Data:   []task{},
		}
	}
	defer db.Close()

	_, salahe = db.Exec("insert into task values (?, ?, ?, ?)", tasks, assignee, deadline, status)
	if salahe != nil {
		return response{
			Status: false,
			Pesan:  "Gagal Query Insert: " + salahe.Error(),
			Data:   []task{},
		}
	}
	return response{
		Status: true,
		Pesan:  "Berhasil Tambah",
		Data:   []task{},
	}
}

func ubah(tasks string, assignee string, deadline string, status string) response {
	db, salahe := koneksi()
	if salahe != nil {
		return response{
			Status: false,
			Pesan:  "Gagal Koneksi: " + salahe.Error(),
			Data:   []task{},
		}
	}
	defer db.Close()

	_, salahe = db.Exec("update task set assignee=?, deadline=?, status=? where tasks=?", assignee, deadline, status, tasks)
	if salahe != nil {
		return response{
			Status: false,
			Pesan:  "Gagal Query Update: " + salahe.Error(),
			Data:   []task{},
		}
	}
	return response{
		Status: true,
		Pesan:  "Berhasil Ubah",
		Data:   []task{},
	}
}

func hapus(tasks string) response {
	db, salahe := koneksi()
	if salahe != nil {
		return response{
			Status: false,
			Pesan:  "Gagal Koneksi: " + salahe.Error(),
			Data:   []task{},
		}
	}
	defer db.Close()

	_, salahe = db.Exec("delete from task where tasks=?", tasks)
	if salahe != nil {
		return response{
			Status: false,
			Pesan:  "Gagal Query Delete: " + salahe.Error(),
			Data:   []task{},
		}
	}
	return response{
		Status: true,
		Pesan:  "Berhasil Hapus",
		Data:   []task{},
	}
}

func kontroler(w http.ResponseWriter, r *http.Request) {

	var tampilHtml, salaheTampil = template.ParseFiles("template/tampil.html")
	if salaheTampil != nil {
		fmt.Println(salaheTampil.Error())
		return
	}

	var tambahHtml, salaheTambah = template.ParseFiles("template/tambah.html")
	if salaheTambah != nil {
		fmt.Println(salaheTambah.Error())
		return
	}

	var ubahHtml, salaheUbah = template.ParseFiles("template/ubah.html")
	if salaheUbah != nil {
		fmt.Println(salaheUbah.Error())
		return
	}

	var hapusHtml, salaheHapus = template.ParseFiles("template/hapus.html")
	if salaheHapus != nil {
		fmt.Println(salaheHapus.Error())
		return
	}

	switch r.Method {

	case "GET":

		aksi := r.URL.Query()["aksi"]
		if len(aksi) == 0 {

			tampilHtml.Execute(w, tampil("Berhasil Tampil"))

		} else if aksi[0] == "tambah" {

			tambahHtml.Execute(w, nil)

		} else if aksi[0] == "ubah" {

			tasks := r.URL.Query()["tasks"]
			ubahHtml.Execute(w, getTsk(tasks[0]))

		} else if aksi[0] == "hapus" {
			tasks := r.URL.Query()["tasks"]
			hapusHtml.Execute(w, getTsk(tasks[0]))

		} else {
			tampilHtml.Execute(w, tampil("Berhasil Tampil"))
		}

	case "POST":
		var salahe = r.ParseForm()
		if salahe != nil {
			fmt.Fprintln(w, "Kesalahan: ", salahe)
			return
		}

		var tasks = r.FormValue("tasks")
		var assignee = r.FormValue("assignee")
		var deadline = r.FormValue("deadline")
		var status = r.FormValue("status")

		var aksi = r.URL.Path
		if aksi == "/tambah" {
			var hasil = tambah(tasks, assignee, deadline, status)
			tampilHtml.Execute(w, tampil(hasil.Pesan))

		} else if aksi == "/ubah" {
			var hasil = ubah(tasks, assignee, deadline, status)
			tampilHtml.Execute(w, tampil(hasil.Pesan))
		} else if aksi == "/hapus" {
			var hasil = hapus(tasks)
			tampilHtml.Execute(w, tampil(hasil.Pesan))
		} else {
			tampilHtml.Execute(w, tampil("Berhasil Tampil"))
		}

	default:
		fmt.Fprint(w, "Maaf. Method yang didukung hanya GET dan POST")
	}

}

func main() {

	http.HandleFunc("/", kontroler)

	fmt.Println("Server berjalan di Port 8080...")
	http.ListenAndServe(":8080", nil)
}
