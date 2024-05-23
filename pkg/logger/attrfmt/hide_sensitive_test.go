package attrfmt_test

import (
	"testing"

	"golang.org/x/exp/slog"

	"bitbucket.org/jurnal/scm/pkg/logger/attrfmt"

	"github.com/stretchr/testify/assert"
)

type logValuer struct {
	v slog.Value
}

func (l logValuer) LogValue() slog.Value {
	return l.v
}

func TestHideSensitiveValueOnKeys(t *testing.T) {
	type TestCase struct {
		Name         string
		Keys         map[string]struct{}
		Group        []string
		InputAttr    slog.Attr
		ExpectedAttr slog.Attr
	}

	testCases := []TestCase{
		{
			Name: "empty attributes",
			Keys: map[string]struct{}{
				"password":      {},
				"client_secret": {},
				"access_token":  {},
				"refresh_token": {},
			},
			Group:        nil,
			InputAttr:    slog.Attr{},
			ExpectedAttr: slog.Attr{},
		},
		{
			Name: "no match value",
			Keys: map[string]struct{}{
				"password":      {},
				"client_secret": {},
				"access_token":  {},
				"refresh_token": {},
			},
			Group: nil,
			InputAttr: slog.Group("my-group",
				slog.String("a", "a value"),
				slog.String("b", "b value"),
				slog.String("c", "c value"),
				slog.String("d", "d value"),
			),
			ExpectedAttr: slog.Group("my-group",
				slog.String("a", "a value"),
				slog.String("b", "b value"),
				slog.String("c", "c value"),
				slog.String("d", "d value"),
			),
		},
		{
			Name: "group key name contains filtered key",
			Keys: map[string]struct{}{
				"password": {},
			},
			Group: nil,
			InputAttr: slog.Group("password",
				slog.String("k1", "actually this is not sensitive value"),
				slog.String("k2", "actually this is not sensitive value"),
				slog.String("k3", "actually this is not sensitive value"),
			),
			ExpectedAttr: slog.Group("password",
				slog.String("k1", "[FILTERED]"),
				slog.String("k2", "[FILTERED]"),
				slog.String("k3", "[FILTERED]"),
			),
		},
		{
			Name: "group value",
			Keys: map[string]struct{}{
				"password":      {},
				"client_secret": {},
				"access_token":  {},
				"refresh_token": {},
			},
			Group: nil,
			InputAttr: slog.Group("my-group",
				slog.String("password", "password value"),
				slog.String("client_secret", "client secret value"),
				slog.String("access_token", "access token value"),
				slog.String("refresh_token", "refresh token value"),
			),
			ExpectedAttr: slog.Group("my-group",
				slog.String("password", "[FILTERED]"),
				slog.String("client_secret", "[FILTERED]"),
				slog.String("access_token", "[FILTERED]"),
				slog.String("refresh_token", "[FILTERED]"),
			),
		},
		{
			Name: "deep group value",
			Keys: map[string]struct{}{
				"password":      {},
				"client_secret": {},
				"access_token":  {},
				"refresh_token": {},
			},
			Group: nil,
			InputAttr: slog.Group("my-group",
				slog.String("password", "password value"),
				slog.Group("second-group",
					slog.String("client_secret", "client secret value"),
					slog.Group("second-group",
						slog.String("access_token", "access token value"),
						slog.Group("third-group",
							slog.String("refresh_token", "refresh token value"),
						),
					),
				),
			),
			ExpectedAttr: slog.Group("my-group",
				slog.String("password", "[FILTERED]"),
				slog.Group("second-group",
					slog.String("client_secret", "[FILTERED]"),
					slog.Group("second-group",
						slog.String("access_token", "[FILTERED]"),
						slog.Group("third-group",
							slog.String("refresh_token", "[FILTERED]"),
						),
					),
				),
			),
		},
		{
			Name: "deep nested log valuer",
			Keys: map[string]struct{}{
				"password":      {},
				"client_secret": {},
				"access_token":  {},
				"refresh_token": {},
			},
			Group: nil,
			InputAttr: slog.Group("my-group",
				slog.Any("password", &logValuer{
					v: slog.AnyValue(&logValuer{
						v: slog.AnyValue(&logValuer{
							v: slog.AnyValue(&logValuer{
								v: slog.StringValue("my password"),
							}),
						}),
					}),
				}),
			),
			ExpectedAttr: slog.Group("my-group",
				slog.Any("password", slog.StringValue("[FILTERED]")),
			),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			actual := attrfmt.HideSensitiveValueOnKeys(tc.Keys)(tc.Group, tc.InputAttr)
			assert.EqualValues(t, tc.ExpectedAttr, actual)
		})
	}
}
