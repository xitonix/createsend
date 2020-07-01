package createsend

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/mail"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/xitonix/createsend/clients"
	"github.com/xitonix/createsend/mock"
)

func TestClientsAPI_Create(t *testing.T) {
	testCases := []struct {
		title                string
		forceHTTPClientError bool
		response             *http.Response
		expectedError        error
		input                clients.Client
		expectedResult       string
		oAuthAuthentication  bool
	}{
		{
			title: "successful execution",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`"client_id""`)),
			},
			expectedResult: "client_id",
		},
		{
			title: "oAuth authentication",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`"client_id"`)),
			},
			expectedResult:      "client_id",
			oAuthAuthentication: true,
		},
		{
			title:                "simulate remote call failure",
			response:             &http.Response{},
			forceHTTPClientError: true,
			expectedError:        mock.ErrDeliberate,
		},
		{
			title: "simulate server side error",
			response: &http.Response{
				StatusCode: 500,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"Message":"msg", "Code":500}`)),
			},
			expectedError: &Error{Code: 500},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.title, func(t *testing.T) {
			client, httpClient := createClient(t, tC.oAuthAuthentication, tC.forceHTTPClientError)
			httpClient.SetResponse(clientsPath, tC.response)
			actual, err := client.Clients().Create(tC.input)
			if err != nil {
				if !checkError(err, tC.expectedError) {
					t.Errorf("Expected '%v' error, actual: '%v'", tC.expectedError, err)
				}
				checkErrorType(t, err, !tC.forceHTTPClientError)
			}
			if actual != tC.expectedResult {
				t.Errorf("Expected: %v, Actual: %v", tC.expectedResult, actual)
			}
		})
	}
}

func TestClientsAPI_Get(t *testing.T) {
	testCases := []struct {
		title                string
		forceHTTPClientError bool
		response             *http.Response
		expected             *clients.ClientDetails
		expectedError        error
		oAuthAuthentication  bool
	}{
		{
			title: "empty json document is returned",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{}`)),
			},
			expected: &clients.ClientDetails{},
		},
		{
			title: "empty server response body is returned",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(&bytes.Buffer{}),
			},
			expected: &clients.ClientDetails{},
		},
		{
			title: "if the billing tier is provided the billing details must be parsed as monthly",
			response: &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(`{
				"ApiKey": "api_key",
				"BasicDetails": {
					"ClientID": "client_id",
					"CompanyName": "company_name",
					"ContactName": "contact_name",
					"EmailAddress": "contact_email_address",
					"Country": "country",
					"TimeZone": "timezone"
				},
				"AccessDetails": {
					"AccessLevel": 10,
					"Username": "username"
				},
				"BillingDetails": {
					"CurrentTier": "current_tier",
					"MonthlyScheme": "monthly_scheme",
					"Currency": "currency",
					"CanPurchaseCredits": true,
					"ClientPays": true,
					"CurrentMonthlyRate": 1.0,
					"MarkupPercentage": 20,
					"Credits": 100,
					"BaseRatePerRecipient": 10.1,
					"MarkupPerRecipient": 20.1,
					"MarkupOnDelivery": 30.1,
					"BaseDeliveryRate": 40.1,
					"MarkupOnDesignSpamTest": 50.1,
					"BaseDesignSpamTestRate": 60.1
				}
			}`)),
			},
			expected: &clients.ClientDetails{
				APIKey:   "api_key",
				Id:       "client_id",
				Company:  "company_name",
				Country:  "country",
				Timezone: "timezone",
				Contact: &clients.ContactDetails{
					Name:         "contact_name",
					EmailAddress: "contact_email_address",
					AccessLevel:  10,
					Username:     "username",
				},
				Billing: &clients.BillingDetails{
					Mode: clients.MonthlyBilling,
					Monthly: &clients.MonthlyBillingDetails{
						Tier:             "current_tier",
						Scheme:           "monthly_scheme",
						Rate:             1.0,
						MarkupPercentage: 20,
						Pending:          nil,
					},
					Currency:   "currency",
					ClientPays: true,
				},
			},
		},
		{
			title: "if the billing tier is not provided the billing details must be parsed as payg",
			response: &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(`{
				"ApiKey": "api_key",
				"BasicDetails": {
					"ClientID": "client_id",
					"CompanyName": "company_name",
					"ContactName": "contact_name",
					"EmailAddress": "contact_email_address",
					"Country": "country",
					"TimeZone": "timezone"
				},
				"AccessDetails": {
					"AccessLevel": 10,
					"Username": "username"
				},
				"BillingDetails": {
					"CurrentTier": "",
					"MonthlyScheme": "monthly_scheme",
					"Currency": "currency",
					"CanPurchaseCredits": true,
					"ClientPays": true,
					"CurrentMonthlyRate": 1.0,
					"MarkupPercentage": 20,
					"Credits": 100,
					"BaseRatePerRecipient": 10.1,
					"MarkupPerRecipient": 20.1,
					"MarkupOnDelivery": 30.1,
					"BaseDeliveryRate": 40.1,
					"MarkupOnDesignSpamTest": 50.1,
					"BaseDesignSpamTestRate": 60.1
				}
			}`)),
			},
			expected: &clients.ClientDetails{
				APIKey:   "api_key",
				Id:       "client_id",
				Company:  "company_name",
				Country:  "country",
				Timezone: "timezone",
				Contact: &clients.ContactDetails{
					Name:         "contact_name",
					EmailAddress: "contact_email_address",
					AccessLevel:  10,
					Username:     "username",
				},
				Billing: &clients.BillingDetails{
					Mode: clients.PAYGBilling,
					PAYG: &clients.PayAsYouGoBillingDetails{
						CanPurchaseCredits:     true,
						Credits:                100,
						BaseRatePerRecipient:   10.1,
						MarkupPerRecipient:     20.1,
						MarkupOnDelivery:       30.1,
						BaseDeliveryRate:       40.1,
						MarkupOnDesignSpamTest: 50.1,
						BaseDesignSpamTestRate: 60.1,
					},
					Currency:   "currency",
					ClientPays: true,
				},
			},
		},
		{
			title: "monthly plan with monthly pending billing details",
			response: &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(`{
				"ApiKey": "api_key",
				"BasicDetails": {
					"ClientID": "client_id",
					"CompanyName": "company_name",
					"ContactName": "contact_name",
					"EmailAddress": "contact_email_address",
					"Country": "country",
					"TimeZone": "timezone"
				},
				"AccessDetails": {
					"AccessLevel": 10,
					"Username": "username"
				},
				"BillingDetails": {
					"CurrentTier": "current_tier",
					"MonthlyScheme": "monthly_scheme",
					"Currency": "currency",
					"CanPurchaseCredits": true,
					"ClientPays": true,
					"CurrentMonthlyRate": 1.0,
					"MarkupPercentage": 20,
					"Credits": 100,
					"BaseRatePerRecipient": 10.1,
					"MarkupPerRecipient": 20.1,
					"MarkupOnDelivery": 30.1,
					"BaseDeliveryRate": 40.1,
					"MarkupOnDesignSpamTest": 50.1,
					"BaseDesignSpamTestRate": 60.1
				},
 				"PendingBillingDetails": {
					"CurrentTier": "pending_current_tier",
					"MonthlyScheme": "pending_monthly_scheme",
					"Currency": "pending_currency",
					"CanPurchaseCredits": true,
					"ClientPays": true,
					"CurrentMonthlyRate": 1.2,
					"MarkupPercentage": 30,
					"Credits": 200,
					"BaseRatePerRecipient": 10.2,
					"MarkupPerRecipient": 20.2,
					"MarkupOnDelivery": 30.2,
					"BaseDeliveryRate": 40.2,
					"MarkupOnDesignSpamTest": 50.2,
					"BaseDesignSpamTestRate": 60.2
				}
			}`)),
			},
			expected: &clients.ClientDetails{
				APIKey:   "api_key",
				Id:       "client_id",
				Company:  "company_name",
				Country:  "country",
				Timezone: "timezone",
				Contact: &clients.ContactDetails{
					Name:         "contact_name",
					EmailAddress: "contact_email_address",
					AccessLevel:  10,
					Username:     "username",
				},
				Billing: &clients.BillingDetails{
					Mode: clients.MonthlyBilling,
					Monthly: &clients.MonthlyBillingDetails{
						Tier:             "current_tier",
						Scheme:           "monthly_scheme",
						Rate:             1.0,
						MarkupPercentage: 20,
						Pending: &clients.BillingDetails{
							Mode: clients.MonthlyBilling,
							Monthly: &clients.MonthlyBillingDetails{
								Tier:             "pending_current_tier",
								Scheme:           "pending_monthly_scheme",
								Rate:             1.2,
								MarkupPercentage: 30,
								Pending:          nil,
							},
							Currency:   "pending_currency",
							ClientPays: true,
						},
					},
					Currency:   "currency",
					ClientPays: true,
				},
			},
		},
		{
			title: "monthly plan with payg pending billing details",
			response: &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(`{
				"ApiKey": "api_key",
				"BasicDetails": {
					"ClientID": "client_id",
					"CompanyName": "company_name",
					"ContactName": "contact_name",
					"EmailAddress": "contact_email_address",
					"Country": "country",
					"TimeZone": "timezone"
				},
				"AccessDetails": {
					"AccessLevel": 10,
					"Username": "username"
				},
				"BillingDetails": {
					"CurrentTier": "current_tier",
					"MonthlyScheme": "monthly_scheme",
					"Currency": "currency",
					"CanPurchaseCredits": true,
					"ClientPays": true,
					"CurrentMonthlyRate": 1.0,
					"MarkupPercentage": 20,
					"Credits": 100,
					"BaseRatePerRecipient": 10.1,
					"MarkupPerRecipient": 20.1,
					"MarkupOnDelivery": 30.1,
					"BaseDeliveryRate": 40.1,
					"MarkupOnDesignSpamTest": 50.1,
					"BaseDesignSpamTestRate": 60.1
				},
 				"PendingBillingDetails": {
					"CurrentTier": "",
					"MonthlyScheme": "",
					"Currency": "pending_currency",
					"CanPurchaseCredits": true,
					"ClientPays": true,
					"CurrentMonthlyRate": 1.2,
					"MarkupPercentage": 30,
					"Credits": 200,
					"BaseRatePerRecipient": 10.2,
					"MarkupPerRecipient": 20.2,
					"MarkupOnDelivery": 30.2,
					"BaseDeliveryRate": 40.2,
					"MarkupOnDesignSpamTest": 50.2,
					"BaseDesignSpamTestRate": 60.2
				}
			}`)),
			},
			expected: &clients.ClientDetails{
				APIKey:   "api_key",
				Id:       "client_id",
				Company:  "company_name",
				Country:  "country",
				Timezone: "timezone",
				Contact: &clients.ContactDetails{
					Name:         "contact_name",
					EmailAddress: "contact_email_address",
					AccessLevel:  10,
					Username:     "username",
				},
				Billing: &clients.BillingDetails{
					Mode: clients.MonthlyBilling,
					Monthly: &clients.MonthlyBillingDetails{
						Tier:             "current_tier",
						Scheme:           "monthly_scheme",
						Rate:             1.0,
						MarkupPercentage: 20,
						Pending: &clients.BillingDetails{
							Mode: clients.PAYGBilling,
							PAYG: &clients.PayAsYouGoBillingDetails{
								CanPurchaseCredits:     true,
								Credits:                200,
								BaseRatePerRecipient:   10.2,
								MarkupPerRecipient:     20.2,
								MarkupOnDelivery:       30.2,
								BaseDeliveryRate:       40.2,
								MarkupOnDesignSpamTest: 50.2,
								BaseDesignSpamTestRate: 60.2,
							},
							Currency:   "pending_currency",
							ClientPays: true,
						},
					},
					Currency:   "currency",
					ClientPays: true,
				},
			},
		},
		{
			title: "contact field will be nil if the server response has no access and contact details",
			response: &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(`{
				"ApiKey": "api_key",
				"BasicDetails": {
					"ClientID": "client_id",
					"CompanyName": "company_name",
					"ContactName": "",
					"EmailAddress": "",
					"Country": "country",
					"TimeZone": "timezone"
				}
			}`)),
			},
			expected: &clients.ClientDetails{
				APIKey:   "api_key",
				Id:       "client_id",
				Company:  "company_name",
				Country:  "country",
				Timezone: "timezone",
				Contact:  nil,
				Billing:  nil,
			},
		},
		{
			title: "contact field should not be nil if the server response only includes the contact name",
			response: &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(`{
				"ApiKey": "api_key",
				"BasicDetails": {
					"ClientID": "client_id",
					"CompanyName": "company_name",
					"ContactName": "name",
					"EmailAddress": "",
					"Country": "country",
					"TimeZone": "timezone"
				}
			}`)),
			},
			expected: &clients.ClientDetails{
				APIKey:   "api_key",
				Id:       "client_id",
				Company:  "company_name",
				Country:  "country",
				Timezone: "timezone",
				Contact: &clients.ContactDetails{
					Name:         "name",
					EmailAddress: "",
					AccessLevel:  -1,
					Username:     "",
				},
				Billing: nil,
			},
		},
		{
			title: "contact field should not be nil if the server response only includes the contact email address",
			response: &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(`{
				"ApiKey": "api_key",
				"BasicDetails": {
					"ClientID": "client_id",
					"CompanyName": "company_name",
					"ContactName": "",
					"EmailAddress": "email",
					"Country": "country",
					"TimeZone": "timezone"
				}
			}`)),
			},
			expected: &clients.ClientDetails{
				APIKey:   "api_key",
				Id:       "client_id",
				Company:  "company_name",
				Country:  "country",
				Timezone: "timezone",
				Contact: &clients.ContactDetails{
					Name:         "",
					EmailAddress: "email",
					AccessLevel:  -1,
					Username:     "",
				},
				Billing: nil,
			},
		},
		{
			title: "contact field should not be nil if the server response only includes access details",
			response: &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(`{
				"ApiKey": "api_key",
				"BasicDetails": {
					"ClientID": "client_id",
					"CompanyName": "company_name",
					"ContactName": "",
					"EmailAddress": "",
					"Country": "country",
					"TimeZone": "timezone"
				},
				"AccessDetails": {
					"AccessLevel": 10,
					"Username": "username"
				}
			}`)),
			},
			expected: &clients.ClientDetails{
				APIKey:   "api_key",
				Id:       "client_id",
				Company:  "company_name",
				Country:  "country",
				Timezone: "timezone",
				Contact: &clients.ContactDetails{
					Name:         "",
					EmailAddress: "",
					AccessLevel:  10,
					Username:     "username",
				},
				Billing: nil,
			},
		},
		{
			title: "oAuth authentication",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{}`)),
			},
			expected:            &clients.ClientDetails{},
			oAuthAuthentication: true,
		},
		{
			title:                "simulate remote call failure",
			response:             &http.Response{},
			forceHTTPClientError: true,
			expectedError:        mock.ErrDeliberate,
		},
		{
			title: "simulate server side error",
			response: &http.Response{
				StatusCode: 500,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"Message":"msg", "Code":500}`)),
			},
			expectedError: &Error{Code: 500},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.title, func(t *testing.T) {
			client, httpClient := createClient(t, tC.oAuthAuthentication, tC.forceHTTPClientError)
			httpClient.SetResponse("clients/client_id.json", tC.response)
			actual, err := client.Clients().Get("client_id")
			if err != nil {
				if !checkError(err, tC.expectedError) {
					t.Errorf("Expected '%v' error, actual: '%v'", tC.expectedError, err)
				}
				checkErrorType(t, err, !tC.forceHTTPClientError)
			}
			if diff := cmp.Diff(tC.expected, actual); diff != "" {
				t.Errorf("Expectations failed (-expected +actual):\n%s", diff)
			}
		})
	}
}

