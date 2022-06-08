package app

import "github.com/spf13/pflag"

// CliOptions 通过配置项读取命令行参数的接口.
type CliOptions interface {
	Flags() (fss NamedFlagSets)
	Validate() []error
}

// NamedFlagSets 存放 pflag.FlagSet 对象的存储结构.
type NamedFlagSets struct {
	// Order 记录 FlagSets 中 name 键的插入顺序.
	Order []string
	// FlagSets 以 name 为键存储 pflag.FlagSet 对象.
	FlagSets map[string]*pflag.FlagSet
}

// FlagSet 根据 name 返回 NamedFlagSets 中存放的 pflag.FlagSet 对象
// 如果 name 不存在则先创建再返回.
func (nfs *NamedFlagSets) FlagSet(name string) *pflag.FlagSet {
	if nfs.FlagSets == nil {
		nfs.FlagSets = map[string]*pflag.FlagSet{}
	}
	if _, ok := nfs.FlagSets[name]; !ok {
		nfs.FlagSets[name] = pflag.NewFlagSet(name, pflag.ExitOnError)
		nfs.Order = append(nfs.Order, name)
	}
	return nfs.FlagSets[name]
}
