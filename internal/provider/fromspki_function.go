package provider

import (
	"context"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"golang.org/x/crypto/ssh"
)

type fromspkiFunction struct{}

func newFromspkiFunction() function.Function {
	return &fromspkiFunction{}
}

func (f *fromspkiFunction) Definition(ctx context.Context, req function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		MarkdownDescription: "Convert an X.509 SubjectPublicKeyInfo structure in PEM format into an SSH-compatible public key string. See provider documentation for details.",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:                "spki",
				MarkdownDescription: "X.509 SubjectPublicKeyInfo structure in PEM format",
			},
		},
		Return: function.StringReturn{},
	}
}

func (f *fromspkiFunction) Metadata(ctx context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "fromspki"
}

func (f *fromspkiFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var spki string

	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &spki))

	if resp.Error != nil {
		return
	}

	result, err := fromspki(spki)

	if err != nil {
		resp.Error = function.NewFuncError(err.Error())

		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, result))
}

type notPemError struct{}

func (notPemError) Error() string {
	return "expected input in PEM format"
}

type wrongBlockTypeError struct{}

func (wrongBlockTypeError) Error() string {
	return "PEM block type should be \"PUBLIC KEY\""
}

func fromspki(spki string) (string, error) {
	block, _ := pem.Decode([]byte(spki))

	if block == nil {
		return "", notPemError{}
	}

	if block.Type != "PUBLIC KEY" {
		return "", wrongBlockTypeError{}
	}

	x509pub, err := x509.ParsePKIXPublicKey(block.Bytes)

	if err != nil {
		return "", err
	}

	sshpub, err := ssh.NewPublicKey(x509pub)

	if err != nil {
		return "", err
	}

	result := sshpub.Type() + " " + base64.StdEncoding.EncodeToString(sshpub.Marshal())

	return result, nil
}
