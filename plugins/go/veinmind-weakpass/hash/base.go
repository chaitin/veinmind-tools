package hash

type Hash interface {
	// Match 密码匹配方式
	Match(hash, guess string) (flag bool, err error)

	// ID 加密算法的ID
	ID() string
}
