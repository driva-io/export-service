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

func searchForExistingObject(client *hubspot.Client, objectType string, filtersMap map[string]any) (any, error) {
	var filters []map[string][]map[string]any

	for key, value := range filtersMap {
		if value == nil {
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

	if len(filters) == 0 {
		return nil, nil
	}

	url := "https://api.hubapi.com/crm/v3/objects/" + objectType + "/search"
	body := map[string]any{
		"filterGroups": filters,
	}

	var res any
	err := client.Post(url, body, &res)
	if err != nil {
		return nil, err
	}

	resMap, ok := res.(map[string]any)
	if !ok {
		return nil, errors.New("unexpected response format")
	}

	total, ok := resMap["total"].(float64)
	if !ok || total == 0 {
		return nil, nil
	}

	results, ok := resMap["results"].([]any)
	if !ok || len(results) == 0 {
		return nil, nil
	}

	return results[0], nil
}

func sendCompany(client *hubspot.Client, mappedCompanyData map[string]any, ownerId string) (ObjectStatus, error) {
	companyEntity, exists := mappedCompanyData["entity"]
	if !exists {
		return ObjectStatus{
			Status:  Failed,
			Message: "company entity not found in mapped company data",
		}, errors.New("company entity not found in mapped company data")
	}
	companyEntityMap, isMap := companyEntity.(map[string]any)
	if !isMap {
		return ObjectStatus{
			Status:  Failed,
			Message: "company entity is not a map",
		}, errors.New("company entity is not a map")
	}

	companyEntityMap["hubspot_owner_id"] = ownerId

	existingCompany, err := searchForExistingObject(client, "companies", map[string]any{"name": companyEntityMap["name"]})
	if err != nil {
		return ObjectStatus{
			Status:  Failed,
			Message: err.Error(),
		}, err
	}

	var company *hubspot.ResponseResource
	var status Status
	var message string
	if existingCompany != nil {
		updatedCompany, err := client.CRM.Company.Update(existingCompany.(map[string]any)["id"].(string), companyEntityMap)
		if err != nil {
			return ObjectStatus{
				Status:  Failed,
				Message: err.Error(),
			}, err
		}
		status = Updated
		company = updatedCompany
		message = "Searched fields: cnpj"
	} else {
		createdCompany, err := client.CRM.Company.Create(companyEntityMap)
		if err != nil {
			return ObjectStatus{
				Status:  Failed,
				Message: err.Error(),
			}, err
		}
		status = Created
		company = createdCompany
	}

	return ObjectStatus{
		CrmId:   company.ID,
		Status:  status,
		Message: message,
	}, nil
}

func sendDeal(client *hubspot.Client, mappedDealData map[string]any, ownerId string, pipelineId string, stageId string) (ObjectStatus, error) {
	dealEntity, exists := mappedDealData["entity"]
	if !exists {
		return ObjectStatus{
			Status:  Failed,
			Message: "deal entity not found in mapped deal data",
		}, errors.New("deal entity not found in mapped deal data")
	}
	dealEntityMap, isMap := dealEntity.(map[string]any)
	if !isMap {
		return ObjectStatus{
			Status:  Failed,
			Message: "deal entity is not a map",
		}, errors.New("deal entity is not a map")
	}

	dealEntityMap["pipeline"] = pipelineId
	dealEntityMap["dealstage"] = stageId
	dealEntityMap["hubspot_owner_id"] = ownerId

	existingDeal, err := searchForExistingObject(client, "deals", map[string]any{"dealname": dealEntityMap["dealname"]})
	if err != nil {
		return ObjectStatus{
			Status:  Failed,
			Message: err.Error(),
		}, err
	}

	var deal *hubspot.ResponseResource
	var status Status
	var message string
	if existingDeal != nil {
		updatedDeal, err := client.CRM.Deal.Update(existingDeal.(map[string]any)["id"].(string), dealEntityMap)
		if err != nil {
			return ObjectStatus{
				Status:  Failed,
				Message: err.Error(),
			}, err
		}
		status = Updated
		deal = updatedDeal
		message = "Searched fields: dealname"
	} else {
		createdDeal, err := client.CRM.Deal.Create(dealEntityMap)
		if err != nil {
			return ObjectStatus{
				Status:  Failed,
				Message: err.Error(),
			}, err
		}
		status = Created
		deal = createdDeal
	}

	return ObjectStatus{
		CrmId:   deal.ID,
		Status:  status,
		Message: message,
	}, nil
}

func sendContact(client *hubspot.Client, mappedContactData map[string]any, ownerId string) (ObjectStatus, error) {
	contactEntity, exists := mappedContactData["entity"]
	if !exists {
		return ObjectStatus{
			Status:  Failed,
			Message: "contact entity not found in mapped contact data",
		}, errors.New("contact entity not found in mapped contact data")
	}
	contactEntityMap, isMap := contactEntity.(map[string]any)
	if !isMap {
		return ObjectStatus{
			Status:  Failed,
			Message: "contact entity is not a map",
		}, errors.New("contact entity is not a map")
	}

	contactEntityMap["hubspot_owner_id"] = ownerId

	existingContact, err := searchForExistingObject(client, "contacts", map[string]any{"email": contactEntityMap["email"]})
	if err != nil {
		return ObjectStatus{
			Status:  Failed,
			Message: err.Error(),
		}, err
	}

	var contact *hubspot.ResponseResource
	var status Status
	var message string
	if existingContact != nil {
		updatedContact, err := client.CRM.Contact.Update(existingContact.(map[string]any)["id"].(string), contactEntityMap)
		if err != nil {
			return ObjectStatus{
				Status:  Failed,
				Message: err.Error(),
			}, err
		}
		status = Updated
		contact = updatedContact
		message = contactEntityMap["email"].(string)
	} else {
		createdContact, err := client.CRM.Contact.Create(contactEntityMap)
		if err != nil {
			return ObjectStatus{
				Status:  Failed,
				Message: err.Error(),
			}, err
		}
		status = Created
		contact = createdContact
		email, exists := contactEntityMap["email"]
		if exists {
			message = email.(string)
		} else {
			message = contactEntityMap["phone"].(string)
		}
	}

	return ObjectStatus{
		CrmId:   contact.ID,
		Status:  status,
		Message: message,
	}, nil
}

func createLeadAssociations(client *hubspot.Client, lead CreatedLead) error {

	if lead.Company != nil && lead.Deal != nil {
		_, err := client.CRM.Deal.AssociateAnotherObj(lead.Deal.CrmId.(string), &hubspot.AssociationConfig{ToObject: hubspot.ObjectTypeCompany, ToObjectID: lead.Company.CrmId.(string), Type: hubspot.AssociationTypeDealToCompany})
		if err != nil {
			return err
		}
	}

	if lead.Deal != nil && lead.Contacts != nil {
		for _, contact := range *lead.Contacts {
			_, err := client.CRM.Deal.AssociateAnotherObj(lead.Deal.CrmId.(string), &hubspot.AssociationConfig{ToObject: hubspot.ObjectTypeContact, ToObjectID: contact.CrmId.(string), Type: hubspot.AssociationTypeDealToContact})
			if err != nil {
				return err
			}
		}
	}

	if lead.Company != nil && lead.Contacts != nil {
		for _, contact := range *lead.Contacts {
			_, err := client.CRM.Company.AssociateAnotherObj(lead.Company.CrmId.(string), &hubspot.AssociationConfig{ToObject: hubspot.ObjectTypeContact, ToObjectID: contact.CrmId.(string), Type: hubspot.AssociationTypeCompanyToContact})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (h HubspotService) SendLead(
	client any,
	mappedLead map[string]any,
	correspondingRawData map[string]any,
	configs map[string]any,
	existingLead map[string]any,
) (CreatedLead, error) {
	husbpotClient, ok := client.(*hubspot.Client)
	if !ok {
		return CreatedLead{}, errors.New("invalid HubSpot client")
	}

	ownerId, err := getConfigValue[string](configs, "owner_id")
	if err != nil {
		return CreatedLead{}, err
	}

	stageId, err := getConfigValue[string](configs, "stage_id")
	if err != nil {
		return CreatedLead{}, err
	}

	pipelineId, err := getConfigValue[string](configs, "pipeline_id")
	if err != nil {
		return CreatedLead{}, err
	}

	createDeal := configs["create_deal"].(bool)
	lead := CreatedLead{}

	if company, exists := mappedLead["company"]; exists {
		companyStatus, err := processCompany(husbpotClient, company, existingLead, correspondingRawData, ownerId)
		if err != nil {
			return lead, err
		}
		lead.Company = companyStatus
	}

	if deal, exists := mappedLead["deal"]; exists && createDeal {
		dealStatus, err := processDeal(husbpotClient, deal, existingLead, correspondingRawData, ownerId, pipelineId, stageId)
		if err != nil {
			return lead, err
		}
		lead.Deal = dealStatus
	}

	if contact, exists := mappedLead["contact"]; exists {
		contactStatus, err := processContact(husbpotClient, contact, existingLead, correspondingRawData, ownerId)
		if err != nil {
			return lead, err
		}
		lead.Contacts = contactStatus
	}

	if contacts, exists := mappedLead["contacts"]; exists {
		contactsStatus, err := processContacts(husbpotClient, contacts, existingLead, correspondingRawData, ownerId)
		if err != nil {
			return lead, err
		}
		lead.Contacts = contactsStatus
	}

	if err := createLeadAssociations(husbpotClient, lead); err != nil {
		return lead, err
	}

	return lead, nil
}

func getConfigValue[T any](configs map[string]any, key string) (T, error) {
	value, exists := configs[key]
	if !exists {
		return *new(T), fmt.Errorf("%s config is required", key)
	}
	castValue, ok := value.(T)
	if !ok {
		return *new(T), fmt.Errorf("invalid type for %s config", key)
	}
	return castValue, nil
}

func processCompany(client *hubspot.Client, company any, existingLead, rawData map[string]any, ownerId string) (*ObjectStatus, error) {
	exportedCompany, exists := existingLead["company"].(map[string]any)
	if exists && exportedCompany["crm_id"] != nil {
		return createExistingStatus(exportedCompany), nil
	}

	companyData, ok := company.(map[string]any)
	if !ok {
		return nil, errors.New("invalid company data")
	}

	sentCompany, err := sendCompany(client, companyData, ownerId)
	if err != nil {
		return nil, err
	}

	if drivaID, exists := rawData["company_contact_id"].(string); exists {
		sentCompany.DrivaContactId = &drivaID
	}

	return &sentCompany, nil
}

func processDeal(client *hubspot.Client, deal any, existingLead, rawData map[string]any, ownerId, pipelineId, stageId string) (*ObjectStatus, error) {
	exportedDeal, exists := existingLead["deal"].(map[string]any)
	if exists && exportedDeal["crm_id"] != nil {
		return createExistingStatus(exportedDeal), nil
	}

	dealData, ok := deal.(map[string]any)
	if !ok {
		return nil, errors.New("invalid deal data")
	}

	sentDeal, err := sendDeal(client, dealData, ownerId, pipelineId, stageId)
	if err != nil {
		return nil, err
	}

	if drivaID, exists := rawData["company_contact_id"].(string); exists {
		sentDeal.DrivaContactId = &drivaID
	}

	return &sentDeal, nil
}

func processContact(client *hubspot.Client, contact any, existingLead, rawData map[string]any, ownerId string) (*[]ObjectStatus, error) {
	exportedContact, exists := existingLead["contact"].(map[string]any)
	if exists && exportedContact["crm_id"] != nil {
		return &[]ObjectStatus{*createExistingStatus(exportedContact)}, nil
	}

	contactData, ok := contact.(map[string]any)
	if !ok {
		return nil, errors.New("invalid contact data")
	}

	sentContact, err := sendContact(client, contactData, ownerId)
	if err != nil {
		return nil, err
	}

	if drivaID, exists := rawData["profile_contact_id"].(string); exists {
		sentContact.DrivaContactId = &drivaID
	}

	return &[]ObjectStatus{sentContact}, nil
}

func processContacts(client *hubspot.Client, contacts any, existingLead, rawData map[string]any, ownerId string) (*[]ObjectStatus, error) {
	contactsData, ok := contacts.([]any)
	if !ok {
		return nil, errors.New("invalid contacts data")
	}

	var statuses []ObjectStatus
	for _, contact := range contactsData {
		contactMap, ok := contact.(map[string]any)
		if !ok {
			continue
		}

		exportedContacts, exists := existingLead["contacts"].([]map[string]any)
		if exists {
			for _, exportedContact := range exportedContacts {
				if exportedContact["driva_contact_id"] == contactMap["profile_company_id"] {
					statuses = append(statuses, *createExistingStatus(exportedContact))
					continue
				}
			}
		}

		sentContact, err := sendContact(client, contactMap, ownerId)
		if err != nil {
			return nil, err
		}

		if drivaID, exists := rawData["profile_contact_id"].(string); exists {
			sentContact.DrivaContactId = &drivaID
		}
		statuses = append(statuses, sentContact)
	}

	return &statuses, nil
}

func createExistingStatus(data map[string]any) *ObjectStatus {
	drivaID := data["driva_contact_id"].(string)
	return &ObjectStatus{
		CrmId:          data["crm_id"].(string),
		Status:         Status(data["status"].(string)),
		Message:        safeString(data, "message"),
		DrivaContactId: &drivaID,
	}
}

func safeString(data map[string]any, key string) string {
	if value, ok := data[key].(string); ok {
		return value
	}
	return ""
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

func (h HubspotService) Authorize(ctx context.Context, workspaceId string) (any, error) {

	company, err := h.companyRepo.Get(ctx, ports.CrmCompanyQueryParams{Crm: "hubspot", WorkspaceId: workspaceId})
	if err != nil {
		return nil, err
	}

	if company.RefreshToken.String == "" {
		return nil, errors.New("refresh token not found for " + workspaceId)
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
func (h HubspotService) Validate(c *fiber.Ctx, client any) bool {

	_, err := h.GetPipelines(client)
	return err == nil
}

func (h HubspotService) Install(installData any) (any, error) {
	baseURL := "https://app.hubspot.com/oauth/authorize"
	clientID := url.QueryEscape(os.Getenv("HUBSPOT_CLIENT_ID"))
	scope := url.QueryEscape(os.Getenv("HUBSPOT_SCOPE"))
	redirectURI := os.Getenv("HUBSPOT_REDIRECT_URI")
	installDataMap, isMap := installData.(map[string]any)
	if !isMap {
		return nil, errors.New("expected install data to be a map")
	}
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
