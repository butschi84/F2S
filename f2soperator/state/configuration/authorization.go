package configuration

import "fmt"

type F2SPrivilege string

const (
	F2SPrivilegeFunctionsList   F2SPrivilege = "functions:list"
	F2SPrivilegeFunctionsCreate F2SPrivilege = "functions:create"
	F2SPrivilegeFunctionsInvoke F2SPrivilege = "functions:invoke"
	F2SPrivilegeFunctionsDelete F2SPrivilege = "functions:delete"
	F2SPrivilegeFunctionsUpdate F2SPrivilege = "functions:update"
	F2SPrivilegeSettingsView    F2SPrivilege = "settings:view"
	F2SPrivilegeSettingsUpdate  F2SPrivilege = "settings:update"
)

// resolve global priileges for a given group name
func ResolveGlobalPrivileges(groupName string) (globalPrivileges []string) {
	logging.Debug(fmt.Sprintf("request to get global privileges for group: %s", groupName))
	authorizationGroups := ActiveConfiguration.Config.F2S.Auth.Authorization
	for _, group := range authorizationGroups {
		if group.Group == groupName {
			logging.Debug(fmt.Sprintf("global privileges of group '%s': %s", groupName, group.Privileges))
			return group.Privileges
		}
	}

	// default return (no privileges)
	logging.Debug(fmt.Sprintf("found no gobal privilege definition for group: %s", groupName))
	return make([]string, 0)
}

// check if a group has a certain privilege globally
func HasGlobalPrivilege(privilege string, groupName string) bool {
	// get global privileges of the group
	groupsGlobalPrivileges := ResolveGlobalPrivileges(groupName)

	// check if specified privilege is in array
	for _, element := range groupsGlobalPrivileges {
		if element == privilege {
			return true
		}
	}
	return false
}
