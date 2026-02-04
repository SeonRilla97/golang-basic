package problem

const (
	BaseURI = "https://api.example.com/errors"

	TypeNotFound      = BaseURI + "/not-found"
	TypeBadRequest    = BaseURI + "/bad-request"
	TypeValidation    = BaseURI + "/validation"
	TypeUnauthorized  = BaseURI + "/unauthorized"
	TypeForbidden     = BaseURI + "/forbidden"
	TypeConflict      = BaseURI + "/conflict"
	TypeInternalError = BaseURI + "/internal-error"
	TypeRateLimited   = BaseURI + "/rate-limited"
)
