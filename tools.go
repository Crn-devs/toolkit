package toolkit

import "crypto/rand"

const RandomStringSource = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789+_"

// Tools is the type used to instantiate this module, any variable of type Tools will have access
// to the reciever methods defined on the *Tools type

type Tools struct{}

// RandomString takes an integer n and returns a string of n random characters defined in the
// RandomStringSource
func (t *Tools) RandomString(n int) string {
	s, r := make([]rune, n), []rune(RandomStringSource)
	for i := range s {
		p, _ := rand.Prime(rand.Reader, len(r))
		x, y := p.Uint64(), uint64(len(r))
		s[i] = r[x%y]
	}
	return string(s)
}
