package hproblem

import "net/http"

// HTTP status code errors as registered with IANA.
// See: https://www.iana.org/assignments/http-status-codes/http-status-codes.xhtml
var (
	ErrStatusContinue                      = Errorf(100, http.StatusText(100)) // RFC 7231, 6.2.1
	ErrStatusSwitchingProtocols            = Errorf(101, http.StatusText(101)) // RFC 7231, 6.2.2
	ErrStatusProcessing                    = Errorf(102, http.StatusText(102)) // RFC 2518, 10.1
	ErrStatusEarlyHints                    = Errorf(103, http.StatusText(103)) // RFC 8297
	ErrStatusOK                            = Errorf(200, http.StatusText(200)) // RFC 7231, 6.3.1
	ErrStatusCreated                       = Errorf(201, http.StatusText(201)) // RFC 7231, 6.3.2
	ErrStatusAccepted                      = Errorf(202, http.StatusText(202)) // RFC 7231, 6.3.3
	ErrStatusNonAuthoritativeInfo          = Errorf(203, http.StatusText(203)) // RFC 7231, 6.3.4
	ErrStatusNoContent                     = Errorf(204, http.StatusText(204)) // RFC 7231, 6.3.5
	ErrStatusResetContent                  = Errorf(205, http.StatusText(205)) // RFC 7231, 6.3.6
	ErrStatusPartialContent                = Errorf(206, http.StatusText(206)) // RFC 7233, 4.1
	ErrStatusMultiStatus                   = Errorf(207, http.StatusText(207)) // RFC 4918, 11.1
	ErrStatusAlreadyReported               = Errorf(208, http.StatusText(208)) // RFC 5842, 7.1
	ErrStatusIMUsed                        = Errorf(226, http.StatusText(226)) // RFC 3229, 10.4.1
	ErrStatusMultipleChoices               = Errorf(300, http.StatusText(300)) // RFC 7231, 6.4.1
	ErrStatusMovedPermanently              = Errorf(301, http.StatusText(301)) // RFC 7231, 6.4.2
	ErrStatusFound                         = Errorf(302, http.StatusText(302)) // RFC 7231, 6.4.3
	ErrStatusSeeOther                      = Errorf(303, http.StatusText(303)) // RFC 7231, 6.4.4
	ErrStatusNotModified                   = Errorf(304, http.StatusText(304)) // RFC 7232, 4.1
	ErrStatusUseProxy                      = Errorf(305, http.StatusText(305)) // RFC 7231, 6.4.5
	ErrStatusTemporaryRedirect             = Errorf(307, http.StatusText(307)) // RFC 7231, 6.4.7
	ErrStatusPermanentRedirect             = Errorf(308, http.StatusText(308)) // RFC 7538, 3
	ErrStatusBadRequest                    = Errorf(400, http.StatusText(400)) // RFC 7231, 6.5.1
	ErrStatusUnauthorized                  = Errorf(401, http.StatusText(401)) // RFC 7235, 3.1
	ErrStatusPaymentRequired               = Errorf(402, http.StatusText(402)) // RFC 7231, 6.5.2
	ErrStatusForbidden                     = Errorf(403, http.StatusText(403)) // RFC 7231, 6.5.3
	ErrStatusNotFound                      = Errorf(404, http.StatusText(404)) // RFC 7231, 6.5.4
	ErrStatusMethodNotAllowed              = Errorf(405, http.StatusText(405)) // RFC 7231, 6.5.5
	ErrStatusNotAcceptable                 = Errorf(406, http.StatusText(406)) // RFC 7231, 6.5.6
	ErrStatusProxyAuthRequired             = Errorf(407, http.StatusText(407)) // RFC 7235, 3.2
	ErrStatusRequestTimeout                = Errorf(408, http.StatusText(408)) // RFC 7231, 6.5.7
	ErrStatusConflict                      = Errorf(409, http.StatusText(409)) // RFC 7231, 6.5.8
	ErrStatusGone                          = Errorf(410, http.StatusText(410)) // RFC 7231, 6.5.9
	ErrStatusLengthRequired                = Errorf(411, http.StatusText(411)) // RFC 7231, 6.5.10
	ErrStatusPreconditionFailed            = Errorf(412, http.StatusText(412)) // RFC 7232, 4.2
	ErrStatusRequestEntityTooLarge         = Errorf(413, http.StatusText(413)) // RFC 7231, 6.5.11
	ErrStatusRequestURITooLong             = Errorf(414, http.StatusText(414)) // RFC 7231, 6.5.12
	ErrStatusUnsupportedMediaType          = Errorf(415, http.StatusText(415)) // RFC 7231, 6.5.13
	ErrStatusRequestedRangeNotSatisfiable  = Errorf(416, http.StatusText(416)) // RFC 7233, 4.4
	ErrStatusExpectationFailed             = Errorf(417, http.StatusText(417)) // RFC 7231, 6.5.14
	ErrStatusTeapot                        = Errorf(418, http.StatusText(418)) // RFC 7168, 2.3.3
	ErrStatusMisdirectedRequest            = Errorf(421, http.StatusText(421)) // RFC 7540, 9.1.2
	ErrStatusUnprocessableEntity           = Errorf(422, http.StatusText(422)) // RFC 4918, 11.2
	ErrStatusLocked                        = Errorf(423, http.StatusText(423)) // RFC 4918, 11.3
	ErrStatusFailedDependency              = Errorf(424, http.StatusText(424)) // RFC 4918, 11.4
	ErrStatusTooEarly                      = Errorf(425, http.StatusText(425)) // RFC 8470, 5.2.
	ErrStatusUpgradeRequired               = Errorf(426, http.StatusText(426)) // RFC 7231, 6.5.15
	ErrStatusPreconditionRequired          = Errorf(428, http.StatusText(428)) // RFC 6585, 3
	ErrStatusTooManyRequests               = Errorf(429, http.StatusText(429)) // RFC 6585, 4
	ErrStatusRequestHeaderFieldsTooLarge   = Errorf(431, http.StatusText(431)) // RFC 6585, 5
	ErrStatusUnavailableForLegalReasons    = Errorf(451, http.StatusText(451)) // RFC 7725, 3
	ErrStatusInternalServerError           = Errorf(500, http.StatusText(500)) // RFC 7231, 6.6.1
	ErrStatusNotImplemented                = Errorf(501, http.StatusText(501)) // RFC 7231, 6.6.2
	ErrStatusBadGateway                    = Errorf(502, http.StatusText(502)) // RFC 7231, 6.6.3
	ErrStatusServiceUnavailable            = Errorf(503, http.StatusText(503)) // RFC 7231, 6.6.4
	ErrStatusGatewayTimeout                = Errorf(504, http.StatusText(504)) // RFC 7231, 6.6.5
	ErrStatusHTTPVersionNotSupported       = Errorf(505, http.StatusText(505)) // RFC 7231, 6.6.6
	ErrStatusVariantAlsoNegotiates         = Errorf(506, http.StatusText(506)) // RFC 2295, 8.1
	ErrStatusInsufficientStorage           = Errorf(507, http.StatusText(507)) // RFC 4918, 11.5
	ErrStatusLoopDetected                  = Errorf(508, http.StatusText(508)) // RFC 5842, 7.2
	ErrStatusNotExtended                   = Errorf(510, http.StatusText(510)) // RFC 2774, 7
	ErrStatusNetworkAuthenticationRequired = Errorf(511, http.StatusText(511)) // RFC 6585, 6
)
