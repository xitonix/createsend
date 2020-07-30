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
	"github.com/xitonix/createsend/order"
)

func TestClientsAPI_Create(t *testing.T) {
	testCases := []struct {
		title                string
		forceHTTPClientError bool
		response             *http.Response
		expectedError        error
		input                clients.BasicDetails
		expectedResult       string
		oAuthAuthentication  bool
	}{
		{
			title: "successful execution",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`"client_id"`)),
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
				"APIKey": "api_key",
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
				ID:       "client_id",
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
				"APIKey": "api_key",
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
				ID:       "client_id",
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
				"APIKey": "api_key",
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
				ID:       "client_id",
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
				"APIKey": "api_key",
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
				ID:       "client_id",
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
				"APIKey": "api_key",
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
				ID:       "client_id",
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
				"APIKey": "api_key",
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
				ID:       "client_id",
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
				"APIKey": "api_key",
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
				ID:       "client_id",
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
				"APIKey": "api_key",
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
				ID:       "client_id",
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
						ID:   "id",
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
						ID:   "id",
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
						ID:   "id",
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
						ID:   "id",
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
						ID:   "id",
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
						ID:   "id",
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
					ID:   "list_id",
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
					ID:   "list_id",
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
						ID:   "list_id",
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
						ID:   "list_id",
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

func TestClientsAPI_Segments(t *testing.T) {
	testCases := []struct {
		title                string
		forceHTTPClientError bool
		response             *http.Response
		expected             []*clients.Segment
		expectedError        error
		oAuthAuthentication  bool
	}{
		{
			title: "no segments",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`[]`)),
			},
			expected: []*clients.Segment{},
		},
		{
			title: "no segments and empty server response body",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(&bytes.Buffer{}),
			},
			expected: []*clients.Segment{},
		},
		{
			title: "client with segments",
			response: &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(`[
				{
					"SegmentID": "segment_id",
					"ListID": "list_id",
					"Title": "segment_title"
    			}
			]`)),
			},
			expected: []*clients.Segment{
				{
					ID:     "segment_id",
					Title:  "segment_title",
					ListID: "list_id",
				},
			},
		},
		{
			title: "oAuth Authentication",
			response: &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(`[
				{
					"SegmentID": "segment_id",
					"ListID": "list_id",
					"Title": "segment_title"
    			}
			]`)),
			},
			expected: []*clients.Segment{
				{
					ID:     "segment_id",
					Title:  "segment_title",
					ListID: "list_id",
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
			httpClient.SetResponse("clients/client_id/segments.json", tC.response)
			actual, err := client.Clients().Segments("client_id")
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

func TestClientsAPI_SuppressionList(t *testing.T) {
	date := time.Date(2020, 12, 1, 20, 21, 22, 0, time.UTC)
	testCases := []struct {
		title                 string
		forceHTTPClientError  bool
		expectClientSideError bool
		response              *http.Response
		expected              *clients.SuppressionList
		expectedError         error
		oAuthAuthentication   bool
	}{
		{
			title: "no lists",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{}`)),
			},
			expected: &clients.SuppressionList{
				Entries: []*clients.SuppressionDetails{},
			},
		},
		{
			title: "no lists and empty server response body",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(&bytes.Buffer{}),
			},
			expected: &clients.SuppressionList{
				Entries: []*clients.SuppressionDetails{},
			},
		},
		{
			title: "client with suppression list order by email asc",
			response: &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(`{
					"Results": [
							{
								"SuppressionReason": "reason",
								"EmailAddress": "email@address.com",
								"Date": "2020-12-01 20:21:22",
								"State": "state"
							}
						],
						"ResultsOrderedBy": "email",
						"OrderDirection": "asc",
						"PageNumber": 1,
						"PageSize": 1000,
						"RecordsOnThisPage": 2,
						"TotalNumberOfRecords": 5,
						"NumberOfPages": 10
					}`)),
			},
			expected: &clients.SuppressionList{
				Entries: []*clients.SuppressionDetails{
					{
						Reason:       "reason",
						EmailAddress: "email@address.com",
						Date:         date,
						State:        "state",
					},
				},
				OrderedBy:            order.BySuppressedEmailAddress,
				OrderDirection:       order.ASC,
				PageNumber:           1,
				PageSize:             1000,
				RecordsOnThisPage:    2,
				TotalNumberOfRecords: 5,
				NumberOfPages:        10,
			},
		},
		{
			title: "client with suppression list order by email desc",
			response: &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(`{
					"Results": [
							{
								"SuppressionReason": "reason",
								"EmailAddress": "email@address.com",
								"Date": "2020-12-01 20:21:22",
								"State": "state"
							}
						],
						"ResultsOrderedBy": "email",
						"OrderDirection": "desc",
						"PageNumber": 1,
						"PageSize": 1000,
						"RecordsOnThisPage": 2,
						"TotalNumberOfRecords": 5,
						"NumberOfPages": 10
					}`)),
			},
			expected: &clients.SuppressionList{
				Entries: []*clients.SuppressionDetails{
					{
						Reason:       "reason",
						EmailAddress: "email@address.com",
						Date:         date,
						State:        "state",
					},
				},
				OrderedBy:            order.BySuppressedEmailAddress,
				OrderDirection:       order.DESC,
				PageNumber:           1,
				PageSize:             1000,
				RecordsOnThisPage:    2,
				TotalNumberOfRecords: 5,
				NumberOfPages:        10,
			},
		},
		{
			title: "client with suppression list order by date asc",
			response: &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(`{
					"Results": [
							{
								"SuppressionReason": "reason",
								"EmailAddress": "email@address.com",
								"Date": "2020-12-01 20:21:22",
								"State": "state"
							}
						],
						"ResultsOrderedBy": "date",
						"OrderDirection": "asc",
						"PageNumber": 1,
						"PageSize": 1000,
						"RecordsOnThisPage": 2,
						"TotalNumberOfRecords": 5,
						"NumberOfPages": 10
					}`)),
			},
			expected: &clients.SuppressionList{
				Entries: []*clients.SuppressionDetails{
					{
						Reason:       "reason",
						EmailAddress: "email@address.com",
						Date:         date,
						State:        "state",
					},
				},
				OrderedBy:            order.BySuppressionDate,
				OrderDirection:       order.ASC,
				PageNumber:           1,
				PageSize:             1000,
				RecordsOnThisPage:    2,
				TotalNumberOfRecords: 5,
				NumberOfPages:        10,
			},
		},
		{
			title: "client with suppression list order by date desc",
			response: &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(`{
					"Results": [
							{
								"SuppressionReason": "reason",
								"EmailAddress": "email@address.com",
								"Date": "2020-12-01 20:21:22",
								"State": "state"
							}
						],
						"ResultsOrderedBy": "date",
						"OrderDirection": "desc",
						"PageNumber": 1,
						"PageSize": 1000,
						"RecordsOnThisPage": 2,
						"TotalNumberOfRecords": 5,
						"NumberOfPages": 10
					}`)),
			},
			expected: &clients.SuppressionList{
				Entries: []*clients.SuppressionDetails{
					{
						Reason:       "reason",
						EmailAddress: "email@address.com",
						Date:         date,
						State:        "state",
					},
				},
				OrderedBy:            order.BySuppressionDate,
				OrderDirection:       order.DESC,
				PageNumber:           1,
				PageSize:             1000,
				RecordsOnThisPage:    2,
				TotalNumberOfRecords: 5,
				NumberOfPages:        10,
			},
		},
		{
			title: "invalid order field",
			response: &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(`{
					"Results": [
							{
								"SuppressionReason": "reason",
								"EmailAddress": "email@address.com",
								"Date": "2020-12-01 20:21:22",
								"State": "state"
							}
						],
						"ResultsOrderedBy": "invalid",
						"OrderDirection": "desc",
						"PageNumber": 1,
						"PageSize": 1000,
						"RecordsOnThisPage": 2,
						"TotalNumberOfRecords": 5,
						"NumberOfPages": 10
					}`)),
			},
			expected:              nil,
			expectedError:         newClientError(ErrCodeInvalidJSON),
			expectClientSideError: true,
		},
		{
			title: "invalid order direction",
			response: &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(`{
					"Results": [
							{
								"SuppressionReason": "reason",
								"EmailAddress": "email@address.com",
								"Date": "2020-12-01 20:21:22",
								"State": "state"
							}
						],
						"ResultsOrderedBy": "email",
						"OrderDirection": "invalid",
						"PageNumber": 1,
						"PageSize": 1000,
						"RecordsOnThisPage": 2,
						"TotalNumberOfRecords": 5,
						"NumberOfPages": 10
					}`)),
			},
			expected:              nil,
			expectedError:         newClientError(ErrCodeInvalidJSON),
			expectClientSideError: true,
		},
		{
			title: "invalid date",
			response: &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(`{
					"Results": [
							{
								"SuppressionReason": "reason",
								"EmailAddress": "email@address.com",
								"Date": "invalid date",
								"State": "state"
							}
						],
						"ResultsOrderedBy": "email",
						"OrderDirection": "asc",
						"PageNumber": 1,
						"PageSize": 1000,
						"RecordsOnThisPage": 2,
						"TotalNumberOfRecords": 5,
						"NumberOfPages": 10
					}`)),
			},
			expected:              nil,
			expectedError:         newClientError(ErrCodeDataProcessing),
			expectClientSideError: true,
		},
		{
			title: "oAuth Authentication",
			response: &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(`{
					"Results": [
							{
								"SuppressionReason": "reason",
								"EmailAddress": "email@address.com",
								"Date": "2020-12-01 20:21:22",
								"State": "state"
							}
						],
						"ResultsOrderedBy": "date",
						"OrderDirection": "desc",
						"PageNumber": 1,
						"PageSize": 1000,
						"RecordsOnThisPage": 2,
						"TotalNumberOfRecords": 5,
						"NumberOfPages": 10
					}`)),
			},
			expected: &clients.SuppressionList{
				Entries: []*clients.SuppressionDetails{
					{
						Reason:       "reason",
						EmailAddress: "email@address.com",
						Date:         date,
						State:        "state",
					},
				},
				OrderedBy:            order.BySuppressionDate,
				OrderDirection:       order.DESC,
				PageNumber:           1,
				PageSize:             1000,
				RecordsOnThisPage:    2,
				TotalNumberOfRecords: 5,
				NumberOfPages:        10,
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
			httpClient.SetResponse("clients/client_id/suppressionlist.json", tC.response)
			actual, err := client.Clients().SuppressionList("client_id", 1000, 1, order.BySuppressedEmailAddress, order.ASC)
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

func TestClientsAPI_Suppress(t *testing.T) {
	testCases := []struct {
		title                string
		forceHTTPClientError bool
		response             *http.Response
		expectedError        error
		oAuthAuthentication  bool
	}{
		{
			title: "successful execution",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(``)),
			},
		},
		{
			title: "oAuth authentication",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(``)),
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
			httpClient.SetResponse("clients/client_id/suppress.json", tC.response)
			err := client.Clients().Suppress("client_id", "email1", "email2")
			if err != nil {
				if !checkError(err, tC.expectedError) {
					t.Errorf("Expected '%v' error, actual: '%v'", tC.expectedError, err)
				}
				checkErrorType(t, err, !tC.forceHTTPClientError)
			}
		})
	}
}

func TestClientsAPI_UnSuppress(t *testing.T) {
	testCases := []struct {
		title                string
		forceHTTPClientError bool
		response             *http.Response
		expectedError        error
		oAuthAuthentication  bool
	}{
		{
			title: "successful execution",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(``)),
			},
		},
		{
			title: "oAuth authentication",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(``)),
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
			httpClient.SetResponse("clients/client_id/unsuppress.json", tC.response)
			err := client.Clients().UnSuppress("client_id", "email")
			if err != nil {
				if !checkError(err, tC.expectedError) {
					t.Errorf("Expected '%v' error, actual: '%v'", tC.expectedError, err)
				}
				checkErrorType(t, err, !tC.forceHTTPClientError)
			}
		})
	}
}

func TestClientsAPI_Templates(t *testing.T) {
	testCases := []struct {
		title                string
		forceHTTPClientError bool
		response             *http.Response
		expected             []*clients.Template
		expectedError        error
		oAuthAuthentication  bool
	}{
		{
			title: "no templates",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`[]`)),
			},
			expected: []*clients.Template{},
		},
		{
			title: "no templates and empty server response body",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(&bytes.Buffer{}),
			},
			expected: []*clients.Template{},
		},
		{
			title: "client with templates",
			response: &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(`[
				{
					"TemplateID": "template_id",
					"Name": "name",
					"PreviewURL": "preview_url",
					"ScreenshotURL": "screenshot_url"
    			}
			]`)),
			},
			expected: []*clients.Template{
				{
					ID:            "template_id",
					Name:          "name",
					PreviewURL:    "preview_url",
					ScreenshotURL: "screenshot_url",
				},
			},
		},
		{
			title: "oAuth Authentication",
			response: &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(`[
				{
					"TemplateID": "template_id",
					"Name": "name",
					"PreviewURL": "preview_url",
					"ScreenshotURL": "screenshot_url"
    			}
			]`)),
			},
			expected: []*clients.Template{
				{
					ID:            "template_id",
					Name:          "name",
					PreviewURL:    "preview_url",
					ScreenshotURL: "screenshot_url",
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
			httpClient.SetResponse("clients/client_id/templates.json", tC.response)
			actual, err := client.Clients().Templates("client_id")
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

func TestClientsAPI_Update(t *testing.T) {
	testCases := []struct {
		title                string
		forceHTTPClientError bool
		response             *http.Response
		expectedError        error
		oAuthAuthentication  bool
	}{
		{
			title: "receiving 200 from the server means success",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{}`)),
			},
		},
		{
			title: "oAuth authentication",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{}`)),
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
			httpClient.SetResponse("clients/client_id/setbasics.json", tC.response)
			err := client.Clients().Update("client_id", clients.BasicDetails{})
			if err != nil {
				if !checkError(err, tC.expectedError) {
					t.Errorf("Expected '%v' error, actual: '%v'", tC.expectedError, err)
				}
				checkErrorType(t, err, !tC.forceHTTPClientError)
			}
		})
	}
}

func TestClientsAPI_SetPaygBilling(t *testing.T) {
	testCases := []struct {
		title                string
		forceHTTPClientError bool
		response             *http.Response
		expectedError        error
		oAuthAuthentication  bool
	}{
		{
			title: "receiving 200 from the server means success",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{}`)),
			},
		},
		{
			title: "oAuth authentication",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{}`)),
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
			httpClient.SetResponse("clients/client_id/setpaygbilling.json", tC.response)
			err := client.Clients().SetPAYGBilling("client_id", clients.PAYGRates{})
			if err != nil {
				if !checkError(err, tC.expectedError) {
					t.Errorf("Expected '%v' error, actual: '%v'", tC.expectedError, err)
				}
				checkErrorType(t, err, !tC.forceHTTPClientError)
			}
		})
	}
}

func TestClientsAPI_SetMonthlyBilling(t *testing.T) {
	testCases := []struct {
		title                string
		forceHTTPClientError bool
		response             *http.Response
		expectedError        error
		oAuthAuthentication  bool
	}{
		{
			title: "receiving 200 from the server means success",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{}`)),
			},
		},
		{
			title: "oAuth authentication",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{}`)),
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
			httpClient.SetResponse("clients/client_id/setmonthlybilling.json", tC.response)
			err := client.Clients().SetMonthlyBilling("client_id", clients.MonthlyRates{})
			if err != nil {
				if !checkError(err, tC.expectedError) {
					t.Errorf("Expected '%v' error, actual: '%v'", tC.expectedError, err)
				}
				checkErrorType(t, err, !tC.forceHTTPClientError)
			}
		})
	}
}

func TestClientsAPI_TransferCredits(t *testing.T) {
	testCases := []struct {
		title                string
		forceHTTPClientError bool
		response             *http.Response
		expected             *clients.CreditTransferResult
		input                clients.CreditTransferRequest
		expectedError        error
		oAuthAuthentication  bool
	}{
		{
			title: "empty JSON response",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{}`)),
			},
			expected: &clients.CreditTransferResult{},
		},
		{
			title: "empty server response",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(&bytes.Buffer{}),
			},
			expected: nil,
		},
		{
			title: "non empty server response",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"AccountCredits":1, "ClientCredits":2}`)),
			},
			expected: &clients.CreditTransferResult{
				AccountCredits: 1,
				ClientCredits:  2,
			},
		},
		{
			title: "oAuth authentication",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"AccountCredits":1, "ClientCredits":2}`)),
			},
			expected: &clients.CreditTransferResult{
				AccountCredits: 1,
				ClientCredits:  2,
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
			httpClient.SetResponse("clients/client_id/credits.json", tC.response)
			actual, err := client.Clients().TransferCredits("client_id", tC.input)
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

func TestClientsAPI_Delete(t *testing.T) {
	testCases := []struct {
		title                string
		forceHTTPClientError bool
		response             *http.Response
		expectedError        error
		oAuthAuthentication  bool
	}{
		{
			title: "successful deletion",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(&bytes.Buffer{}),
			},
		},
		{
			title: "oAuth authentication",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(&bytes.Buffer{}),
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
			httpClient.SetResponse("clients/client_id.json", tC.response)
			err := client.Clients().Delete("client_id")
			if err != nil {
				if !checkError(err, tC.expectedError) {
					t.Errorf("Expected '%v' error, actual: '%v'", tC.expectedError, err)
				}
				checkErrorType(t, err, !tC.forceHTTPClientError)
			}
		})
	}
}

