package deployer

type configStruct struct {
	AwsAccessID           string `json:"aws_access_id"`
	AwsSecretKey          string `json:"aws_secret_key"`
	DigitaloceanToken     string `json:"digitalocean_token"`
	AzureTenantID         string `json:"azure_tenant_id"`
	AzureClientID         string `json:"azure_client_id"`
	AzureClientSecret     string `json:"azure_client_secret"`
	AzureLocation         string `json:"azure_location"`
	AzureSubscriptionID   string `json:"azure_subscription_id"`
	GoogleCredentialsPath string `json:"google_credentials_path"`
	PublicKey             string `json:"public_key"`
	PrivateKey            string `json:"private_key"`
}
