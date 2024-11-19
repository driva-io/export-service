package crm_exporter

import (
	"context"
	"encoding/json"
	"errors"
	"export-service/internal/core/ports"
	"export-service/internal/repositories/crm_company_repo"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/belong-inc/go-hubspot"
	"github.com/gofiber/fiber/v2"
)

type HubspotService struct {
	companyRepo *crm_company_repo.PgCrmCompanyRepository
}

func NewHubspotService(companyRepo *crm_company_repo.PgCrmCompanyRepository) *HubspotService {
	return &HubspotService{companyRepo: companyRepo}
}

func searchForExistingCompany(client *hubspot.Client, fields ...map[string]any) (any, error) {
	var filters []map[string][]map[string]any
	for _, field := range fields {
		var isNil bool
		var key string
		var value any
		for k, v := range field {
			if v == nil {
				isNil = true
				break
			}
			key, value = k, v
			break
		}
		if isNil {
			continue
		}
		filters = append(filters, map[string][]map[string]any{
			"filters": {
				{
					"propertyName": key,
					"operator":     "EQ",
					"value":        value,
				},
			},
		})
	}

	if filters == nil {
		return nil, nil
	}

	url := "https://api.hubapi.com/crm/v3/objects/companies/search"
	body := map[string]any{
		"filterGroups": filters,
	}

	var res any
	err := client.Post(url, body, &res)
	if err != nil {
		return nil, err
	}

	if res.(map[string]any)["total"].(float64) == 0 {
		return nil, nil
	}

	return res.(map[string]any)["results"].([]any)[0], nil
}

func searchForExistingDeal(client *hubspot.Client, fields ...map[string]any) (any, error) {
	var filters []map[string][]map[string]any
	for _, field := range fields {
		var isNil bool
		var key string
		var value any
		for k, v := range field {
			if v == nil {
				isNil = true
				break
			}
			key, value = k, v
			break
		}
		if isNil {
			continue
		}
		filters = append(filters, map[string][]map[string]any{
			"filters": {
				{
					"propertyName": key,
					"operator":     "EQ",
					"value":        value,
				},
			},
		})
	}

	if filters == nil {
		return nil, nil
	}

	url := "https://api.hubapi.com/crm/v3/objects/deals/search"
	body := map[string]any{
		"filterGroups": filters,
	}

	var res any
	err := client.Post(url, body, &res)
	if err != nil {
		return nil, err
	}

	if res.(map[string]any)["total"].(float64) == 0 {
		return nil, nil
	}

	return res.(map[string]any)["results"].([]any)[0], nil
}

func searchForExistingContact(client *hubspot.Client, fields ...map[string]any) (any, error) {
	var filters []map[string][]map[string]any
	for _, field := range fields {
		var isNil bool
		var key string
		var value any
		for k, v := range field {
			if v == nil {
				isNil = true
				break
			}
			key, value = k, v
			break
		}
		if isNil {
			continue
		}
		filters = append(filters, map[string][]map[string]any{
			"filters": {
				{
					"propertyName": key,
					"operator":     "EQ",
					"value":        value,
				},
			},
		})
	}

	if filters == nil {
		return nil, nil
	}

	url := "https://api.hubapi.com/crm/v3/objects/contacts/search"
	body := map[string]any{
		"filterGroups": filters,
	}

	var res any
	err := client.Post(url, body, &res)
	if err != nil {
		return nil, err
	}

	if res.(map[string]any)["total"].(float64) == 0 {
		return nil, nil
	}

	return res.(map[string]any)["results"].([]any)[0], nil
}

func sendCompany(client *hubspot.Client, mappedCompanyData map[string]any) (ObjectStatus, error) {
	companyEntity, exists := mappedCompanyData["entity"]
	if !exists {
		return ObjectStatus{}, errors.New("company entity not found in mapped company data")
	}
	companyEntityMap, isMap := companyEntity.(map[string]any)
	if !isMap {
		return ObjectStatus{}, errors.New("company entity is not a map")
	}

	existingCompany, err := searchForExistingCompany(client, map[string]any{"cnpj_empresa": companyEntityMap["cnpj_empresa"]})
	if err != nil {
		return ObjectStatus{}, err
	}

	var company *hubspot.ResponseResource
	var status Status
	var message string
	if existingCompany != nil {
		updatedCompany, err := client.CRM.Company.Update(existingCompany.(map[string]any)["id"].(string), companyEntityMap)
		if err != nil {
			return ObjectStatus{}, err
		}
		status = Updated
		company = updatedCompany
	} else {
		createdCompany, err := client.CRM.Company.Create(companyEntityMap)
		if err != nil {
			return ObjectStatus{}, err
		}
		status = Created
		company = createdCompany
	}

	return ObjectStatus{
		Id:      company.ID,
		Status:  status,
		Message: message,
	}, nil
}