func TestClientsAPI_SentCampaigns(t *testing.T) {
	date := time.Date(2020, 12, 1, 20, 21, 22, 0, time.UTC)
	testCases := []struct {
		title                 string
		forceHTTPClientError  bool
		response              *http.Response
		expected              []*clients.SentCampaign
		expectedError         error
		oAuthAuthentication   bool
		expectClientSideError bool
	}{
		{
			title: "no campaigns",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`[]`)),
			},
			expected: []*clients.SentCampaign{},
		},
		{
			title: "no sent campaigns and empty server response body",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(&bytes.Buffer{}),
			},
			expected: []*clients.SentCampaign{},
		},
		{
			title: "client with campaigns",
			response: &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(`[
				{
					"FromName": "from_name",
					"FromEmail": "from@email.com",
					"ReplyTo": "reply_to",
					"WebVersionURL": "web_version",
					"WebVersionTextURL": "web_version_text",
					"CampaignID": "id",
					"Subject": "subject",
					"Name": "name",
					"SentDate": "2020-12-01 20:21:22",
					"TotalRecipients": 100
    			}
			]`)),
			},
			expected: []*clients.SentCampaign{
				{
					Campaign: clients.Campaign{
						Id:   "id",
						Name: "name",
						From: mail.Address{
							Name:    "from_name",
							Address: "from@email.com",
						},
						ReplyTo:           "reply_to",
						WebVersionURL:     "web_version",
						WebVersionTextURL: "web_version_text",
						Subject:           "subject",
					},
					SentDate:   date,
					Recipients: 100,
				},
			},
		},
		{
			title: "oAuth Authentication",
			response: &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(`[
				{
					"FromName": "from_name",
					"FromEmail": "from@email.com",
					"ReplyTo": "reply_to",
					"WebVersionURL": "web_version",
					"WebVersionTextURL": "web_version_text",
					"CampaignID": "id",
					"Subject": "subject",
					"Name": "name",
					"SentDate": "2020-12-01 20:21:22",
					"TotalRecipients": 100
    			}
			]`)),
			},
			expected: []*clients.SentCampaign{
				{
					Campaign: clients.Campaign{
						Id:   "id",
						Name: "name",
						From: mail.Address{
							Name:    "from_name",
							Address: "from@email.com",
						},
						ReplyTo:           "reply_to",
						WebVersionURL:     "web_version",
						WebVersionTextURL: "web_version_text",
						Subject:           "subject",
					},
					SentDate:   date,
					Recipients: 100,
				},
			},
			oAuthAuthentication: true,
		},
		{
			title: "invalid send date",
			response: &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(`[
				{
					"FromName": "from_name",
					"FromEmail": "from@email.com",
					"ReplyTo": "reply_to",
					"WebVersionURL": "web_version",
					"WebVersionTextURL": "web_version_text",
					"CampaignID": "id",
					"Subject": "subject",
					"Name": "name",
					"SentDate": "invalid date",
					"TotalRecipients": 100
    			}
			]`)),
			},
			expected:              nil,
			expectClientSideError: true,
			expectedError:         newClientError(ErrCodeDataProcessing),
		},
		{
			title:                "simulate remote call failure",
			response:             &http.Response{},
			forceHTTPClientError: true,
			expectedError:        mock.ErrDeliberate,
		},
		{
			title: "simulate server side error",
			response: &http.Response{
				StatusCode: 500,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"Message":"msg", "Code":500}`)),
			},
			expectedError: &Error{Code: 500},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.title, func(t *testing.T) {
			client, httpClient := createClient(t, tC.oAuthAuthentication, tC.forceHTTPClientError)
			httpClient.SetResponse("clients/client_id/campaigns.json", tC.response)
			actual, err := client.Clients().SentCampaigns("client_id")
			if err != nil {
				if !checkError(err, tC.expectedError) {
					t.Errorf("Expected '%v' error, actual: '%v'", tC.expectedError, err)
				}
				checkErrorType(t, err, !tC.expectClientSideError && !tC.forceHTTPClientError)
			}
			if diff := cmp.Diff(tC.expected, actual); diff != "" {
				t.Errorf("Expectations failed (-expected +actual):\n%s", diff)
			}
		})
	}
}

func TestClientsAPI_ScheduledCampaigns(t *testing.T) {
	scheduleDate := time.Date(2020, 12, 1, 20, 21, 22, 0, time.UTC)
	createdDate := time.Date(2020, 12, 2, 20, 21, 22, 0, time.UTC)

	testCases := []struct {
		title                 string
		forceHTTPClientError  bool
		response              *http.Response
		expected              []*clients.ScheduledCampaign
		expectedError         error
		oAuthAuthentication   bool
		expectClientSideError bool
	}{
		{
			title: "no campaigns",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`[]`)),
			},
			expected: []*clients.ScheduledCampaign{},
		},
		{
			title: "no campaigns and empty server response body",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(&bytes.Buffer{}),
			},
			expected: []*clients.ScheduledCampaign{},
		},
		{
			title: "client with campaigns",
			response: &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(`[
				{
					"FromName": "from_name",
					"FromEmail": "from@email.com",
					"ReplyTo": "reply_to",
					"WebVersionURL": "web_version",
					"WebVersionTextURL": "web_version_text",
					"CampaignID": "id",
					"Subject": "subject",
					"Name": "name",
					"DateScheduled": "2020-12-01 20:21:22",
					"DateCreated": "2020-12-02 20:21:22",
					"ScheduledTimeZone": "tz"
    			}
			]`)),
			},
			expected: []*clients.ScheduledCampaign{
				{
					Campaign: clients.Campaign{
						Id:   "id",
						Name: "name",
						From: mail.Address{
							Name:    "from_name",
							Address: "from@email.com",
						},
						ReplyTo:           "reply_to",
						WebVersionURL:     "web_version",
						WebVersionTextURL: "web_version_text",
						Subject:           "subject",
					},
					DateScheduled: scheduleDate,
					DateCreated:   createdDate,
					Timezone:      "tz",
				},
			},
		},
		{
			title: "invalid date created",
			response: &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(`[
				{
					"FromName": "from_name",
					"FromEmail": "from@email.com",
					"ReplyTo": "reply_to",
					"WebVersionURL": "web_version",
					"WebVersionTextURL": "web_version_text",
					"CampaignID": "id",
					"Subject": "subject",
					"Name": "name",
					"DateScheduled": "2020-12-01 20:21:22",
					"DateCreated": "invalid date",
					"ScheduledTimeZone": "tz"
    			}
			]`)),
			},
			expected:              nil,
			expectClientSideError: true,
			expectedError:         newClientError(ErrCodeDataProcessing),
		},
		{
			title: "invalid date scheduled",
			response: &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(`[
				{
					"FromName": "from_name",
					"FromEmail": "from@email.com",
					"ReplyTo": "reply_to",
					"WebVersionURL": "web_version",
					"WebVersionTextURL": "web_version_text",
					"CampaignID": "id",
					"Subject": "subject",
					"Name": "name",
					"DateScheduled": "invalid date",
					"DateCreated": "2020-12-02 20:21:22",
					"ScheduledTimeZone": "tz"
    			}
			]`)),
			},
			expected:              nil,
			expectClientSideError: true,
			expectedError:         newClientError(ErrCodeDataProcessing),
		},
		{
			title: "oAuth Authentication",
			response: &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(`[
				{
					"FromName": "from_name",
					"FromEmail": "from@email.com",
					"ReplyTo": "reply_to",
					"WebVersionURL": "web_version",
					"WebVersionTextURL": "web_version_text",
					"CampaignID": "id",
					"Subject": "subject",
					"Name": "name",
					"DateScheduled": "2020-12-01 20:21:22",
					"DateCreated": "2020-12-02 20:21:22",
					"ScheduledTimeZone": "tz"
    			}
			]`)),
			},
			expected: []*clients.ScheduledCampaign{
				{
					Campaign: clients.Campaign{
						Id:   "id",
						Name: "name",
						From: mail.Address{
							Name:    "from_name",
							Address: "from@email.com",
						},
						ReplyTo:           "reply_to",
						WebVersionURL:     "web_version",
						WebVersionTextURL: "web_version_text",
						Subject:           "subject",
					},
					DateScheduled: scheduleDate,
					DateCreated:   createdDate,
					Timezone:      "tz",
				},
			},
			oAuthAuthentication: true,
		},
		{
			title:                "simulate remote call failure",
			response:             &http.Response{},
			forceHTTPClientError: true,
			expectedError:        mock.ErrDeliberate,
		},
		{
			title: "simulate server side error",
			response: &http.Response{
				StatusCode: 500,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"Message":"msg", "Code":500}`)),
			},
			expectedError: &Error{Code: 500},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.title, func(t *testing.T) {
			client, httpClient := createClient(t, tC.oAuthAuthentication, tC.forceHTTPClientError)
			httpClient.SetResponse("clients/client_id/scheduled.json", tC.response)
			actual, err := client.Clients().ScheduledCampaigns("client_id")
			if err != nil {
				if !checkError(err, tC.expectedError) {
					t.Errorf("Expected '%v' error, actual: '%v'", tC.expectedError, err)
				}
				checkErrorType(t, err, !tC.expectClientSideError && !tC.forceHTTPClientError)
			}
			if diff := cmp.Diff(tC.expected, actual); diff != "" {
				t.Errorf("Expectations failed (-expected +actual):\n%s", diff)
			}
		})
	}
}

