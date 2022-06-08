package hash

type Hash interface {
	// 密码匹配方式
	Match(hash, guess string) (plain string, flag bool)
	// 获取密码明文
	Plain() (plain string)
	// 加密算法的ID
	ID() string
}