func sendDeal(client *hubspot.Client, mappedDealData map[string]any) (ObjectStatus, error) {
	dealEntity, exists := mappedDealData["entity"]
	if !exists {
		return ObjectStatus{}, errors.New("deal entity not found in mapped deal data")
	}
	dealEntityMap, isMap := dealEntity.(map[string]any)
	if !isMap {
		return ObjectStatus{}, errors.New("deal entity is not a map")
	}

	existingDeal, err := searchForExistingDeal(client, map[string]any{"dealname": dealEntityMap["dealname"]})
	if err != nil {
		return ObjectStatus{}, err
	}

	var deal *hubspot.ResponseResource
	var status Status
	var message string
	if existingDeal != nil {
		updatedDeal, err := client.CRM.Deal.Update(existingDeal.(map[string]any)["id"].(string), dealEntityMap)
		if err != nil {
			return ObjectStatus{}, err
		}
		status = Updated
		deal = updatedDeal
	} else {
		createdDeal, err := client.CRM.Deal.Create(dealEntityMap)
		if err != nil {
			return ObjectStatus{}, err
		}
		status = Created
		deal = createdDeal
	}

	return ObjectStatus{
		Id:      deal.ID,
		Status:  status,
		Message: message,
	}, nil
}

func sendContact(client *hubspot.Client, mappedContactData map[string]any) (ObjectStatus, error) {
	contactEntity, exists := mappedContactData["entity"]
	if !exists {
		return ObjectStatus{}, errors.New("contact entity not found in mapped contact data")
	}
	contactEntityMap, isMap := contactEntity.(map[string]any)
	if !isMap {
		return ObjectStatus{}, errors.New("contact entity is not a map")
	}

	existingContact, err := searchForExistingContact(client, map[string]any{"firstname": contactEntityMap["firstname"]})
	if err != nil {
		return ObjectStatus{}, err
	}

	var contact *hubspot.ResponseResource
	var status Status
	var message string
	if existingContact != nil {
		updatedContact, err := client.CRM.Contact.Update(existingContact.(map[string]any)["id"].(string), contactEntityMap)
		if err != nil {
			return ObjectStatus{}, err
		}
		status = Updated
		contact = updatedContact
	} else {
		createdContact, err := client.CRM.Contact.Create(contactEntityMap)
		if err != nil {
			return ObjectStatus{}, err
		}
		status = Created
		contact = createdContact
	}

	return ObjectStatus{
		Id:      contact.ID,
		Status:  status,
		Message: message,
	}, nil
}