func TestClientsAPI_DraftCampaigns(t *testing.T) {
	date := time.Date(2020, 12, 1, 20, 21, 22, 0, time.UTC)
	testCases := []struct {
		title                 string
		forceHTTPClientError  bool
		expectClientSideError bool
		response              *http.Response
		expected              []*clients.DraftCampaign
		expectedError         error
		oAuthAuthentication   bool
	}{
		{
			title: "no campaigns",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`[]`)),
			},
			expected: []*clients.DraftCampaign{},
		},
		{
			title: "no campaigns and empty server response body",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(&bytes.Buffer{}),
			},
			expected: []*clients.DraftCampaign{},
		},
		{
			title: "client with campaigns",
			response: &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(`[
				{
					"FromName": "from_name",
					"FromEmail": "from@email.com",
					"ReplyTo": "reply_to",
					"WebVersionURL": "web_version",
					"WebVersionTextURL": "web_version_text",
					"CampaignID": "id",
					"Subject": "subject",
					"Name": "name",
					"DateCreated": "2020-12-01 20:21:22"
    			}
			]`)),
			},
			expected: []*clients.DraftCampaign{
				{
					Campaign: clients.Campaign{
						Id:   "id",
						Name: "name",
						From: mail.Address{
							Name:    "from_name",
							Address: "from@email.com",
						},
						ReplyTo:           "reply_to",
						WebVersionURL:     "web_version",
						WebVersionTextURL: "web_version_text",
						Subject:           "subject",
					},
					DateCreated: date,
				},
			},
		},
		{
			title: "invalid date value",
			response: &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(`[
				{
					"FromName": "from_name",
					"FromEmail": "from@email.com",
					"ReplyTo": "reply_to",
					"WebVersionURL": "web_version",
					"WebVersionTextURL": "web_version_text",
					"CampaignID": "id",
					"Subject": "subject",
					"Name": "name",
					"DateCreated": "invalid date"
    			}
			]`)),
			},
			expected:              nil,
			expectedError:         newClientError(ErrCodeDataProcessing),
			expectClientSideError: true,
		},
		{
			title: "oAuth Authentication",
			response: &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(`[
				{
					"FromName": "from_name",
					"FromEmail": "from@email.com",
					"ReplyTo": "reply_to",
					"WebVersionURL": "web_version",
					"WebVersionTextURL": "web_version_text",
					"CampaignID": "id",
					"Subject": "subject",
					"Name": "name",
					"DateCreated": "2020-12-01 20:21:22"
    			}
			]`)),
			},
			expected: []*clients.DraftCampaign{
				{
					Campaign: clients.Campaign{
						Id:   "id",
						Name: "name",
						From: mail.Address{
							Name:    "from_name",
							Address: "from@email.com",
						},
						ReplyTo:           "reply_to",
						WebVersionURL:     "web_version",
						WebVersionTextURL: "web_version_text",
						Subject:           "subject",
					},
					DateCreated: date,
				},
			},
			oAuthAuthentication: true,
		},
		{
			title:                "simulate remote call failure",
			response:             &http.Response{},
			forceHTTPClientError: true,
			expectedError:        mock.ErrDeliberate,
		},
		{
			title: "simulate server side error",
			response: &http.Response{
				StatusCode: 500,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"Message":"msg", "Code":500}`)),
			},
			expectedError: &Error{Code: 500},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.title, func(t *testing.T) {
			client, httpClient := createClient(t, tC.oAuthAuthentication, tC.forceHTTPClientError)
			httpClient.SetResponse("clients/client_id/drafts.json", tC.response)
			actual, err := client.Clients().DraftCampaigns("client_id")
			if err != nil {
				if !checkError(err, tC.expectedError) {
					t.Errorf("Expected '%v' error, actual: '%v'", tC.expectedError, err)
				}
				checkErrorType(t, err, !tC.expectClientSideError && !tC.forceHTTPClientError)
			}
			if diff := cmp.Diff(tC.expected, actual); diff != "" {
				t.Errorf("Expectations failed (-expected +actual):\n%s", diff)
			}
		})
	}
}

func TestClientsAPI_Lists(t *testing.T) {
	testCases := []struct {
		title                string
		forceHTTPClientError bool
		response             *http.Response
		expected             []*clients.List
		expectedError        error
		oAuthAuthentication  bool
	}{
		{
			title: "no lists",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`[]`)),
			},
			expected: []*clients.List{},
		},
		{
			title: "no lists and empty server response body",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(&bytes.Buffer{}),
			},
			expected: []*clients.List{},
		},
		{
			title: "client with lists",
			response: &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(`[
				{
					"ListID": "list_id",
					"Name": "list_name"
    			}
			]`)),
			},
			expected: []*clients.List{
				{
					Id:   "list_id",
					Name: "list_name",
				},
			},
		},
		{
			title: "oAuth Authentication",
			response: &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(`[
				{
					"ListID": "list_id",
					"Name": "list_name"
    			}
			]`)),
			},
			expected: []*clients.List{
				{
					Id:   "list_id",
					Name: "list_name",
				},
			},
			oAuthAuthentication: true,
		},
		{
			title:                "simulate remote call failure",
			response:             &http.Response{},
			forceHTTPClientError: true,
			expectedError:        mock.ErrDeliberate,
		},
		{
			title: "simulate server side error",
			response: &http.Response{
				StatusCode: 500,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"Message":"msg", "Code":500}`)),
			},
			expectedError: &Error{Code: 500},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.title, func(t *testing.T) {
			client, httpClient := createClient(t, tC.oAuthAuthentication, tC.forceHTTPClientError)
			httpClient.SetResponse("clients/client_id/lists.json", tC.response)
			actual, err := client.Clients().Lists("client_id")
			if err != nil {
				if !checkError(err, tC.expectedError) {
					t.Errorf("Expected '%v' error, actual: '%v'", tC.expectedError, err)
				}
				checkErrorType(t, err, !tC.forceHTTPClientError)
			}
			if diff := cmp.Diff(tC.expected, actual); diff != "" {
				t.Errorf("Expectations failed (-expected +actual):\n%s", diff)
			}
		})
	}
}

