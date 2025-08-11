package provider

import (
	"testing"

	"gotest.tools/assert"
)

func TestFromspkiP256(t *testing.T) {
	pem := `
-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE3m3JQFAEc2DqBGhA/3K68q5HwjSp
KASwcbiLtht5ne7CmyVRUT7qKVCphcmm81Hy6bzUR6PZLaMhToq8dTnXAA==
-----END PUBLIC KEY-----
	`

	expect := "ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBN5tyUBQBHNg6gRoQP9yuvKuR8I0qSgEsHG4i7YbeZ3uwpslUVE+6ilQqYXJpvNR8um81Eej2S2jIU6KvHU51wA="

	ssh, err := fromspki(pem)
	assert.NilError(t, err)
	assert.Equal(t, ssh, expect)
}

func TestFromspkiRsa(t *testing.T) {
	pem := `
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAsdNeyJP9gOZKytYOeLin
c0oXmRqnir1yPOWnytv7EqKurWS5Tz8v3cvavTeTspFIVYfwvPiOuy1MSbOEw6c/
LcJpDyJDHmAnFHQ5sNA7rJy8DbOaJBXUK8z0FG/y0T639nVQVywfWJZVofdpbYl3
Dvh1E9rsvc4VDAxaXClZ0KfLOy5uJXrZyLUMMzNyaL9/L/eWQ+RNOOwiI0I+ExNB
4LVnI16TYbS0il4XrFHqoIlkgs45Inyae20DvqrYID6OrgGmE48Su77m81vddDRn
YNdCLePUAXIggDbrTOLR16mWfk5+140bTUyrsflQXqMxRIh/7x4RBT4fQWZMSh/4
iwIDAQAB
-----END PUBLIC KEY-----
	`

	expect := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCx017Ik/2A5krK1g54uKdzSheZGqeKvXI85afK2/sSoq6tZLlPPy/dy9q9N5OykUhVh/C8+I67LUxJs4TDpz8twmkPIkMeYCcUdDmw0DusnLwNs5okFdQrzPQUb/LRPrf2dVBXLB9YllWh92ltiXcO+HUT2uy9zhUMDFpcKVnQp8s7Lm4letnItQwzM3Jov38v95ZD5E047CIjQj4TE0HgtWcjXpNhtLSKXhesUeqgiWSCzjkifJp7bQO+qtggPo6uAaYTjxK7vubzW910NGdg10It49QBciCANutM4tHXqZZ+Tn7XjRtNTKux+VBeozFEiH/vHhEFPh9BZkxKH/iL"

	ssh, err := fromspki(pem)
	assert.NilError(t, err)
	assert.Equal(t, ssh, expect)
}

func TestInvalidInput(t *testing.T) {
	var err error

	_, err = fromspki("-----BEGIN CERTIFICATE-----\nAAAA\n-----END CERTIFICATE-----")
	assert.ErrorType(t, err, wrongBlockTypeError{})

	_, err = fromspki("asdf")
	assert.ErrorType(t, err, notPemError{})
}