func createLeadAssociations(client *hubspot.Client, lead CreatedLead) error {

	if lead.company != nil && lead.deal != nil {
		_, err := client.CRM.Deal.AssociateAnotherObj(lead.deal.Id.(string), &hubspot.AssociationConfig{ToObject: hubspot.ObjectTypeCompany, ToObjectID: lead.company.Id.(string), Type: hubspot.AssociationTypeDealToCompany})
		if err != nil {
			return err
		}
	}

	if lead.deal != nil && lead.contacts != nil {
		for _, contact := range *lead.contacts {
			_, err := client.CRM.Deal.AssociateAnotherObj(lead.deal.Id.(string), &hubspot.AssociationConfig{ToObject: hubspot.ObjectTypeContact, ToObjectID: contact.Id.(string), Type: hubspot.AssociationTypeDealToContact})
			if err != nil {
				return err
			}
		}
	}

	if lead.company != nil && lead.contacts != nil {
		for _, contact := range *lead.contacts {
			_, err := client.CRM.Company.AssociateAnotherObj(lead.company.Id.(string), &hubspot.AssociationConfig{ToObject: hubspot.ObjectTypeContact, ToObjectID: contact.Id.(string), Type: hubspot.AssociationTypeCompanyToContact})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (h HubspotService) SendLead(client any, mappedLead map[string]any) (CreatedLead, error) {
	husbpotClient := client.(*hubspot.Client)
	lead := CreatedLead{
		company:  nil,
		contacts: nil,
		deal:     nil,
		other:    nil,
	}

	company, exists := mappedLead["company"]
	if exists {
		companyData, isMap := company.(map[string]any)
		if !isMap {
			return CreatedLead{}, errors.New("invalid company data to send crm. must be a map")
		}
		sentCompany, err := sendCompany(husbpotClient, companyData)
		if err != nil {
			return CreatedLead{}, err
		}
		lead.company = &sentCompany
	}

	deal, exists := mappedLead["deal"]
	if exists {
		dealData, isMap := deal.(map[string]any)
		if !isMap {
			return CreatedLead{}, errors.New("invalid deal data to send to crm. must be a map")
		}
		sentDeal, err := sendDeal(husbpotClient, dealData)
		if err != nil {
			return CreatedLead{}, err
		}
		lead.deal = &sentDeal
	}

	contact, exists := mappedLead["contact"]
	if exists {
		contactData, isMap := contact.(map[string]any)
		if !isMap {
			return CreatedLead{}, errors.New("invalid contact data to send to crm. must be a map")
		}
		sentContact, err := sendContact(husbpotClient, contactData)
		if err != nil {
			return CreatedLead{}, err
		}
		if lead.contacts == nil {
			lead.contacts = &[]ObjectStatus{}
		}
		*lead.contacts = append(*lead.contacts, sentContact)
	}

	contacts, exists := mappedLead["contacts"]
	if exists {
		contactsData, isArray := contacts.([]map[string]any)
		if !isArray {
			return CreatedLead{}, errors.New("invalid contacts data to send to crm. must be an array of maps")
		}
		for _, contact := range contactsData {
			sentContact, err := sendContact(husbpotClient, contact)
			if err != nil {
				return CreatedLead{}, err
			}
			if lead.contacts == nil {
				lead.contacts = &[]ObjectStatus{}
			}
			*lead.contacts = append(*lead.contacts, sentContact)
		}
	}

	err := createLeadAssociations(husbpotClient, lead)
	if err != nil {
		return CreatedLead{}, err
	}

	return lead, nil
}

func (h HubspotService) GetPipelines(client any) ([]Pipeline, error) {
	url := "https://api.hubapi.com/crm/v3/pipelines/deals"

	var res any
	err := client.(*hubspot.Client).Get(url, &res, nil)
	if err != nil {
		return nil, err
	}

	var pipelines []Pipeline

	for _, value := range res.(map[string]any)["results"].([]any) {
		var stages []Stage
		for _, stageValue := range value.(map[string]any)["stages"].([]any) {

			pipelineStage := Stage{
				Id:   stageValue.(map[string]any)["id"].(string),
				Name: stageValue.(map[string]any)["label"].(string),
			}
			stages = append(stages, pipelineStage)

		}

		pipeline := Pipeline{
			Id:     value.(map[string]any)["id"].(string),
			Name:   value.(map[string]any)["label"].(string),
			Stages: stages,
		}

		pipelines = append(pipelines, pipeline)
	}

	return pipelines, nil
}

func buildFields(fields *hubspot.CrmPropertiesList) []CrmField {
	var builtFields []CrmField
	for _, value := range fields.Results {
		if strings.HasPrefix(value.Name.String(), "hs_") {
			continue
		}

		var fieldOptions []FieldOptions
		for _, optionValue := range value.Options {
			fieldOption := FieldOptions{
				Id:    optionValue.Value.String(),
				Label: optionValue.Label.String(),
			}

			fieldOptions = append(fieldOptions, fieldOption)
		}

		crmField := CrmField{
			Id:      value.Name.String(),
			Label:   value.Label.String(),
			Type:    value.Type.String(),
			Options: &fieldOptions,
		}
		builtFields = append(builtFields, crmField)
	}

	return builtFields
}

func (h HubspotService) GetFields(client any) (CrmFields, error) {
	dealFields, err := client.(*hubspot.Client).CRM.Properties.List("deals")
	if err != nil {
		return CrmFields{}, err
	}
	companyFields, err := client.(*hubspot.Client).CRM.Properties.List("companies")
	if err != nil {
		return CrmFields{}, err
	}
	contactFields, err := client.(*hubspot.Client).CRM.Properties.List("contacts")
	if err != nil {
		return CrmFields{}, err
	}

	builtDealFields := buildFields(dealFields)
	builtCompanyFields := buildFields(companyFields)
	builtContactFields := buildFields(contactFields)

	return CrmFields{
		Deals:     &builtDealFields,
		Companies: &builtCompanyFields,
		Contacts:  &builtContactFields,
	}, nil
}

func (h HubspotService) GetOwners(client any) ([]Owner, error) {
	url := "https://api.hubapi.com/crm/v3/owners"
	options := struct {
		After    *string `url:"after,omitempty"`
		Limit    int     `url:"limit"`
		Archived bool    `url:"archived"`
	}{
		After:    nil,
		Limit:    100,
		Archived: false,
	}

	var res any
	err := client.(*hubspot.Client).Get(url, &res, options)
	if err != nil {
		return nil, err
	}

	//TODO: pagination
	var owners []Owner
	for _, value := range res.(map[string]any)["results"].([]any) {
		owner := Owner{
			Id:   value.(map[string]any)["id"].(string),
			Name: value.(map[string]any)["firstName"].(string) + " " + value.(map[string]any)["lastName"].(string),
		}
		owners = append(owners, owner)
	}

	return owners, nil
}

func (h HubspotService) Authorize(ctx context.Context, companyName string) (any, error) {

	company, err := h.companyRepo.Get(ctx, ports.CrmCompanyQueryParams{Crm: "hubspot", Company: companyName})
	if err != nil {
		return nil, err
	}

	client, _ := hubspot.NewClient(hubspot.SetOAuth(&hubspot.OAuthConfig{
		GrantType:    hubspot.GrantTypeRefreshToken,
		ClientID:     os.Getenv("HUBSPOT_CLIENT_ID"),
		ClientSecret: os.Getenv("HUBSPOT_CLIENT_SECRET"),
		RefreshToken: company.RefreshToken.String,
	}))

	log.Printf("Authenticated hubspot for company %s - workspaceId: %s", company.Name.String, company.WorkspaceId.String)

	return client, nil
}

func (h HubspotService) Install(installData any) (any, error) {
	baseURL := "https://app.hubspot.com/oauth/authorize"
	clientID := url.QueryEscape(os.Getenv("HUBSPOT_CLIENT_ID"))
	scope := url.QueryEscape(os.Getenv("HUBSPOT_SCOPE"))
	redirectURI := os.Getenv("HUBSPOT_REDIRECT_URI")
	installDataMap := installData.(map[string]any)
	state := fmt.Sprintf("%s|%s|%s", installDataMap["workspace_id"], installDataMap["user_id"], installDataMap["company"])

	authURL := fmt.Sprintf("%s?client_id=%s&scope=%s&redirect_uri=%s&state=%s", baseURL, clientID, scope, url.QueryEscape(redirectURI), url.QueryEscape(state))

	return map[string]string{"url": authURL}, nil
}

func (h HubspotService) OAuthCallback(c *fiber.Ctx, params ...any) (any, error) {
	if len(params) != 3 {
		return nil, errors.New("expected 4 parms in oauth callback")
	}

	hubspotCode := c.Query("code")
	workspaceId := params[0].(string)
	userId := params[1].(string)
	company := params[2].(string)

	formData := url.Values{}
	formData.Set("grant_type", "authorization_code")
	formData.Set("client_id", os.Getenv("HUBSPOT_CLIENT_ID"))
	formData.Set("client_secret", os.Getenv("HUBSPOT_CLIENT_SECRET"))
	formData.Set("redirect_uri", os.Getenv("HUBSPOT_REDIRECT_URI"))
	formData.Set("code", hubspotCode)

	resp, err := http.PostForm("https://api.hubapi.com/oauth/v1/token", formData)
	if err != nil {
		log.Fatalf("Error making request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response: %v", err)
	}

	var responseData map[string]any
	if err := json.Unmarshal(body, &responseData); err != nil {
		log.Fatalf("Error unmarshaling JSON: %v", err)
	}

	h.companyRepo.AddHubspot(context.Background(), ports.CrmAddHubspotCompanyQueryParams{Company: company, WorkspaceId: workspaceId, UserId: userId, RefreshToken: responseData["refresh_token"].(string), AccessToken: responseData["access_token"].(string)})

	return nil, nil
}

// func refreshToken(ctx context.Context, company any) (string, error) {

// 	client, _ := hubspot.NewClient(hubspot.SetOAuth(&hubspot.OAuthConfig{
// 		GrantType:    hubspot.GrantTypeRefreshToken,
// 		ClientID:     os.Getenv("HUBSPOT_CLIENT_ID"),
// 		ClientSecret: os.Getenv("HUBSPOT_CLIENT_SECRET"),
// 		RefreshToken: company["refresh_token"],
// 	}))

// 	company.RefreshToken = refreshToken
// 	company.AccessToken = accessToken
// 	company.ExpiresIn = fmt.Sprintf("%d", expiresIn)
// 	company.RefreshedAt = time.Now().UTC().Format(time.RFC3339)

// 	if err := companyRepo.Update("hubspot", company.Name, *company); err != nil {
// 		return "", fmt.Errorf("failed to update company: %w", err)
// 	}

// 	return accessToken, nil
// }
