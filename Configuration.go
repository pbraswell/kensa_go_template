package main

type Configuration struct {
	KENSA_TEMPLATE_MYADDON_URL string
}

type AddonResource struct {
	Id     string        `json:"id"`
	Config Configuration `json:"config"`
}
