package hash

type Hash interface {
	// ID 加密算法的ID
	ID() string

	// Match 密码匹配方式
	Match(hash, guess string) (flag bool, err error)
}
