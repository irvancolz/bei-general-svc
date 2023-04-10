package httpresponse

const (
	ERR_GENERAL_400     string = "kesalahan dalam sistem"
	ERR_REQUESTBODY_400 string = "kesalahan pada request"

	DATANOTFOUND_400    string = "data tidak ditemukan"
	CONTENTNOTFOUND_404 string = "konten tidak ditemukan"
	WRONGPASSWORD_400   string = "password salah"

	CREATESUCCESS_200   string = "sukses menambahkan data"
	CREATEFAILED_400    string = "gagal menambahkan data"
	CREATEDUPLICATE_400 string = "data sudah ada"

	READSUCCESS_200 string = "sukses menampilkan data"
	READFAILED_400  string = "gagal menampilkan data"

	UPDATESUCCESS_200 string = "sukses mengubah data"
	UPDATEFAILED_400  string = "gagal mengubah data"

	DELETESUCCESS_200 string = "sukses menghapus data"
	DELETEFAILED_400  string = "gagal menghapus data"

	UPLOADSUCCESS_200 string = "sukses mengunggah data"
	UPLOADFAILED_400  string = "gagal mengunggah data"

	DOWNLOADSUCCESS_200 string = "sukses mengunduh data"
	DOWNLOADFAILED_400  string = "gagal mengunduh data"

	AUTHISLOGIN_400 string = "user sedang login di device lain"
)
