package trace

type Domain struct {
	Ip          string    `json:"ip"`
	Name        string    `json:"name"`
	Created     string    `json:"created"`
	Expires     string    `json:"expires"`
	Changed     string    `json:"changed"`
	IdnName     string    `json:"idn_name"`
	AskWhois    string    `json:"ask_whois"`
	Contacts    Contacts  `json:"contacts"`
	Registrar   Registrar `json:"registrar"`
	Nameservers []string  `json:"nameservers"`
}

type Contacts struct {
	Tech  Contact `json:"tech"`
	Admin Contact `json:"admin"`
	Owner Contact `json:"owner"`
}

type Contact struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	Country      string `json:"country"`
	Address      string `json:"address"`
	Organization string `json:"organization"`
}

type Registrar struct {
	Url   string `json:"url"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
