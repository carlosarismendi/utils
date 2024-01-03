package udatabase

import "github.com/carlosarismendi/utils/uerr"

var PqErrors = map[string]*uerr.UError{
	"unique_violation":   uerr.NewError(uerr.ResourceAlreadyExistsError, "Resource already exists."),
	"not_null_violation": uerr.NewError(uerr.WrongInputParameterError, "Missing required value."),
}
