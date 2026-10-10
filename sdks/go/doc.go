// Package alldebrid is a client for the AllDebrid API (v4 and v4.1).
//
// The models in types.gen.go are generated from openapi.yaml at the root of
// the repository; the client methods are hand-written on top of them, one per
// operation of the spec.
//
//	client := alldebrid.New(os.Getenv("ALLDEBRID_API_KEY"), alldebrid.WithUserAgent("my-app"))
//	user, err := client.User(ctx)
//
// API errors are *Error values; compare their code with errors.Is:
//
//	_, err := client.Magnet(ctx, 42)
//	if errors.Is(err, alldebrid.CodeMagnetInvalidID) { ... }
//
// Without an account, use one of the demo keys (DemoPremiumKey, DemoTrialKey,
// DemoFreeKey): the API answers them with static sample data.
package alldebrid

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.8.0 -config oapi-codegen.yaml ../../openapi.yaml
