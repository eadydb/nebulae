package testutil

// SetEnvs takes a map of key values to set using t.Setenv and restore
// the environment variable to its original value after the test.
func (t *T) SetEnvs(envs map[string]string) {
	for key, value := range envs {
		t.Setenv(key, value)
	}
}