func TestClientsAPI_AddPerson(t *testing.T) {
	testCases := []struct {
		title                string
		forceHTTPClientError bool
		response             *http.Response
		expectedError        error
		oAuthAuthentication  bool
		expected             string
	}{
		{
			title: "empty JSON response",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{}`)),
			},
			expected: "",
		},
		{
			title: "empty server response",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(&bytes.Buffer{}),
			},
			expected: "",
		},
		{
			title: "non empty server response",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"EmailAddress":"e@d.com"}`)),
			},
			expected: "e@d.com",
		},
		{
			title: "oAuth authentication",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"EmailAddress":"e@d.com"}`)),
			},
			expected:            "e@d.com",
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
			httpClient.SetResponse("clients/client_id/people.json", tC.response)
			actual, err := client.Clients().AddPerson("client_id", clients.Person{})
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

func TestClientsAPI_UpdatePerson(t *testing.T) {
	testCases := []struct {
		title                string
		forceHTTPClientError bool
		response             *http.Response
		expectedError        error
		oAuthAuthentication  bool
		expected             string
	}{
		{
			title: "empty JSON response",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{}`)),
			},
			expected: "",
		},
		{
			title: "empty server response",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(&bytes.Buffer{}),
			},
			expected: "",
		},
		{
			title: "non empty server response",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"EmailAddress":"e@d.com"}`)),
			},
			expected: "e@d.com",
		},
		{
			title: "oAuth authentication",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"EmailAddress":"e@d.com"}`)),
			},
			expected:            "e@d.com",
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
			httpClient.SetResponse("clients/client_id/people.json", tC.response)
			actual, err := client.Clients().UpdatePerson("client_id", "e@d.com", clients.Person{})
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

