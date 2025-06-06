package api

func (p *SignupParams) ConfigureDefaultsYuzha() {
	if p.Data == nil {
		p.Data = make(map[string]interface{})
		p.Data["user_language"] = "en"
		p.Data["country"] = "US"
	} else {
		_, ok := p.Data["user_language"]
		if !ok {
			p.Data["user_language"] = "en"
		}
		_, ok = p.Data["country"]
		if !ok {
			p.Data["country"] = "US"
		}
	}
}
