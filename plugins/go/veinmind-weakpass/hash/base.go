package hash

type Hash interface {
	// 密码匹配方式
	Match(hash, guess string) (flag bool, err error)

	// 加密算法的ID
	ID() string
}
