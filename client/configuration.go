package client

import (
	"github.com/spf13/viper"
    
)

const (

	//
	// All properties needed in the configuration file
	//

	//Paths to personal config files
	GPG_PUBRING_PATH string = "gpg_pubring_path"
	GPG_SECRING_PATH string = "gpg_secring_path"
	GPG_CONF_PATH    string = "gpg_conf_path"
	SSH_ID_PATH      string = "ssh_id_path"

	// GitHub Access
	GITHUB_LOGIN          string = "github_login"
	GITHUB_FULLNAME       string = "github_fullname"
	GITHUB_EMAIL          string = "github_email"
	GITHUB_SSH_PASSPHRASE string = "github_ssh_passphrase"

	//Nexus Access
	NEXUS_LOGIN          string = "exo_login"
	NEXUS_PASSWORD       string = "exo_password"
	NEXUS_GPG_PASSPHRASE string = "gpg_passphrase"
	NEXUS_GPG_KEYNAME    string = "gpg_keyname"

	// Jira Access
	JIRA_LOGIN    string = "exo_jira_login"
	JIRA_PASSWORD string = "exo_jira_password"

	// eXo Infra
	// git
	GIT_HOST string = "GIT_HOST"

	// Nexus eXo
	NEXUS_REPO_URL                   string = "NEXUS_REPO_URL"
	NEXUS_STAGING_SERVER_ID          string = "NEXUS_STAGING_SERVER_ID"
	NEXUS_STAGING_PROFILE_PRIVATE_ID string = "NEXUS_STAGING_PROFILE_PRIVATE_ID"
	NEXUS_STAGING_PROFILE_ADDONS_ID  string = "NEXUS_STAGING_PROFILE_ADDONS_ID"
	NEXUS_STAGING_PROFILE_PUBLIC_ID  string = "NEXUS_STAGING_PROFILE_PUBLIC_ID"

	// Nexus JBoss
	NEXUS_JBOSS_REPO_URL           string = "NEXUS_JBOSS_REPO_URL"
	NEXUS_JBOSS_STAGING_SERVER_ID  string = "NEXUS_JBOSS_STAGING_SERVER_ID"
	NEXUS_STAGING_PROFILE_JBOSS_ID string = "NEXUS_STAGING_PROFILE_JBOSS_ID"

	// eXo JIRA
	JIRA_API_URL     string = "JIRA_API_URL"
	JIRA_WEBHOOK_URL string = "JIRA_WEBHOOK_URL"

	// eXo Release Catalog
	CATALOG_BASE_URL string = "CATALOG_BASE_URL"

	EXOR_USER string = "exo-release"
)



// Read the config files and create string array to pass all as env values
func ConfigFileAsEnv() []string {

	return []string{
		GITHUB_LOGIN + "=" + viper.GetString(GITHUB_LOGIN), GITHUB_FULLNAME + "=" + viper.GetString(GITHUB_FULLNAME),
		GITHUB_EMAIL + "=" + viper.GetString(GITHUB_EMAIL), GITHUB_SSH_PASSPHRASE + "=" + viper.GetString(GITHUB_SSH_PASSPHRASE),
		NEXUS_LOGIN + "=" + viper.GetString(NEXUS_LOGIN), NEXUS_PASSWORD + "=" + viper.GetString(NEXUS_PASSWORD),
		NEXUS_GPG_PASSPHRASE + "=" + viper.GetString(NEXUS_GPG_PASSPHRASE), NEXUS_GPG_KEYNAME + "=" + viper.GetString(NEXUS_GPG_KEYNAME),
		JIRA_LOGIN + "=" + viper.GetString(JIRA_LOGIN), JIRA_PASSWORD + "=" + viper.GetString(JIRA_PASSWORD),
		GIT_HOST + "=" + viper.GetString(GIT_HOST), NEXUS_REPO_URL + "=" + viper.GetString(NEXUS_REPO_URL),
		NEXUS_STAGING_SERVER_ID + "=" + viper.GetString(NEXUS_STAGING_SERVER_ID), NEXUS_STAGING_PROFILE_PRIVATE_ID + "=" + viper.GetString(NEXUS_STAGING_PROFILE_PRIVATE_ID),
		NEXUS_STAGING_PROFILE_ADDONS_ID + "=" + viper.GetString(NEXUS_STAGING_PROFILE_ADDONS_ID), NEXUS_STAGING_PROFILE_PUBLIC_ID + "=" + viper.GetString(NEXUS_STAGING_PROFILE_PUBLIC_ID),
		NEXUS_JBOSS_REPO_URL + "=" + viper.GetString(NEXUS_JBOSS_REPO_URL), NEXUS_JBOSS_STAGING_SERVER_ID + "=" + viper.GetString(NEXUS_JBOSS_STAGING_SERVER_ID),
		JIRA_WEBHOOK_URL + "=" + viper.GetString(JIRA_WEBHOOK_URL), JIRA_API_URL + "=" + viper.GetString(JIRA_API_URL),
		NEXUS_STAGING_PROFILE_JBOSS_ID + "=" + viper.GetString(NEXUS_STAGING_PROFILE_JBOSS_ID), CATALOG_BASE_URL + "=" + viper.GetString(CATALOG_BASE_URL),
	}
}

func BindFiles(name string) []string {

	var binds = []string{
		viper.GetString(GPG_PUBRING_PATH) + ":/home/exo-release/.gnupg/pubring.gpg:ro",
		viper.GetString(GPG_SECRING_PATH) + ":/home/exo-release/.gnupg/secring.gpg:ro",
		viper.GetString(GPG_CONF_PATH) + ":/home/exo-release/.gnupg/gpg.conf:ro",
		viper.GetString(SSH_ID_PATH) + ":/home/exo-release/.ssh/id_rsa:ro",
		"workspace-" + name + ":/opt/plf-release/workspace",
	}

	return binds

}
