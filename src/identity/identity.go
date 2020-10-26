package identity

import (
	msalgo "github.com/AzureAD/microsoft-authentication-library-for-go/src/msal"
	log "github.com/sirupsen/logrus"
)

const authority = "https://login.microsoftonline.com/73a99466-ad05-4221-9f90-e7142aa2f6c1/oauth2/v2.0/authorize"
const clientId = "824ce29e-a0bb-4ad5-9d29-90664358f3a2"
const notSoSecretSecret = "AxxIG_4O6feYwj_ueCN-6-ca8rN24U52sN"
var scopes []string = []string{"user.read"}

func tryClientSecretFlow(confidentialClientApp *msalgo.ConfidentialClientApplication) {
	clientSecretParams := msalgo.CreateAcquireTokenClientCredentialParameters(scopes)
	result, err := confidentialClientApp.AcquireTokenByClientCredential(clientSecretParams)
	if err != nil {
		log.Fatal(err)
	}
	accessToken := result.GetAccessToken()
	log.Info("Access token is: " + accessToken)
}

func AcquireTokenClientSecret() {
	secret, err := msalgo.CreateClientCredentialFromSecret(notSoSecretSecret)
	if err != nil {
		log.Fatal(err)
	}
	confidentialClientApp, err := msalgo.CreateConfidentialClientApplication(
		clientId, authority, secret)
	if err != nil {
		log.Fatal(err)
	}
	//Comes from sample_cache_accessor.go. Disabling cache for now.
	//confidentialClientApp.SetCacheAccessor(cacheAccessor)
	silentParams := msalgo.CreateAcquireTokenSilentParameters(scopes)
	result, err := confidentialClientApp.AcquireTokenSilent(silentParams)
	if err != nil {
		log.Info(err)
		tryClientSecretFlow(confidentialClientApp)
	} else {
		accessToken := result.GetAccessToken()
		log.Info("Access token is: " + accessToken)
	}

}