func TestClientsAPI_ListsByEmailAddress(t *testing.T) {
	date := time.Date(2020, 12, 1, 20, 21, 22, 0, time.UTC)
	testCases := []struct {
		title                 string
		forceHTTPClientError  bool
		expectClientSideError bool
		response              *http.Response
		expected              []*clients.SubscriberList
		expectedError         error
		oAuthAuthentication   bool
	}{
		{
			title: "no lists",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`[]`)),
			},
			expected: []*clients.SubscriberList{},
		},
		{
			title: "no lists and empty server response body",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(&bytes.Buffer{}),
			},
			expected: []*clients.SubscriberList{},
		},
		{
			title: "client with lists",
			response: &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(`[
				{
					"ListID": "list_id",
					"ListName": "list_name",
					"SubscriberState": "state", 
					"DateSubscriberAdded": "2020-12-01 20:21:22"
    			}
			]`)),
			},
			expected: []*clients.SubscriberList{
				{
					List: clients.List{
						Id:   "list_id",
						Name: "list_name",
					},
					Subscriber: clients.Subscriber{
						State:     "state",
						DateAdded: date,
					},
				},
			},
		},
		{
			title: "invalid date value",
			response: &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(`[
				{
					"ListID": "list_id",
					"ListName": "list_name",
					"SubscriberState": "state", 
					"DateSubscriberAdded": "invalid date"
    			}
			]`)),
			},
			expected:              nil,
			expectedError:         newClientError(ErrCodeDataProcessing),
			expectClientSideError: true,
		},
		{
			title: "oAuth Authentication",
			response: &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(`[
				{
					"ListID": "list_id",
					"ListName": "list_name",
					"SubscriberState": "state", 
					"DateSubscriberAdded": "2020-12-01 20:21:22"
    			}
			]`)),
			},
			expected: []*clients.SubscriberList{
				{
					List: clients.List{
						Id:   "list_id",
						Name: "list_name",
					},
					Subscriber: clients.Subscriber{
						State:     "state",
						DateAdded: date,
					},
				},
			},
			oAuthAuthentication: true,
		},
		{
			title:                "simulate remote call failure",
			response:             &http.Response{},
			forceHTTPClientError: true,
			expectedError:        mock.ErrDeliberate,
		},
		{
			title: "simulate server side error",
			response: &http.Response{
				StatusCode: 500,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"Message":"msg", "Code":500}`)),
			},
			expectedError: &Error{Code: 500},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.title, func(t *testing.T) {
			client, httpClient := createClient(t, tC.oAuthAuthentication, tC.forceHTTPClientError)
			httpClient.SetResponse("clients/client_id/listsforemail.json", tC.response)
			actual, err := client.Clients().ListsByEmailAddress("client_id", "email@address.com")
			if err != nil {
				if !checkError(err, tC.expectedError) {
					t.Errorf("Expected '%v' error, actual: '%v'", tC.expectedError, err)
				}
				checkErrorType(t, err, !tC.expectClientSideError && !tC.forceHTTPClientError)
			}
			if diff := cmp.Diff(tC.expected, actual); diff != "" {
				t.Errorf("Expectations failed (-expected +actual):\n%s", diff)
			}
		})
	}
}
