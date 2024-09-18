package fetch

import (
	"fmt"
	"strings"

	"github.com/anti-raid/evil-befall/assets"
	"github.com/anti-raid/evil-befall/types/bigint"
	"github.com/anti-raid/evil-befall/types/bitflag"
	"github.com/anti-raid/evil-befall/types/silverpelt"
)

func getMapKeys[V any](m map[string]V) []string {
	var keys = make([]string, 0)

	for k := range m {
		keys = append(keys, k)
	}

	return keys
}

func fmtMap[K comparable, V any](m map[K]V) string {
	var sb strings.Builder

	sb.WriteString("{")

	var i = 0
	for k, v := range m {
		sb.WriteString(fmt.Sprintf("%v: %v", k, v))

		if i < len(m)-1 {
			sb.WriteString(", ")
		}

		i++
	}

	sb.WriteString("}")

	return sb.String()
}

type PermissionCheckFormatter struct {
	PermCheck silverpelt.PermissionCheck
}

func NewPermissionCheckFormatter(permCheck silverpelt.PermissionCheck) *PermissionCheckFormatter {
	return &PermissionCheckFormatter{PermCheck: permCheck}
}

func (pcf *PermissionCheckFormatter) NativePerms() []bigint.BigInt {
	return pcf.PermCheck.NativePerms
}

func (pcf *PermissionCheckFormatter) KittycatPerms() []string {
	return pcf.PermCheck.KittycatPerms
}

func (pcf *PermissionCheckFormatter) InnerAnd() bool {
	return pcf.PermCheck.InnerAnd
}

func (pcf *PermissionCheckFormatter) OuterAnd() bool {
	return pcf.PermCheck.OuterAnd
}

func (pcf *PermissionCheckFormatter) String() string {
	var result strings.Builder

	if len(pcf.NativePerms()) > 0 {
		result.WriteString("\t- Discord: ")
		for i, perm := range pcf.NativePerms() {
			if i != 0 {
				result.WriteString(" ")
			}
			permsBf := bitflag.NewBitFlag(assets.DiscordPermissions, perm.String())
			result.WriteString(fmt.Sprintf("%v (%s)", getMapKeys(permsBf.GetSetFlags()), perm.String()))
			if i < len(pcf.NativePerms())-1 {
				if pcf.InnerAnd() {
					result.WriteString(" AND")
				} else {
					result.WriteString(" OR")
				}
			}
		}
	}

	if len(pcf.KittycatPerms()) > 0 {
		result.WriteString("\n\t- Custom Permissions (kittycat): ")
		for i, perm := range pcf.KittycatPerms() {
			if i != 0 {
				result.WriteString(" ")
			}
			result.WriteString(perm)
			if i < len(pcf.KittycatPerms())-1 {
				if pcf.InnerAnd() {
					result.WriteString(" AND")
				} else {
					result.WriteString(" OR")
				}
			}
		}
	}

	return result.String()
}

type PermissionResultFormatter struct {
	Result silverpelt.PermissionResult
}

func NewPermissionResultFormatter(result silverpelt.PermissionResult) *PermissionResultFormatter {
	return &PermissionResultFormatter{Result: result}
}

func (prf *PermissionResultFormatter) ToMarkdown() string {
	switch prf.Result.Var {
	case "Ok":
		return "No message/context available"
	case "OkWithMessage":
		return prf.Result.Message
	case "MissingKittycatPerms", "MissingNativePerms", "MissingAnyPerms":
		if prf.Result.Check == nil {
			panic("Missing checks for permission result")
		}
		checksFmt1 := NewPermissionCheckFormatter(*prf.Result.Check)
		return fmt.Sprintf("You do not have the required permissions to perform this action. Try checking that you have the below permissions: %s", checksFmt1.String())
	case "CommandDisabled":
		return fmt.Sprintf("You cannot perform this action because the command ``%s`` is disabled on this server", prf.Result.CommandConfig.Command)
	case "UnknownModule":
		return fmt.Sprintf("The module ``%s`` does not exist", prf.Result.ModuleConfig.Module)
	case "ModuleNotFound":
		return "The module corresponding to this command could not be determined!"
	case "ModuleDisabled":
		return fmt.Sprintf("The module ``%s`` is disabled on this server", prf.Result.ModuleConfig.Module)
	case "NoChecksSucceeded":
		if prf.Result.Check == nil {
			panic("Missing checks for permission result")
		}
		checksFmt2 := NewPermissionCheckFormatter(*prf.Result.Check)
		return fmt.Sprintf("You do not have the required permissions to perform this action. You need at least one of the following permissions to execute this command:\n\n**Required Permissions**:\n\n%s", checksFmt2.String())
	case "DiscordError":
		return fmt.Sprintf("A Discord-related error seems to have occurred: %s.\n\nPlease try again later, it might work!", prf.Result.Error)
	case "SudoNotGranted":
		return "This module is only available for root (staff) and/or developers of the bot"
	case "GenericError":
		return prf.Result.Error
	default:
		return "Unknown error"
	}
}

