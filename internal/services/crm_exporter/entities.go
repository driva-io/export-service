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

var DrivaTestLead = map[string]any{
	"possui_email":                 true,
	"data_ultima_movimentacao_qsa": "2021-08-30",
	"data_situacao_cadastral":      "2020-01-10",
	"capital_uf":                   false,
	"socios_em_comum": []any{
		map[string]any{
			"cnpjs_com_mesmo_socio":      []int64{9074301000134, 24448877000108, 42436658000190, 26032328000183},
			"socio":                      "PATRICK DE CESAR FRANCISCO | ***302639**",
			"qtde_cnpjs_com_mesmo_socio": 4,
		},
	},
	"emails": []any{
		map[string]any{
			"email":             "FINANCEIRO@DRIVA.COM.BR",
			"nome":              "DRIVA TECNOLOGIA LTDA",
			"pertence_contador": false,
		},
	},
	"telefones": []any{
		map[string]any{
			"ddd_correto":       false,
			"fixo_movel":        "MOVEL",
			"pertence_contador": false,
			"telefone_completo": "41 989042906",
			"validado_discador": false,
			"whatsapp":          false,
		},
		map[string]any{
			"ddd_correto":       true,
			"fixo_movel":        "MOVEL",
			"pertence_contador": false,
			"telefone_completo": "41 985198510",
			"validado_discador": false,
			"whatsapp":          false,
		},
		map[string]any{
			"ddd_correto":       true,
			"fixo_movel":        "MOVEL",
			"pertence_contador": false,
			"telefone_completo": "41 985198510",
			"validado_discador": false,
			"whatsapp":          true,
		},
		map[string]any{
			"ddd_correto":       false,
			"fixo_movel":        "MOVEL",
			"pertence_contador": false,
			"telefone_completo": "41 985198510",
			"validado_discador": false,
			"whatsapp":          false,
		},
		map[string]any{
			"ddd_correto":       false,
			"fixo_movel":        "MOVEL",
			"pertence_contador": false,
			"telefone_completo": "41 997204336",
			"validado_discador": true,
			"whatsapp":          true,
		},
		map[string]any{
			"ddd_correto":       true,
			"fixo_movel":        "MOVEL",
			"pertence_contador": false,
			"telefone_completo": "41 997204336",
			"validado_discador": true,
			"whatsapp":          true,
		},
		map[string]any{
			"ddd_correto":       true,
			"fixo_movel":        "MOVEL",
			"pertence_contador": false,
			"telefone_completo": "41 989042906",
			"validado_discador": false,
			"whatsapp":          false,
		},
		map[string]any{
			"ddd_correto":       false,
			"fixo_movel":        "MOVEL",
			"pertence_contador": false,
			"telefone_completo": "41 985198510",
			"validado_discador": false,
			"whatsapp":          true,
		},
	},
	"cnae_principal_subclasse":      6311900,
	"data_atualizacao":              "2024-10-31",
	"sigla_uf":                      "PR",
	"regiao_saude":                  "2a RS METROPOLITANA",
	"nivel_de_atividade":            "MÉDIA",
	"cnae_principal_desc_subclasse": "6311900 - TRATAMENTO DE DADOS, PROVEDORES DE SERVICOS DE APLICACAO E SERVICOS DE HOSPEDAGEM NA INTERNET",
	"faturamento":                   744120,
	"cod_regiao_intermediaria":      4101,
	"uf_norm":                       "PARANA",
	"movimentacoes_qsa": []any{
		map[string]any{
			"nome_socio":             "PATRICK DE CESAR FRANCISCO",
			"data_nascimento":        576460800000,
			"data_entrada_sociedade": 1578614400000,
			"sexo":                   "M",
		},
		map[string]any{
			"nome_socio":             "WAGNER RODRIGUES ULIAN AGOSTINHO",
			"data_nascimento":        889401600000,
			"data_entrada_sociedade": 1630281600000,
			"sexo":                   "M",
		},
	},
	"empresa_publico_privada": "PRIVADA",
	"qtde_socios":             2,
	"municipio":               "SAO JOSE DOS PINHAIS",
	"bairro":                  "PARQUE DA FONTE",
	"qtde_filiais":            2,
	"possui_pat":              true,
	"microrregiao":            "CURITIBA",
	"possui_site":             true,
	"raiz_cnpj":               35965725,
	"qtde_beneficiarios_pat":  3,
	"capital_social":          110000,
	"cod_municipio_bcb":       14443,
	"cnae_principal_divisao":  63,
	"nome_fantasia":           "DRIVA TECNOLOGIA LTDA",
	"cnae_principal_classe":   63119,
	"nbs": []any{
		map[string]any{
			"nbs_desc": "01.03 - PROCESSAMENTO, ARMAZENAMENTO OU HOSPEDAGEM DE DADOS, TEXTOS, IMAGENS, VIDEOS, PAGINAS ELETRONICAS, APLICATIVOS E SISTEMAS DE INFORMACAO, ENTRE OUTROS FORMATOS, E TRATAMENTO DE DADOS, PROVEDORES DE SERVICOS DE APLICACAO E SERVICOS DE HOSPEDAGEM NA CONGENERES.",
			"nbs":      103,
		},
	},
	"numero":                            211,
	"mesorregiao":                       "METROPOLITANA DE CURITIBA",
	"faixa_funcionarios_grupo":          "50 A 99",
	"lon":                               -49.1697494,
	"qualificacao_responsavel":          "SOCIO-ADMINISTRADOR",
	"forma_de_tributacao":               "LUCRO REAL",
	"data_atualizacao_dados_cadastrais": "2024-10-16",
	"matriz":                            true,
	"cep":                               83050610,
	"regiao_intermediaria":              "CURITIBA",
	"cod_municipio_ibge_6":              412550,
	"cnae_principal_grupo":              631,
	"faixa_funcionarios":                "02 A 05",
	"regiao_imediata":                   "CURITIBA",
	"porte":                             "MICRO EMPRESA",
	"opcao_pelo_mei":                    false,
	"forma_tributacao": []any{
		map[string]any{
			"ano":              2022,
			"forma_tributacao": "LUCRO REAL",
		},
	},
	"possui_linkedin":           true,
	"cod_municipio_ibge":        4125506,
	"ordem_cnpj":                1,
	"motivo_situacao_cadastral": "SEM MOTIVO",
	"coords":                    "-25.496327,-49.169749",
	"possui_socio_jovem":        true,
	"opcao_pelo_simples":        false,
	"ufiso31662":                "BR-PR",
	"possui_facebook":           false,
	"segmento":                  "SERVICOS",
	"sucessoes_qsa":             []any{},
	"populacao_municipio":       334620,
	"macrorregiao":              "SUL",
	"logradouro":                "LAURA NUNES FERNANDES",
	"municipio_norm":            "SAO JOSE DOS PINHAIS",
	"logradouro_norm":           "RUA LAURA NUNES FERNANDES",
	"ramo_de_atividade":         "OUTROS SERVICOS",
	"cod_municipio_rfb":         7885,
	"cod_uf":                    41,
	"dv_cnpj":                   7,
	"cnpj":                      35965725000107,
	"data_exclusao_simples":     "2021-12-31",
	"fonte_coords":              "google",
	"possui_instagram":          false,
	"uf":                        "PARANA",
	"faixa_faturamento_grupo":   "7M A 10M",
	"crescimento_por_ano":       33.0,
	"cod_regiao_saude":          41002,
	"razao_social":              "DRIVA TECNOLOGIA LTDA",
	"sigla_uf_norm":             "PR",
	"parentes_qsa":              []any{},
}
