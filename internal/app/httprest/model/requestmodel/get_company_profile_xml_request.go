package requestmodel

// type: request/responnse
// data: profile
// extermnal type: "ab|participant|pjsppa|du
// company code: kalo mau spesifik
// syarat company code di isi, harus external typenya diisi juga
type CompanyProfileXml struct {
	Type         string `xml:"Type"`
	Data         string `xml:"Data"`
	ExternalType string `xml:"ExternalType"`
	CompanyCode  string `xml:"CompanyCode"`
}
