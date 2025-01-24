package crm_exporter

type Stage struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type Pipeline struct {
	Id     string  `json:"id,omitempty"`
	Name   string  `json:"name,omitempty"`
	Stages []Stage `json:"stages,omitempty"`
}

type FieldOptions struct {
	Id    string `json:"id,omitempty"`
	Label string `json:"label,omitempty"`
}

type CrmField struct {
	Id       interface{}     `json:"id"` // Can be string or number
	Label    string          `json:"label"`
	Type     string          `json:"type"`
	Options  *[]FieldOptions `json:"options,omitempty"`
	Required *bool           `json:"required,omitempty"`
}

type CrmFields struct {
	Deals     *[]CrmField `json:"deals,omitempty"`
	Companies *[]CrmField `json:"companies,omitempty"`
	Contacts  *[]CrmField `json:"contacts,omitempty"`
}

type Owner struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

var DrivaTestLead = map[string]interface{}{
	"address":                   "Rua Laura Nunes Fernandes 211",
	"bairro":                    "PARQUE DA FONTE",
	"banner_url":                "https://media.licdn.com/dms/image/v2/D4D3DAQG3SSw49BniiA/image-scale_127_750/image-scale_127_750/0/1693338849355/driva_tech_cover?e=1738342800&v=beta&t=KY6Bmq-XFUVcvzCMBSGCHg7QjsCG_lCdZRS9IfEaXQQ",
	"capital_social":            110000,
	"cep":                       83050610,
	"city":                      "São José dos Pinhais",
	"cnpj":                      35965725000107,
	"company_contact_id":        "24e3bf31-2695-4e13-b70c-992aeae0be61",
	"complemento":               "",
	"country":                   "BR",
	"data_inicio_atividade":     "2020-01-10",
	"descricao_tipo_logradouro": "RUA",
	"description":               "Driva é a solução de Inteligência de Mercado perfeita que otimiza todos os aspectos da sua jornada de vendas. Mapeie seu mercado, descubra as melhores oportunidades e conecte-se com seus clientes com uma tecnologia inovadora de big data e inteligência de mercado.",
	"empresa_publico_privada":   "PRIVADA",
	"endereco":                  "RUA LAURA NUNES FERNANDES 211 - PARQUE DA FONTE - 83050610",
	"follower_count":            7276,
	"founded_on":                2020,
	"funding_type":              "",
	"funding_value":             "",
	"id":                        52163831,
	"industries":                "Atividades dos serviços de tecnologia da informação",
	"logradouro":                "LAURA NUNES FERNANDES",
	"matriz":                    true,
	"name":                      "Driva",
	"natureza_juridica":         "SOCIEDADE EMPRESARIA LIMITADA",
	"nome_fantasia":             "DRIVA TECNOLOGIA LTDA",
	"numero":                    211,
	"opcao_pelo_mei":            false,
	"opcao_pelo_simples":        false,
	"operates_in_brazil":        true,
	"phone_number":              "(41) 985198510",
	"phones": []interface{}{
		map[string]interface{}{
			"has_whatsapp":      false,
			"phone":             "41 989042906",
			"type":              "MOVEL",
			"validado_discador": false,
			"validation":        "NUMERO_MUDOU",
		},
		map[string]interface{}{
			"has_whatsapp":      true,
			"phone":             "41 997204336",
			"type":              "MOVEL",
			"validado_discador": true,
			"validation":        "ATENDIDA",
		},
	},
	"picture_url": "https://media.licdn.com/dms/image/v2/D4D0BAQEwbXBBAHrajg/company-logo_200_200/company-logo_200_200/0/1691499793056/driva_tech_logo?e=1746057600&v=beta&t=3HU8jo022BPx9kIt4mCZTisHURLjY7jgbJfzqdUFhNI",
	"porte":       "MICRO EMPRESA",
	"postal_code": "Curitiba, PR",
	"profiles": []interface{}{
		map[string]interface{}{
			"profile_contact_id": "f0810792-30f0-4e9f-981b-2095cf073933",
			"company_contact_id": "24e3bf31-2695-4e13-b70c-992aeae0be61",
			"id":                 "ACwAABgXfPYBAqk7BZbjzGe6pYy0ScA1NktVmwM",
			"area":               "",
			"headline":           "",
			"location":           "Santa Catarina, Brazil",
			"name":               "Elimar Sanches Kauffmann",
			"picture_url":        "https://media.licdn.com/dms/image/v2/D4D03AQGd8Sw5T-kDPg/profile-displayphoto-shrink_200_200/profile-displayphoto-shrink_200_200/0/1675382565872?e=1743033600&v=beta&t=n1ZoUqxKtlj3qBGLwjRDhk4x8_g3DcQDrVPH_6czbUw",
			"profile_url":        "linkedin.com/in/ACwAABgXfPYBAqk7BZbjzGe6pYy0ScA1NktVmwM",
			"role":               "Business Development Manager LATAM",
			"seniority":          "GERENTE",
			"emails": []interface{}{
				map[string]interface{}{
					"email":      "elimar@datadriva.com",
					"validation": "ENTREGAVEL",
				},
			},
		},
		map[string]interface{}{
			"profile_contact_id": "0bf3bb87-1c00-47ea-ac3b-a91b66d6679b",
			"company_contact_id": "24e3bf31-2695-4e13-b70c-992aeae0be61",
			"id":                 "ACwAABPJzKcBKpjiLedI4t6DutjJ5xOh8f6wnBQ",
			"area":               "",
			"headline":           "",
			"location":           "Curitiba, Paraná, Brazil",
			"name":               "Gabriel Galvão",
			"picture_url":        "https://media.licdn.com/dms/image/v2/C4D03AQGpJih1-IUNlA/profile-displayphoto-shrink_200_200/profile-displayphoto-shrink_200_200/0/1649304531050?e=1743033600&v=beta&t=UmM0W0OxU-A0zmeRiR_3ET-ldRQatbGsEeGWJtcajyE",
			"profile_url":        "linkedin.com/in/ACwAABPJzKcBKpjiLedI4t6DutjJ5xOh8f6wnBQ",
			"role":               "COO",
			"seniority":          "C-SUITE / DIRETOR",
			"emails": []interface{}{
				map[string]interface{}{
					"email":      "gabriel.galvao@driva.com.br",
					"validation": "ENTREGAVEL",
				},
			},
		},
	},
	"profiles_count":           0,
	"public_id":                "driva-tech",
	"qualificacao_responsavel": "SOCIO-ADMINISTRADOR",
	"raiz_cnpj":                35965725,
	"razao_social":             "DRIVA TECNOLOGIA LTDA",
	"situacao_cadastral":       "ATIVA",
	"specialities":             "",
	"staff_count":              78,
	"staff_count_range":        "51-200",
	"state":                    "Paraná",
	"tagline":                  "Venda mais e melhor com a empresa de Inteligência Comercial que mais cresce no Brasil!",
	"url":                      "linkedin.com/company/driva-tech",
	"website":                  "driva.io",
}