func TestClientsAPI_People(t *testing.T) {
	testCases := []struct {
		title                string
		forceHTTPClientError bool
		response             *http.Response
		expected             []*clients.PersonDetails
		expectedError        error
		oAuthAuthentication  bool
	}{
		{
			title: "no persons",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`[]`)),
			},
			expected: []*clients.PersonDetails{},
		},
		{
			title: "no persons and empty server response body",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(&bytes.Buffer{}),
			},
			expected: []*clients.PersonDetails{},
		},
		{
			title: "client with people",
			response: &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(`[
				{
					"EmailAddress": "e@d.com",
					"Name":         "name",
					"AccessLevel":  10,
					"Status": "status"
    			}
			]`)),
			},
			expected: []*clients.PersonDetails{
				{
					PersonBasicDetails: clients.PersonBasicDetails{
						EmailAddress: "e@d.com",
						Name:         "name",
						AccessLevel:  10,
					},
					Status: "status",
				},
			},
		},
		{
			title: "oAuth authentication",
			response: &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(`[
				{
					"EmailAddress": "e@d.com",
					"Name":         "name",
					"AccessLevel":  10,
					"Status": "status"
    			}
			]`)),
			},
			expected: []*clients.PersonDetails{
				{
					PersonBasicDetails: clients.PersonBasicDetails{
						EmailAddress: "e@d.com",
						Name:         "name",
						AccessLevel:  10,
					},
					Status: "status",
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
			httpClient.SetResponse("clients/client_id/people.json", tC.response)
			actual, err := client.Clients().People("client_id")
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

func TestClientsAPI_Person(t *testing.T) {
	testCases := []struct {
		title                string
		forceHTTPClientError bool
		response             *http.Response
		expected             *clients.PersonDetails
		expectedError        error
		oAuthAuthentication  bool
	}{
		{
			title: "no person",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{}`)),
			},
			expected: &clients.PersonDetails{},
		},
		{
			title: "no person and empty server response body",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(&bytes.Buffer{}),
			},
			expected: nil,
		},
		{
			title: "person found",
			response: &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(`
				{
					"EmailAddress": "e@d.com",
					"Name":         "name",
					"AccessLevel":  10,
					"Status": "status"
    			}
			`)),
			},
			expected: &clients.PersonDetails{
				PersonBasicDetails: clients.PersonBasicDetails{
					EmailAddress: "e@d.com",
					Name:         "name",
					AccessLevel:  10,
				},
				Status: "status",
			},
		},
		{
			title: "oAuth authentication",
			response: &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(`
				{
					"EmailAddress": "e@d.com",
					"Name":         "name",
					"AccessLevel":  10,
					"Status": "status"
    			}
			`)),
			},
			expected: &clients.PersonDetails{
				PersonBasicDetails: clients.PersonBasicDetails{
					EmailAddress: "e@d.com",
					Name:         "name",
					AccessLevel:  10,
				},
				Status: "status",
			},
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
			httpClient.SetResponse("clients/client_id/people.json", tC.response)
			actual, err := client.Clients().Person("client_id", "e@d.com")
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

func TestClientsAPI_DeletePerson(t *testing.T) {
	testCases := []struct {
		title                string
		forceHTTPClientError bool
		response             *http.Response
		expectedError        error
		oAuthAuthentication  bool
	}{
		{
			title: "successful deletion",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(&bytes.Buffer{}),
			},
		},
		{
			title: "oAuth authentication",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(&bytes.Buffer{}),
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
			httpClient.SetResponse("clients/client_id/people.json", tC.response)
			err := client.Clients().DeletePerson("client_id", "e@d.com")
			if err != nil {
				if !checkError(err, tC.expectedError) {
					t.Errorf("Expected '%v' error, actual: '%v'", tC.expectedError, err)
				}
				checkErrorType(t, err, !tC.forceHTTPClientError)
			}
		})
	}
}

func TestClientsAPI_SetPrimaryContact(t *testing.T) {
	testCases := []struct {
		title                string
		forceHTTPClientError bool
		response             *http.Response
		expectedError        error
		oAuthAuthentication  bool
		expected             string
	}{
		{
			title: "empty JSON response",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{}`)),
			},
			expected: "",
		},
		{
			title: "empty server response",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(&bytes.Buffer{}),
			},
			expected: "",
		},
		{
			title: "non empty server response",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"EmailAddress":"e@d.com"}`)),
			},
			expected: "e@d.com",
		},
		{
			title: "oAuth authentication",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"EmailAddress":"e@d.com"}`)),
			},
			expected:            "e@d.com",
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
			httpClient.SetResponse("clients/client_id/primarycontact.json", tC.response)
			actual, err := client.Clients().SetPrimaryContact("client_id", "e@d.com")
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

func TestClientsAPI_PrimaryContact(t *testing.T) {
	testCases := []struct {
		title                string
		forceHTTPClientError bool
		response             *http.Response
		expected             string
		expectedError        error
		oAuthAuthentication  bool
	}{
		{
			title: "no primary contact",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{}`)),
			},
			expected: "",
		},
		{
			title: "no primary contact and empty server response body",
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(&bytes.Buffer{}),
			},
			expected: "",
		},
		{
			title: "primary contact found",
			response: &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(`
				{
					"EmailAddress": "e@d.com"
    			}
			`)),
			},
			expected: "e@d.com",
		},
		{
			title: "oAuth authentication",
			response: &http.Response{
				StatusCode: 200,
				Body: ioutil.NopCloser(bytes.NewBufferString(`
				{
					"EmailAddress": "e@d.com"
    			}
			`)),
			},
			expected: "e@d.com",
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
			httpClient.SetResponse("clients/client_id/primarycontact.json", tC.response)
			actual, err := client.Clients().PrimaryContact("client_id")
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