type SettingsErrorFormatter struct {
	Error silverpelt.CanonicalSettingsError
}

func NewSettingsErrorFormatter(error silverpelt.CanonicalSettingsError) *SettingsErrorFormatter {
	return &SettingsErrorFormatter{Error: error}
}

func (sef *SettingsErrorFormatter) ToMarkdown() string {
	if sef.Error.Generic != nil {
		return fmt.Sprintf("An error occurred: `%s` from src `%s` of type `%s`", sef.Error.Generic.Message, sef.Error.Generic.Src, sef.Error.Generic.Typ)
	} else if sef.Error.OperationNotSupported != nil {
		return fmt.Sprintf("Operation `%s` is not supported", sef.Error.OperationNotSupported.Operation)
	} else if sef.Error.SchemaTypeValidationError != nil {
		return fmt.Sprintf("Column `%s` expected type `%s`, got type `%s`", sef.Error.SchemaTypeValidationError.Column, sef.Error.SchemaTypeValidationError.ExpectedType, sef.Error.SchemaTypeValidationError.GotType)
	} else if sef.Error.SchemaNullValueValidationError != nil {
		return fmt.Sprintf("Column `%s` is not nullable, yet value is null", sef.Error.SchemaNullValueValidationError.Column)
	} else if sef.Error.SchemaCheckValidationError != nil {
		return fmt.Sprintf("Column `%s` failed check `%s`, accepted range: `%s`, error: `%s`", sef.Error.SchemaCheckValidationError.Column, sef.Error.SchemaCheckValidationError.Check, sef.Error.SchemaCheckValidationError.AcceptedRange, sef.Error.SchemaCheckValidationError.Error)
	} else if sef.Error.MissingOrInvalidField != nil {
		return fmt.Sprintf("Missing (or invalid) field `%s` with src: `%s`", sef.Error.MissingOrInvalidField.Field, sef.Error.MissingOrInvalidField.Src)
	} else if sef.Error.RowExists != nil {
		return fmt.Sprintf("A row with the same column `%s` already exists. Count: `%d`", sef.Error.RowExists.ColumnId, sef.Error.RowExists.Count)
	} else if sef.Error.RowDoesNotExist != nil {
		return fmt.Sprintf("A row with the same column `%s` does not exist", sef.Error.RowDoesNotExist.ColumnId)
	} else if sef.Error.MaximumCountReached != nil {
		return fmt.Sprintf("The maximum number of entities this server may have (`%d`) has been reached. This server currently has `%d`.", sef.Error.MaximumCountReached.Max, sef.Error.MaximumCountReached.Current)
	} else {
		return fmt.Sprintf("Unknown error: %+v", sef.Error)
	}
}

func (sef *SettingsErrorFormatter) Code() string {
	if sef.Error.Generic != nil {
		return "Generic"
	} else if sef.Error.OperationNotSupported != nil {
		return "OperationNotSupported"
	} else if sef.Error.SchemaTypeValidationError != nil {
		return "SchemaTypeValidationError"
	} else if sef.Error.SchemaNullValueValidationError != nil {
		return "SchemaNullValueValidationError"
	} else if sef.Error.SchemaCheckValidationError != nil {
		return "SchemaCheckValidationError"
	} else if sef.Error.MissingOrInvalidField != nil {
		return "MissingOrInvalidField"
	} else if sef.Error.RowExists != nil {
		return "RowExists"
	} else if sef.Error.RowDoesNotExist != nil {
		return "RowDoesNotExist"
	} else if sef.Error.MaximumCountReached != nil {
		return "MaximumCountReached"
	} else {
		return "Unknown"
	}
}
