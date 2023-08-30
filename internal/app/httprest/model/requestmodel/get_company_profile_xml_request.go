package requestmodel

//type: request/responnse
//data: profile
//extermnal type: "ab|participant|pjsppa|du
//company code: kalo mau spesifik
//syarat company code di isi, harus external typenya diisi juga
type CompanyProfileXml struct {
	Type string
	Data string
	ExternalType string
	CompanyCode string
}

